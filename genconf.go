package genconf

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// Load a directory of JSON config into a target struct. Which config is loaded will
// depend on the provided env. Conflicting config will be merged from less specific to
// more specific, and then by the provided order.
func Load(root string, order []string, env map[string]string, target interface{}) error {
	filePaths, err := loadFilePaths(root)
	if err != nil {
		return err
	}

	overrides := make([]override, 0, len(filePaths))
	for _, path := range filePaths {
		tuple, err := parseTuple(root, path)
		if err != nil {
			return err
		}
		// Only apply overrides which match the provided env.
		if subset(env, tuple) {
			overrides = append(overrides, override{tuple: tuple, path: path})
		}
	}

	orderLookup := map[string]int{}
	for i, key := range order {
		orderLookup[key] = i
	}

	// Sort by cardinality of dimensions first (where "default" has a length of 0 and is
	// therefore going to be applied first).
	// Then sort by provided order if cardinality is equal.
	sort.Slice(overrides, func(i, j int) bool {
		if len(overrides[i].tuple) == len(overrides[j].tuple) {
			iKeys := make([]int, 0, len(overrides[i].tuple))
			for key := range overrides[i].tuple {
				iKeys = append(iKeys, orderLookup[key])
			}
			jKeys := make([]int, 0, len(overrides[j].tuple))
			for key := range overrides[j].tuple {
				jKeys = append(jKeys, orderLookup[key])
			}
			sort.Ints(iKeys)
			sort.Ints(jKeys)
			for k := 0; k < len(iKeys); k++ {
				if iKeys[k] == jKeys[k] {
					continue
				}
				return iKeys[k] < jKeys[k]
			}
			return false
		}
		return len(overrides[i].tuple) < len(overrides[j].tuple)
	})

	// This merge is relatively naive. Alternatives would include recursion, and possibly
	// supporting other encodings like YAML or protobuf.
	for _, override := range overrides {
		data, err := ioutil.ReadFile(override.path)
		if err != nil {
			return err
		}
		err = json.Unmarshal(data, target)
		if err != nil {
			return err
		}
	}

	return nil
}

func loadFilePaths(root string) ([]string, error) {
	filePaths := []string{}
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		filePaths = append(filePaths, path)
		return nil
	})
	if err != nil {
		return nil, err
	}
	return filePaths, nil
}

type override struct {
	tuple map[string]string
	path  string
}

func parseTuple(root, path string) (map[string]string, error) {
	rel, err := filepath.Rel(root, path)
	if err != nil {
		return nil, err
	}
	rel = rel[0 : len(rel)-len(filepath.Ext(rel))]
	if rel == "default" {
		return map[string]string{}, nil
	}
	parts := strings.Split(rel, "/")
	if len(parts)%2 != 0 {
		return nil, fmt.Errorf("Path '%s' must be even", path)
	}
	tuple := map[string]string{}
	for i := 0; i < len(parts); i += 2 {
		tuple[parts[i]] = parts[i+1]
	}
	return tuple, nil
}

func subset(env, tuple map[string]string) bool {
	for key, value := range tuple {
		v, ok := env[key]
		if !ok || value != v {
			return false
		}
	}
	return true
}
