/*******************************************************************************
regulare expression REPL tool is an application to construct regular expression
statements in real time

This file holds functions related to printing to a modern unix terminal. Massive
assumptions are made about how modern the terminal is and what can be done.

Created for https://earthdata.nasa.gov
Created by thomas.a.cherry@nasa.gov
Created: January 2023

License: Public Domain or what ever the default NASA license is.
*******************************************************************************/

package regex_tool

import (
    "fmt"
    "os"
    "syscall"
    "unsafe"
)

/******************************************************************************/
// #MARK: Variables and Structs

/* returned by os for terminal width and height */
type winsize struct {
    Row    uint16
    Col    uint16
    Xpixel uint16
    Ypixel uint16
}

/* Console Terminal Codes For style */
type ConsoleCodes int
const (
    ClearAll ConsoleCodes = iota //0
    Bold
    Light
    Italic      //3
    Underline
    Blink
    _           //6
    Reverse
    Invisible   //not working
    Strike      //9
    DoubleUnder ConsoleCodes = iota + (21-10)
    Black ConsoleCodes = iota + (30-11)
    Red
    Green
    Yellow      //33
    Blue
    Magenta
    Cyan        //36
    Gray
    White ConsoleCodes = iota + (97-19)
    DefaultCode
)

/*
Apparently there is not a one-one matching of control start codes and end codes
so this map contains the matching values
*/
var (
    ends = map[ConsoleCodes]int{ClearAll:0,
        Bold:22,
        Light:22,
        Italic:23,
        Underline:24,
        DoubleUnder:24,
        Blink:25,
        Reverse:27,
        Invisible:28,
        Strike:29,
        /*Black:39,Red:39,Green:39,Yellow:39,Blue:39,Magenta:39,Cyan:39,Gray:39,White:39,*/
        DefaultCode: 39}
)

/*
Find the correct end code to match the, assume that the ending code is for a
color and return the "DefaultCode"
*/
func (cc ConsoleCodes) End() int {
    if value, okay := ends[cc] ; okay {
        return value
    } else {
        return ends[DefaultCode]
    }
}

func (cc ConsoleCodes) Wrap(texts ...any) string {
    text := fmt.Sprint(texts...)
    return fmt.Sprintf("\033[%dm%s\033[%dm", cc, text, cc.End())
}

func (cc ConsoleCodes) W(texts ...any) string {
    return cc.Wrap(texts...)
}


func init() {
    //fmt.Printf(Red.W("Test ", Bold.W(Blink.W("B"), "old"), " Test") + ".\n")
}

/* Other console control codes not related to style */
const (
    ESC_SAVE_SCREEN = "?47h"
    ESC_RESTORE_SCREEN = "?47l"

    ESC_SAVE_CURSOR = "s"
    ESC_RESTORE_CURSOR = "u"

    ESC_CURSOR_ON = "?25h"
    ESC_CURSOR_OFF = "?25l"

    ESC_CLEAR_SCREEN = "2J"
    ESC_CLEAR_LINE = "2K"
)

/******************************************************************************/
// #MARK: - Console Functions

func GetWidth() int {
    ws := &winsize{}
    retCode, _, errno := syscall.Syscall(syscall.SYS_IOCTL,
        uintptr(syscall.Stdin),
        uintptr(syscall.TIOCGWINSZ),
        uintptr(unsafe.Pointer(ws)))

    if int(retCode) == -1 {
        panic(errno)
    }
    return int(ws.Col)
}

func GetHeight() int {
    ws := &winsize{}
    retCode, _, errno := syscall.Syscall(syscall.SYS_IOCTL,
        uintptr(syscall.Stdin),
        uintptr(syscall.TIOCGWINSZ),
        uintptr(unsafe.Pointer(ws)))

    if int(retCode) == -1 {
        panic(errno)
    }
    return int(ws.Row)
}

/* A add console style to text */
func Style(color ConsoleCodes, text string) string {
    return fmt.Sprintf("\033[%dm%s\033[%dm", color, text, ends[color])
}

/*
Short hand for adding console style to text
Example usage:
    fmt.Printf(S(Red, "Test "+S(Bold, S(Blink, "B")+"old") + " Test") + ".\n")
*/
func S(color ConsoleCodes, text string) string {
    return Style(color, text)
}

func PrintCtrOnOut(esc string) {
    fmt.Fprintf(os.Stdout, "\033[%s", esc)
}

func PrintCtrOnErr(esc string) {
    fmt.Fprintf(os.Stderr, "\033[%s", esc)
}

func PrintStrOnOutAt(msg string, y, x int) {
    fmt.Fprintf(os.Stdout, "\033[%d;%dH%s", y, x, msg)
}

func PrintStrOnErrAt(msg string, y, x int) {
    fmt.Fprintf(os.Stderr, "\033[%d;%dH%s", y, x, msg)
}
