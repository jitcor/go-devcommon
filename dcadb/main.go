package dcadb



import (
	"fmt"
	"io/ioutil"
	"log"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
)

type DcAdb struct {
	mAdbPath    string
	mCurrentCmd string
	mResult     string
}

func NewDcAdb(adbPath string) *DcAdb {
	return &DcAdb{mAdbPath: adbPath}
}

func (that *DcAdb) Adb(params string) *DcAdb {
	that.mCurrentCmd = that.mAdbPath + " " + params
	that.mResult = ExecWrap2(that.mCurrentCmd)
	return that
}
func (that *DcAdb) GetModel() *DcAdb {
	that.mCurrentCmd = that.mAdbPath + " shell getprop ro.product.model"
	that.mResult = ExecWrap2(that.mCurrentCmd)
	return that
}
func (that *DcAdb) ClearCRLF() *DcAdb {
	reg := regexp.MustCompile(`[\r\n]`)
	that.mResult=reg.ReplaceAllString(that.mResult, "")
	//that.mResult = strings.ReplaceAll(that.mResult,"[\r\n]","")
	return that
}
func (that *DcAdb)GetSystemProperties(key,def string)*DcAdb {
	that.mCurrentCmd = that.mAdbPath + " shell getprop "+key
	that.mResult = ExecWrap2(that.mCurrentCmd)
	if that.mResult==""{
		that.mResult=def
	}
	that.ClearCRLF()
	return that
}
func (that *DcAdb) GetAndroidVersion() *DcAdb {
	that.mCurrentCmd = that.mAdbPath + " shell getprop ro.build.version.sdk"
	that.mResult = ExecWrap2(that.mCurrentCmd)
	return that
}
func (that *DcAdb) Pull(remoteFileOrFolder string, localFileOrFolder string) *DcAdb {
	that.mCurrentCmd = that.mAdbPath + " pull " + remoteFileOrFolder + " " + localFileOrFolder
	that.mResult = ExecWrap2(that.mCurrentCmd)
	return that
}
func (that *DcAdb) Delete(remoteFileOrFolder string) *DcAdb {
	that.mCurrentCmd = that.mAdbPath + " shell rm -rf " + remoteFileOrFolder
	that.mResult = ExecWrap2(that.mCurrentCmd)
	return that
}
func (that *DcAdb) Install(apkPath string) *DcAdb {
	that.mCurrentCmd = that.mAdbPath + " install -r -d " + apkPath
	that.mResult = ExecWrap2(that.mCurrentCmd)
	return that
}
func (that *DcAdb) Uninstall(packageName string) *DcAdb {
	that.mCurrentCmd = that.mAdbPath + " uninstall " + packageName
	that.mResult = ExecWrap2(that.mCurrentCmd)
	return that
}
func (that *DcAdb) GetVersionCode(packageName string) int {
	that.mCurrentCmd = that.mAdbPath +" shell dumpsys package " + packageName + " | grep versionCode"
	that.mResult = ExecWrap2(that.mCurrentCmd)
	that.ClearCRLF()
	if !that.IsResultEmpty(){
		re, _ := regexp.Compile("versionCode=(\\d*?) ")
		versionCode,_:=strconv.Atoi(re.FindStringSubmatch(that.mResult)[1])
		return versionCode
	}
	return -1
}
func (that *DcAdb) CheckFileExist(remoteFilePath string) bool {
	that.mCurrentCmd = that.mAdbPath + " shell ls "+remoteFilePath
	that.mResult = ExecWrap2(that.mCurrentCmd)
	that.ClearCRLF()
	return that.mResult!=""&&!strings.Contains(that.mResult,"No such file or directory")
}
func (that *DcAdb) CheckApkExist(packageName string) bool {
	that.mCurrentCmd = that.mAdbPath + " shell pm list package"
	that.mResult = ExecWrap2(that.mCurrentCmd)
	that.ClearCRLF()
	return that.mResult!=""&&strings.Contains(that.mResult,packageName)
}

