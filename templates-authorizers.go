package main

import (
	"bytes"
	"fmt"
	"github.com/Odania-IT/terraless/schema"
	"github.com/Odania-IT/terraless/support"
	"github.com/Odania-IT/terraless/templates"
)

const (
	helperFunctionCommand = `
terraless auth --auth-provider {{ .TeamName }}:{{ .ProfileName }}:role={{ .RoleName }}
export AWS_PROFILE={{ .AwsProfile }}
`
)

func (provider *ProviderAws) RenderAuthorizerTemplates(config schema.TerralessConfig) string {
	var buffer bytes.Buffer
	for _, authorizer := range config.Authorizers {
		if authorizer.Type == "aws" {
			logger.Debug(fmt.Sprintf("Generating authorizer template for %s\n", authorizer.Name))
			authorizer.TerraformName = "terraless-authorizer-" + support.SanitizeString(authorizer.Name)

			buffer = templates.RenderTemplateToBuffer(authorizer, buffer, awsTemplates("authorizer.tf.tmpl"), "terraless-authorizer")
		}
	}

	return buffer.String()
}

func (provider *ProviderAws) GenerateHelperFunctionCommand(teamName string, ProviderName string, roleName string) string {
	buffer := bytes.Buffer{}

	data := map[string]string{
		"AwsProfile": ProviderName + "-" + roleName,
		"TeamName": teamName,
		"ProfileName": ProviderName,
		"RoleName": roleName,
	}

	buffer = templates.RenderTemplateToBuffer(data, buffer, helperFunctionCommand, "terraless-helper-function")

	return buffer.String()
}
