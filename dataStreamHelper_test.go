package main

import (
    "testing"
    //"fmt"
    "strconv"
)

var passwdTestPath string = "./testFiles/passwdTest"
var malformedPasswdTestPath string = "./testFiles/malformedPasswdTest"
var groupTestPath string = "./testFiles/groupTest"
var malformedGroupTestPath string = "./testFiles/malformedGroupTest"
var fileNotExistPath string = "./testFiles/fileNotExist"

func TestGetFileData(t* testing.T) {

    passwdData, err := getFileData(fileNotExistPath)
    if err == nil {
        t.Errorf("Processed file "+ fileNotExistPath  +", which should not exist")
    }
    if (len(passwdData) != 0 ) {
        t.Errorf("Total file line size mismatch. Expected: 0, Actual:"+ strconv.Itoa(len(passwdData)))
    }

    passwdData ,err = getFileData(passwdTestPath)
    if err != nil {
        t.Errorf ("Encountered error reading file " + passwdTestPath + ", which should exist")
    }
    if (len(passwdData) != 4) {
        t.Errorf ("Total file line size mismatch. Expected: 4, Actual:" + strconv.Itoa(len(passwdData)))
    }

}

func TestDecodePasswd(t* testing.T) {
    passwdData, err := getFileData(passwdTestPath)
    allEntries, err := decodePasswd(passwdData)
    if err != nil {
        t.Errorf ("Unexpected error hit when decoding test data: " + err.Error())
    }
    if len(allEntries) != 2 {
        t.Errorf ("passwdTest file gave back more entries than expected. Expected: 2, Actual:" + strconv.Itoa(len(allEntries)))
    }

    passwdData, err = getFileData(malformedPasswdTestPath)
    allEntries, err = decodePasswd(passwdData)
    if err == nil {
        t.Errorf ("Error about corrupted passwd did not hit as expected")
    }
    if len(allEntries) != 0 {
        t.Errorf ("Number of entries mismatched. Expected: 0, Actual: " + strconv.Itoa(len(allEntries)))
    }
}

func TestDecodePasswdWithQuery(t* testing.T) {
    var queriedParams UserInfo
    queriedParams.Name = "nobody"
    queriedParams.Uid = "-2"
    queriedParams.Gid = "-2"
    queriedParams.Comment = "Unprivileged User"
    queriedParams.Home = "/var/empty"
    queriedParams.Shell = "/usr/bin/false"

    passwdData, err := getFileData(passwdTestPath)

    queriedEntries, err := decodePasswdWithQuery(passwdData, queriedParams)
    if err != nil || queriedEntries == nil{
        t.Errorf(err.Error())
    }

    if (len(queriedEntries) != 1) {
        t.Errorf("Number of query matches incorrect. Expected: 1, Actual:" + strconv.Itoa(len(queriedEntries)))
    }
}

func TestCompareUserInfo(t* testing.T) {
    var oneInfo UserInfo
    var twoInfo UserInfo
    var twoInfoCopy UserInfo

    oneInfo.Name = "name1"
    oneInfo.Uid = "1"
    oneInfo.Gid = "1"
    oneInfo.Comment = "Unprivileged User"
    oneInfo.Home = "/var/empty"
    oneInfo.Shell = "/usr/bin/false"

    twoInfo.Name = "name2"
    twoInfo.Uid = "2"
    twoInfo.Gid = "2"
    twoInfo.Comment = "number2"
    twoInfo.Home = "/var/two"
    twoInfo.Shell = "/usr/bin/two"

    twoInfoCopy.Name = "name2"
    twoInfoCopy.Uid = "2"
    twoInfoCopy.Gid = "2"
    twoInfoCopy.Comment = "number2"
    twoInfoCopy.Home = "/var/two"
    twoInfoCopy.Shell = "/usr/bin/two"

    isMatch := compareUserInfo(oneInfo, twoInfo)
    if (isMatch == true) {
        t.Errorf("Received inappropriate match with two different UserInfo structs")
    }
    isMatch = compareUserInfo(twoInfo, twoInfoCopy)
    if (isMatch == false) {
        t.Errorf("Received inappropriate mismatch with two of the same UserInfo structs")
    }

}

func TestRetrieveUserInfoFromUid(t* testing.T) {
    passwdData, _ := getFileData(passwdTestPath)
    uid := "-2"
    myUserInfo, _ := retrieveUserInfoFromUid(passwdData, uid)

    if myUserInfo.Name != "nobody"{
        t.Errorf("Incorrect name picked up from passwdTestPath. Expected: nobody, Actual: " + myUserInfo.Name)
    }
    if myUserInfo.Uid != "-2" {
        t.Errorf("Incorrect UID picked up from passwdTestPath. Expected: -2, Actual: " + myUserInfo.Uid)
    }
    if myUserInfo.Gid != "-2"{
        t.Errorf("Incorrect GID picked up from passwdTestPath. Expected: -2, Actual: " + myUserInfo.Gid)
    }
    if myUserInfo.Comment != "Unprivileged User" {
        t.Errorf("Incorrect Comment picked up from passwdTestPath. Expected: Unprivileged User, Actual: " + myUserInfo.Comment)
    }
    if myUserInfo.Home != "/var/empty"{
        t.Errorf("Incorrect Home picked up from passwdTestPath. Expected: /var/empty, Actual: " + myUserInfo.Home)
    }
    if myUserInfo.Shell != "/usr/bin/false" {
        t.Errorf("Incorrect Shell picked up from passwdTestPath. Expected: /usr/bin/false, Actual: " + myUserInfo.Shell)
    }
}

