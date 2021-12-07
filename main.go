package devcommon

import (
	"archive/zip"
	"bufio"
	"crypto/md5"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"io"
	"io/ioutil"
	"log"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
)

var NilErr = new(error)

func Min(a, b int) int {
	if a < b {
		return a
	} else {
		return b
	}
}

func GetUUID() ([]byte, error) {
	u := uuid.New()
	if bUUID, err := u.MarshalBinary(); err != nil {
		return nil, err
	} else {
		return bUUID, nil
	}
}
func CreateFileAndDelIfExists(path string,retErr *error) *os.File {
	if *retErr!=nil{
		return nil
	}
	if Exists(path) || IsDir(path) {
		if err := os.RemoveAll(path); err != nil {
			*retErr=err
			return nil
		}
	}
	if err := os.MkdirAll(filepath.Dir(path), os.ModePerm); err != nil {
		*retErr=err
		return nil
	}
	if file, err := os.OpenFile(path, os.O_CREATE, 0777); err != nil {
		*retErr=err
		return nil
	} else {
		return file
	}
}
func CreateFileIfNotExists(path string,retErr *error) *os.File {
	if *retErr!=nil{
		return nil
	}
	if Exists(path) && IsFile(path) {
		if file, err := os.OpenFile(path, os.O_APPEND, 0777); err != nil {
			*retErr=err
			return nil
		} else {
			return file
		}
	}
	if Exists(path) || IsDir(path) {
		if err := os.RemoveAll(path); err != nil {
			*retErr=err
			return nil
		}
	}
	if err := os.MkdirAll(filepath.Dir(path), os.ModePerm); err != nil {
		*retErr=err
		return nil
	}
	if file, err := os.OpenFile(path, os.O_CREATE, 0777); err != nil {
		*retErr=err
		return nil
	} else {
		return file
	}
}

// 判断所给路径文件/文件夹是否存在
func Exists(path string) bool {
	_, err := os.Stat(path) //os.Stat获取文件信息
	if err != nil {
		if os.IsExist(err) {
			return true
		}
		return false
	}
	return true
}

// 判断所给路径是否为文件夹
func IsDir(path string) bool {
	s, err := os.Stat(path)
	if err != nil {
		return false
	}
	return s.IsDir()
}

