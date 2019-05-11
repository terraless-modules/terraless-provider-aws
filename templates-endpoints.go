package main

import (
	"bytes"
	"github.com/Odania-IT/terraless/schema"
	"github.com/Odania-IT/terraless/support"
	"github.com/Odania-IT/terraless/templates"
)

func (provider *ProviderAws) RenderEndpointTemplates(config schema.TerralessConfig, buffer bytes.Buffer) bytes.Buffer {
	for _, endpoint := range config.Endpoints {
		if endpoint.Type == "apigateway" {
			logger.Debug("Generating certificate template for %s\n", endpoint.Domain)
			endpoint.TerralessCertificate = config.Certificates[endpoint.Certificate]
			endpoint.TerraformName = "terraless-endpoint-" + support.SanitizeString(endpoint.Domain)

			buffer = templates.RenderTemplateToBuffer(endpoint, buffer, awsTemplates("endpoint.tf.tmpl"), "terraless-endpoint")
		}
	}

	return buffer
}
