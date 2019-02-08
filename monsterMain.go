package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"github.com/gorilla/mux" //for extensibility purposes
)

type UserInfo struct {
	Name    string `json:"name"`
	Uid     string `json:"uid"`
	Gid     string `json:"gid"`
	Comment string `json:"comment"`
	Home    string `json:"home"`
	Shell   string `json:"shell"`
}

type GroupInfo struct {
	Name    string   `json:"name"`
	Gid     string   `json:"gid"`
	Members []string `json:"members"`
}

type UserInfos []UserInfo
type GroupInfos []GroupInfo

func printJSON(w http.ResponseWriter, allEntries UserInfos) {
	//json.NewEncoder(w).Encode(allEntries)   //use this for non-pretty print
	jsonEntry := json.NewEncoder(w)
	jsonEntry.SetIndent("", "    ")
	jsonEntry.Encode(allEntries)
}

func getFileData(w http.ResponseWriter, filePath string) (csvData [][]string) {
	csvFile, err := os.Open(filePath)
	if err != nil {
		fmt.Println(err)
	}

	defer csvFile.Close()

	reader := csv.NewReader(csvFile)
	reader.Comma = ':'
	reader.FieldsPerRecord = -1

	csvData, err = reader.ReadAll()
	if err != nil {
		errorMsg := filePath + " may not have read access rights or does not exist" +
			" on this system"
		fmt.Fprintf(w, errorMsg)
		fmt.Println(err)
	}

	return csvData
}

