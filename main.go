package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
)

type APIs struct {
	APIs []API `json:"apis"`
}

type API struct {
	Name        string   `json:"name"`
	URIs        string   `json:"uris"`
	UpstreamUrl string   `json:"upstream_url"`
	Checks      []Checks `json:"checks"`
}

type Checks struct {
	CheckType string `json:"checkType,omitempty"`
	Key       string `json:"key,omitempty"`
	Token     string `json:"token,omitempty"`
	Body      string `json:"body,omitempty"`
	Header    string `json:"header,omitempty"`
}

func RoleCheck(roleCheck Checks, fp *os.File) {
	fp.WriteString("  api_roles := [")
	switch roleCheck.Key {
	case "token":
		if roleCheck.Token != "" {
			length := len(strings.Split(roleCheck.Token, ","))
			for k, v := range strings.Split(roleCheck.Token, ",") {
				fp.WriteString(strconv.Quote(strings.TrimSpace(v)))
				if k < length-1 {
					fp.WriteString(", ")
				}
			}
			fp.WriteString("]\n")
			fp.WriteString("  some i\n")
			fp.WriteString("  api_roles[_] == token.payload.roles[i].role\n")
		}
	}
}

// func OrgCheck(orgCheck Checks, roleCheck bool, fp *os.File) {
// 	count := strings.Count(orgCheck.Body, "*")
// 	var placeHolders []rune
// 	fp.WriteString("  some j")
// 	for i, v := 0, 'j'; i < count; i++ {
// 		v++
// 		placeHolders = append(placeHolders, v)
// 	}
// 	for _, v := range placeHolders {
// 		fp.WriteString(string(", " + string(v)))
// 	}
// 	fp.WriteString("\n")
// 	for _, v := range placeHolders {
// 		orgCheck.Body = strings.Replace(orgCheck.Body, "*", string(v), 1)
// 	}
// 	switch orgCheck.Key {
// 	case "body":
// 		if orgCheck.Body != "" {
// 			if roleCheck {
// 				fp.WriteString("  token.payload.roles[i].scope[j].orgId == input.parsed_body." + orgCheck.Body + "\n")
// 			} else {
// 				fp.WriteString("  token.payload.roles[_].scope[j].orgId == input.parsed_body." + orgCheck.Body + "\n")
// 			}
// 		}
// 	case "header":
// 		if orgCheck.Header != "" {
// 			if roleCheck {
// 				fp.WriteString("  token.payload.roles[i].scope[j].orgId == http_request.header." + orgCheck.Header + "\n")
// 			} else {
// 				fp.WriteString("  token.payload.roles[_].scope[j].orgId == http_request.header." + orgCheck.Header + "\n")
// 			}
// 		}
// 	}
// }

func OrgCheck(orgCheck Checks, roleCheck bool, fp *os.File) {
	count := strings.Count(orgCheck.Body, "*")
	var placeHolders []rune
	fp.WriteString("  some j")
	for i, v := 0, 'j'; i < count; i++ {
		v++
		placeHolders = append(placeHolders, v)
	}
	for _, v := range placeHolders {
		fp.WriteString(string(", " + string(v)))
	}
	fp.WriteString("\n")
	if strings.Contains(orgCheck.Key, "&&") {
		values := strings.Split(orgCheck.Key, "&&")
		for k, v := range values {
			values[k] = strings.TrimSpace(v)
		}
		for _, v := range placeHolders {
			orgCheck.Body = strings.Replace(orgCheck.Body, "*", string(v), 1)
		}
		for _, v := range values {
			switch v {
			case "body":
				if orgCheck.Body != "" {
					if roleCheck {
						fp.WriteString("  token.payload.roles[i].scope[j].orgId == input.parsed_body." + orgCheck.Body + "\n")
					} else {
						fp.WriteString("  token.payload.roles[_].scope[j].orgId == input.parsed_body." + orgCheck.Body + "\n")
					}
				}
			case "header":
				if orgCheck.Header != "" {
					if roleCheck {
						fp.WriteString("  token.payload.roles[i].scope[j].orgId == http_request.header." + orgCheck.Header + "\n")
					} else {
						fp.WriteString("  token.payload.roles[_].scope[j].orgId == http_request.header." + orgCheck.Header + "\n")
					}
				}
			}
		}
	} else {
		for _, v := range placeHolders {
			orgCheck.Body = strings.Replace(orgCheck.Body, "*", string(v), 1)
			switch orgCheck.Key {
			case "body":
				if orgCheck.Body != "" {
					fp.WriteString("  token.payload.roles[i].scope[j].orgId == input.parsed_body." + orgCheck.Body + "\n")
				}
			case "header":
				if orgCheck.Header != "" {
					fp.WriteString("  token.payload.roles[i].scope[j].orgId == http_request.header." + orgCheck.Header + "\n")
				}
			}
		}
	}
}

