package modulemeta

import (
	"encoding/json"
	"io"
	"strings"
)

type Module struct {
	Key     string `json:"Key"`
	Source  string `json:"Source"`
	Version string `json:"Version,omitempty"`
	Dir     string `json:"Dir"`
}

type ModuleMeta struct {
	Modules []Module `json:"Modules"`
}

func Decode(r io.Reader) (*ModuleMeta, error) {
	mm := &ModuleMeta{}
	decoder := json.NewDecoder(r)
	if err := decoder.Decode(mm); err != nil {
		return nil, err
	}
	return mm, nil
}

func (mm ModuleMeta) GetModuleMap() map[string]Module {
	mMap := make(map[string]Module)
	for _, m := range mm.Modules {
		if m.Key != "" {
			mMap[mm.format(m.Key)] = m
		}
	}
	return mMap
}
func (mm ModuleMeta) format(key string) string {
	parts := strings.Split(key, ".")
	return strings.Join(parts, ".module.")
}
