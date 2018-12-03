package main

import (
	"fmt"
	"net/http"
	"os"
	"time"
)

//定义struct,作为http的handler
type MessageHandler struct {
	//处理get/post数据的channel
	ch chan string
}

//MessageHandler需要实现ServeHTTP方法，才能正式成为一个handler
func (m *MessageHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	if len(r.FormValue("message")) > 0 {
		fmt.Println("message:", r.FormValue("message"))
		fmt.Fprintf(w, r.FormValue("message"))
		m.ch <- r.FormValue("message")
	}
}

//启动http server
func startHttpServer(ch chan string, f *os.File) {
	//初始化MessageHandler，并注册到http中
	http.Handle("/test", &MessageHandler{ch})
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println("http server start error!", err)
		return
	}
}

//写文件的loop
func fileWriterLoop(ch chan string, f *os.File) {

	//声明并初始化一个capacity=2的slice
	msgBuffer := make([]string, 0, 2)
	for {
		//从ch中读取数据，并附上换行符，放入slice
		msgBuffer = append(msgBuffer, <-ch+"\r\n")
		if len(msgBuffer) == 2 { //slice长度==2，则写文件
			for _, msg := range msgBuffer {
				//把slice中的字符串，一个个写入文件
				_, err := f.WriteString(msg)
				if err != nil {
					fmt.Println("write message to file error：", err)
				}
			}
			//重置slice，内容清空，长度设置为0
			msgBuffer = msgBuffer[0:0]
		}
	}
}

func main3() {
	//声明并初始化一个带buffer的channel
	ch := make(chan string, 99)
	f, err := os.OpenFile("d:\\test.txt", os.O_RDWR|os.O_CREATE, 0766)
	if err != nil {
		fmt.Println("open file failed")
		return
	}
	defer f.Close()

	//1st. start fileWriter goroutine
	go fileWriterLoop(ch, f)

	//2nd. start receiver goroutine
	go startHttpServer(ch, f)

	//waiting....
	time.Sleep(time.Duration(10000) * time.Second)

}

func main() {
	ch := make(chan int)
	param := 3
	var param1 int
	go func() {
		select {
		case param1 = <-ch:
			fmt.Printf("received %d from ch\n", param1)
			//default:
			//fmt.Printf("no communication222\n")
		}
	}()

	go func() {
		select {
		case ch <- param:
			fmt.Printf("sent %d to ch\n", param)
		default:
			fmt.Println("no communication111")
		}
	}()

	time.Sleep(time.Duration(2000))
}
