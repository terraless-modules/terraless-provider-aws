package main

import (
	"bytes"
	"fmt"
	"github.com/Odania-IT/terraless/schema"
	"github.com/Odania-IT/terraless/support"
	"github.com/Odania-IT/terraless/templates"
	"strconv"
)

var lambdaFunctionsTemplate = `
# Lambda Function {{.FunctionName}}

resource "aws_cloudwatch_log_group" "lambda-log-{{.FunctionName}}" {
  name = "/aws/lambda/{{ .ProjectName }}-{{.FunctionName}}"
  retention_in_days = 14

  tags = "${merge(var.terraless-default-tags, map("FunctionName", "{{.FunctionName}}"))}"
}

resource "aws_lambda_function" "lambda-{{.FunctionName}}" {
  filename = "${data.archive_file.lambda-archive.output_path}"
  function_name = "{{ .ProjectName }}-{{.FunctionName}}"
  role = "{{.RoleArn}}"
  handler = "{{.Handler}}"
  source_code_hash = "${data.archive_file.lambda-archive.output_base64sha256}"
  runtime = "{{.Runtime}}"

  {{ if .RenderEnvironment }}
  environment {
    variables = {
      {{ range $key, $val := .Environment }}{{ $key }} = "{{ $val }}"
      {{ end }}
    }
  }
  {{ end }}

  tags = "${merge(var.terraless-default-tags, map("FunctionName", "{{.FunctionName}}"))}"
}

{{ if .AddApiGatewayPermission }}
resource "aws_lambda_permission" "lambda-{{.FunctionName}}" {
  action = "lambda:InvokeFunction"
  function_name = "${aws_lambda_function.lambda-{{.FunctionName}}.function_name}"
  principal = "apigateway.amazonaws.com"
}
{{ end }}

`

var addTerralessLambdaRole bool
func renderBaseFunction(functionConfig schema.TerralessFunction, functionName string, config schema.TerralessConfig) string {
	var buffer bytes.Buffer
	logger.Debug(fmt.Sprintf("Rendering Template for Lambda Function: %s\n", functionName))
	functionConfig.RenderEnvironment = len(functionConfig.Environment) > 0
	functionConfig.FunctionName = functionName
	functionConfig.ProjectName = config.ProjectName

	// Set default runtime if none is specified for the function
	if functionConfig.Runtime == "" {
		functionConfig.Runtime = config.Settings.Runtime
	}

	if functionConfig.RoleArn == "" {
		functionConfig.RoleArn = "${aws_iam_role.terraless-lambda-iam-role.arn}"
		addTerralessLambdaRole = true
	}

	for _, event := range functionConfig.Events {
		if event.Type == "http" {
			functionConfig.AddApiGatewayPermission = true
		}
	}

	buffer = templates.RenderTemplateToBuffer(functionConfig, buffer, lambdaFunctionsTemplate, "aws-lambda-function")
	return buffer.String()
}

func (provider *ProviderAws) RenderFunctionTemplates(resourceType string, functionEvents schema.FunctionEvents, terralessData *schema.TerralessData) string {
	if !provider.CanHandle(resourceType) {
		return ""
	}

	var buffer bytes.Buffer
	buffer.WriteString("## Terraless Functions AWS\n\n")
	functionsToRender := map[string]bool{}
	for eventType, functionEventArray := range functionEvents.Events {
		baseTemplate, err := getAwsTemplate("function-event/" + eventType + ".tf.tmpl", false)

		if err != nil {
			baseTemplate = "## Terraless " + eventType + "\n"
		}

		buffer = templates.RenderTemplateToBuffer(terralessData, buffer, baseTemplate, "function-event-" + eventType)

		// Events
		pathsRendered := map[string]string{}
		for key, event := range functionEventArray {
			logger.Debug(fmt.Sprintf("[EventType %s][AWS %s] Rendering Event %s\n", eventType, event.FunctionName, event))
			functionsToRender[event.FunctionName] = true

			// Render function template
			functionEvent := event.FunctionEvent
			functionEvent.FunctionName = event.FunctionName
			functionEvent.Idx = strconv.FormatInt(int64(key), 10)
			functionEvent.ProjectName = terralessData.Config.ProjectName
			functionEvent.PathsRendered = pathsRendered
			functionEvent.ResourceNameForPath = functionEvent.Idx

			// Authorization
			if functionEvent.Authorizer == "" {
				functionEvent.Authorization = "NONE"
			} else {
				authorizer := terralessData.Config.Authorizers[functionEvent.Authorizer]
				functionEvent.Authorization = authorizer.AuthorizerType
				authorizer.TerraformName = "terraless-authorizer-" + support.SanitizeString(authorizer.Name)
				functionEvent.AuthorizerId = "${aws_api_gateway_authorizer." + authorizer.TerraformName + ".id}"
			}

			// Integration Template
			integrationTemplate, err := getAwsTemplate("function-event/integration/" + functionEvent.Type + ".tf.tmpl", true)

			if err != nil {
				fatal(fmt.Sprintf("Event Type %s unknown! Function: %s", functionEvent.Type, event.FunctionName))
			}

			buffer = templates.RenderTemplateToBuffer(functionEvent, buffer, integrationTemplate, "function-event-" + functionEvent.Type + "-" + functionEvent.Idx)
			pathsRendered[functionEvent.Path] = support.SanitizeString(functionEvent.Path)
		}
	}

	// Render function base
	for functionName := range functionsToRender {
		functionConfig := terralessData.Config.Functions[functionName]
		buffer.WriteString(renderBaseFunction(functionConfig, functionName, terralessData.Config))

		terralessData.Config.Runtimes = append(terralessData.Config.Runtimes, functionConfig.Runtime)
	}

	return buffer.String()
}
