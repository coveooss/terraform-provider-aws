package aws

import (
	"regexp"
)

var (
	ec2ResourceIdRegexp = regexp.MustCompile(`\A([a-z]+)-([a-z0-9]+)\z`)
)

// parseEc2ResourceId parses an EC2 resource ID into prefix and unique part.
func parseEc2ResourceId(s string) (string, string) {
	matches := ec2ResourceIdRegexp.FindStringSubmatch(s)
	if matches == nil {
		return "", ""
	}
	return matches[1], matches[2]
}
