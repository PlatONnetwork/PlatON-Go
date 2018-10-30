package core

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

func Send(params interface{}, action, url string) (string, error) {
	param := JsonParam{
		Jsonrpc: "2.0",
		Method:  action,
		Params:  params,
		Id:      1,
	}
	resp, err := HttpPost(param, url)
	if err != nil {
		panic(fmt.Sprintf("send http post error .\n %s" + err.Error()))
	}
	return resp, err

}

func HttpPost(param JsonParam, url string) (string, error) {

	client := &http.Client{}
	req, _ := json.Marshal(param)
	req_new := bytes.NewBuffer(req)

	//fmt.Printf("request json data：%s\n",string(req))

	request, _ := http.NewRequest("POST", url, req_new)
	request.Header.Set("Content-type", "application/json")
	response, err := client.Do(request)
	if err == nil && response.StatusCode == 200 {
		body, _ := ioutil.ReadAll(response.Body)
		return string(body), nil
	} else {
		panic(fmt.Sprintf("http response status :%s", response.Status))
	}
	return "", err
}
