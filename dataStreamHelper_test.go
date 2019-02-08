package main

import (
    "testing"
    //"fmt"
    "strconv"
)

var passwdTestPath string = "./testFiles/passwdTest"

func TestGetFileData(t* testing.T){

    fileNotExistPath := "./testFiles/fileNotExist"
    //passwdTestPath := "./testFiles/passwdTest"

    passwdData, err := getFileData(fileNotExistPath)
    if err == nil {
        t.Errorf("Processed file "+ fileNotExistPath  +" which should not exist")
    }
    if (len(passwdData) != 0 ) {
        t.Errorf("Total file line size mismatch. Expected: 0, Actual:"+ strconv.Itoa(len(passwdData)))
    }

    passwdData ,err = getFileData(passwdTestPath)
    if err != nil {
        t.Errorf ("Encountered error reading file " + passwdTestPath + "which should exist")
    }
    if (len(passwdData) != 4) {
        t.Errorf ("Total file line size mismatch. Expected: 4, Actual:" + strconv.Itoa(len(passwdData)))
    }

}

func TestDecodePasswd(t* testing.T){
    passwdData, err := getFileData(passwdTestPath)
    allEntries, err := decodePasswd(passwdData)
    if err != nil {
        t.Errorf ("Unexpected error hit when decoding test data: " + err.Error())
    }
    if len(allEntries) != 2 {
        t.Errorf ("passwdTest file gave back more entries than expected. Expected: 2, Actual:" + strconv.Itoa(len(allEntries)))
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

    csvData, err := getFileData(passwdTestPath)

    queriedEntries, err := decodePasswdWithQuery(csvData, queriedParams)
    if err != nil || queriedEntries == nil{
        t.Errorf(err.Error())
    }

    if (len(queriedEntries) != 1) {
        t.Errorf("number of matches incorrect. Expected: 1, Actual:" + strconv.Itoa(len(queriedEntries)))
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
