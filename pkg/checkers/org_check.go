package checkers

import (
	"os"
	"rego-go-parser/pkg/helpers"
	"rego-go-parser/pkg/schema"
	"strings"
)

/* This functions generates a rego policy as below with proper indentations:

	With role check enabled on the API, compare token orgIds and request body orgIds:
		some j, k
    	token.payload.roles[i].scope[j].orgId == input.parsed_body.request.content.createdFor[k]

	Without role check enabled on the API, compare token orgIds and request body orgIds:
		some j, k
    	token.payload.roles[_].scope[j].orgId == input.parsed_body.request.content.createdFor[k]

	With role check enabled on the API, compare token orgIds and request header orgIds:
		some j
    	token.payload.roles[i].scope[j].orgId == http_request.header.X-Channel-Id

	Without role check enabled on the API, compare token orgIds and request header orgIds:
		some j
    	token.payload.roles[_].scope[j].orgId == http_request.header.X-Channel-Id
*/

func OrgCheck(checks schema.Checks, needsRoleCheck bool, f *os.File) {
	var loopVariablesArray []rune

	countOfAsterisk := strings.Count(checks.Body, "*")
	f.WriteString(helpers.SomeJ)

	for i, listOfLoopVarsRequired := 0, 'j'; i < countOfAsterisk; i++ {
		listOfLoopVarsRequired++
		loopVariablesArray = append(loopVariablesArray, listOfLoopVarsRequired)
	}

	for _, value := range loopVariablesArray {
		f.WriteString(helpers.CommaSpace + string(value))
	}
	f.WriteString(helpers.NewLine)

	if strings.Contains(checks.Key, "&&") {
		jsonKeysArray := strings.Split(checks.Key, "&&")

		for key, value := range jsonKeysArray {
			jsonKeysArray[key] = strings.TrimSpace(value)
		}

		for _, value := range loopVariablesArray {
			checks.Body = strings.Replace(checks.Body, "*", string(value), 1)
		}

		for _, value := range jsonKeysArray {
			switch value {
			case "body":
				if checks.Body != "" {
					if needsRoleCheck {
						f.WriteString(helpers.OrgCheckRequestBodyWithRoleCheck + checks.Body)
						f.WriteString(helpers.NewLine)
					} else {
						f.WriteString(helpers.OrgCheckRequestBodyWithoutRoleCheck + checks.Body)
						f.WriteString(helpers.NewLine)
					}
				}
			case "header":
				if checks.Header != "" {
					if needsRoleCheck {
						f.WriteString(helpers.OrgCheckRequestHeaderWithRoleCheck + checks.Header)
						f.WriteString(helpers.NewLine)
					} else {
						f.WriteString(helpers.OrgCheckRequestHeaderWithoutRoleCheck + checks.Header)
						f.WriteString(helpers.NewLine)
					}
				}
			}
		}
	} else {
		for _, value := range loopVariablesArray {
			checks.Body = strings.Replace(checks.Body, "*", string(value), 1)
			switch checks.Key {
			case "body":
				if checks.Body != "" {
					if needsRoleCheck {
						f.WriteString(helpers.OrgCheckRequestBodyWithRoleCheck + checks.Body)
						f.WriteString(helpers.NewLine)
					} else {
						f.WriteString(helpers.OrgCheckRequestBodyWithoutRoleCheck + checks.Body)
						f.WriteString(helpers.NewLine)
					}
				}
			case "header":
				if checks.Header != "" {
					if needsRoleCheck {
						f.WriteString(helpers.OrgCheckRequestHeaderWithRoleCheck + checks.Header)
						f.WriteString(helpers.NewLine)
					} else {
						f.WriteString(helpers.OrgCheckRequestHeaderWithoutRoleCheck + checks.Header)
						f.WriteString(helpers.NewLine)
					}
				}
			}
		}
	}
}
