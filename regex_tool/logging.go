/*******************************************************************************
regulare expression REPL tool is an application to construct regular expression
statements in real time.

This file holds functions related starting the application and managing command
line flags, and direct interfacing with the user.

Created for https://earthdata.nasa.gov
Created by thomas.a.cherry@nasa.gov
Created: January 2023

License: Public Domain or what ever the default NASA license is.
*******************************************************************************/
/*
*/

package regex_tool

import (
    "io/ioutil"
    "log"
    "os"
)

/******************************************************************************/
// #MARK: Variables and Structs

var (
    LogError *log.Logger
    LogWarning *log.Logger
    LogWarn *log.Logger
    LogInfo *log.Logger
    LogDebug *log.Logger
)

func init() {
    file := os.Stderr

    LogError = log.New(file, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
    LogWarning = log.New(file, "WARNING: ", log.Ldate|log.Ltime|log.Lshortfile)
    LogWarn = LogWarning    //an alias
    LogInfo = log.New(file, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
    LogDebug = log.New(file, "DEBUG: ", log.Ldate|log.Ltime|log.Lshortfile)

    LogInfo.SetOutput(ioutil.Discard)
    LogDebug.SetOutput(ioutil.Discard)
}

func EnableInfo() {
    LogInfo.SetOutput(os.Stderr)
}

func EnableDebug() {
    LogDebug.SetOutput(os.Stderr)
}
