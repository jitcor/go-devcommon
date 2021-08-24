package dcmd

import (
	"fmt"
	"io/ioutil"
	"log"
	"os/exec"
	"smsonline/devcommon/dctask"
	"strings"
)

type DCmd struct {
	mAdbPath    string
	mCurrentCmd string
	mResult     string
}

func NewDCmd(adbPath string) *DCmd {
	return &DCmd{mAdbPath: adbPath}
}
func (that *DCmd) GetResult() string {
	return that.mResult
}
func (that *DCmd) IsSuccess(successKeyword string) bool {
	return that.mResult != "" && strings.Contains(that.mResult, successKeyword)
}
func (that *DCmd) IsResultEmpty() bool {
	result:=strings.ReplaceAll(that.mResult,"[\r\n]","")
	return result== ""
}
func (that *DCmd) Printf(tag string) *DCmd {
	fmt.Println("DCmd", tag+" "+that.mResult)
	return that
}
func Exec_(cmdAndParams string,err *error) string {
	if cmdAndParams == "" {
		return "Please input cmd And Params to call execWrap2"
	}
	localParams := strings.Split(cmdAndParams, " ")
	cmd := ""
	cmd = localParams[0]
	if len(localParams) == 1 {
		return Exec(cmd, "",err)
	}
	params := ""
	for i := 1; i < len(localParams); i++ {
		params += " " + localParams[i]
	}
	return Exec(cmd, params[1:],err)
}
func Exec(cmd string, params string,err *error) string {
	var result = ""
	localParams := strings.Split(params, " ")
	cmdPtr := exec.Command(cmd, localParams...)
	log.Println("Exec:",cmdPtr)
	stdout, e := cmdPtr.StdoutPipe()
	if e != nil {
		*err=e
		return ""
	}
	defer stdout.Close()
	stderr, e := cmdPtr.StderrPipe()
	if e != nil {
		*err=e
		return ""
	}
	defer stderr.Close()

	if err := cmdPtr.Start(); err != nil {
		return err.Error()
	}
	dctask.NewDcTask().Go(func(err *error) {
		if opBytes, e := ioutil.ReadAll(stderr); e != nil {
			*err=e
			return
		} else {
			result += "\n" + string(opBytes)
		}
	},err).Go(func(err *error) {
		if opBytes, e := ioutil.ReadAll(stdout); e != nil {
			*err=e
			return
		} else {
			result = string(opBytes)
		}
	},err).Wait()
	if *err!=nil{
		return ""
	}
	cmdPtr.Wait()
	//OnInfo(result)
	return result
}