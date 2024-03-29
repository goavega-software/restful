package restful

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"

	"github.com/PaesslerAG/jsonpath"
	"github.com/dop251/goja"
	"github.com/qntfy/kazaam"
)

type Options struct {
	Method      string
	Payload     string
	Headers     map[string]string
	Transformer string
	XPath       string
	JS          string
}

func interfacify(input []string) []interface{} {
	vals := make([]interface{}, len(input))
	for i, v := range input {
		vals[i] = v
	}
	return vals
}

func tokenize(input string) (string, []string) {
	re := regexp.MustCompile(`(?m)\$\{(.+?)\}`)
	substitution := "%s"
	var variables []string
	for _, variable := range re.FindAllStringSubmatch(input, -1) {
		variables = append(variables, os.Getenv(variable[1]))
	}
	return re.ReplaceAllString(input, substitution), variables
}

func cleanString(format string, variables ...any) string {
	return fmt.Sprintf(format, variables...)
}

func Call(url string, options *Options) (string, int) {
	method := "GET"
	var data []byte
	format, tokens := tokenize(url)
	parsedUrl := cleanString(format, interfacify(tokens)...)

	if options.Payload != "" {
		format, tokens := tokenize(options.Payload)
		data = []byte(cleanString(format, interfacify(tokens)...))
	}
	if options.Method != "" {
		method = options.Method
	}
	log.Println("Performing Http ...", parsedUrl, strings.ToUpper(method))
	client := &http.Client{}
	var req *http.Request
	if data != nil {
		req, _ = http.NewRequest(method, parsedUrl, bytes.NewBuffer(data))
	} else {
		req, _ = http.NewRequest(method, parsedUrl, nil)
	}

	if options.Headers != nil {
		if _, ok := options.Headers["Content-Type"]; !ok {
			req.Header.Set("Content-Type", "application/json")
		}
		for key, value := range options.Headers {
			format, tokens := tokenize(value)
			req.Header.Set(key, cleanString(format, interfacify(tokens)...))
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
	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return bodyString, resp.StatusCode
	}
	if !kazaam.IsJsonFast(bodyBytes) {
		return "Invalid JSON", 400
	}
	v := interface{}(nil)
	json.Unmarshal([]byte(bodyString), &v)
	if options.JS != "" {
		vm := goja.New()
		vm.SetFieldNameMapper(goja.TagFieldNameMapper("json", true))
		vm.Set("data", v)
		val, err := vm.RunString(options.JS)
		if err != nil {
			panic(err)
		}
		v = val.Export()
		// convert the map to JSON
		b, err := json.Marshal(v)
		if err != nil {
			return err.Error(), 500
		}
		bodyString = string(b)
	}
	if options.XPath != "" {
		log.Println("Performing xpath with ", options.XPath)

		doc, err := jsonpath.Get(options.XPath, v)
		if err != nil {
			return err.Error(), 500
		}
		// convert the map to JSON
		b, err := json.Marshal(doc)
		if err != nil {
			return err.Error(), 500
		}
		bodyString = string(b)
		log.Println("xpath result ", bodyString)
	}
	if options.Transformer != "" {
		log.Println("Performing transformation with ", options.Transformer)
		k, _ := kazaam.NewKazaam(options.Transformer)
		json, e := k.TransformJSONStringToString(bodyString)
		if e != nil {
			return e.Error(), 500
		}
		log.Println("Performed transformation with {0}", json)
		return json, resp.StatusCode
	}
	return bodyString, resp.StatusCode
}
