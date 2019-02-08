
package main

import (
    "fmt"
    "encoding/json"
    "encoding/csv"
    "os"
    "errors"
    "net/http"
    "strconv"
    "strings"
)
func printJSON(w http.ResponseWriter, allEntries interface{}) {
	jsonEntry := json.NewEncoder(w)
	jsonEntry.SetIndent("", "    ")    //For pretty print
	jsonEntry.Encode(allEntries)
}

func getFileData( filePath string) (csvData [][]string, err error) {
	csvFile, err := os.Open(filePath)
	if err != nil {
		fmt.Println(err)
	}

	defer csvFile.Close()

	reader := csv.NewReader(csvFile)
	reader.Comma = ':'
	reader.FieldsPerRecord = -1

	return reader.ReadAll()
}


func decodePasswd(csvData [][]string) (allEntries UserInfos, err error) {

	var oneEntry UserInfo
	for index, each := range csvData {

		if each[0][0] == '#' {
			continue
		}

		if len(each) != 7 {
            errString := "Error! passwd file may be corrupt! Found entry with " +
                          strconv.Itoa(len(each)) + " fields on line:" + strconv.Itoa(index+1)
			err = errors.New(errString)
			allEntries = nil
			break
		}
		oneEntry.Name = each[0]
		oneEntry.Uid = each[2]
		oneEntry.Gid = each[3]
		oneEntry.Comment = each[4]
		oneEntry.Home = each[5]
		oneEntry.Shell = each[6]
		allEntries = append(allEntries, oneEntry)
	}
	return allEntries, err
}

func decodePasswdWithQuery( csvData [][]string, params UserInfo) (queriedEntries UserInfos, err error) {
	var oneEntry UserInfo

	for index, each := range csvData {

		if each[0][0] == '#' {
			continue
		}

        err = validateEntryInPasswdFile(len(each), index)
        if err != nil {
            queriedEntries = nil
            break
        }

		oneEntry.Name = each[0]
		oneEntry.Uid = each[2]
		oneEntry.Gid = each[3]
		oneEntry.Comment = each[4]
		oneEntry.Home = each[5]
		oneEntry.Shell = each[6]

		isMatch := compareUserInfo(params, oneEntry)

		if isMatch {
			queriedEntries = append(queriedEntries, oneEntry)
		}
	}
	return queriedEntries, err

}

func compareUserInfo(params UserInfo, dataRecord UserInfo) (isMatch bool) {

	if params.Name != "" {
		if params.Name != dataRecord.Name {
			return false
		}
	}
	if params.Uid != "" {
		if params.Uid != dataRecord.Uid {
			return false
		}
	}
	if params.Gid != "" {
		if params.Gid != dataRecord.Uid {
			return false
		}
	}
	if params.Comment != "" {
		if params.Comment != dataRecord.Comment {
			return false
		}
	}
	if params.Home != "" {
		if params.Home != dataRecord.Home {
			return false
		}
	}
	if params.Shell != "" {
		if params.Shell != dataRecord.Shell {
			return false
		}
	}

	return true
}

func retrieveUserInfoFromUid(csvData [][]string, uid string) (matchingEntryPtr *UserInfo, err error) {

	var matchingEntry UserInfo
	matchingEntryPtr = nil

	for index, each := range csvData {

		if each[0][0] == '#' {
			continue
		}

        err = validateEntryInPasswdFile(len(each), index)
        if err != nil {
            matchingEntryPtr = nil
            break
        }

		if each[2] == uid {
			matchingEntry.Name = each[0]
			matchingEntry.Uid = each[2]
			matchingEntry.Gid = each[3]
			matchingEntry.Comment = each[4]
			matchingEntry.Home = each[5]
			matchingEntry.Shell = each[6]
			matchingEntryPtr = &matchingEntry
			break
		}
	}
	return matchingEntryPtr, err
}

func retrieveGroupsFromUser( csvData [][]string, userName string) (groupEntries GroupInfos, err error) {
	var oneEntry GroupInfo
	var foundMatch bool

	for index, each := range csvData {

		foundMatch = false
		if each[0][0] == '#' {
			continue
		}

        err = validateEntryInGroupFile(len(each), index)
        if err != nil {
            groupEntries = nil
            break
        }

		oneEntry.Name = each[0]
		oneEntry.Gid = each[2]
		// Not all Gid matches from /etc/passwd show up in /etc/group file...
		// if we want to add the primary Gid from /etc/passwd, we will need to compare
		// it with this Gid and if they match, include it in the groupEntries

		oneEntry.Members = strings.Split(each[3], ",")

		for _, element := range oneEntry.Members {
			if element == userName {
				foundMatch = true
			}
		}

		if foundMatch == true {
			groupEntries = append(groupEntries, oneEntry)
		}

	}
	return groupEntries, err
}

func decodeGroup( csvData [][]string) (groupEntries GroupInfos, err error) {
	var oneEntry GroupInfo

	for index, each := range csvData {

		if each[0][0] == '#' {
			continue
		}

        err = validateEntryInGroupFile(len(each), index)
        if err != nil {
            groupEntries = nil
            break
        }

		oneEntry.Name = each[0]
		oneEntry.Gid = each[2]
		// Not all Gid matches from /etc/passwd show up in /etc/group file...
		// if we want to add the primary Gid from /etc/passwd, we will need to compare
		// it with this Gid and if they match, include it in the groupEntries

		oneEntry.Members = strings.Split(each[3], ",")

		groupEntries = append(groupEntries, oneEntry)

	}
	return groupEntries, err
}
