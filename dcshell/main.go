package dcshell

import (
	"io/ioutil"
	"log"
	"os/exec"
	"regexp"
	"strings"
	"time"
)

type DcShell struct {
	root bool
	debug bool
	result string
	currentCmd string
}

func NewDcShell(root,debug bool) *DcShell {
	return &DcShell{root: root,debug: debug}
}
func (that*DcShell)CheckAppIsRunning(packageName string) bool {
	that.currentCmd = "ps -A"
	that.result = that.execWrap2(that.currentCmd)
	if re,err:=regexp.Compile(strings.ReplaceAll(packageName,`.`,`\.`)+`(?:\n|$|\r)`);err!=nil{
		log.Println("[-] reg compile error:",err)
		return false
	}else {
		return that.result!=""&&strings.Contains(that.result,packageName)&&re.MatchString(that.result)
	}
}
func (that*DcShell)CheckPortIsListen(port string) bool {
	that.currentCmd = "netstat -anp"
	that.result = that.execWrap2(that.currentCmd)
	if re,err:=regexp.Compile(`:::`+port+`[^\r\n]*?LISTEN`);err!=nil{
		log.Println("[-] reg compile error:",err)
		return false
	}else {
		return that.result!=""&&strings.Contains(that.result,":::"+port)&&re.MatchString(that.result)
	}
}

func (that *DcShell) LaunchApp(packageName string) bool {
	that.currentCmd = "dumpsys package "+packageName
	that.result = that.execWrap2(that.currentCmd)
	that.ClearCRLF()
	//Category: "android.intent.category.LAUNCHER"
	if that.result!=""&&strings.Contains(that.result,"Activity Resolver Table:"){
		if re, err := regexp.Compile(`[a-z0-9]{3,10} ([a-zA-Z0-9\\._]*?)/([a-zA-Z0-9\\._]*?) filter [a-z0-9]{3,10}[^/]*?Category: "android\.intent\.category\.LAUNCHER"`);err!=nil{
			log.Println("[-] reg compile error:",err)
			return false
		}else if match:=re.FindStringSubmatch(that.result);match==nil{
			log.Println("[-] not found activity on package:",packageName)
			return false
		}else{
			pkgName:=match[1]
			activity:=match[2]
			log.Println("[+] find package:",pkgName,",activity:",activity)
			that.currentCmd = "am start -n "+pkgName+"/"+activity
			that.result = that.execWrap2(that.currentCmd)
			that.ClearCRLF()
			time.Sleep(2*time.Second)
			return that.result!=""&&strings.Contains(that.result,"cmp="+pkgName+"/"+activity)
		}

	}
	return false
}
func (that *DcShell) LaunchAppWhenStopped(packageName string) bool {
	if that.CheckAppIsRunning(packageName){
		return true
	}
	that.currentCmd = "dumpsys package "+packageName
	that.result = that.execWrap2(that.currentCmd)
	that.ClearCRLF()
	//Category: "android.intent.category.LAUNCHER"
	if that.result!=""&&strings.Contains(that.result,"Activity Resolver Table:"){
		if re, err := regexp.Compile(`[a-z0-9]{3,10} ([a-zA-Z0-9\\._]*?)/([a-zA-Z0-9\\._]*?) filter [a-z0-9]{3,10}[^/]*?Category: "android\.intent\.category\.LAUNCHER"`);err!=nil{
			log.Println("[-] reg compile error:",err)
			return false
		}else if match:=re.FindStringSubmatch(that.result);match==nil{
			log.Println("[-] not found activity on package:",packageName)
			return false
		}else{
			pkgName:=match[1]
			activity:=match[2]
			log.Println("[+] find package:",pkgName,",activity:",activity)
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
	if that.debug{
		log.Println(cmdAndParams)
	}
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
	if that.debug{
		log.Println("exec ret:",ret)
	}
	return ret
}
func (that *DcShell) execWrap(cmd string, params string) string {
	var result = ""
	localParams := strings.Split(params, " ")
	cmdPtr := exec.Command(cmd, localParams...)
	if that.debug{
		log.Println(cmdPtr.String())
	}
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
	if that.debug{
		log.Println(result)
	}
	//OnInfo(result)
	return result
}
