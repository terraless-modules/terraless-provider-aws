package main

import (
	"bytes"
	"github.com/Odania-IT/terraless/support"
	"github.com/Odania-IT/terraless/templates"
	"github.com/sirupsen/logrus"
)

var templateRoute53AliasRecord = `# This is generated by Terraless

resource "aws_route53_record" "{{ .TerraformName }}" {
  name = "{{ .Domain }}"
  type = "A"
  zone_id = "{{ .ZoneId }}"

  alias {
    evaluate_target_health = false
    name = "${aws_cloudfront_distribution.terraless-default.domain_name}"
    zone_id = "${aws_cloudfront_distribution.terraless-default.hosted_zone_id}"
  }
}

`

func Route53AliasRecordFor(domain string, zoneId string, buffer bytes.Buffer) bytes.Buffer {
	if zoneId == "" {
		logrus.Warnf("Not making route53 alias record for domain %s cause no zone id is set!\n", domain)
		return buffer
	}

	data := map[string]string {
		"Domain": domain,
		"TerraformName": "terraless-cloudfront-target-" + support.SanitizeString(domain),
		"ZoneId": zoneId,
	}

	return templates.RenderTemplateToBuffer(data, buffer, templateRoute53AliasRecord, "route53-alias")
}
