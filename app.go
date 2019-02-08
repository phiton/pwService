package main

import (
    "net/http"
    "log"
    "fmt"
    "strings"
    "github.com/gorilla/mux" //for extensibility purposes
)

type App struct {
    Router *mux.Router
    PasswordPath string
    GroupPath string
}

func (a *App) Run(addr string) {
	log.Fatal(http.ListenAndServe(addr, a.Router))
}

func (a *App) Initialize(passwdPath string, groupPath string) {

	a.Router = mux.NewRouter().StrictSlash(false)
    a.PasswordPath = passwdPath
    a.GroupPath = groupPath
    a.InitializeRoutes()


}

func (a *App) InitializeRoutes() {
    a.Router.HandleFunc("/", a.getHomePage)
    a.Router.HandleFunc("/users", a.getAllUserInfos).Methods("GET")
    a.Router.HandleFunc("/users/query", a.getQueryUserInfos).Methods("GET")
    // a.Router.HandleFunc("/users/{uid}", uidUser).Methods("GET")
    //
    // a.Router.HandleFunc("/users/{uid}/groups", uidGroupInfo).Methods("GET")
    // a.Router.HandleFunc("/groups", allGroupInfos).Methods("GET")
    // a.Router.HandleFunc("/groups/query", queryGroupInfos).Methods("GET")
    // a.Router.HandleFunc("/groups/{gid}", gidGroup).Methods("GET")
}

func (a *App) getHomePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Homepage Endpoint Hit")
	fmt.Println("Endpoint Hit: homePage")
}

func (a *App)getAllUserInfos(w http.ResponseWriter, r *http.Request) {

	fmt.Println("GET Endpoint Hit: /users")

	passwdData ,err := getFileData(a.PasswordPath)
    if err != nil {
		errorMsg := a.PasswordPath + " file does not exist or can't be read" +
			" on this system"
		fmt.Fprintf(w, errorMsg)
		fmt.Println(err)
        return
	}

	allEntries, err := decodePasswd(passwdData)
    if err != nil {
        fmt.Fprintf(w, err.Error())
    }else if allEntries != nil {
		printJSON(w, allEntries)
	}
}

func (a *App)getQueryUserInfos(w http.ResponseWriter, r *http.Request) {
	fmt.Println("GET Endpoint Hit: /users/query")

	urlQueryParams := r.URL.Query()

	invalidParams := validateUserParams(urlQueryParams)
	if len(invalidParams) != 0 {
        invalidParamsString := strings.Join(invalidParams, " ")
		fmt.Fprintf(w, "Error! invalid query parameters given: " + invalidParamsString)
		fmt.Println("Error! invalid query parameters given: " + invalidParamsString)
        return
	}

	var queriedParams UserInfo
	queriedParams.Name = urlQueryParams.Get("name")
	queriedParams.Uid = urlQueryParams.Get("uid")
	queriedParams.Gid = urlQueryParams.Get("gid")
	queriedParams.Comment = urlQueryParams.Get("comment")
	queriedParams.Home = urlQueryParams.Get("home")
	queriedParams.Shell = urlQueryParams.Get("shell")

	csvData, err := getFileData(a.PasswordPath)
    if err != nil {
		errorMsg := a.PasswordPath + " file does not exist or can't be read" +
			" on this system"
		fmt.Fprintf(w, errorMsg)
		fmt.Println(err)
        return
	}

	queriedEntries, err := decodePasswdWithQuery(csvData, queriedParams)
    if err != nil {
        fmt.Fprintf(w, err.Error())
    } else if queriedEntries != nil && len(queriedEntries) != 0 {
		printJSON(w, queriedEntries)
	} else {
		fmt.Fprintf(w, " Unable to find matching entry with given query")
	}

}
