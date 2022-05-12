package myHttp

import (
	"bytes"
	"fmt"
	"net/http"
	"time"
)

func HttpRequest(method string, url string, jsonData []byte, auth string) (*http.Response, error) {
	//http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	client := http.Client{
		Timeout: time.Second * 5, // Timeout after 5 seconds
	}

	request, err := http.NewRequest(method, url, bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Println(err)
	}

	if auth != "" {
		request.Header.Set("Authorization", auth)
	}
	request.Header.Set("Content-Type", "application/json")

	res, err := client.Do(request)
	if err != nil {
		fmt.Println(err)
	}

	return res, nil
}
