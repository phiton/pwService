package main

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
