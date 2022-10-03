package restful

import (
	"bytes"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/qntfy/kazaam"
)

type Options struct {
	Method      string
	Payload     string
	Headers     map[string]string
	Transformer string
}

func Hello() string {
	return "Hello world"
}

func Call(url string, options *Options) (string, int) {
	method := "GET"
	var data []byte

	if options.Payload != "" {
		data = []byte(options.Payload)
	}
	if options.Method != "" {
		method = options.Method
	}
	log.Println("{0} Performing Http {1}...", url, strings.ToUpper(method))
	client := &http.Client{}
	var req *http.Request
	if data != nil {
		req, _ = http.NewRequest(method, url, bytes.NewBuffer(data))
	} else {
		req, _ = http.NewRequest(method, url, nil)
	}

	if options.Headers != nil {
		if _, ok := options.Headers["Content-Type"]; !ok {
			req.Header.Set("Content-Type", "application/json")
		}
		for key, value := range options.Headers {
			req.Header.Set(key, value)
		}
	}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalln(err)
		return err.Error(), resp.StatusCode
	}

	defer resp.Body.Close()
	bodyBytes, _ := ioutil.ReadAll(resp.Body)

	// Convert response body to string
	bodyString := string(bodyBytes)
	if options.Transformer != "" && resp.StatusCode >= 200 && resp.StatusCode <= 299 {
		if !kazaam.IsJsonFast(bodyBytes) {
			return "Invalid JSON", 400
		}
		k, _ := kazaam.NewKazaam(options.Transformer)
		json, e := k.TransformJSONStringToString(bodyString)
		if e != nil {
			return e.Error(), 500
		}
		return json, resp.StatusCode
	}
	return bodyString, resp.StatusCode
}
