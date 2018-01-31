package util

import (
	"bytes"
	"crypto/tls"
	"errors"
	"fmt"
	"net/http"
  "io/ioutil"
)

// SendReq is used to send a http/https request to a URL defind by u, using
// method m, optional headers h and an optional body b
func SendReq(m, u string, h http.Header, b []byte) (body []byte, err error) {
	var req *http.Request
	req, err = http.NewRequest(m, u, bytes.NewBuffer(b))
	if err != nil {
		fmt.Println("Invalid Request error: ", err)
		return
	}

	req.Header = h

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Request error: ", err)
		return
	}
	defer resp.Body.Close()

	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Response error: ", err)
		err = fmt.Errorf("Error: %v; Status: %d", err, resp.StatusCode)
		return
	}

	if resp.StatusCode > 399 {
		fmt.Println("Response code: ", resp.StatusCode)
		err = errors.New("Status " + resp.Status)
	}

	return
}
