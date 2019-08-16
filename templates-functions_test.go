package main

import (
	"github.com/Odania-IT/terraless/schema"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestTemplatesFunctions_RenderFunctionTemplates_DoesNotHandleWrongType(t *testing.T) {
	// given
	provider := ProviderAws{}

	// when
	result := provider.RenderFunctionTemplates("dummy", schema.FunctionEvents{}, &schema.TerralessData{})

	// then
	expected := ``
	assert.Equal(t, expected, result)
}

func TestTemplatesFunctions_RenderFunctionTemplates_HttpEvent(t *testing.T) {
	// given
	provider := ProviderAws{}
	functionEvents := schema.FunctionEvents{
		Events: map[string][]schema.FunctionEvent{
			"http": {
				{
					FunctionName: "DummyFunction",
					FunctionEvent: schema.TerralessFunctionEvent{
						Type: "http",
						Path: "dummy",
					},
				},
			},
		},
	}
	terralessData := schema.TerralessData{
		Arguments: schema.Arguments{
			Environment: "DummyEnvironment",
		},
		Config: schema.TerralessConfig{
			ProjectName: "DummyProjectName",
		},
	}

	// when
	result := provider.RenderFunctionTemplates("aws", functionEvents, &terralessData)

	// then
	expected := `resource "aws_api_gateway_rest_api" "terraless-api-gateway" {
  name        = "DummyProjectName-DummyEnvironment"
  description = "Terraless Api Gateway for DummyProjectName-DummyEnvironment"
}`
	assert.Contains(t, result, expected)
	assert.Contains(t, result, `output "api-gateway-invoke-url"`)
	assert.Contains(t, result, `resource "aws_api_gateway_resource" "terraless-lambda-DummyFunction-0"`)
	assert.Contains(t, result, `resource "aws_cloudwatch_log_group" "lambda-log-DummyFunction"`)
}

func TestTemplatesFunctions_RenderFunctionTemplates_SqsEvent(t *testing.T) {
	// given
	provider := ProviderAws{}
	functionEvents := schema.FunctionEvents{
		Events: map[string][]schema.FunctionEvent{
			"sqs": {
				{
					FunctionName: "SpecificFunction",
					FunctionEvent: schema.TerralessFunctionEvent{
						Type: "sqs",
						Arn: "arn:aws::::sqs",
					},
				},
			},
		},
	}
	terralessData := schema.TerralessData{
		Arguments: schema.Arguments{
			Environment: "DummyEnvironment",
		},
		Config: schema.TerralessConfig{
			ProjectName: "DummyProjectName",
		},
	}

	// when
	result := provider.RenderFunctionTemplates("aws", functionEvents, &terralessData)

	// then
	assert.Contains(t, result, `# Function DummyProjectName SpecificFunction EventKey: 0`)
}