func decodePasswd(w http.ResponseWriter, csvData [][]string, filePath string) (allEntries UserInfos) {

	var oneEntry UserInfo
	var lineNumber = 0
	for _, each := range csvData {

		lineNumber++
		if each[0][0] == '#' {
			continue
		}

		if len(each) != 7 {
			fmt.Fprintf(w, "Error! passwd file may be corrupt!"+
				" Found entry with %d fields on line:%d.", len(each), lineNumber)
			fmt.Println("Error!:", filePath, "file may be corrupt")
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
	return allEntries
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

func decodePasswdWithQuery(w http.ResponseWriter, csvData [][]string, filePath string, params UserInfo) (queriedEntries UserInfos) {
	var oneEntry UserInfo
	//var allEntries UserInfos
	var lineNumber = 0

	for _, each := range csvData {

		lineNumber++
		if each[0][0] == '#' {
			continue
		}

		if len(each) != 7 {
			fmt.Fprintf(w, "Error! passwd file may be corrupt!"+
				" Found entry with %d fields on line:%d.", len(each), lineNumber)
			fmt.Println("Error!:", filePath, "file may be corrupt")
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
	return queriedEntries

}

func allUserInfos(w http.ResponseWriter, r *http.Request) {

	fmt.Println("GET Endpoint Hit: /users")

	filePath := "/etc/passwd"
	csvData := getFileData(w, filePath)
	allEntries := decodePasswd(w, csvData, filePath)
	if allEntries != nil {
		printJSON(w, allEntries)
	}
}

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
		// isValid = false
	}

	return invalidStrings
}

func queryUserInfos(w http.ResponseWriter, r *http.Request) {
	fmt.Println("GET Endpoint Hit: /users/query")
	filePath := "/etc/passwd"

	urlQueryParams := r.URL.Query()

	invalidParams := validateUserParams(urlQueryParams)
	if len(invalidParams) != 0 {
		fmt.Fprintf(w, "Error! invalid query parameters given:", invalidParams)
		fmt.Println("Error! invalid query parameters given:", invalidParams)
        return
	}

	var queriedParams UserInfo
	queriedParams.Name = urlQueryParams.Get("name")
	queriedParams.Uid = urlQueryParams.Get("uid")
	queriedParams.Gid = urlQueryParams.Get("gid")
	queriedParams.Comment = urlQueryParams.Get("comment")
	queriedParams.Home = urlQueryParams.Get("home")
	queriedParams.Shell = urlQueryParams.Get("shell")

	csvData := getFileData(w, filePath)
	queriedEntries := decodePasswdWithQuery(w, csvData, filePath, queriedParams)
	if queriedEntries != nil && len(queriedEntries) != 0 {
		printJSON(w, queriedEntries)
	} else {
		fmt.Fprintf(w, " Unable to find matching entry with given query")
	}

}

func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Homepage Endpoint Hit")
	fmt.Println("Endpoint Hit: homePage")
}

func retrieveUserInfoFromUid(w http.ResponseWriter, csvData [][]string, filePath string, uid string) (matchingEntryPtr *UserInfo) {

	var lineNumber = 0
	var matchingEntry UserInfo
	matchingEntryPtr = nil
	for _, each := range csvData {

		lineNumber++
		if each[0][0] == '#' {
			continue
		}

		//NEED TO IMPLEMENT CHECK HERE FOR CORRUPT FILE!!
		// if len(each) != 7 {
		//     fmt.Fprintf(w, "Error! passwd file may be corrupt!" +
		//         " Found entry with %d fields on line:%d.", len(each), lineNumber)
		//     fmt.Println("Error!:", filePath,  "file may be corrupt")
		//     matchingEntries = nil
		//     break
		// }

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
	return matchingEntryPtr
}

func uidUser(w http.ResponseWriter, r *http.Request) {
	fmt.Println("GET Endpoint Hit: /users/{uid}")
	filePath := "/etc/passwd"

	csvData := getFileData(w, filePath)
	vars := mux.Vars(r)
	uid := vars["uid"]

	myUserInfo := retrieveUserInfoFromUid(w, csvData, filePath, uid)
	if myUserInfo != nil {
		// printJSON(w, myUserInfo)   //TODO - CLEANUP?
		jsonEntry := json.NewEncoder(w)
		jsonEntry.SetIndent("", "    ")
		jsonEntry.Encode(myUserInfo)
	} else {
		fmt.Fprintf(w, "404 page not found. \nUnable to find matching entry with uid="+uid)
	}
}

func retrieveGroupsFromUser(w http.ResponseWriter, csvData [][]string, filePath string, userName string) (groupEntries GroupInfos) {
	var oneEntry GroupInfo
	var foundMatch bool
	//var allEntries UserInfos
	var lineNumber = 0

	for _, each := range csvData {

		foundMatch = false
		lineNumber++
		if each[0][0] == '#' {
			continue
		}

		if len(each) != 4 {
			fmt.Fprintf(w, "Error! group file may be corrupt!"+
				" Found entry with %d fields on line:%d.", len(each), lineNumber)
			fmt.Println("Error!:", filePath, "file may be corrupt")
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
			fmt.Println(element)
			if element == userName {
				foundMatch = true
			}
		}

		if foundMatch == true {
			groupEntries = append(groupEntries, oneEntry)
		}

	}
	return groupEntries
}

func uidGroupInfo(w http.ResponseWriter, r *http.Request) {
	fmt.Println("GET Endpoint Hit: /users/{uid}/groups")
	passwdFilePath := "/etc/passwd"
	groupFilePath := "/etc/group"

	csvDataPasswd := getFileData(w, passwdFilePath)
	csvDataGroups := getFileData(w, groupFilePath)
	vars := mux.Vars(r)
	uid := vars["uid"]

	myUserInfo := retrieveUserInfoFromUid(w, csvDataPasswd, passwdFilePath, uid)
	if myUserInfo != nil {
		myUserName := (*myUserInfo).Name

		myUserGroupInfos := retrieveGroupsFromUser(w, csvDataGroups, groupFilePath, myUserName)
        fmt.Println(myUserGroupInfos)
		if myUserGroupInfos != nil && len(myUserGroupInfos) > 0{
			// printJSON(w, myUserInfo)   //TODO - CLEANUP?
			jsonEntry := json.NewEncoder(w)
			jsonEntry.SetIndent("", "    ")
			jsonEntry.Encode(myUserGroupInfos)
		} else {
            fmt.Fprintf(w, "No Groups found for given user")
        }
	} else {
		fmt.Fprintf(w, "404 page not found. \nUnable to find matching entry with uid="+uid + " in " + passwdFilePath)
	}
}

func decodeGroup(w http.ResponseWriter, csvData [][]string, filePath string) (groupEntries GroupInfos) {
	var oneEntry GroupInfo
	var lineNumber = 0

	for _, each := range csvData {

		lineNumber++
		if each[0][0] == '#' {
			continue
		}

		if len(each) != 4 {
			fmt.Fprintf(w, "Error! group file may be corrupt!"+
				" Found entry with %d fields on line:%d.", len(each), lineNumber)
			fmt.Println("Error!:", filePath, "file may be corrupt")
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
	return groupEntries
}

func allGroupInfos(w http.ResponseWriter, r *http.Request) {
	fmt.Println("GET Endpoint Hit: /groups")

	filePath := "/etc/group"
	csvData := getFileData(w, filePath)
	allGroupEntries := decodeGroup(w, csvData, filePath)
	if allGroupEntries != nil {
		//printJSON(w, allEntries) // TODO -CLEANUP
		jsonEntry := json.NewEncoder(w)
		jsonEntry.SetIndent("", "    ")
		jsonEntry.Encode(allGroupEntries)
	}
}

func retrieveGroupInfoFromGid(w http.ResponseWriter, csvData [][]string, filePath string, gid string) (matchingEntryPtr *GroupInfo) {

	var lineNumber = 0
	var matchingEntry GroupInfo
	matchingEntryPtr = nil
	for _, each := range csvData {

		lineNumber++
		if each[0][0] == '#' {
			continue
		}

		//NEED TO IMPLEMENT CHECK HERE FOR CORRUPT FILE!!
		// if len(each) != 7 {
		//     fmt.Fprintf(w, "Error! passwd file may be corrupt!" +
		//         " Found entry with %d fields on line:%d.", len(each), lineNumber)
		//     fmt.Println("Error!:", filePath,  "file may be corrupt")
		//     matchingEntries = nil
		//     break
		// }

		if each[2] == gid {
			matchingEntry.Name = each[0]
			matchingEntry.Gid = each[2]
			matchingEntry.Members = strings.Split(each[3], ",")
			matchingEntryPtr = &matchingEntry
			break
		}
	}
	return matchingEntryPtr
}

func gidGroup(w http.ResponseWriter, r *http.Request) {
	fmt.Println("GET Endpoint Hit: /groups/{gid}")
	groupFilePath := "/etc/group"
	csvDataGroups := getFileData(w, groupFilePath)
	vars := mux.Vars(r)
	gid := vars["gid"]

	myGroupInfo := retrieveGroupInfoFromGid(w, csvDataGroups, groupFilePath, gid)
	if myGroupInfo != nil {
		// printJSON(w, myUserInfo)   //TODO - CLEANUP?
		jsonEntry := json.NewEncoder(w)
		jsonEntry.SetIndent("", "    ")
		jsonEntry.Encode(myGroupInfo)
	} else {
		fmt.Fprintf(w, "404 page not found. \nUnable to find matching entry with gid="+gid)
	}
}
func validateGroupParams(params map[string][]string) (invalidStrings []string) {
	validParams := [3]string{"name", "gid", "member"}
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

func compareGroupQuery (params GroupInfo, dataRecord GroupInfo) (isMatch bool) {

    if (params.Name != ""){
        if (params.Name != dataRecord.Name) {
            return false
        }
    }
    if (params.Gid != ""){
        if (params.Gid != dataRecord.Gid) {
            return false
        }
    }

    isFound := false
    if (len(params.Members) != 0){
        for _, paramMember := range params.Members {
            isFound = false
            for _, dataRecordMember := range dataRecord.Members {
                if (dataRecordMember == paramMember) {
                    isFound = true
                }
            }
            if isFound == false{
                return false
            }
        }
    }

    return true
}

func decodeGroupWithQuery(w http.ResponseWriter, csvData [][]string, filePath string, params GroupInfo) (queriedEntries GroupInfos) {
	var oneEntry GroupInfo
	//var allEntries UserInfos
	var lineNumber = 0

	for _, each := range csvData {

		lineNumber++
		if each[0][0] == '#' {
			continue
		}

		if len(each) != 4 {
			fmt.Fprintf(w, "Error! group file may be corrupt!"+
				" Found entry with %d fields on line:%d.", len(each), lineNumber)
			fmt.Println("Error!:", filePath, "file may be corrupt")
			queriedEntries = nil
			break
		}
		oneEntry.Name = each[0]
		oneEntry.Gid = each[2]
		oneEntry.Members = strings.Split(each[3], ",")

		isMatch := compareGroupQuery(params, oneEntry)

		if (isMatch){
		    queriedEntries = append(queriedEntries, oneEntry)
		}
	}
	return queriedEntries

}

func queryGroupInfos(w http.ResponseWriter, r *http.Request) {
	fmt.Println("GET Endpoint Hit: /groups/query")
	groupFilePath := "/etc/group"

	urlQueryParams := r.URL.Query()

	invalidParams := validateGroupParams(urlQueryParams)
	if len(invalidParams) != 0 {
		fmt.Fprintf(w, "Error! invalid query parameters given:", invalidParams)
		fmt.Println("Error! invalid query parameters given:", invalidParams)
		return
	}

	var queriedParams GroupInfo
    var memberValues []string
	queriedParams.Name = urlQueryParams.Get("name")
	queriedParams.Gid = urlQueryParams.Get("gid")

    // fmt.Println(urlQueryParams)
    for mapKey, mapValue := range urlQueryParams {
        if mapKey == "member"{
            memberValues = mapValue
        }
    }
    queriedParams.Members = memberValues
    // fmt.Println(queriedParams.Members)

	fmt.Println("Hello?")

	csvData := getFileData(w, groupFilePath)
	queriedEntries := decodeGroupWithQuery(w, csvData, groupFilePath, queriedParams)
	if queriedEntries != nil && len(queriedEntries) != 0 {
		//printJSON(w, queriedEntries)  //TODO - CLEANUP?
		jsonEntry := json.NewEncoder(w)
		jsonEntry.SetIndent("", "    ")
		jsonEntry.Encode(queriedEntries)
	} else {
		fmt.Fprintf(w, " Unable to find matching entry with given query")
	}

}

func handleRequests() {

	myRouter := mux.NewRouter().StrictSlash(false)

	myRouter.HandleFunc("/", homePage)
	myRouter.HandleFunc("/users", allUserInfos).Methods("GET")
	myRouter.HandleFunc("/users/query", queryUserInfos).Methods("GET")
	myRouter.HandleFunc("/users/{uid}", uidUser).Methods("GET")


	myRouter.HandleFunc("/users/{uid}/groups", uidGroupInfo).Methods("GET")
	myRouter.HandleFunc("/groups", allGroupInfos).Methods("GET")
    myRouter.HandleFunc("/groups/query", queryGroupInfos).Methods("GET")
	myRouter.HandleFunc("/groups/{gid}", gidGroup).Methods("GET")


	log.Fatal(http.ListenAndServe(":8081", myRouter))
}

func main() {
	handleRequests()
}
