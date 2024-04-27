package aws

import (
	"github.com/gainings/tfirmg/internal/model/resourceid"
	"github.com/gainings/tfirmg/internal/rules"
)

func init() {
	rules.Register("aws_acm_certificate_validation", AWSAcmCertificateValidationRule{})
}

type AWSAcmCertificateValidationRule struct{}

func (r AWSAcmCertificateValidationRule) Generate(attributes map[string]interface{}) *resourceid.ResourceID {
	//import unsupported
	return nil
}
