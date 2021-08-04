package helpers

var (
	ApiRolesStart                         string = "  api_roles := ["
	ApiRolesEnd                           string = "]\n"
	SomeI                                 string = "  some i"
	SomeJ                                 string = "  some j"
	RolesCheckWithoutOrgCheck             string = "  api_roles[_] == token.payload.roles[_].role\n"
	RolesCheckWithOrgCheck                string = "  api_roles[_] == token.payload.roles[i].role\n"
	CommaSpace                            string = ", "
	NewLine                               string = "\n"
	OrgCheckRequestBodyWithRoleCheck      string = "  token.payload.roles[i].scope[j].orgId == input.parsed_body."
	OrgCheckRequestBodyWithoutRoleCheck   string = "  token.payload.roles[_].scope[j].orgId == input.parsed_body."
	OrgCheckRequestHeaderWithRoleCheck    string = "  token.payload.roles[i].scope[j].orgId == http_request.header."
	OrgCheckRequestHeaderWithoutRoleCheck string = "  token.payload.roles[_].scope[j].orgId == http_request.header."
	OwnerCheckRequestBody                 string = "  sub[2] == input.parsed_body."
	OwnerCheckRequestHeader               string = "  sub[2] == http_request.header."
)
