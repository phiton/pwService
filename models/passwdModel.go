package main

import (
    "fmt"

)

type UserInfo struct {
    Name string `json:"name"`
    Uid string `json:"uid"`
    Gid string `json:"gid"`
    Comment string `json:"comment"`
    Home string `json:"home"`
    Shell string `json:"shell"`
}
