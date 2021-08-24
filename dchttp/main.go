package dchttp

import (
	"errors"
	"net/http"
	"time"
)

func Error(w http.ResponseWriter, error string, code int)  {

}

func TryGet(url string,timeoutMS int64) (*http.Response, error) {
	curTime := time.Now().UnixNano() / 1000/1000
	for true {
		if res,e:=http.Get(url);e==nil{
			return res,nil
		}
		newTime:=time.Now().UnixNano() / 1000/1000
		if newTime-curTime>timeoutMS{
			return nil,errors.New("timeout get:"+url)
		}
		time.Sleep(500*time.Millisecond)
	}
	return nil,errors.New("unknown error:TryGet:"+url)
}