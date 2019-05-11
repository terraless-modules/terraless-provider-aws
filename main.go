package main

import (
	"bytes"
	"github.com/Odania-IT/terraless/schema"
	"github.com/Odania-IT/terraless/templates"
	"github.com/gobuffalo/packr/v2"
	"github.com/hashicorp/go-plugin"
	"github.com/sirupsen/logrus"
)

const (
	ProviderName = "terraless-provider-aws"
	VERSION      = "0.1.0"
)

type ProviderAws struct {}

func (provider *ProviderAws) Info(logLevel string) schema.PluginInfo {
	// Set log level
	level, _ := logrus.ParseLevel(logLevel)
	logrus.SetLevel(level)

	return schema.PluginInfo{
		Name:    ProviderName,
		Version: VERSION,
	}
}

func (provider *ProviderAws) Exec(data schema.TerralessData) error {
	logrus.Infof("[%s] Executing", ProviderName)

	return nil
}

var handshakeConfig = plugin.HandshakeConfig{
	ProtocolVersion:  1,
	MagicCookieKey:   "provider-plugin",
	MagicCookieValue: "terraless",
}

func main() {
	logrus.Info("Running Plugin: Terraless Provider AWS")
	provider := &ProviderAws{}

	// pluginMap is the map of plugins we can dispense.
	var pluginMap = map[string]plugin.Plugin{
		"Provider": &schema.ProviderPlugin{Impl: provider},
	}

	plugin.Serve(&plugin.ServeConfig{
		HandshakeConfig: handshakeConfig,
		Plugins:         pluginMap,
	})
}

func awsTemplates(name string) string {
	pckr := packr.New("awsTemplates", "./templates")

	tpl, err := pckr.FindString(name)
	if err != nil {
		logrus.Fatal("Failed retrieving template: ", err)
	}

	return tpl
}

func (provider *ProviderAws) CanHandle(resourceType string) bool {
	return resourceType == "aws"
}

func (provider *ProviderAws) FinalizeTemplates(terralessData schema.TerralessData, buffer bytes.Buffer) bytes.Buffer {
	if addTerralessLambdaRole {
		buffer = templates.RenderTemplateToBuffer(terralessData, buffer, awsTemplates("iam.tf.tmpl"), "aws-lambda-iam")
	}

	return buffer
}
