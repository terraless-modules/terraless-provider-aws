package main

import (
	"bytes"
	"github.com/Odania-IT/terraless/schema"
	"github.com/Odania-IT/terraless/support"
	"github.com/Odania-IT/terraless/templates"
	"github.com/sirupsen/logrus"
)

func (provider *ProviderAws) RenderCertificateTemplates(config schema.TerralessConfig, buffer bytes.Buffer) bytes.Buffer {
	for _, certificate := range config.Certificates {
		if provider.CanHandle(certificate.Type) {
			logrus.Debugf("Generating certificate template for %s\n", certificate.Domain)
			certificate.ProjectName = config.ProjectName
			certificate.TerraformName = "terraless-certificate-" + support.SanitizeString(certificate.Domain)

			buffer = templates.RenderTemplateToBuffer(certificate, buffer, awsTemplates("certificate.tf.tmpl"), "terraless-certificate")
		}
	}

	return buffer
}