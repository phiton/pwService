package main

import (
    "flag"
    "fmt"
    "os"
)

func main() {
    argsWithoutProg := os.Args[1:]
    if len(argsWithoutProg) > 2 {
        fmt.Println("Received too many arguments, only passwd and group file paths are configurable")
        return
    }

    if len(argsWithoutProg) == 1 && argsWithoutProg[0] == "-h" {
        fmt.Println("This tool is used to expose a systems passwd and group via HTTP servuce")
        fmt.Println("command to run with specified locations:  ./main.go -p=<passwdFileLocation> -h=<groupFileLocation>")
        fmt.Println("       ex: ./main.go -p=/etc/myPasswd -h=/etc/myGroup")
        fmt.Println("If inputs are not specified, the default locations will be /etc/passwd and /etc/groups")

        return
    }

    passwdLocationPtr := flag.String("p", "/etc/passwd", "location of the passwd file")
    groupLocationPtr := flag.String("g", "/etc/group", "location of the group file")
    flag.Parse()

    var passwd string = "/etc/passwd"
    var group string = "/etc/group"

    if *passwdLocationPtr != ""{
        passwd = *passwdLocationPtr
    }

    if *groupLocationPtr != ""{
        group = *groupLocationPtr
    }

    a := App{}
    fmt.Println("passwd location set to:" +  passwd +
               "\ngroup location set to:" + group)
    a.Initialize(passwd, group)

    fmt.Println ("HTTP Service is up and running on localhost:8081")
    a.Run(":8081")
}
