package dcwindow

import (
	"errors"
	"github.com/Humenger/go-devcommon"
	"github.com/jthmath/winapi"
	"golang.org/x/image/bmp"
	"image"
	"image/jpeg"
	"image/png"
	"os"
	"runtime"
	"strconv"
)

// 全局变量
var hBitmap winapi.HBITMAP

type ImageWindow struct {
	Title   string
	Path    string
	BmpPath string
	Width   int32
	Height  int32
}

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())
}
func ShowImage(title, path string) error {
	inst, err := winapi.GetModuleHandle("")
	if err != nil {
		return err
	}

	imageWindow := &ImageWindow{
		Title: title,
		Path:  path,
	}
	if err = imageWindow.initParams(); err != nil {
		return err
	}
	if r := imageWindow.winMain(inst, "", 0); r != 0 {
		return errors.New("[ShowImage] err=" + strconv.Itoa(int(r)))
	}
	if imageWindow.BmpPath!=imageWindow.Path&&imageWindow.BmpPath!=""{
		devcommon.FileDelete(imageWindow.BmpPath,new(error))
	}
	return nil
}

func (that *ImageWindow) initParams() error {
	that.BmpPath = ""
	that.Width = 0
	that.Height = 0
	if originalF, err := os.Open(that.Path); err != nil {
		return err
	} else {
		var decodeData image.Image
		if decodeData, err = png.Decode(originalF); err != nil {
			_ = originalF.Close()
			originalF, _ = os.Open(that.Path)
			if decodeData, err = jpeg.Decode(originalF); err != nil {
				_ = originalF.Close()
				originalF, _ = os.Open(that.Path)
				if decodeData, err = bmp.Decode(originalF); err != nil {
					return err
				} else {
					that.BmpPath = that.Path
				}
			}
		}

		that.Width = int32(decodeData.Bounds().Max.X)
		that.Height = int32(decodeData.Bounds().Max.Y)
		if that.BmpPath != "" {
			return nil
		}
		nilErr := new(error)
		f := devcommon.FileCreate(devcommon.FileJoins(devcommon.FileDir(that.Path, nilErr), devcommon.FileNamePrefix(that.Path)+".bmp"), nilErr)
		if *nilErr != nil {
			return *nilErr
		}
		if err = bmp.Encode(f, decodeData); err != nil {
			return err
		} else {
			that.BmpPath = f.Name()
			return nil
		}
	}
}

func (that *ImageWindow) winMain(Inst winapi.HINSTANCE, Cmd string, nCmdShow int32) int32 {
	var err error

	// 1. 注册窗口类
	_, err = that.myRegisterClass(Inst)
	if err != nil {
		return 0
	}

	// 2. 创建窗口
	wnd, err := winapi.CreateWindow("Main Window Class", that.Title,
		winapi.WS_OVERLAPPEDWINDOW, 0,
		winapi.CW_USEDEFAULT, winapi.CW_USEDEFAULT, winapi.CW_USEDEFAULT, winapi.CW_USEDEFAULT,
		0, 0, Inst, 0)
	if err != nil {
		return 0
	}
	winapi.ShowWindow(wnd, winapi.SW_SHOW)
	winapi.UpdateWindow(wnd)

	// 3. 主消息循环
	var msg winapi.MSG
	msg.Message = winapi.WM_QUIT + 1 // 让它不等于 winapi.WM_QUIT

	for winapi.GetMessage(&msg, 0, 0, 0) > 0 {
		winapi.TranslateMessage(&msg)
		winapi.DispatchMessage(&msg)
	}

	return int32(msg.WParam)
}

func (that *ImageWindow) wndProc(hWnd winapi.HWND, message uint32, wParam uintptr, lParam uintptr) uintptr {
	var hTemp winapi.HANDLE

	switch message {
	case winapi.WM_CREATE:
		hTemp, _ = winapi.LoadImageByName(0, that.BmpPath,
			winapi.IMAGE_BITMAP, 0, 0, winapi.LR_LOADFROMFILE)
		hBitmap = winapi.HBITMAP(hTemp)
	case winapi.WM_PAINT:
		that.onPaint(hWnd)
	case winapi.WM_DESTROY:
		winapi.PostQuitMessage(0)
	case winapi.WM_COMMAND:
		that.onCommand(hWnd, wParam, lParam)
	default:
		return winapi.DefWindowProc(hWnd, message, wParam, lParam)
	}
	return 0
}

func (that *ImageWindow) onPaint(hWnd winapi.HWND) {
	var err error
	var ps winapi.PAINTSTRUCT

	hdc, err := winapi.BeginPaint(hWnd, &ps)
	if err != nil {
		return
	}
	defer winapi.EndPaint(hWnd, &ps) // defer 终于有用武之地了

	// HDC mdc = CreateCompatibleDC(hdc);
	mdc, err := winapi.CreateCompatibleDC(hdc)
	if err != nil {
		return
	}
	defer winapi.DeleteDC(mdc)

	winapi.SelectObject(mdc, winapi.HGDIOBJ(hBitmap))

	// 这个函数的第4、5个参数分别是图片的宽、高
	// 为了简便起见，我直接写在了这里
	// 实际项目中当然要用过程序获取一下
	winapi.BitBlt(hdc, 0, 0, that.Width, that.Height, mdc, 0, 0, winapi.SRCCOPY)
}

func (that *ImageWindow) onCommand(hWnd winapi.HWND, wParam uintptr, lParam uintptr) {
	// 暂时不需要特殊处理 WM_COMMAND
	winapi.DefWindowProc(hWnd, winapi.WM_COMMAND, wParam, lParam)
}

func (that *ImageWindow) myRegisterClass(hInstance winapi.HINSTANCE) (atom uint16, err error) {
	var wc winapi.WNDCLASS
	wc.Style = winapi.CS_HREDRAW | winapi.CS_VREDRAW
	wc.PfnWndProc = that.wndProc
	wc.CbClsExtra = 0
	wc.CbWndExtra = 0
	wc.HInstance = hInstance
	wc.HIcon = 0
	wc.HCursor, err = winapi.LoadCursorById(0, winapi.IDC_ARROW)
	if err != nil {
		return
	}
	wc.HbrBackground = winapi.COLOR_WINDOW + 1
	wc.Menu = uint16(0)
	wc.PszClassName = "Main Window Class"
	wc.HIconSmall = 0

	return winapi.RegisterClass(&wc)
}
