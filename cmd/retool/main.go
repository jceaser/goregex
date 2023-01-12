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

package main

import (
    "flag"
    "fmt"
    "github.com/peterh/liner"
    "os"
    "path/filepath"
    reg "earthdata.nasa.gov/regex_tool/regex_tool"
    "strings"
)

/******************************************************************************/
// #MARK: Variables and Structs

var (
    history_fn = filepath.Join(home(), ".reg_history") //used by liner
    names = []string{"exit", "quit", "dump", "help"} //used by liner and prompt
)

type PNames int
const (
    PExit PNames = iota
    PQuit
    PDump
    PHelp
)

func (pname PNames) Text() string {
    return names[pname]
}

/*
Codes for any follow up action that must take place after a user requests a
command.
*/
type FollowUp int
const (
    FollowNone FollowUp = iota
    FollowExit
    FollowDump
)

/******************************************************************************/
// #MARK: Application Functions

func Prompt() string {
    out := ""
    //first letter of each item
    for i := range names {
        out = out + string(names[i][0])
    }
    return out
}

/* Find the users home folder, if that can not be found, then a temp directory. */
func home() string {
    if home, err := os.UserHomeDir() ; err==nil {
        return home
    } else {
        reg.LogWarning.Println ("Request for user home directory failed, using temp directory")
        return os.TempDir()
    }
}

/* Process a command from the user */
func process_line(line string, state *reg.AppState) FollowUp {

    if len(line)>0 {
        switch text := strings.ToLower(line); {
        case strings.HasPrefix(PExit.Text(), text) ||
                strings.HasPrefix(PQuit.Text(), text):
            return FollowExit
        case strings.HasPrefix(PDump.Text(), text):
            return FollowDump
        case strings.HasPrefix(PHelp.Text(), text):
            reg.Help()
            return FollowNone
        }
    }
    reg.Process(state, line)
    reg.Draw(*state)

    return FollowNone
}

/* Initialize liner */
func setup_liner(line *liner.State) {
    line.SetCtrlCAborts(true)

    line.SetTabCompletionStyle(liner.TabPrints)
    line.SetCompleter(func(line string) (c []string) {
        for _, name := range append(names, reg.Names...) {
            if strings.HasPrefix(name, strings.ToLower(line)) {
                c = append(c, name)
            }
        }
        return
    })
    if file_pointer, err := os.Open(history_fn); err!=nil {
        reg.LogError.Print("Could not open history file:", err)
    } else {
        line.ReadHistory(file_pointer)
        file_pointer.Close()
    }
}

/*
Enter interactive mode where all app state info is displayed and a prompt is
shown asking for commands
*/
func Interactivity(state reg.AppState) {
    reg.InitUserInterface()
    defer reg.DestroyUserInterface()

    //draw initial interface
    reg.Process(&state, "")
    reg.Draw(state)

    line := liner.NewLiner()
    defer line.Close()
    setup_liner(line)
    running := true
    dump := false
    for running {
        prompt := fmt.Sprintf("%s%s>", reg.Prompt(), Prompt())
        if text, err := line.Prompt(prompt) ; err == nil {
            line.AppendHistory(text)
            follow_up := process_line(text, &state)
            switch follow_up {
            case FollowExit:
                running = false
            case FollowDump:
                reg.LogInfo.Print("dump")
                running = false
                dump = true
            }
        } else if (err == liner.ErrPromptAborted) {
            reg.LogWarn.Println ("Prompt was exited")
        } else {
            reg.LogError.Println ("Reading line:", err)
        }
        if file_pointer, err := os.Create(history_fn); err != nil {
            reg.LogError.Print("Could not Create history file: ", err)
        } else {
            line.WriteHistory(file_pointer)
            file_pointer.Close()
        }
    }
    reg.DestroyUserInterface()
    if dump {
        Work(state)
    }
}

/* Non-interactive mode, just dump the value out */
func Work(state reg.AppState) {
    reg.Process(&state, "")
    reg.Dump(state)
}

func main() {
    pattern := flag.String("pattern", " [a-z]",
        "Regular Expression Pattern")
    input := flag.String("input", "The quick brown fox jumps over the lazy dogs",
        "Input text to apply pattern on")
    replace := flag.String("replace", "*",
        "Regular Expression replacement test, groups start at $1")
    version := flag.Bool("version", false,
        "Print out the application version and exit")
    interactive := flag.Bool("interactive", false,
        "Enter interactive mode. At prompt: p <text>, i <text>, r <text>, exit, dump")
    flag.Parse()

    state := reg.AppState{}
    state.Init(*pattern, *input, *replace)

    if *version {
        fmt.Println ("Written by thomas.a.cherry@nasa.gov")
        fmt.Println (reg.REG_TOOL_VERSION)
    }

    if *interactive {
        Interactivity(state)
    } else {
        Work(state)
    }
}
