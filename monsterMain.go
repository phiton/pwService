package main

import (
    "fmt"
    "log"
    "net/http"
    "encoding/json"
    "encoding/csv"
    "github.com/gorilla/mux"  //for extensibility purposes
    "os"
    // "strconv"
)

type UserInfo struct {
    Name string `json:"name"`
    Uid string `json:"uid"`
    Gid string `json:"gid"`
    Comment string `json:"comment"`
    Home string `json:"home"`
    Shell string `json:"shell"`
}

type UserInfos [] UserInfo

func printJSON(w http.ResponseWriter, allEntries UserInfos ) {
    //json.NewEncoder(w).Encode(allEntries)   //use this for non-pretty print
    jsonEntry := json.NewEncoder(w)
    jsonEntry.SetIndent("", "    ")
    jsonEntry.Encode(allEntries)
}

func getFileData(w http.ResponseWriter, filePath string) (csvData [][]string) {
    csvFile, err := os.Open (filePath)
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

func marshalPasswd(w http.ResponseWriter, csvData [][]string, filePath string) (allEntries UserInfos){

    var oneEntry UserInfo
    var lineNumber = 0
    for _, each := range csvData {

        lineNumber++
        if each[0][0] == '#' {
            continue
        }

        if len(each) != 7 {
            fmt.Fprintf(w, "Error! passwd file may be corrupt!" +
                " Found entry with %d fields on line:%d.", len(each), lineNumber)
            fmt.Println("Error!:", filePath,  "file may be corrupt")
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

func compareWithQuery (params UserInfo, dataRecord UserInfo) (isMatch bool) {

    if (params.Name != ""){
        if (params.Name != dataRecord.Name) {
            return false
        }
    }
    if (params.Uid != ""){
        if (params.Uid != dataRecord.Uid) {
            return false
        }
    }
    if (params.Gid != ""){
        if (params.Gid != dataRecord.Uid) {
            return false
        }
    }
    if (params.Comment != ""){
        if (params.Comment != dataRecord.Comment) {
            return false
        }
    }
    if (params.Home != ""){
        if (params.Home != dataRecord.Home) {
            return false
        }
    }
    if (params.Shell != ""){
        if (params.Shell != dataRecord.Shell) {
            return false
        }
    }

    return true
}

func marshalPasswdWithQuery(w http.ResponseWriter, csvData [][]string, filePath string ,params UserInfo) (queriedEntries UserInfos) {
    var oneEntry UserInfo
    //var allEntries UserInfos
    var lineNumber = 0

    for _, each := range csvData {

        lineNumber++
        if each[0][0] == '#' {
            continue
        }

        if len(each) != 7 {
            fmt.Fprintf(w, "Error! passwd file may be corrupt!" +
                " Found entry with %d fields on line:%d.", len(each), lineNumber)
            fmt.Println("Error!:", filePath,  "file may be corrupt")
            queriedEntries = nil
            break
        }
        oneEntry.Name = each[0]
        oneEntry.Uid = each[2]
        oneEntry.Gid = each[3]
        oneEntry.Comment = each[4]
        oneEntry.Home = each[5]
        oneEntry.Shell = each[6]

        isMatch := compareWithQuery(params, oneEntry)

        if (isMatch){
            queriedEntries = append(queriedEntries, oneEntry)
        }
    }
    return queriedEntries

}

func allUserInfos(w http.ResponseWriter, r *http.Request) {

    fmt.Println("GET Endpoint Hit: /users")

    filePath := "/etc/passwd"
    csvData := getFileData(w, filePath)
    allEntries := marshalPasswd(w, csvData, filePath)
    if (allEntries != nil) {
        printJSON(w, allEntries)
    }
}

func validateParams(params map[string][]string) (invalidStrings []string) {
    validParams := [6]string {"name", "uid", "gid", "comment", "home", "shell"}
    isValid := false

    for mapKey := range params {
        for _, validParam := range validParams {
            if mapKey == validParam {
                isValid = true
            }
        }
        if isValid == false {
            invalidStrings = append(invalidStrings, mapKey)
        }
        isValid = false
    }

    return invalidStrings
}

func queryUserInfos(w http.ResponseWriter, r *http.Request) {
    fmt.Println("GET Endpoint Hit: /users/query")
    filePath := "/etc/passwd"

    urlQueryParams := r.URL.Query()
    // fmt.Println(urlQueryParams)
    // fmt.Println(urlQueryParams.Get("name"))
    // fmt.Println("uid =" + urlQueryParams.Get("uid"))

    invalidParams := validateParams(urlQueryParams)
    if len(invalidParams) !=0 {
        fmt.Fprintf(w, "Error! invalid query parameters given:", invalidParams)
        fmt.Println("Error! invalid query parameters given:" , invalidParams)
    }

    var queriedParams UserInfo
    queriedParams.Name = urlQueryParams.Get("name")
    queriedParams.Uid = urlQueryParams.Get("uid")
    queriedParams.Gid = urlQueryParams.Get("gid")
    queriedParams.Comment = urlQueryParams.Get("comment")
    queriedParams.Home = urlQueryParams.Get("home")
    queriedParams.Shell = urlQueryParams.Get("shell")

    csvData := getFileData(w, filePath)
    queriedEntries := marshalPasswdWithQuery(w, csvData, filePath, queriedParams)
    if (queriedEntries != nil && len(queriedEntries) != 0) {
        printJSON(w, queriedEntries)
    } else {
        fmt.Fprintf(w, " Unable to find matching entry with given query")
    }

}

func homePage(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "Homepage Endpoint Hit")
    fmt.Println("Endpoint Hit: homePage")
}

func retrieveUserInfo (w http.ResponseWriter, csvData [][]string, filePath string ,uid string) (matchingEntries UserInfos) {

    var lineNumber = 0
    var matchingEntry UserInfo
    for _, each := range csvData {

        lineNumber++
        if each[0][0] == '#' {
            continue
        }

        if len(each) != 7 {
            fmt.Fprintf(w, "Error! passwd file may be corrupt!" +
                " Found entry with %d fields on line:%d.", len(each), lineNumber)
            fmt.Println("Error!:", filePath,  "file may be corrupt")
            matchingEntries = nil
            break
        }

        if (each[2] == uid){
            matchingEntry.Name = each[0]
            matchingEntry.Uid = each[2]
            matchingEntry.Gid = each[3]
            matchingEntry.Comment = each[4]
            matchingEntry.Home = each[5]
            matchingEntry.Shell = each[6]
            matchingEntries = append(matchingEntries, matchingEntry)
            break
        }
    }
    return matchingEntries
}

func uidUser(w http.ResponseWriter, r *http.Request) {
    fmt.Println("GET Endpoint Hit: /users/{uid}")
    filePath := "/etc/passwd"

    csvData := getFileData(w, filePath)
    vars := mux.Vars(r)
    uid := vars["uid"]

    myUserInfo := retrieveUserInfo(w, csvData, filePath, uid )
    if (myUserInfo != nil) && (len(myUserInfo) > 0){
        printJSON(w, myUserInfo)
    } else {
        fmt.Fprintf(w, " Unable to find matching entry with uid=" + uid)
    }
}

func handleRequests() {

    myRouter  := mux.NewRouter().StrictSlash(true)

    myRouter.HandleFunc("/", homePage)
    myRouter.HandleFunc("/users", allUserInfos).Methods("GET")
    myRouter.HandleFunc("/users/{uid}", uidUser).Methods("GET")
    myRouter.HandleFunc("/users/query", queryUserInfos).Methods("GET")
    log.Fatal(http.ListenAndServe(":8081", myRouter))
}

func main() {
    handleRequests()
}
