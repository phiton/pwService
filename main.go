package main

import
(
    "fmt"
    "log"
    "net/http"
    "encoding/json"
    // "encoding/csv"
    "github.com/gorilla/mux"
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

func allUserInfos(w http.ResponseWriter, r *http.Request) {
    UserInfos :=UserInfos {
        UserInfo {Name:"Test Title", Uid: "Test Description", Gid: "Hello World", Comment: "hah", Home: "home", Shell: "shell"},
    }

    fmt.Println("Endpoint Hit: All UserInfos Endpoint")
    json.NewEncoder(w).Encode(UserInfos)
}

func homePage(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "Homepage Endpoint Hit")
    fmt.Println("Endpoint Hit: homePage")
}

func handleRequests() {

    myRouter  := mux.NewRouter().StrictSlash(true)

    myRouter.HandleFunc("/", homePage)
    myRouter.HandleFunc("/UserInfos", allUserInfos).Methods("GET")
    log.Fatal(http.ListenAndServe(":8081", myRouter))
}

func main() {
    handleRequests()
}
