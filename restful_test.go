package restful

import (
	"fmt"
	"testing"
)

func TestCall(t *testing.T) {
	options := Options{}
	options.Method = "GET"
	options.Headers = make(map[string]string)
	options.Headers["Content-Type"] = "application/json"
	options.XPath = "$.results[0].name"
	options.Transformer = `
	[
		{
			"operation": "delete",
			"spec": {
			  "paths": ["title"]
			}
		  }
	  ]	
	`
	message, _ := Call("https://randomuser.me/api/", &options)
	fmt.Print(message)
	t.Log(message)
}
