package checkers

import (
	"os"
	"rego-go-parser/pkg/schema"
	"strings"
)

func NothingContainsOr(api schema.API, checks schema.API, orgCheckKey string, ownerCheckKey string, f *os.File) {
	flag, roleCheck := false, false
	for j := 0; j < len(checks.Checks); j++ {
		for _, v := range checks.Checks {
			if strings.Contains(v.Key, "&&") && !strings.Contains(v.Key, "||") {
				flag = true
				break
			}
		}
		for _, v := range checks.Checks {
			if v.CheckType == "roleCheck" {
				roleCheck = true
				break
			}
		}
		if flag {
			if j == 0 {
				f.WriteString(api.Name)
				f.WriteString(" {\n")
			}
			switch checks.Checks[j].CheckType {
			case "roleCheck":
				RoleCheck(checks.Checks[j], f)
			case "orgCheck":
				OrgCheck(checks.Checks[j], roleCheck, f)
			case "ownerCheck":
				OwnerCheck(checks.Checks[j], f)
				f.WriteString("}\n\n")
			}
		}
	}
}