func (that *DcAdb) CheckApkIsRunning(packageName string) bool {
	that.mCurrentCmd = that.mAdbPath + " shell ps | findstr "+packageName
	that.mResult = ExecWrap2(that.mCurrentCmd)
	that.ClearCRLF()
	return that.mResult!=""&&strings.Contains(that.mResult,packageName)
}
func (that *DcAdb) LaunchApp(packageName string) bool {
	if that.CheckApkIsRunning(packageName){
		return true
	}
	that.mCurrentCmd = that.mAdbPath + " shell dumpsys package "+packageName
	that.mResult = ExecWrap2(that.mCurrentCmd)
	that.ClearCRLF()
	if that.mResult!=""&&strings.Contains(that.mResult,"Activity Resolver Table:"){
		if re, err := regexp.Compile(`android\.intent\.action\.MAIN:.*? ([a-zA-Z0-9\\._]*?)/([a-zA-Z0-9\\._]*?) filter .*?Action: "android\.intent\.action\.MAIN"`);err!=nil{
			return false
		}else if match:=re.FindStringSubmatch(that.mResult);match==nil{
			return false
		}else{
			pkgName:=match[1]
			activity:=match[2]
			that.mCurrentCmd = that.mAdbPath + " shell am start -n "+pkgName+"/"+activity
			that.mResult = ExecWrap2(that.mCurrentCmd)
			that.ClearCRLF()
			return that.mResult!=""&&strings.Contains(that.mResult,"cmp="+pkgName+"/"+activity)
		}

	}
	return false
}

func (that *DcAdb) CheckDevicesExist() bool {
	that.mCurrentCmd = that.mAdbPath + " devices"
	that.mResult = ExecWrap2(that.mCurrentCmd)
	that.ClearCRLF()
	offset:=strings.LastIndex(that.mResult,"List of devices attached")+len("List of devices attached")-1
	subStr:=that.mResult[offset:len(that.mResult)]
	//subStr:= main.Substr(that.mResult,offset, len(that.mResult)-offset)
	//List of devices attached
	return that.mResult!=""&&strings.Contains(subStr,"device")
}
func (that *DcAdb) GetExternalStorage() *DcAdb {
	that.mCurrentCmd = that.mAdbPath + " shell echo $EXTERNAL_STORAGE"
	that.mResult = ExecWrap2(that.mCurrentCmd)
	return that
}
func (that *DcAdb) ForwardTcp(localPort int, remotePort int) *DcAdb {
	that.mCurrentCmd = that.mAdbPath + " forward tcp:" + strconv.Itoa(localPort) + " tcp:" + strconv.Itoa(remotePort)
	that.mResult = ExecWrap2(that.mCurrentCmd)
	return that
}
func (that *DcAdb) GetResult() string {
	return that.mResult
}
func (that *DcAdb) IsSuccess(successKeyword string) bool {
	return that.mResult != "" && strings.Contains(that.mResult, successKeyword)
}
func (that *DcAdb) IsFailed(failedKeyword string) bool {
	return that.mResult != "" && strings.Contains(that.mResult, failedKeyword)
}

func (that *DcAdb) IsResultEmpty() bool {
	result:=strings.TrimSpace(that.mResult)
	return result== ""
}
func (that *DcAdb) Printf(tag string) *DcAdb {
	fmt.Println("DcAdb", tag+" "+that.mResult)
	return that
}
func ExecWrap2(cmdAndParams string) string {
	if cmdAndParams == "" {
		return "Please input cmd And Params to call execWrap2"
	}
	log.Println(cmdAndParams)
	localParams := strings.Split(cmdAndParams, " ")
	cmd := ""
	cmd = localParams[0]
	if len(localParams) == 1 {
		return ExecWrap(cmd, "")
	}
	params := ""
	for i := 1; i < len(localParams); i++ {
		params += " " + localParams[i]
	}
	ret:= ExecWrap(cmd, params[1:])
	log.Println("exec ret:",ret)
	return ret
}
func ExecWrap(cmd string, params string) string {
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
