package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"
)

const (
	message = "shanghai"
)

var dataList = []string{"test1", "test2", "test3", "test4", "test5", "test6", "test7", "test8", "test9", "test10"}

func TestSendMessage(t *testing.T) {

	for _, data := range dataList {
		t.Run(fmt.Sprintf("send:%s", data), func(t *testing.T) {

			resp, err := http.Get("http://localhost:8080/test?message=" + data)

			if err != nil {
				// handle error
			}
			defer resp.Body.Close()

			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				// handle error
			}
			fmt.Println(string(body))
			// 检测返回的数据
			if string(body) != data {
				t.Errorf("handler returned unexpected body: got %v want %v", string(body), data)
			}
		})
	}
}
