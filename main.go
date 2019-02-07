package main

import (
    "fmt"
    "log"
    "net/http"
    "encoding/json"
    "encoding/csv"
    "github.com/gorilla/mux"  //for extensibility purposes
    "os"
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

func marshalPasswdFromFile(w http.ResponseWriter, filePath string) (allEntries UserInfos){
    csvFile, err := os.Open (filePath)
    if err != nil {
        fmt.Println(err)
    }

    defer csvFile.Close()

    reader := csv.NewReader(csvFile)
    reader.Comma = ':'
    reader.FieldsPerRecord = -1

    csvData, err := reader.ReadAll()
    if err != nil {
        errorMsg := filePath + " may not have read access rights or does not exist" +
                    " on this system"
        fmt.Fprintf(w, errorMsg)
        fmt.Println(err)
    }

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

func allUserInfos(w http.ResponseWriter, r *http.Request) {

    fmt.Println("GET Endpoint Hit: /users")
    allEntries := marshalPasswdFromFile(w, "/etc/passwd")
    if (allEntries != nil) {
        printJSON(w, allEntries)
    }
}

func homePage(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "Homepage Endpoint Hit")
    fmt.Println("Endpoint Hit: homePage")
}

func handleRequests() {

    myRouter  := mux.NewRouter().StrictSlash(true)

    myRouter.HandleFunc("/", homePage)
    myRouter.HandleFunc("/users", allUserInfos).Methods("GET")

    log.Fatal(http.ListenAndServe(":8081", myRouter))
}

func main() {
    handleRequests()
}
