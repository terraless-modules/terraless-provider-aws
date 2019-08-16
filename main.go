package main

import (
	"bytes"
	"fmt"
	"github.com/Odania-IT/terraless/schema"
	"github.com/Odania-IT/terraless/templates"
	"github.com/gobuffalo/packr/v2"
	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-plugin"
	"github.com/pkg/errors"
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
	template, err := getAwsTemplate(name, true)

	if err != nil {
		fatal(err.Error())
	}

	return template
}

func getAwsTemplate(name string, doFatal bool) (string, error) {
	logger.Debug(fmt.Sprintf("Looking up template %s\n", name))
	pckr := packr.New("awsTemplates", "./templates")

	tpl, err := pckr.FindString(name)
	if err != nil {
		if doFatal {
			fatal(fmt.Sprintf("Failed retrieving template: %s", err))
		}

		return "", errors.New(fmt.Sprintf("Failed retrieving template: %s", err))
	}

	return tpl, nil
}

// func listTemplates() {
// 	pckr := packr.New("awsTemplates", "./templates")
//
// 	list := pckr.List()
// 	for _, entry := range list {
// 		logger.Warn(entry)
// 	}
// }

func (provider *ProviderAws) CanHandle(resourceType string) bool {
	return resourceType == "aws"
}

func (provider *ProviderAws) FinalizeTemplates(terralessData schema.TerralessData) string {
	var buffer bytes.Buffer
	if addTerralessLambdaRole {
		buffer = templates.RenderTemplateToBuffer(terralessData, buffer, awsTemplates("iam.tf.tmpl"), "aws-lambda-iam")
	}

	return buffer.String()
}
