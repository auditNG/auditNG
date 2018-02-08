package source

import (
	"fmt"
	"github.com/buger/jsonparser"
	"io/ioutil"
	"net/http"

	"github.com/pre-processink/util"
)

func NewESSource() Source {
	return ESSource{}
}

type ESSource struct{}

func (es ESSource) Fetch() (string, error) {

	// Read config
	f := "es.json"
	config, err := ioutil.ReadFile(f)
	if err != nil {
		fmt.Println("Error reading file: ", err, ", File: ", f)
		return "", err
	}

	uri, err := jsonparser.GetString(config, "es_config", "uri")
	if err != nil {
		fmt.Println("Error parsing config: ", err)
		return "", err
	}

	body, err := jsonparser.GetString(config, "es_config", "payload")
	if err != nil {
		fmt.Println("Error parsing config: ", err)
		return "", err
	}

	header := http.Header{
		"Content-Type": []string{"application/json"},
	}

	resp, err := util.SendReq(
		"GET",
		uri,
		header,
		[]byte(body))

	if err != nil {
		fmt.Println("ES Request Error: ", err)
		fmt.Println("ES Error Response body: ", resp)
		return "", err
	}

	return string(resp), nil
}
