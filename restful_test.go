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
	options.JS = `const x = data.results.map(u => u.name)[0];
	const { first, last } = x;` +
		"const result = {\"name\": `${last}, ${first}`};result;"
	message, _ := Call("https://randomuser.me/api/", &options)
	fmt.Print(message)
	t.Log(message)
}
