package checkers

import (
	"os"
	"rego-go-parser/pkg/helpers"
	"rego-go-parser/pkg/schema"
	"strconv"
	"strings"
)

/* This functions generates a rego policy as below with proper indentations:

	Role check with org check
		api_roles := ["CONTENT_CREATOR", "COURSE_CREATOR"]
  		some i
  		api_roles[_] == token.payload.roles[i].role

Todo:
	Role check without org check
		api_roles := ["CONTENT_CREATOR", "COURSE_CREATOR"]
  		some i
  		api_roles[_] == token.payload.roles[_].role
*/
func RoleCheck(checks schema.Checks, f *os.File) {
	f.WriteString(helpers.ApiRolesStart)
	switch checks.Key {
	case "token":
		if checks.Token != "" {
			length := len(strings.Split(checks.Token, ","))
			for k, v := range strings.Split(checks.Token, ",") {
				f.WriteString(strconv.Quote(strings.TrimSpace(v)))
				if k < length-1 {
					f.WriteString(helpers.CommaSpace)
				}
			}
			f.WriteString(helpers.ApiRolesEnd)
			f.WriteString(helpers.SomeI)
			f.WriteString(helpers.ApiRolesCheck)
		}
	}
}
