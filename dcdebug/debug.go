package dcdebug

import (
	"errors"
	"fmt"
	"path/filepath"
	"runtime"
	"strings"
)

import "github.com/sirupsen/logrus"

func initSetting() {
	logrus.SetLevel(logrus.InfoLevel)
}
func LogI(args ...interface{}) {
	initSetting()
	logrus.Info(args...)

}
func LogW(args ...interface{}) {
	initSetting()
	logrus.Warn(args...)
}
func LogD(args ...interface{}) {
	initSetting()
	logrus.Debug(args...)
}

func LogE(args ...interface{}) {
	initSetting()
	logrus.Error(GetLine())
	logrus.Error(args...)

}
func LogNotImplement() {
	filename, line, funcname := "???", 0, "???"
	pc, filename, line, ok := runtime.Caller(1)
	// fmt.Println(reflect.TypeOf(pc), reflect.ValueOf(pc))
	if ok {
		funcname = runtime.FuncForPC(pc).Name()      // main.(*MyStruct).foo
		funcname = filepath.Ext(funcname)            // .foo
		funcname = strings.TrimPrefix(funcname, ".") // foo
		filename = filepath.Base(filename)           // /full/path/basename.go => basename.go
	}

	LogE("not implement ", filename, ":", line, ":", funcname, "()")
}
func addCallFunctionName(args ...interface{}) []interface{}{
	_, _, funcname := "???", 0, "???"
	pc, _, _, ok := runtime.Caller(2)
	// fmt.Println(reflect.TypeOf(pc), reflect.ValueOf(pc))
	if ok {
		funcname = runtime.FuncForPC(pc).Name()      // main.(*MyStruct).foo
		funcname = filepath.Ext(funcname)            // .foo
		funcname = strings.TrimPrefix(funcname, ".") // foo
	}
	newArgs:=[]interface{}{funcname+"():"}
	newArgs=append(newArgs,args...)
	return newArgs
}

func GetLine() string {
	filename, line, funcname := "???", 0, "???"
	pc, filename, line, ok := runtime.Caller(2)
	// fmt.Println(reflect.TypeOf(pc), reflect.ValueOf(pc))
	if ok {
		funcname = runtime.FuncForPC(pc).Name()      // main.(*MyStruct).foo
		funcname = filepath.Ext(funcname)            // .foo
		funcname = strings.TrimPrefix(funcname, ".") // foo
		filename = filepath.Base(filename)           // /full/path/basename.go => basename.go
	}
	return fmt.Sprint(filename, ":", line, ":", funcname, "()")
}
func getCallLine() string {
	pc, filename, line, ok := runtime.Caller(2)
	if ok{
		funcName := runtime.FuncForPC(pc).Name()
		return filename+":"+ itoa(line)+":"+funcName+"()"
	}
	return ""
}
func NewError(msg string) error {
	return errors.New(GetLine()+" MSG:"+msg)
}