package aws

import (
	"github.com/gainings/tfirmg/internal/model/resourceid"
	"github.com/gainings/tfirmg/internal/rules"
)

func init() {
	rules.Register("aws_default_route_table", AWSDefaultRouteTable{})
}

type AWSDefaultRouteTable struct{}

func (r AWSDefaultRouteTable) Generate(attributes map[string]interface{}) *resourceid.ResourceID {
	//import unsupported
	return nil
}
