package parser

import (
	"errors"
	"fmt"
	"io"
	"sort"

	"gopkg.in/yaml.v3"
)

// Workload is the normalised view of one resource-bearing unit found in a
// Helm values file. It is intentionally cloud-agnostic and decoupled from
// the Kubernetes API: the same shape supports `Deployment`, `StatefulSet`,
// `CronJob`, etc.
type Workload struct {
	Name     string       // YAML key path joined by "." (e.g. "api", "subchart.worker")
	Kind     string       // best-effort: "Deployment" by default
	Requests ResourceList //
	Limits   ResourceList //
}

// ResourceList captures the CPU and memory of either requests or limits.
type ResourceList struct {
	CPU    Quantity
	Memory Quantity
}

// ParseValues reads a Helm values.yaml stream and returns every workload
// candidate found. A "workload" is any nested map containing a
// `resources` key with `requests` and/or `limits`.
//
// Phase 1 supports the common pattern where chart authors expose a
// `<workload>.resources` block per service. It does not (yet) render
// templates or evaluate functions; sub-chart and template support land
// in Phase 2.
func ParseValues(r io.Reader) ([]Workload, error) {
	raw, err := io.ReadAll(r)
	if err != nil {
		return nil, fmt.Errorf("parser: read: %w", err)
	}
	var root yaml.Node
	if err := yaml.Unmarshal(raw, &root); err != nil {
		return nil, fmt.Errorf("parser: yaml: %w", err)
	}
	if root.Kind != yaml.DocumentNode || len(root.Content) == 0 {
		return nil, errors.New("parser: empty document")
	}
	if root.Content[0].Kind != yaml.MappingNode {
		return nil, errors.New("parser: top-level must be a map")
	}

	var workloads []Workload
	walk(root.Content[0], "", &workloads)
	sort.Slice(workloads, func(i, j int) bool { return workloads[i].Name < workloads[j].Name })
	return workloads, nil
}

// walk descends the YAML tree, emitting a Workload whenever it sees a
// `resources` map under any named key.
func walk(n *yaml.Node, path string, out *[]Workload) {
	if n.Kind != yaml.MappingNode {
		return
	}
	for i := 0; i+1 < len(n.Content); i += 2 {
		k, v := n.Content[i], n.Content[i+1]
		if k.Kind != yaml.ScalarNode {
			continue
		}
		childPath := joinPath(path, k.Value)

		// Treat any mapping with a `resources.requests` or
		// `resources.limits` child as a workload.
		if v.Kind == yaml.MappingNode {
			if res := findChild(v, "resources"); res != nil && res.Kind == yaml.MappingNode {
				wl := Workload{Name: childPath, Kind: "Deployment"}
				if reqs := findChild(res, "requests"); reqs != nil {
					wl.Requests = readResourceList(reqs)
				}
				if lims := findChild(res, "limits"); lims != nil {
					wl.Limits = readResourceList(lims)
				}
				*out = append(*out, wl)
			}
			walk(v, childPath, out)
		}
	}
}

// readResourceList extracts cpu/memory under a requests/limits mapping.
// Unparseable values are dropped silently — the rule engine surfaces
// missing fields, not malformed ones (those are a chart-author bug).
func readResourceList(n *yaml.Node) ResourceList {
	var rl ResourceList
	if n.Kind != yaml.MappingNode {
		return rl
	}
	for i := 0; i+1 < len(n.Content); i += 2 {
		k, v := n.Content[i], n.Content[i+1]
		if k.Kind != yaml.ScalarNode || v.Kind != yaml.ScalarNode {
			continue
		}
		switch k.Value {
		case "cpu":
			if q, err := ParseCPU(v.Value); err == nil {
				rl.CPU = q
			}
		case "memory":
			if q, err := ParseMemory(v.Value); err == nil {
				rl.Memory = q
			}
		}
	}
	return rl
}

func findChild(n *yaml.Node, key string) *yaml.Node {
	if n.Kind != yaml.MappingNode {
		return nil
	}
	for i := 0; i+1 < len(n.Content); i += 2 {
		k, v := n.Content[i], n.Content[i+1]
		if k.Kind == yaml.ScalarNode && k.Value == key {
			return v
		}
	}
	return nil
}

func joinPath(parent, child string) string {
	if parent == "" {
		return child
	}
	return parent + "." + child
}
