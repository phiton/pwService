package main

import (

)

func validateUserParams(params map[string][]string) (invalidStrings []string) {
	validParams := [6]string{"name", "uid", "gid", "comment", "home", "shell"}
	isValid := false

	for mapKey := range params {
        isValid = false
		for _, validParam := range validParams {
			if mapKey == validParam {
				isValid = true
			}
		}
		if isValid == false {
			invalidStrings = append(invalidStrings, mapKey)
		}
	}

	return invalidStrings
}
