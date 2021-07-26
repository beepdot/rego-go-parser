package helpers

import (
	"rego-go-parser/pkg/schema"
	"strings"
)

func IdentifyKeys(apis *schema.APIs) ([]string, []string, bool, bool) {
	var orgCheckKeys, ownerCheckKeys []string
	var orgCheckContainsAnd, ownerCheckContainsAnd = false, false

	for i := 0; i < len(apis.APIs); i++ {
		for j := 0; j < len(apis.APIs[i].Checks); j++ {
			if strings.Contains(apis.APIs[i].Checks[j].Key, "||") {
				if apis.APIs[i].Checks[j].CheckType == "orgCheck" {
					orgCheckKeys = strings.Split(apis.APIs[i].Checks[j].Key, "||")
					for key, value := range orgCheckKeys {
						orgCheckKeys[key] = strings.TrimSpace(value)
					}
				} else if apis.APIs[i].Checks[j].CheckType == "ownerCheck" {
					ownerCheckKeys = strings.Split(apis.APIs[i].Checks[j].Key, "||")
					for key, value := range ownerCheckKeys {
						ownerCheckKeys[key] = strings.TrimSpace(value)
					}
				}
			} else if strings.Contains(apis.APIs[i].Checks[j].Key, "&&") {
				if apis.APIs[i].Checks[j].CheckType == "orgCheck" {
					orgCheckKeys = strings.Split(apis.APIs[i].Checks[j].Key, "&&")
					for k, v := range orgCheckKeys {
						orgCheckKeys[k] = strings.TrimSpace(v)
					}
					orgCheckContainsAnd = true
				} else if apis.APIs[i].Checks[j].CheckType == "ownerCheck" {
					ownerCheckKeys = strings.Split(apis.APIs[i].Checks[j].Key, "||")
					for key, value := range ownerCheckKeys {
						ownerCheckKeys[key] = strings.TrimSpace(value)
					}
					ownerCheckContainsAnd = true
				}
			} else {
				if apis.APIs[i].Checks[j].CheckType == "orgCheck" {
					orgCheckKeys = append(orgCheckKeys, apis.APIs[i].Checks[j].Key)
				} else if apis.APIs[i].Checks[j].CheckType == "ownerCheck" {
					ownerCheckKeys = append(ownerCheckKeys, apis.APIs[i].Checks[j].Key)
				}
			}
		}
	}
	return orgCheckKeys, ownerCheckKeys, orgCheckContainsAnd, ownerCheckContainsAnd
}
