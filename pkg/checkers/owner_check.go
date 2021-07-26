package checkers

import (
	"os"
	"rego-go-parser/pkg/helpers"
	"rego-go-parser/pkg/schema"
	"strings"
)

func OwnerCheck(checks schema.Checks, f *os.File) {
	if strings.Contains(checks.Key, "&&") {
		jsonKeysArray := strings.Split(checks.Key, "&&")

		for key, value := range jsonKeysArray {
			jsonKeysArray[key] = strings.TrimSpace(value)
		}

		for _, value := range jsonKeysArray {
			switch value {
			case "body":
				if checks.Body != "" {
					f.WriteString(helpers.OwnerCheckRequestBody + checks.Body)
					f.WriteString(helpers.NewLine)
				}
			case "header":
				if checks.Header != "" {
					f.WriteString(helpers.OwnerCheckRequestHeader + checks.Header)
					f.WriteString(helpers.NewLine)
				}
			}
		}
	} else {
		switch checks.Key {
		case "body":
			if checks.Body != "" {
				f.WriteString(helpers.OwnerCheckRequestBody + checks.Body)
				f.WriteString(helpers.NewLine)
			}
		case "header":
			if checks.Header != "" {
				f.WriteString(helpers.OwnerCheckRequestHeader + checks.Header)
				f.WriteString(helpers.NewLine)
			}
		}
	}
}
