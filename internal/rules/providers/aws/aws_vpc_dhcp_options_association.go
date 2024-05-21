package aws

import (
	"github.com/gainings/tfirmg/internal/model/resourceid"
	"github.com/gainings/tfirmg/internal/rules"
)

func init() {
	rules.Register("aws_vpc_dhcp_options_association", AWSDhcpOptionsAssociationRule{})
}

type AWSDhcpOptionsAssociationRule struct{}

func (r AWSDhcpOptionsAssociationRule) Generate(attributes map[string]interface{}) *resourceid.ResourceID {
	id := resourceid.ResourceID(attributes["vpc_id"].(string))
	return &id
}
