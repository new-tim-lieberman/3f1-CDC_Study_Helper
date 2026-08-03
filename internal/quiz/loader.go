package quiz

import (
	"embed"
	"encoding/json"
	"fmt"
	"sort"
)

// Drop additional module files (e.g. module2.json) in data/ to add more study sets.
//
//go:embed data/*.json
var dataFS embed.FS

// LoadModules reads every embedded module file and returns them sorted by ID.
func LoadModules() ([]Module, error) {
	entries, err := dataFS.ReadDir("data")
	if err != nil {
		return nil, err
	}

	var modules []Module
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		raw, err := dataFS.ReadFile("data/" + e.Name())
		if err != nil {
			return nil, fmt.Errorf("reading %s: %w", e.Name(), err)
		}
		var m Module
		if err := json.Unmarshal(raw, &m); err != nil {
			return nil, fmt.Errorf("parsing %s: %w", e.Name(), err)
		}
		for i := range m.Questions {
			m.Questions[i].ModuleID = m.ID
		}
		modules = append(modules, m)
	}

	sort.Slice(modules, func(i, j int) bool { return modules[i].ID < modules[j].ID })
	return modules, nil
}