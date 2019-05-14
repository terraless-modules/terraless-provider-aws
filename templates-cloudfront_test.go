package main

import (
	"github.com/Odania-IT/terraless/schema"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestTemplatesFunctions_RenderUploadTemplates(t *testing.T) {
	// given
	provider := ProviderAws{}
	terralessData := schema.TerralessData{
		Config: schema.TerralessConfig{
			ProjectName: "DummyProjectName",
			Uploads: []schema.TerralessUpload{
				{
					Type: "s3",
					Cloudfront: schema.TerralessCloudfront{
						Handler: "singleEntryPointHandler",
					},
				},
			},
		},
	}

	// when
	result := provider.RenderUploadTemplates(terralessData)

	// then
	assert.NotContains(t, result, `## Terraless Lambda@Edge`)
	assert.NotContains(t, result, `resource "aws_cloudwatch_log_group" "lambda-log-terraless-lambda-cloudfront"`)
	assert.NotContains(t, result, `resource "aws_lambda_function" "terraless-lambda-cloudfront"`)
}

func TestTemplatesFunctions_RenderUploadTemplates_WithDomain(t *testing.T) {
	// given
	provider := ProviderAws{}
	terralessData := schema.TerralessData{
		Config: schema.TerralessConfig{
			ProjectName: "DummyProjectName",
			Uploads: []schema.TerralessUpload{
				{
					Type: "s3",
					Cloudfront: schema.TerralessCloudfront{
						Domain: "my-dummy-domain.org",
						Handler: "redirectToWww",
					},
				},
			},
		},
	}

	// when
	result := provider.RenderUploadTemplates(terralessData)

	// then
	assert.Contains(t, result, `## Terraless Lambda@Edge`)
	assert.Contains(t, result, `resource "aws_cloudwatch_log_group" "lambda-log-terraless-lambda-cloudfront"`)
	assert.Contains(t, result, `resource "aws_lambda_function" "terraless-lambda-cloudfront"`)
	assert.Contains(t, result, `resource "aws_cloudfront_origin_access_identity" "terraless-default"`)
	assert.Contains(t, result, `resource "aws_cloudfront_distribution" "terraless-default"`)
}