// 判断所给路径是否为文件
func IsFile(path string) bool {
	return !IsDir(path)
}
func TcpToFile(conn net.Conn, file *os.File, length int) (err error) {
	writer := bufio.NewWriter(file)
	if wn, err := io.CopyN(writer, conn, int64(length)); err != nil {
		return err
	} else {
		log.Println("wn:", wn)
	}
	//if n,err:=writer.Write(dcsTcp.ReadBytes(uint(length)));err!=nil{
	//	panic(err)
	//}else if n!=length{
	//	panic("n!=length")
	//}else if err:=writer.Flush();err!=nil{
	//	panic(err)
	//}
	return nil
}
func GetBytes(conn net.Conn, fileName string, fileLength int) {
	log.Println("-->getBytes(fileName:" + fileName + ",fileLength:" + strconv.Itoa(fileLength) + ")")
	_ = os.RemoveAll(fileName)
	if !MakeSureFileExists(fileName) {
		log.Println("无法确保文件存在:" + fileName)
		return
	}
	f, _ := os.OpenFile(fileName, os.O_CREATE, 0777)
	defer f.Close()
	w := bufio.NewWriter(f)
	var buf = make([]byte, 1024)
	var localFileLength = 0
	for {
		var readLengthTmp int
		if fileLength-localFileLength > len(buf) {
			readLengthTmp = len(buf)
		} else {
			readLengthTmp = fileLength - localFileLength
		}
		buf = make([]byte, readLengthTmp)
		if !makeSureReadAll(conn, buf) {
			log.Println("无法确保读到完整数据")
			return
		}
		w.Write(buf)
		w.Flush()
		localFileLength += readLengthTmp
		if localFileLength == fileLength {
			break
		}
	}
	log.Println("file receive success")

}
func makeSureReadAll(conn net.Conn, buf []byte) bool {
	length := len(buf)
	if length <= 0 || cap(buf) <= 0 {
		return false
	}
	n := 0
	offset := 0
	for length > 0 && offset < cap(buf) {
		n, _ := conn.Read(buf[offset : offset+length])
		length -= n
		offset += n
	}
	return n == 0 && length == 0 && offset == len(buf)
}
func MakeSureFileExists(path string) bool {
	if Exists(path) {
		return true
	}
	if IsDir(path) {
		err := os.MkdirAll(path, os.ModePerm)
		return err == nil
	}
	if !Exists(GetParentDirectory(path)) {
		os.MkdirAll(GetParentDirectory(path), os.ModePerm)
	}
	_, err := os.Create(path)
	return err == nil
}
func GetParentDirectory(filePath string) string {
	return filepath.Dir(filePath)
}
func Zip(src, target string, err *error) {
	if *err != nil {
		return
	}
	src = FilePath(src, err)
	target = FilePath(target, err)
	// 预防：旧文件无法覆盖
	_ = os.RemoveAll(target)

	// 创建：zip文件
	zipfile, _ := os.Create(target)
	defer zipfile.Close()

	// 打开：zip文件
	archive := zip.NewWriter(zipfile)
	defer archive.Close()

	// 遍历路径信息
	filepath.Walk(src, func(path string, info os.FileInfo, _ error) error {
		path = FilePath(path, err)
		// 如果是源路径，提前进行下一个遍历
		if path == src {
			return nil
		}

		// 获取：文件头信息
		header, _ := zip.FileInfoHeader(info)
		header.Name = strings.TrimPrefix(path, src+"/")

		// 判断：文件是不是文件夹
		if info.IsDir() {
			header.Name += `/`
		} else {
			// 设置：zip的文件压缩算法
			header.Method = zip.Deflate
		}

		// 创建：压缩包头部信息
		writer, _ := archive.CreateHeader(header)
		if !info.IsDir() {
			file, _ := os.Open(path)
			defer file.Close()
			io.Copy(writer, file)
		}
		return nil
	})
}
func Unzip(src, out string) error {
	or, err := zip.OpenReader(src)
	if err != nil {
		return err
	}
	for _, file := range or.Reader.File {
		if file.FileInfo().IsDir() {
			if err := os.MkdirAll(OptSeparator(filepath.Join(out, file.Name)), 0644); err != nil {
				return err
			}
			continue
		}
		f, err := file.Open()
		if err != nil {
			return err
		}
		defer f.Close()
		cfile, err := os.Create(OptSeparator(filepath.Join(out, file.Name)))
		if err != nil {
			return err
		}
		if _, err := io.Copy(cfile, f); err != nil {
			return err
		}
		cfile.Close()
	}
	or.Close()
	return nil
}

func OptSeparator(path string) string {
	return strings.ReplaceAll(path, "\\", "/")
}
func MD5Bytes(data []byte) string {
	m := md5.New()
	m.Write(data)
	return fmt.Sprintf("%x", m.Sum(nil))
}
func MD5String(text string) string {
	m := md5.New()
	m.Write([]byte(text))
	return fmt.Sprintf("%x", m.Sum(nil))
}
func FilePath(path string, err *error) string {
	if *err != nil {
		return path
	}
	if absFile, e := filepath.Abs(path); e != nil {
		log.Println("Error:", e)
		*err = e
		return OptSeparator(path)
	} else {
		return OptSeparator(absFile)
	}

}
func FileNamePrefix(path string) string {
	return filepath.Base(path)[0 : len(filepath.Base(path))-len(FileNameSuffix(path))]
}
func FileNameSuffix(path string) string {
	return filepath.Ext(path)
}
func DirCreate(path string, err *error) {
	if *err != nil {
		return
	}
	path = OptSeparator(path)
	if Exists(path) {
		return
	}
	if !Exists(filepath.Dir(path)) {
		if e := os.MkdirAll(filepath.Dir(path), 0644); e != nil {
			*err = e
			return
		}
	}
}
func FileCreate(path string, err *error) *os.File {
	if *err != nil {
		return nil
	}
	path = OptSeparator(path)
	if Exists(path) {
		return nil
	}
	if !Exists(filepath.Dir(path)) {
		if e := os.MkdirAll(filepath.Dir(path), 0644); e != nil {
			*err = e
			return nil
		}
	}
	file, e := os.OpenFile(path, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0660)
	*err = e
	return file
}