func TestRetrieveGroupsFromUser(t* testing.T) {
    groupData, _ := getFileData(groupTestPath)
    myUserName := "root"
    myUserGroupInfos, err := retrieveGroupsFromUser(groupData, myUserName)
    if (len(myUserGroupInfos) != 2){
        t.Errorf("Did not receive appropriate count for username root. Expected: 2, Actual: " + strconv.Itoa(len(myUserGroupInfos)))
    }

    groupData, _ = getFileData(malformedGroupTestPath)
    myUserGroupInfos, err = retrieveGroupsFromUser(groupData, myUserName)
    if err == nil {
        t.Errorf("Error was not hit with malformedGroupTest file Expected:Error! passwd file may be corrupt!, Actual: nil ")
    }

}

func TestDecodeGroup(t* testing.T) {
    groupData, err := getFileData(groupTestPath)
    allEntries, err := decodeGroup(groupData)
    if err != nil {
        t.Errorf ("Unexpected error hit when decoding test data: " + err.Error())
    }
    if len(allEntries) != 4 {
        t.Errorf ("passwdTest file gave back more entries than expected. Expected: 4, Actual:" + strconv.Itoa(len(allEntries)))
    }

    groupData, err = getFileData(malformedGroupTestPath)
    allEntries, err = decodeGroup(groupData)
    if err == nil {
        t.Errorf("Error about corrupted group did not hit as expected")
    }
    if len(allEntries) != 0 {
        t.Errorf ("Number of entries mismatched. Expected: 0, Actual: " + strconv.Itoa(len(allEntries)))
    }
}

func TestDecodeGroupWithQuery(t* testing.T) {
    groupData, err := getFileData(groupTestPath)

    var queriedParams GroupInfo
    queriedParams.Name = "wheel"
    queriedParams.Gid = "0"
    queriedParams.Members = []string{"root"}
    queriedEntries, err := decodeGroupWithQuery(groupData, queriedParams)

    if err != nil {
        t.Errorf("Unexpected error hit: " + err.Error() )
    }
    if len(queriedEntries) != 1 {
        t.Errorf ("Found wrong number of matches. Expected: 1, Actual: " + strconv.Itoa(len(queriedEntries)))
    }

    if queriedEntries[0].Name != queriedParams.Name ||
        queriedEntries[0].Gid != queriedParams.Gid ||
        queriedEntries[0].Members[0] != queriedParams.Members[0] {
            t.Errorf ("Found mismatch in queried match for groups. Expected: Name = " + queriedParams.Name +
            "; Gid = " + queriedParams.Gid + "; Members = " + queriedParams.Members[0] + " Actual: Name = " + queriedEntries[0].Name +
            "; Gid = " + queriedEntries[0].Gid + "; Members = " + queriedEntries[0].Members[0])
        }

}

func TestCompareGroupInfo (t* testing.T) {
    var oneGroup GroupInfo
    var twoGroup GroupInfo
    var twoGroupCopy GroupInfo

    oneGroup.Name = "name1"
    oneGroup.Gid = "1"
    oneGroup.Members = []string{"member1"}


    twoGroup.Name = "name2"
    twoGroup.Gid = "2"
    twoGroup.Members = []string{"member2"}


    twoGroupCopy.Name = "name2"
    twoGroupCopy.Gid = "2"
    twoGroupCopy.Members = []string{"member2"}


    isMatch := compareGroupInfo(oneGroup, twoGroup)
    if (isMatch == true) {
        t.Errorf("Received inappropriate match with two different GroupInfo structs")
    }
    isMatch = compareGroupInfo(twoGroup, twoGroupCopy)
    if (isMatch == false) {
        t.Errorf("Received inappropriate mismatch with two of the same GroupInfo structs")
    }
}
func TestRetrieveGroupInfoFromGid (t* testing.T) {
    groupData, _ := getFileData(groupTestPath)
    gid := "3"
    myGroupInfo, _ := retrieveGroupInfoFromGid(groupData, gid)

    if myGroupInfo.Name != "chow"{
        t.Errorf("Incorrect name picked up from groupTestPath. Expected: chow, Actual: " + myGroupInfo.Name)
    }
    if myGroupInfo.Gid != "3" {
        t.Errorf("Incorrect GID picked up from groupTestPath. Expected: -2, Actual: " + myGroupInfo.Gid)
    }
    if myGroupInfo.Members[0] != "root"{
        t.Errorf("Incorrect members picked up from groupTestPath. Expected: {\"root\"}, Actual: " + myGroupInfo.Members[0])
    }
}
