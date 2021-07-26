package checkers

import (
	"os"
	"rego-go-parser/pkg/schema"
	"strconv"
	"strings"
)

func AllChecksContainingOnlyOr(api schema.API, checks schema.Checks, orgCheckKey string, ownerCheckKey string, f *os.File) {
	switch checks.CheckType {
	case "roleCheck":
		switch checks.Key {
		case "token":
			if checks.Token != "" {
				f.WriteString(api.Name)
				f.WriteString(" {\n")
				f.WriteString("  api_roles := [")
				length := len(strings.Split(checks.Token, ","))
				for k, v := range strings.Split(checks.Token, ",") {
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
			if checks.Body != "" {
				count := strings.Count(checks.Body, "*")
				var placeHolders []rune
				f.WriteString("  some j")
				for i, v := 0, 'j'; i < count; i++ {
					v++
					placeHolders = append(placeHolders, v)
				}
				for _, v := range placeHolders {
					f.WriteString(string(", " + string(v)))
					checks.Body = strings.Replace(checks.Body, "*", string(v), 1)
				}
				f.WriteString("\n")

				f.WriteString("  token.payload.roles[i].scope[j].orgId == input.parsed_body." + checks.Body + "\n")
			}
		case "header":
			if checks.Header != "" {
				f.WriteString("  some j")
				f.WriteString("\n")
				f.WriteString("  token.payload.roles[i].scope[j].orgId == http_request.header." + checks.Header + "\n")
			}
		}
	case "ownerCheck":
		switch ownerCheckKey {
		case "body":
			f.WriteString("  sub[2] == input.parsed_body." + checks.Body + "\n")
		case "header":
			f.WriteString("  sub[2] == http_request.header." + checks.Header + "\n")
		}
		f.WriteString("}\n\n")
	}
}
