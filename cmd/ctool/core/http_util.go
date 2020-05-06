// Copyright 2018-2020 The PlatON Network Authors
// This file is part of the PlatON-Go library.
//
// The PlatON-Go library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The PlatON-Go library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the PlatON-Go library. If not, see <http://www.gnu.org/licenses/>.

package core

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

func Send(params interface{}, action string) (string, error) {
	param := JsonParam{
		Jsonrpc: "2.0",
		Method:  action,
		Params:  params,
		Id:      1,
	}
	resp, err := HttpPost(param)
	if err != nil {
		panic(fmt.Sprintf("send http post error .\n %s" + err.Error()))
	}

	return resp, err
}

func HttpPost(param JsonParam) (string, error) {

	client := &http.Client{}
	req, _ := json.Marshal(param)
	reqNew := bytes.NewBuffer(req)

	request, _ := http.NewRequest("POST", config.Url, reqNew)
	request.Header.Set("Content-type", "application/json")
	response, err := client.Do(request)
	if response == nil && err != nil {
		panic(fmt.Sprintf("no response from node,%s", err.Error()))
	}
	if err == nil && response.StatusCode == 200 {
		body, _ := ioutil.ReadAll(response.Body)
		return string(body), nil
	} else {
		panic(fmt.Sprintf("http response status :%s", response.Status))
	}
}

func parseResponse(r string) *Response {
	var resp = Response{}
	err := json.Unmarshal([]byte(r), &resp)

	if err != nil {
		panic(fmt.Sprintf("parse result error ! error:%s \n", err.Error()))
	}

	if resp.Error.Code != 0 {
		panic(fmt.Sprintf("send transaction error ,error:%v \n", resp.Error.Message))
	}
	return &resp
}
