package resource

import (
	"fmt"
	"github.com/gainings/tfirg/internal/model/resourceid"
	"strings"
)

type Resource struct {
	Address  Address
	Name     string
	ID       *resourceid.ResourceID
	IndexKey string
	Module   *struct {
		Name string
	}
}

type Resources []Resource

type Address string

func (a Address) String() string {
	return string(a)
}

type ResourceCreator struct {
	rig resourceid.GeneratorFactory
}

func NewResourceCreator(rig resourceid.GeneratorFactory) ResourceCreator {
	return ResourceCreator{
		rig: rig,
	}
}

func (rf ResourceCreator) Create(typeName, resourceName, indexKey string, moduleName *string, attributes map[string]interface{}) Resource {
	address := rf.newAddress(typeName, resourceName, indexKey, moduleName)
	rID := rf.rig.Generator(typeName).Generate(attributes)
	r := Resource{
		Address:  address,
		Name:     fmt.Sprintf("%s.%s", typeName, resourceName),
		ID:       rID,
		IndexKey: indexKey,
	}
	if moduleName != nil {
		mn := strings.Split(*moduleName, ".")
		r.Module = &struct{ Name string }{Name: mn[0]}
	}
	return r
}

func (rf ResourceCreator) newAddress(typeName, resourceName, indexKey string, moduleName *string) Address {
	if moduleName != nil {
		return Address(fmt.Sprintf("%s.%s.%s", *moduleName, typeName, resourceName))
	} else if indexKey != "" {
		return Address(fmt.Sprintf("%s.%s[%s]", typeName, resourceName, indexKey))
	} else {
		return Address(fmt.Sprintf("%s.%s", typeName, resourceName))
	}
}
