package main

import (
	"bytes"
	"github.com/Odania-IT/terraless/schema"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestTemplatesFunctions_Provider(t *testing.T) {
	// given

	// when
	provider := ProviderAws{}

	// then
	assert.Equal(t, true, provider.CanHandle("aws"))
	assert.Equal(t, false, provider.CanHandle("aws2"))
	assert.Equal(t, false, provider.CanHandle("dummy"))
	assert.Equal(t, "terraless-provider-aws", provider.Info("none").Name)
}

func TestTemplatesFunctions_Info(t *testing.T) {
	// given
	provider := ProviderAws{}

	// when
	pluginInfo := provider.Info("debug")

	// then
	assert.Equal(t, "terraless-provider-aws", pluginInfo.Name)
	assert.Equal(t, VERSION, pluginInfo.Version)
}

func TestTemplatesFunctions_AwsTemplates(t *testing.T) {
	// given

	// when
	template := awsTemplates("iam.tf.tmpl")

	// then
	assert.Contains(t, template, `resource "aws_iam_role" "terraless-lambda-iam-role"`)
}

func TestTemplatesFunctions_FinalizeTemplates(t *testing.T) {
	// given
	provider := ProviderAws{}
	addTerralessLambdaRole = true
	terralessData := schema.TerralessData{
		Arguments: schema.Arguments{
			Environment: "DummyEnvironment",
		},
		Config: schema.TerralessConfig{
			ProjectName: "DummyProjectName",
		},
	}

	// when
	result := provider.FinalizeTemplates(terralessData)

	// then
	assert.Contains(t, result, `resource "aws_iam_role" "terraless-lambda-iam-role"`)
}
