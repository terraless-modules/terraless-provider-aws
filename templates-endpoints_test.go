package main

import (
	"github.com/Odania-IT/terraless/schema"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestTemplatesFunctions_RenderEndpointTemplates(t *testing.T) {
	// given
	provider := ProviderAws{}
	config := schema.TerralessConfig{
			ProjectName: "DummyProjectName",
			Endpoints: []schema.TerralessEndpoint{
				{
					Type: "dummy",
					Domain: "my-secret-dummy-domain.org",
				},
				{
					Type: "apigateway",
					Domain: "my-secret-domain.org",
				},
			},
		}

	// when
	result := provider.RenderEndpointTemplates(config)

	// then
	assert.Contains(t, result, `domain_name     = "my-secret-domain.org"`)
	assert.Contains(t, result, `resource "aws_api_gateway_base_path_mapping" "terraless-endpoint-my-secret-domain-org"`)
	assert.Contains(t, result, `resource "aws_route53_record" "terraless-endpoint-my-secret-domain-org"`)
}
