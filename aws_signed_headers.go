package s3

import (
	"net/http"
	"sort"
	"strings"
)

func canonicalAndSignedHeaders(original http.Header) (canonical, signed string) {
	lowercaseKeys := map[string]string{} // map[lowercase]original
	for key := range original {
		if headerEligibleForSigning(key) {
			lowercaseKeys[strings.ToLower(key)] = key
		}
	}

	var sortedKeys []string
	for key := range lowercaseKeys {
		sortedKeys = append(sortedKeys, key)
	}
	sort.Strings(sortedKeys)

	canonicalBuilder := new(strings.Builder)
	for _, lowerKey := range sortedKeys {
		titleKey := lowercaseKeys[lowerKey]
		var values []string
		for _, value := range original[titleKey] {
			if titleKey == "Host" && strings.Contains(value, ":") {
				value = strings.Split(value, ":")[0] // AWS does not include port in signing request.
			}
			values = append(values, trimHeaderValue(value))
		}
		canonicalBuilder.WriteString(lowerKey)
		canonicalBuilder.WriteString(":")
		canonicalBuilder.WriteString(strings.Join(values, ","))
		canonicalBuilder.WriteString("\n")
	}
	return canonicalBuilder.String(), strings.Join(sortedKeys, ";")
}

func trimHeaderValue(value string) string {
	value = strings.TrimSpace(value)
	for strings.Contains(value, "  ") {
		value = strings.ReplaceAll(value, "  ", " ")
	}
	return value
}

func headerEligibleForSigning(key string) bool {
	if runningAWSTestSuite {
		return true
	}
	switch key {
	case "Content-Type", "Content-Md5", "Host":
		return true
	default:
		return strings.HasPrefix(key, "X-Amz")
	}
}

var runningAWSTestSuite bool