package node

import (
	"github.com/TobiasYin/go_web_ui/logs"
	"github.com/TobiasYin/go_web_ui/vdom"
	"time"
)

func init() {
	allowRerender <- 1
	SetMaxFPS(30)
	go renderLoop()
}

var (
	allowRerender = make(chan int, 1)
	needRerender  = true
	mainPage      main
	frameTime     time.Duration
)

type main struct {
	page *Page
}

func panicCatch() {
	if r := recover(); r != nil {
		logs.Errorf("panic found!, catch: %v\n", r)
	}
}

func recoverReverse() {
	if r := recover(); r != nil {
		logs.Errorf("panic found! try to recover! %v\n", r)
		page := mainPage.page
		if stack.size == 0 {
			go newApp(page)
			return
		}
		newStack := pageStack{}
		for stack.size > 0 {
			newStack.Add(stack.Pop())
		}
		go newApp(newStack.Pop())
		for newStack.size > 0 {
			top := newStack.Pop()
			if page.path != "" {
				_ = PushByPathWithPathParams(top.path)
			} else {
				PushToPage(top)
			}
		}
	}
}

func (m main) GetPage() *Page {
	return m.page
}

func newApp(page *Page) {
	defer recoverReverse()
	mainPage = main{page}
	RegisterRoute("/", func(m map[string]interface{}) PageGetter {
		return mainPage
	})
	page.path = "/"
	PushToPage(page)
	initPush()
}

func NewApp(page *Page) {
	defer recoverReverse()
	newApp(page)
	//在这里等待，防止wasm退出。
	c := make(chan struct{})
	<-c
}

// 设置每秒最大帧率
func SetMaxFPS(fps int) {
	if fps > 144 {
		fps = 144
	}
	if fps < 10 {
		fps = 10
	}
	frameTime = time.Second / time.Duration(fps)
}

func renderLoop() {
	defer recoverReverse()
	for {
		time.Sleep(frameTime)
		if !needRerender {
			continue
		}
		select {
		case <-allowRerender:
			go func() {
				//start := time.Now()
				defer recoverReverse()
				needRerender = false
				rerender()
				allowRerender <- 1
				//end := time.Now()
				//logs.Println("render page, using: ", end.Sub(start))
			}()
		default:
		}
	}
}

func FlashApp() {
	needRerender = true
}

func rerender() {
	start := time.Now()
	top := stack.Top()
	d := top.pack()
	vdom.MergeTwoTree(&d, top.oldTree)
	top.oldTree = &d
	end := time.Now()
	logs.Infof("Re Render Page, Using: %v\n", end.Sub(start))
}

func rerenderTree(newTree *vdom.JsDomElement, oldTree *vdom.JsDomElement) {
	start := time.Now()
	vdom.MergeTwoContext(newTree, oldTree)
	end := time.Now()
	logs.Infof("Re Render Page, Using: %v\n", end.Sub(start))
}
