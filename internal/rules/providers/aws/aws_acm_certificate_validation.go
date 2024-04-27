package aws

import (
	"github.com/gainings/tfirg/internal/model/resourceid"
	"github.com/gainings/tfirg/internal/rules"
)

func init() {
	rules.Register("aws_acm_certificate_validation", AWSAcmCertificateValidationRule{})
}

type AWSAcmCertificateValidationRule struct{}

func (r AWSAcmCertificateValidationRule) Generate(attributes map[string]interface{}) *resourceid.ResourceID {
	//import unsupported
	return nil
}