// func OwnerCheck(ownerCheck Checks, fp *os.File) {
// 	switch ownerCheck.Key {
// 	case "body":
// 		fp.WriteString("  sub[2] == input.parsed_body." + ownerCheck.Body + "\n")
// 	case "header":
// 		fp.WriteString("  sub[2] == http_request.header." + ownerCheck.Header + "\n")
// 	}
// }

func OwnerCheck(ownerCheck Checks, fp *os.File) {
	if strings.Contains(ownerCheck.Key, "&&") {
		values := strings.Split(ownerCheck.Key, "&&")
		for k, v := range values {
			values[k] = strings.TrimSpace(v)
		}
		for _, v := range values {
			switch v {
			case "body":
				if ownerCheck.Body != "" {
					fp.WriteString("  sub[2] == input.parsed_body." + ownerCheck.Body + "\n")
				}
			case "header":
				if ownerCheck.Header != "" {
					fp.WriteString("  sub[2] == http_request.header." + ownerCheck.Header + "\n")
				}
			}
		}
	} else {
		switch ownerCheck.Key {
		case "body":
			if ownerCheck.Body != "" {
				fp.WriteString("  sub[2] == input.parsed_body." + ownerCheck.Body + "\n")
			}
		case "header":
			if ownerCheck.Header != "" {
				fp.WriteString("  sub[2] == http_request.header." + ownerCheck.Header + "\n")
			}
		}
	}
}

