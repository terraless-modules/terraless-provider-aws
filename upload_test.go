package main

import (
	"github.com/Odania-IT/terraless/schema"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

var uploadedFiles int
func uploadFileMock(svc *s3manager.Uploader, uploadInput s3manager.UploadInput) (*s3manager.UploadOutput, error) {
	uploadedFiles += 1
	response := s3manager.UploadOutput{
		Location: "only-mock",
	}
	return &response, nil
}

func TestTemplatesFunctions_RecursiveUpload(t *testing.T) {
	// given
	provider := ProviderAws{}
	dir, _ := os.Getwd()
	uploadFileFunc = uploadFileMock
	terralessData := schema.TerralessData{
		Config: schema.TerralessConfig{
			SourcePath: dir,
		},
	}
	upload := schema.TerralessUpload{
		Type: "s3",
		Source: "templates",
	}

	// when
	uploadedFilenames := provider.ProcessUpload(terralessData, upload)

	// then
	expected := []string{
		"authorizer.tf.tmpl",
		"certificate.tf.tmpl",
		"cloudfront.tf.tmpl",
		"endpoint.tf.tmpl",
		"function-event/http.tf.tmpl",
		"function-event/integration/CodeCommit.tf.tmpl",
		"function-event/integration/http.tf.tmpl",
		"function-event/integration/sqs.tf.tmpl",
		"iam.tf.tmpl",
		"lambda-at-edge.js",
		"lambda-at-edge.tf.tmpl",
	}
	assert.Equal(t, 11, uploadedFiles)
	assert.Equal(t, expected, uploadedFilenames)
}
