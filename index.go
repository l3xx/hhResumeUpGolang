package main

import (
	"fmt"
	"net/http"
	"strings"
	"io/ioutil"
	"io"
	"encoding/json"
	"os"
	"log"
)

type SampleParamReq struct {
	url string
	method string
	token string
	data io.Reader
}

func downloadFromUrl(dataParamReq SampleParamReq) (status int, body string) {
	tokens := strings.Split(dataParamReq.url, "/")
	fileName := tokens[len(tokens)-1]
	fmt.Println("Downloading", dataParamReq.url, "to", fileName)

	client := &http.Client{}
	//req, err := http.NewRequest("POST", "http://example.com", bytes.NewReader(postData))
	req, err := http.NewRequest(dataParamReq.method, dataParamReq.url, dataParamReq.data)
	if dataParamReq.token != "" {
		req.Header.Add("User-Agent", "Letunovskiymn/1.0 (miha-1221@inbox.ru)")
		req.Header.Add("Accept", "*/*")
		req.Header.Add("Authorization", "Bearer "+ token)
	}
	response, err := client.Do(req)


	if err != nil {
		fmt.Println("Error while downloading", dataParamReq.url, "-", err)
	}
	defer response.Body.Close()

	fmt.Println( dataParamReq.url, response)
	bodyBytes, _ := ioutil.ReadAll(response.Body)
	return response.StatusCode, string(bodyBytes)
}

var token,dataBody string

func main() {
	configUrl :="https://opendevelopers.ru/getKey.php"
	baseUrl :="https://api.hh.ru"

	statusCode, body :=downloadFromUrl(SampleParamReq{url: configUrl, method: "GET"})
	if statusCode == http.StatusOK {
		token = body
	}
	//fmt.Println(token)

	statusCode, body = downloadFromUrl(SampleParamReq{url: baseUrl + "/resumes/mine", method: "GET", token: token})
	if statusCode == http.StatusOK {
		dataBody = body
	}
	fmt.Println(dataBody)

	var dat map[string]interface{}

	if err := json.Unmarshal([]byte(dataBody), &dat); err != nil {
		panic(err)
	}

	f, err := os.OpenFile("log", os.O_RDWR | os.O_CREATE | os.O_APPEND, 0775)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	log.SetOutput(f)
	for _,val := range dat["items"].([]interface{}){
		item :=val.(map[string]interface{})
		var id  = item["id"].(string)
		statusCode, body = downloadFromUrl(SampleParamReq{url: baseUrl + "/resumes/" + id +"/publish", method: "POST", token: token})
		log.Println(statusCode)
		log.Println(body)
	}
}