package main

import (
	"bytes"
	"github.com/Odania-IT/terraless/schema"
	"github.com/Odania-IT/terraless/templates"
	"github.com/gobuffalo/packr/v2"
	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-plugin"
	"os"
)

const (
	ProviderName = "terraless-provider-aws"
	VERSION      = "0.1.0"
)

type ProviderAws struct {
	logger hclog.Logger
}

func (provider *ProviderAws) Info() schema.PluginInfo {
	return schema.PluginInfo{
		Name:    ProviderName,
		Version: VERSION,
	}
}

func (provider *ProviderAws) Exec(data schema.TerralessData) error {
	provider.logger.Info("[%s] Executing", ProviderName)

	return nil
}

var handshakeConfig = plugin.HandshakeConfig{
	ProtocolVersion:  1,
	MagicCookieKey:   "provider-plugin",
	MagicCookieValue: "terraless",
}

var logger hclog.Logger

func main() {
	logger = hclog.New(&hclog.LoggerOptions{
		Level:      hclog.Trace,
		Output:     os.Stderr,
		JSONFormat: true,
	})
	logger.Info("Running Plugin: Terraless Provider AWS")

	provider := &ProviderAws{
		logger: logger,
	}

	// pluginMap is the map of plugins we can dispense.
	var pluginMap = map[string]plugin.Plugin{
		"Provider": &schema.ProviderPlugin{Impl: provider},
	}

	plugin.Serve(&plugin.ServeConfig{
		HandshakeConfig: handshakeConfig,
		Plugins:         pluginMap,
	})
}

func fatal(msg string, args ...interface{}) {
	logger.Error(msg, args...)
	os.Exit(1)
}

func awsTemplates(name string) string {
	pckr := packr.New("awsTemplates", "./templates")

	tpl, err := pckr.FindString(name)
	if err != nil {
		fatal("Failed retrieving template: ", err)
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
