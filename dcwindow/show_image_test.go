package dcwindow

import (
	"github.com/Humenger/go-devcommon"
	"golang.org/x/image/bmp"
	"image"
	"image/jpeg"
	"os"
	"sync"
	"testing"
)

func TestShowImage(t *testing.T) {
	var lock sync.WaitGroup
	go func() {
		lock.Add(1)
		if err:=ShowImage("Demo","F:\\project\\go\\wechat_test\\debug\\qrcode.png");err!=nil{
			t.Error(err)
		}
		lock.Done()
	}()
	lock.Wait()
}
func TestShowImage2(t *testing.T) {
	//if f,err:=os.Open("F:\\project\\go\\wechat_test\\debug\\qrcode.png");err!=nil{
	//	log.Println("err:",err.Error())
	//}else {
	//	var decodeData image.Image
	//	if decodeData,err=jpeg.Decode(f);err!=nil{
	//		log.Println("err2:",err.Error())
	//	}else{
	//		log.Println("err2:",decodeData)
	//		nilErr := new(error)
	//		f := devcommon.FileCreate(devcommon.FileJoins(devcommon.FileDir("F:\\project\\go\\wechat_test\\debug\\qrcode.png", nilErr), devcommon.FileNamePrefix("F:\\project\\go\\wechat_test\\debug\\qrcode.png")+".bmp"), nilErr)
	//		if *nilErr != nil {
	//			log.Println("err2:",*nilErr)
	//		}
	//		if err:=bmp.Encode(f,decodeData);err!=nil{
	//			log.Println("err3:",err.Error())
	//		}
	//	}
	//}
	Path:="F:\\project\\go\\wechat_test\\debug\\qrcode.png"
	bmpPath:=""
	if originalF, err := os.Open(Path); err != nil {
		t.Error(err)
	} else {
		var decodeData image.Image
		//if decodeData, err = png.Decode(originalF); err != nil {
			if decodeData, err = jpeg.Decode(originalF); err != nil {
				//if decodeData, err = bmp.Decode(originalF); err != nil {
				//	t.Error(err)
				//} else {
					bmpPath=Path
				}
			//}
		//}
		if bmpPath != "" {
			return
		}
		nilErr := new(error)
		f := devcommon.FileCreate(devcommon.FileJoins(devcommon.FileDir(Path, nilErr), devcommon.FileNamePrefix(Path)+".bmp"), nilErr)
		if *nilErr != nil {
			t.Error(*nilErr)
		}
		if err = bmp.Encode(f, decodeData); err != nil {
			t.Error(err)
		} else {
			bmpPath = f.Name()
			return
		}
	}
}