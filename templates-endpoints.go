package main

import (
	"bytes"
	"fmt"
	"github.com/Odania-IT/terraless/schema"
	"github.com/Odania-IT/terraless/support"
	"github.com/Odania-IT/terraless/templates"
)

func (provider *ProviderAws) RenderEndpointTemplates(config schema.TerralessConfig) string {
	var buffer bytes.Buffer
	for _, endpoint := range config.Endpoints {
		if endpoint.Type == "apigateway" {
			logger.Debug(fmt.Sprintf("Generating certificate template for %s\n", endpoint.Domain))
			endpoint.TerralessCertificate = config.Certificates[endpoint.Certificate]
			endpoint.TerraformName = "terraless-endpoint-" + support.SanitizeString(endpoint.Domain)

			buffer = templates.RenderTemplateToBuffer(endpoint, buffer, awsTemplates("endpoint.tf.tmpl"), "terraless-endpoint")
		}
	}

	return buffer.String()
}
