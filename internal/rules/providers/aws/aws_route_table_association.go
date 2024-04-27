package aws

import (
	"fmt"
	"github.com/gainings/tfirmg/internal/model/resourceid"
	"github.com/gainings/tfirmg/internal/rules"
)

func init() {
	rules.Register("aws_route_table_association", AWSRouteTableAssociationRule{})
}

type AWSRouteTableAssociationRule struct{}

func (r AWSRouteTableAssociationRule) Generate(attributes map[string]interface{}) *resourceid.ResourceID {
	id := resourceid.ResourceID(fmt.Sprintf("%s/%s", attributes["subnet_id"].(string), attributes["route_table_id"].(string)))
	return &id
}