func main() {
	jsonFile, err := os.Open("checks.json")
	if err != nil {
		panic(err)
	}

	fmt.Println("Successfully Opened checks.json")

	defer jsonFile.Close()
	byteValue, _ := ioutil.ReadAll(jsonFile)

	var apis APIs
	json.Unmarshal(byteValue, &apis)

	f, err := os.OpenFile("policy.rego", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	// Basic json
	// for i := 0; i < len(apis.APIs); i++ {
	// 	flag, roleCheck := true, false
	// 	for j := 0; j < len(apis.APIs[i].Checks); j++ {
	// 		for _, v := range apis.APIs[i].Checks {
	// 			if strings.Contains(v.Key, "||") || strings.Contains(v.Key, "&&") {
	// 				// flag = false
	// 				// temporary comment out
	// 				flag = true
	// 				break
	// 			}
	// 		}
	// 		for _, v := range apis.APIs[i].Checks {
	// 			if v.CheckType == "roleCheck" {
	// 				roleCheck = true
	// 				break
	// 			}
	// 		}
	// 		if flag {
	// 			if j == 0 {
	// 				f.WriteString(apis.APIs[i].Name)
	// 				f.WriteString(" {\n")
	// 			}
	// 			switch apis.APIs[i].Checks[j].CheckType {
	// 			case "roleCheck":
	// 				RoleCheck(apis.APIs[i].Checks[j], f)
	// 			case "orgCheck":
	// 				OrgCheck(apis.APIs[i].Checks[j], roleCheck, f)
	// 			case "ownerCheck":
	// 				OwnerCheck(apis.APIs[i].Checks[j], f)
	// 			}
	// 		}
	// 	}
	// 	if flag {
	// 		f.WriteString("}\n\n")
	// 	}
	// }

	// && json
	// for i := 0; i < len(apis.APIs); i++ {
	// 	flag, roleCheck := false, false
	// 	for j := 0; j < len(apis.APIs[i].Checks); j++ {
	// 		for _, v := range apis.APIs[i].Checks {
	// 			if strings.Contains(v.Key, "&&") && !strings.Contains(v.Key, "||") {
	// 				flag = true
	// 				break
	// 			}
	// 		}
	// 		for _, v := range apis.APIs[i].Checks {
	// 			if v.CheckType == "roleCheck" {
	// 				roleCheck = true
	// 				break
	// 			}
	// 		}
	// 		if flag {
	// 			if j == 0 {
	// 				f.WriteString(apis.APIs[i].Name)
	// 				f.WriteString(" {\n")
	// 			}
	// 			switch apis.APIs[i].Checks[j].CheckType {
	// 			case "roleCheck":
	// 				RoleCheck(apis.APIs[i].Checks[j], f)
	// 			case "orgCheck":
	// 				OrgCheck(apis.APIs[i].Checks[j], roleCheck, f)
	// 			case "ownerCheck":
	// 				OwnerCheck(apis.APIs[i].Checks[j], f)
	// 			}
	// 		}
	// 	}
	// 	if flag {
	// 		f.WriteString("}\n\n")
	// 	}
	// }

	// All 3 scenarios
	var ownerCheck []string
	var orgCheck []string
	var orgAnd, ownerAnd = false, false
	for i := 0; i < len(apis.APIs); i++ {
		for j := 0; j < len(apis.APIs[i].Checks); j++ {
			if strings.Contains(apis.APIs[i].Checks[j].Key, "||") {
				if apis.APIs[i].Checks[j].CheckType == "ownerCheck" {
					ownerCheck = strings.Split(apis.APIs[i].Checks[j].Key, "||")
					for k, v := range ownerCheck {
						ownerCheck[k] = strings.TrimSpace(v)
					}
				} else if apis.APIs[i].Checks[j].CheckType == "orgCheck" {
					orgCheck = strings.Split(apis.APIs[i].Checks[j].Key, "||")
					for k, v := range orgCheck {
						orgCheck[k] = strings.TrimSpace(v)
					}
				}
			} else if strings.Contains(apis.APIs[i].Checks[j].Key, "&&") {
				if apis.APIs[i].Checks[j].CheckType == "ownerCheck" {
					ownerCheck = strings.Split(apis.APIs[i].Checks[j].Key, "||")
					for k, v := range ownerCheck {
						ownerCheck[k] = strings.TrimSpace(v)
					}
					ownerAnd = true
				} else if apis.APIs[i].Checks[j].CheckType == "orgCheck" {
					orgCheck = strings.Split(apis.APIs[i].Checks[j].Key, "&&")
					for k, v := range orgCheck {
						orgCheck[k] = strings.TrimSpace(v)
					}
					orgAnd = true
				}
			} else {
				if apis.APIs[i].Checks[j].CheckType == "ownerCheck" {
					ownerCheck = append(ownerCheck, apis.APIs[i].Checks[j].Key)
				} else if apis.APIs[i].Checks[j].CheckType == "orgCheck" {
					orgCheck = append(orgCheck, apis.APIs[i].Checks[j].Key)
				}
			}
		}
	}

	fmt.Println(orgAnd, ownerAnd)

	for i := 0; i < len(apis.APIs); i++ {
		for oc := 0; oc < len(ownerCheck); oc++ {
			if !ownerAnd {
				for ogc := 0; ogc < len(orgCheck); ogc++ {
					if !orgAnd {
						for j := 0; j < len(apis.APIs[i].Checks); j++ {
							switch apis.APIs[i].Checks[j].CheckType {
							case "roleCheck":
								switch apis.APIs[i].Checks[j].Key {
								case "token":
									if apis.APIs[i].Checks[j].Token != "" {
										f.WriteString(apis.APIs[i].Name)
										f.WriteString(" {\n")
										f.WriteString("  api_roles := [")
										length := len(strings.Split(apis.APIs[i].Checks[j].Token, ","))
										for k, v := range strings.Split(apis.APIs[i].Checks[j].Token, ",") {
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
								switch orgCheck[ogc] {
								case "body":
									if apis.APIs[i].Checks[j].Body != "" {
										count := strings.Count(apis.APIs[i].Checks[j].Body, "*")
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

										var replacedBody string
										for _, v := range placeHolders {
											replacedBody = strings.Replace(apis.APIs[i].Checks[j].Body, "*", string(v), 1)

											f.WriteString("  token.payload.roles[i].scope[j].orgId == input.parsed_body." + replacedBody + "\n")
										}
									}
								case "header":
									if apis.APIs[i].Checks[j].Header != "" {
										f.WriteString("  some j")
										f.WriteString("\n")
										f.WriteString("  token.payload.roles[i].scope[j].orgId == http_request.header." + apis.APIs[i].Checks[j].Header + "\n")
									}
								}
							case "ownerCheck":
								switch ownerCheck[oc] {
								case "body":
									f.WriteString("  sub[2] == input.parsed_body." + apis.APIs[i].Checks[j].Body + "\n")
								case "header":
									f.WriteString("  sub[2] == http_request.header." + apis.APIs[i].Checks[j].Header + "\n")
								}
								f.WriteString("}\n\n")
							}
						}
					} else {
						flag, roleCheck := false, false
						for j := 0; j < len(apis.APIs[i].Checks); j++ {
							for _, v := range apis.APIs[i].Checks {
								if strings.Contains(v.Key, "&&") && !strings.Contains(v.Key, "||") {
									flag = true
									break
								}
							}
							for _, v := range apis.APIs[i].Checks {
								if v.CheckType == "roleCheck" {
									roleCheck = true
									break
								}
							}
							if flag {
								if j == 0 {
									f.WriteString(apis.APIs[i].Name)
									f.WriteString(" {\n")
								}
								switch apis.APIs[i].Checks[j].CheckType {
								case "roleCheck":
									RoleCheck(apis.APIs[i].Checks[j], f)
								case "orgCheck":
									OrgCheck(apis.APIs[i].Checks[j], roleCheck, f)
								case "ownerCheck":
									switch ownerCheck[oc] {
									case "body":
										f.WriteString("  sub[2] == input.parsed_body." + apis.APIs[i].Checks[j].Body + "\n")
									case "header":
										f.WriteString("  sub[2] == http_request.header." + apis.APIs[i].Checks[j].Header + "\n")
									}
									f.WriteString("}\n\n")
								}
							}
						}
						break
					}
				}
			} else {
				for ogc := 0; ogc < len(orgCheck); ogc++ {
					if !orgAnd {
						for j := 0; j < len(apis.APIs[i].Checks); j++ {
							switch apis.APIs[i].Checks[j].CheckType {
							case "roleCheck":
								switch apis.APIs[i].Checks[j].Key {
								case "token":
									if apis.APIs[i].Checks[j].Token != "" {
										f.WriteString(apis.APIs[i].Name)
										f.WriteString(" {\n")
										f.WriteString("  api_roles := [")
										length := len(strings.Split(apis.APIs[i].Checks[j].Token, ","))
										for k, v := range strings.Split(apis.APIs[i].Checks[j].Token, ",") {
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
								switch orgCheck[ogc] {
								case "body":
									if apis.APIs[i].Checks[j].Body != "" {
										count := strings.Count(apis.APIs[i].Checks[j].Body, "*")
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

										var replacedBody string
										for _, v := range placeHolders {
											replacedBody = strings.Replace(apis.APIs[i].Checks[j].Body, "*", string(v), 1)

											f.WriteString("  token.payload.roles[i].scope[j].orgId == input.parsed_body." + replacedBody + "\n")
										}
									}
								case "header":
									if apis.APIs[i].Checks[j].Header != "" {
										f.WriteString("  some j")
										f.WriteString("\n")
										f.WriteString("  token.payload.roles[i].scope[j].orgId == http_request.header." + apis.APIs[i].Checks[j].Header + "\n")
									}
								}
							case "ownerCheck":
								OwnerCheck(apis.APIs[i].Checks[j], f)
								f.WriteString("}\n\n")
							}
						}
					} else {
						flag, roleCheck := false, false
						for j := 0; j < len(apis.APIs[i].Checks); j++ {
							for _, v := range apis.APIs[i].Checks {
								if strings.Contains(v.Key, "&&") && !strings.Contains(v.Key, "||") {
									flag = true
									break
								}
							}
							for _, v := range apis.APIs[i].Checks {
								if v.CheckType == "roleCheck" {
									roleCheck = true
									break
								}
							}
							if flag {
								if j == 0 {
									f.WriteString(apis.APIs[i].Name)
									f.WriteString(" {\n")
								}
								switch apis.APIs[i].Checks[j].CheckType {
								case "roleCheck":
									RoleCheck(apis.APIs[i].Checks[j], f)
								case "orgCheck":
									OrgCheck(apis.APIs[i].Checks[j], roleCheck, f)
								case "ownerCheck":
									switch ownerCheck[oc] {
									case "body":
										f.WriteString("  sub[2] == input.parsed_body." + apis.APIs[i].Checks[j].Body + "\n")
									case "header":
										f.WriteString("  sub[2] == http_request.header." + apis.APIs[i].Checks[j].Header + "\n")
									}
									f.WriteString("}\n\n")
								}
							}
						}
						break
					}
				}
			}
		}
	}
}
