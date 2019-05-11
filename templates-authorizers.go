package main

import (
	"bytes"
	"github.com/Odania-IT/terraless/schema"
	"github.com/Odania-IT/terraless/support"
	"github.com/Odania-IT/terraless/templates"
)

func (provider *ProviderAws) RenderAuthorizerTemplates(config schema.TerralessConfig, buffer bytes.Buffer) bytes.Buffer {
	for _, authorizer := range config.Authorizers {
		if authorizer.Type == "aws" {
			logger.Debug("Generating authorizer template for %s\n", authorizer.Name)
			authorizer.TerraformName = "terraless-authorizer-" + support.SanitizeString(authorizer.Name)

			buffer = templates.RenderTemplateToBuffer(authorizer, buffer, awsTemplates("authorizer.tf.tmpl"), "terraless-authorizer")
		}
	}

	return buffer
}
