package main

import (
	"github.com/Odania-IT/terraless/schema"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestTerralessFunctions_RenderCertificateTemplates(t *testing.T) {
	// given
	provider := ProviderAws{}
	config := schema.TerralessConfig{
		Certificates: map[string]schema.TerralessCertificate{
			"TestCert": {
				Type: "aws",
				Domain: "example.com",
				ZoneId: "my-test-zone",
			},
		},
	}

	// when
	result := provider.RenderCertificateTemplates(config)

	// then
	expected := `## Terraless Certificate

resource "aws_acm_certificate" "terraless-certificate-example-com" {
  domain_name = "example.com"
  validation_method = "DNS"
  subject_alternative_names = [
    
  ]

  provider = "aws.us-east"

  tags = "${var.terraless-default-tags}"
}


resource "aws_route53_record" "terraless-certificate-example-com-validation" {
  name = "${aws_acm_certificate.terraless-certificate-example-com.domain_validation_options.0.resource_record_name}"
  type = "${aws_acm_certificate.terraless-certificate-example-com.domain_validation_options.0.resource_record_type}"
  zone_id = "my-test-zone"
  records = [
    "${aws_acm_certificate.terraless-certificate-example-com.domain_validation_options.0.resource_record_value}"
  ]
  ttl = 60
}

resource "aws_acm_certificate_validation" "terraless-certificate-example-com-validation" {
  certificate_arn = "${aws_acm_certificate.terraless-certificate-example-com.arn}"
  validation_record_fqdns = [
    "${aws_route53_record.terraless-certificate-example-com-validation.fqdn}"
  ]

  provider = "aws.us-east"
  depends_on = [
    "aws_route53_record.terraless-certificate-example-com-validation",
    "aws_acm_certificate.terraless-certificate-example-com"
  ]
}

`
	assert.Equal(t, expected, result)
}
