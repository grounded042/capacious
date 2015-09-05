package utils

import (
	"log"
	"runtime"
	"strconv"
)

// Capacious Error Interface
// Anything that satisfies this interface also satisfies the error
// interface and can thus be used in those situations.
// We add the Code() func since we want to be able to pass error codes.
// We add the Location() func since we want to get the location at which the
// error happened.
// These codes can be anything from http status codes to some internal meaning.
type Error interface {
	Code() int
	Error() string
	Location() string
}

// Struct that implements the Error interface
type ApiError struct {
	c int
	e string
	l string
}

func (err ApiError) Code() int {
	return err.c
}

func (err ApiError) Error() string {
	return err.e
}

func (err ApiError) Location() string {
	return err.l
}

// Build a new ApiError
func NewApiError(code int, err string) ApiError {
	// get the location info
	pc, _, _, _ := runtime.Caller(1)
	loc := runtime.FuncForPC(pc).Name()

	return ApiError{c: code, e: err, l: loc}
}

// Log errors. When you pass an obj that implements Error, we log
// the code as well.
func LogError(logme error) {
	if e, ok := logme.(Error); ok {
		PrintRed("Error: " + strconv.Itoa(e.Code()) + " - " + e.Error())
		PrintRed("Error Location: " + e.Location())
	} else if _, ok = logme.(error); ok {
		PrintRed("Error: " + logme.Error())
	}
}

// This function is for when there might be an error, but
// we don't really care about. It will check and see if
// there is an error and, if there is, log that error.
func CheckErr(err error, msg string) {
	if err != nil {
		PrintRed(msg + ":")
		PrintRed(err.Error())
	}
}

func Print(msg string) {
	log.Println(msg)
}

func PrintRed(msg string) {
	log.Println("\033[31;1m" + msg + "\033[0m")
}
