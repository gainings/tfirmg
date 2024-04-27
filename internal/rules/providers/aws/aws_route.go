package aws

import (
	"fmt"
	"github.com/gainings/tfirg/internal/model/resourceid"
	"github.com/gainings/tfirg/internal/rules"
)

func init() {
	rules.Register("aws_route", AWSRouteRule{})
}

type AWSRouteRule struct{}

func (r AWSRouteRule) Generate(attributes map[string]interface{}) *resourceid.ResourceID {
	id := resourceid.ResourceID(fmt.Sprintf("%s_%s", attributes["route_table_id"].(string), attributes["destination_cidr_block"].(string)))
	return &id
}
