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
    //passwdTestPath := "./
    passwdData, err := getFileData(passwdTestPath)
    allEntries, err := decodePasswd(passwdData)
    if err != nil {
        t.Errorf ("Unexpected error hit when decoding test data: " + err.Error())
    }
    if len(allEntries) != 2 {
        t.Errorf ("passwdTest file gave back more entries than expected. Expected: 2, Actual:" + strconv.Itoa(len(allEntries)))
    }
}
