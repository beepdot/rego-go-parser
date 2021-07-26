package checkers

import (
	"os"
	"rego-go-parser/pkg/schema"
	"strconv"
	"strings"
)

func OnlyOrgCheckContainsOr(api schema.API, checks schema.API, orgCheckKey string, ownerCheckKey string, f *os.File) {
	//flag, roleCheck := false, false
	for j := 0; j < len(checks.Checks); j++ {
		switch checks.Checks[j].CheckType {
		case "roleCheck":
			switch checks.Checks[j].Key {
			case "token":
				if checks.Checks[j].Token != "" {
					f.WriteString(api.Name)
					f.WriteString(" {\n")
					f.WriteString("  api_roles := [")
					length := len(strings.Split(checks.Checks[j].Token, ","))
					for k, v := range strings.Split(checks.Checks[j].Token, ",") {
						f.WriteString(strconv.Quote(strings.TrimSpace(v)))
						if k < length-1 {
							f.WriteString(", ")
						}
					}
					f.WriteString("]\n")
					f.WriteString("  some i\n")
					f.WriteString("  api_roles[_] == token.payload.roles[i].role\n")
				}
			}
		case "orgCheck":
			switch orgCheckKey {
			case "body":
				if checks.Checks[j].Body != "" {
					count := strings.Count(checks.Checks[j].Body, "*")
					var placeHolders []rune
					f.WriteString("  some j")
					for i, v := 0, 'j'; i < count; i++ {
						v++
						placeHolders = append(placeHolders, v)
					}
					for _, v := range placeHolders {
						f.WriteString(string(", " + string(v)))
					}
					f.WriteString("\n")

					for _, v := range placeHolders {
						checks.Checks[j].Body = strings.Replace(checks.Checks[j].Body, "*", string(v), 1)
					}
					f.WriteString("  token.payload.roles[i].scope[j].orgId == input.parsed_body." + checks.Checks[j].Body + "\n")
				}
			case "header":
				if checks.Checks[j].Header != "" {
					f.WriteString("  some j")
					f.WriteString("\n")
					f.WriteString("  token.payload.roles[i].scope[j].orgId == http_request.header." + checks.Checks[j].Header + "\n")
				}
			}
		case "ownerCheck":
			OwnerCheck(checks.Checks[j], f)
			f.WriteString("}\n\n")
		}
	}
}
