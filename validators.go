package main

import (
    "strconv"
    "errors"
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

func validateEntryInPasswdFile (entryLength int, index int) (err error) {
    err = nil
    if entryLength != 7 {
        errString :=  "Error! passwd file may be corrupt! Found entry with " +
                       strconv.Itoa(entryLength) + "fields on line:" + strconv.Itoa(index+1)
        err = errors.New(errString)
    }
    return err
}
