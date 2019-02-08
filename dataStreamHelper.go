
package main

import (
    "fmt"
    "encoding/json"
    "encoding/csv"
    "os"
    "errors"
    "net/http"
    "strconv"
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

		if len(each) != 7 {
			errString :=  "Error! passwd file may be corrupt! Found entry with " +
                           strconv.Itoa(len(each)) + "fields on line:" + strconv.Itoa(index+1)
			err = errors.New(errString)
			queriedEntries = nil
			break
		}
		oneEntry.Name = each[0]
		oneEntry.Uid = each[2]
		oneEntry.Gid = each[3]
		oneEntry.Comment = each[4]
		oneEntry.Home = each[5]
		oneEntry.Shell = each[6]

		isMatch := compareUserQuery(params, oneEntry)

		if isMatch {
			queriedEntries = append(queriedEntries, oneEntry)
		}
	}
	return queriedEntries, err

}

func compareUserQuery(params UserInfo, dataRecord UserInfo) (isMatch bool) {

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