func FileReCreate(path string, err *error) *os.File {
	if *err != nil {
		return nil
	}
	path = OptSeparator(path)
	if Exists(path) {
		if *err = os.Remove(path); *err != nil {
			return nil
		}
	}
	return FileCreate(path, err)
}
func FileJoins(elem ...string) string {
	path := filepath.Join(elem...)
	return FilePath(path, new(error))
}
func ToFile(path string, err *error) *os.File {
	if *err != nil {
		return nil
	}
	if file, e := os.OpenFile(path, os.O_RDWR|os.O_CREATE, os.ModePerm); e != nil {
		*err = e
		return nil
	} else {
		return file
	}
}
func FileSuffixRename(path, newSuffix string, err *error) string {
	path = FilePath(path, err)
	return FilePath(fmt.Sprintf("%s/%s%s", filepath.Dir(path), FileNamePrefix(path), newSuffix),err)
}
func FileStreamCopy(targetFile *os.File,reader io.Reader,count int64,err *error)*os.File  {
	if *err!=nil{
		return targetFile
	}
	if wn,e:=io.Copy(targetFile,reader);e!=nil{
		*err=fmt.Errorf("FileStreamCopy:%w",e)
		return targetFile
	}else if wn!=count{
		*err=errors.New("FileStreamCopy:wn!=count")
		return targetFile
	}else {
		for  {
			if fi,e:=targetFile.Stat();e!=nil{
				*err=fmt.Errorf("FileStreamCopy:%w",e)
				return targetFile
			}else if fi.Size()==count{
				return targetFile
			}else {
				log.Println("No Size :",fi.Size())
			}
		}
	}
}
func FileCopy(src, target string, err *error) {
	if *err != nil {
		return
	}
	src = FilePath(src, err)
	target = FilePath(target, err)

	FileCreate(target, err)

	if srcFile, e := os.Open(src); e != nil {
		*err = fmt.Errorf("open srcFile:%w", e)
		return
	} else if targetFile, e := os.OpenFile(target, os.O_WRONLY, 0660); e != nil {
		*err = fmt.Errorf("open targetFile:%w", e)
		return
	} else if _, e := io.Copy(targetFile, srcFile); e != nil {
		*err = fmt.Errorf("copy target to src:%w", e)
		return
	}else {
		FileClose(srcFile,err)
		FileClose(targetFile,err)
	}
}
func FileTempDir(path string, err *error) string {
	if *err!=nil{
		return path
	}
	if tempDir, e := ioutil.TempDir(path, "devcom-"); e != nil {
		*err = e
		return path
	} else {
		return tempDir
	}

}
func FileDelete(path string, err *error) {
	if *err != nil {
		return
	}
	path = FilePath(path, err)
	if !Exists(path) {
		return
	}
	if FileIsDir(path, err) {
		dir, e := ioutil.ReadDir(path)
		if e!=nil{
			*err=e
			return
		}
		for _, d := range dir {
			log.Println("dName:",d.Name())
			for e := os.RemoveAll(d.Name()); e != nil; e = os.Remove(path) {
				log.Println("removeAll error:", e)
			}
		}
		_=os.RemoveAll(path)

	} else {
		for e := os.Remove(path); e != nil; e = os.Remove(path) {
			log.Println("remove error:", e)
		}
	}
}
func FileIsDir(path string, err *error) bool {
	if *err != nil {
		return false
	}
	if stat, e := os.Stat(path); e != nil {
		*err = e
		return false
	} else {
		return stat.IsDir()
	}

}
func FileDir(path string,err *error) string {
	return FilePath(filepath.Dir(path),err)
}
func FileClose(file *os.File,err *error) {
	if *err!=nil{
		return
	}
	if e:=file.Close();e!=nil{
		*err=e
		return
	}
	log.Println("Close:",file.Name())
}
func IoCopy(dst io.Writer, src io.Reader,err *error) {
	if *err!=nil{
		return
	}
	if wn,e:=io.Copy(dst,src);e!=nil{
		*err=fmt.Errorf("IoCopy writeLen:%d,error:%w",wn,e)
		return
	}else if v,ok:=src.(*os.File);ok{
		if e:=v.Close();e!=nil{
			*err=e
			return
		}
		return
	}else {
		return
	}
}
func GetBoolEnv(key string) bool {
	return strings.ToLower(os.Getenv(key))=="true"
}
func GetAppDir() string {
	if GetBoolEnv("APP_DEBUG"){
		return OptSeparator("./")
	}
	file, _ := exec.LookPath(os.Args[0])
	return FileDir(file,new(error))
}