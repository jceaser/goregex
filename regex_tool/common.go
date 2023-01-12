/*******************************************************************************
regulare expression REPL tool is an application to construct regular expression
statements in real time

This file holds general functions for the application. Most behavior code is in
this file and not main().

Created for https://earthdata.nasa.gov
Created by thomas.a.cherry@nasa.gov
Created: January 2023

License: Public Domain or what ever the default NASA license is.
*******************************************************************************/

package regex_tool

import (
    "fmt"
    "regexp"
    "strings"
)

/******************************************************************************/
// #MARK: Variables and Structs

const (
    REG_TOOL_VERSION = "0.0.1"
)

var (
    screen_setup bool
    Names = []string{"pattern", "input", "replace"} //used by liner and prompt
)

type PNames int
const (
    PPattern PNames = iota
    PInput
    PReplace
)

func (pname PNames) Text() string {
    return Names[pname]
}

/* App state for passing from one function to another */
type AppState struct {
    Pattern string
    InputText string
    ReplaceText string
    ResultText string
}

/* Populate values in the app state, overwriting existing values */
func (state *AppState) Init(pattern, input, replace string) {
    (*state).Pattern = pattern
    (*state).InputText = input
    (*state).ReplaceText = replace
    (*state).ResultText = ""
}

/******************************************************************************/
// #MARK: Application Functions

func Prompt() string {
    out := ""
    for _, name := range Names {
        out = out + string(name[0])
    }
    return out
}

func init() {
    screen_setup = false
}

/*
Called externally when the Screen is to be saved. ScrRestore() will revert the
console. It is safe to call this function multiple times.
*/
func ScrSave() {
    if !screen_setup {
        PrintCtrOnOut(ESC_SAVE_SCREEN)
        PrintCtrOnOut(ESC_SAVE_CURSOR)
        PrintCtrOnOut(ESC_CLEAR_SCREEN)
        screen_setup = true
    }
}

/**
Called externally when the console is to be restored from from SrcSave(). It is
safe to call this function multiple times.
*/
func ScrRestore() {
    if screen_setup {
        PrintCtrOnOut(ESC_RESTORE_CURSOR)
        PrintCtrOnOut(ESC_RESTORE_SCREEN)
        screen_setup = false
    }
}

/*
Called externally when the user interface is to be initialized. This will Save
the screen state with ScrSave()
*/
func InitUserInterface() {
    ScrSave()
    PrintStrOnOutAt("", GetHeight(), 0)
}

/*
Called externally when the user interface is to be restored after calling
InitUserInterface(). Will call ScrRestore()
*/
func DestroyUserInterface() {
    ScrRestore()
}

/*
Update one field and recalculate the result. Fields to be updated are prefixed
by a command followed by the data which is used to do the update with, for
example:

    in The quick brown fox

"in" will match the field "input" where as everything else is the data.

Fields to match are "patern, input, replace. Partial matches are accepted.
*/
func Process(state *AppState, text string) {
    if parts:=strings.Split(text, " ") ; len(parts)>1 && len(parts[0])>0 {
        data := strings.Join(parts[1:], " ")
        switch command := parts[0] ; {
            case strings.HasPrefix(PPattern.Text(), command):
                state.Pattern = data
            case strings.HasPrefix(PInput.Text(), command):
                state.InputText = data
            case strings.HasPrefix(PReplace.Text(), command):
                state.ReplaceText = data
        }
    }

    Calculate(state)
}

/*
This is everything, the entire application is just a wrapper around this one
function which takes a pattern and does a Regular Expression Replace All action
*/
func Calculate(state *AppState) {
    var re = regexp.MustCompile(state.Pattern)
    state.ResultText = re.ReplaceAllString(state.InputText, state.ReplaceText)
}

/*
Draw the interactive screen using as much or as little space as is provided. If
there are less then 2 lines of text, then a warning will be generated and only
the prompt will display.
*/
func Draw(state AppState) {
    height := GetHeight()
    PrintCtrOnOut(ESC_CLEAR_SCREEN)
    space := (height)/4
    if height>=9 {
        PrintStrOnOutAt(Green.W(Bold.W("Pattern")+":"), 1, 0)
        PrintStrOnOutAt(state.Pattern, 2, 0)

        PrintStrOnOutAt(Green.W(Bold.W("Input")+":"), space+1, 0)
        PrintStrOnOutAt(state.InputText, space+2, 0)

        PrintStrOnOutAt(Green.W(Bold.W("Replacement text")+":"), space*2, 0)
        PrintStrOnOutAt(state.ReplaceText, space*2+1, 0)

        PrintStrOnOutAt(Green.W(Bold.W("Result")+":"), space*3, 0)
        PrintStrOnOutAt(state.ResultText, space*3+1, 0)
    } else if height >= 5 {
        PrintStrOnOutAt(state.Pattern, 1, 0)
        PrintStrOnOutAt(state.InputText, space+1, 0)
        PrintStrOnOutAt(state.ReplaceText, space*2, 0)
        PrintStrOnOutAt(state.ResultText, space*3, 0)
    } else if height >= 2 {
        msg := fmt.Sprintf("%s%s ; %s%s ; %s%s %s %s",
            Style(Green, "P:"), state.Pattern,
            Style(Green, "I:"), state.InputText,
            Style(Green, "R:"), state.ReplaceText,
            Style(Green, "=="), state.ResultText)
        PrintStrOnOutAt(msg, 1, 0)
    } else {
        LogWarn.Print("Not enough room for output")
    }
    PrintStrOnOutAt("", height, 0)
}

/* Dump out the display text for non-interactive modes */
func Dump(state AppState) {
    format := "%7s: \"%s\"\n"
    fmt.Printf(format, "Pattern", state.Pattern)
    fmt.Printf(format, "Input", state.InputText)
    fmt.Printf(format, "Replace", state.ReplaceText)
    fmt.Printf(format, "Result", state.ResultText)
}

/* Print out a help screen, assume the full screen can be used */
func Help() {
    PrintCtrOnOut(ESC_CLEAR_SCREEN)
    PrintStrOnOutAt("", 1, 0)

    format := "%8s %5s %s\n"
    fmt.Printf(format, "command", "input", "Description")
    fmt.Printf(format, "-------", "-----", "-----")
    fmt.Printf(format, "pattern", "text", "Update Regular Expression Pattern.")
    fmt.Printf(format, "input", "text", "Update Input text.")
    fmt.Printf(format, "replace", "text", "Update Regular Expression Replacement Text")
    fmt.Printf(format, "-------", "-----", "-----")
    fmt.Printf(format, "exit", "", "Exit application")
    fmt.Printf(format, "help", "", "Print this help message")
    fmt.Printf("Partial commands accepted on all, p==pattern\n")

    PrintStrOnOutAt("", GetHeight(), 0)
}
