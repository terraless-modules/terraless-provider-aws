package main

import (
	"github.com/Odania-IT/terraless/schema"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestTemplatesFunctions_RenderAuthorizerTemplates(t *testing.T) {
	// given
	provider := ProviderAws{}
	config := schema.TerralessConfig{
		ProjectName: "DummyProjectName",
		Authorizers: map[string]schema.TerralessAuthorizer{
			"unsed": {
				Type: "dummy",
				Name: "UnusupportedAuthorizer",
			},
			"dummy": {
				Type: "aws",
				Name: "SupportedAuthorizer",
				ProviderArns: []string{
					"arn1",
					"arn2",
				},
			},
		},
	}

	// when
	result := provider.RenderAuthorizerTemplates(config)

	// then
	assert.Contains(t, result, `## Terraless Authorizer`)
	assert.Contains(t, result, `resource "aws_api_gateway_authorizer" "terraless-authorizer-SupportedAuthorizer"`)
	assert.Contains(t, result, `"arn1"`)
	assert.Contains(t, result, `"arn2"`)
}
