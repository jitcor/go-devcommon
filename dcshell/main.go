package dcshell

import (
	"io/ioutil"
	"log"
	"os/exec"
	"regexp"
	"strings"
)

type DcShell struct {
	root bool
	result string
	currentCmd string
}

func NewDcShell(root bool) *DcShell {
	return &DcShell{root: root}
}
func (that*DcShell)CheckAppIsRunning(packageName string) bool {
	that.currentCmd = "ps -A | grep "+packageName
	that.result = that.execWrap2(that.currentCmd)
	that.ClearCRLF()
	return that.result!=""&&strings.Contains(that.result,packageName)
}

func (that *DcShell) LaunchApp(packageName string) bool {
	if that.CheckAppIsRunning(packageName){
		return true
	}
	that.currentCmd = "dumpsys package "+packageName
	that.result = that.execWrap2(that.currentCmd)
	that.ClearCRLF()
	if that.result!=""&&strings.Contains(that.result,"Activity Resolver Table:"){
		if re, err := regexp.Compile(`android\.intent\.action\.MAIN:.*? ([a-zA-Z0-9\\._]*?)/([a-zA-Z0-9\\._]*?) filter .*?Action: "android\.intent\.action\.MAIN"`);err!=nil{
			return false
		}else if match:=re.FindStringSubmatch(that.result);match==nil{
			return false
		}else{
			pkgName:=match[1]
			activity:=match[2]
			that.currentCmd = "am start -n "+pkgName+"/"+activity
			that.result = that.execWrap2(that.currentCmd)
			that.ClearCRLF()
			return that.result!=""&&strings.Contains(that.result,"cmp="+pkgName+"/"+activity)
		}

	}
	return false
}
func (that *DcShell) ClearCRLF() *DcShell {
	reg := regexp.MustCompile(`[\r\n]`)
	that.result=reg.ReplaceAllString(that.result, "")
	//that.mResult = strings.ReplaceAll(that.mResult,"[\r\n]","")
	return that
}
func (that *DcShell) execWrap2(cmdAndParams string) string {
	if cmdAndParams == "" {
		return "Please input cmd And Params to call execWrap2"
	}
	log.Println(cmdAndParams)
	localParams := strings.Split(cmdAndParams, " ")
	cmd := ""
	cmd = localParams[0]
	if len(localParams) == 1 {
		return that.execWrap(cmd, "")
	}
	params := ""
	for i := 1; i < len(localParams); i++ {
		params += " " + localParams[i]
	}
	ret:= that.execWrap(cmd, params[1:])
	log.Println("exec ret:",ret)
	return ret
}
func (that *DcShell) execWrap(cmd string, params string) string {
	var result = ""
	localParams := strings.Split(params, " ")
	cmdPtr := exec.Command(cmd, localParams...)
	log.Println(cmdPtr.String())
	stdout, err := cmdPtr.StdoutPipe()
	if err != nil {
		return err.Error()
	}
	defer stdout.Close()
	stderr, err := cmdPtr.StderrPipe()
	if err != nil {
		return err.Error()
	}
	defer stderr.Close()

	if err := cmdPtr.Start(); err != nil {
		return err.Error()
	}
	if opBytes, err := ioutil.ReadAll(stdout); err != nil {
		return err.Error()
	} else {
		result = string(opBytes)
	}

	if opBytes, err := ioutil.ReadAll(stderr); err != nil {
		return err.Error()
	} else {
		result += "\n" + string(opBytes)
	}
	cmdPtr.Wait()
	log.Println(result)
	//OnInfo(result)
	return result
}
