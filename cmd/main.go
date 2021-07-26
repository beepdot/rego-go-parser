package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"rego-go-parser/pkg/checkers"
	"rego-go-parser/pkg/helpers"
	"rego-go-parser/pkg/schema"
)

func main() {
	jsonFile, err := os.Open("../input/checks.json")
	if err != nil {
		panic(err)
	}

	fmt.Println("Successfully Opened checks.json")

	defer jsonFile.Close()
	byteValue, _ := ioutil.ReadAll(jsonFile)

	var apis schema.APIs
	json.Unmarshal(byteValue, &apis)

	f, err := os.OpenFile("../output/policy.rego", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	orgCheckKeys, ownerCheckKeys, orgCheckContainsAnd, ownerCheckContainsAnd := helpers.IdentifyKeys(&apis)

	for i := 0; i < len(apis.APIs); i++ {
		for j := 0; j < len(orgCheckKeys); j++ {
			if !orgCheckContainsAnd {
				for k := 0; k < len(ownerCheckKeys); k++ {
					if !ownerCheckContainsAnd {
						for l := 0; l < len(apis.APIs[i].Checks); l++ {
							checkers.AllChecksContainingOnlyOr(apis.APIs[i], apis.APIs[i].Checks[l], orgCheckKeys[j], ownerCheckKeys[k], f)
						}
					} else {
						checkers.OnlyOrgCheckContainsOr(apis.APIs[i], apis.APIs[i], orgCheckKeys[j], ownerCheckKeys[k], f)
					}
				}
			} else {
				for k := 0; k < len(ownerCheckKeys); k++ {
					if !ownerCheckContainsAnd {
						checkers.OnlyOwnerCheckContainsOr(apis.APIs[i], apis.APIs[i], orgCheckKeys[j], ownerCheckKeys[k], f)
					} else {
						checkers.NothingContainsOr(apis.APIs[i], apis.APIs[i], orgCheckKeys[j], ownerCheckKeys[k], f)
					}
				}
				break
			}
		}
	}
}
