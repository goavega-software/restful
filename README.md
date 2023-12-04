# REST API Manager 

Allows you to call any API and then transform the returned json data and return back. Restful also supports using JSONPath style xpath to filter the JSON data. Currently, if both XPath and Transformer are provided, restful pipes the output of XPath to Transformer and then returns the final result. 

Example use

```golang
func main() {
	// Get a message and print it.
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
}
```

The above ex. makes the API call, extracts the name object from JSON, removes the `title` and returns back.

JSON Transformation supports Jolt like transformations using [Kazaam](https://github.com/qntfy/kazaam)
XPath provides complete implementation of http://goessner.net/articles/JsonPath/ using [JsonPath](https://github.com/PaesslerAG/jsonpath)
Supports Basic authorization, and token/key based authentication.

There is basic string interpolation support in url, raw body and headers. Tokens in the form of `${var}` are replaced with `os.Getenv(var)`. Below is one example:

```golang
func main() {
	options := restful.Options{}
	options.Method = "GET"
	options.Headers = make(map[string]string)
	options.Headers["Content-Type"] = "application/json"
	options.Transformer = `[
		{
			"operation": "shift", 
			  "spec": {
				"data.title": "joke",
				"data.subtitle": "category"
			}
		},
		{
		"operation": "default",
		"spec": {
		  "event": "jotd"
		}
	  }
	  ]`
	os.Setenv("type", "single")
	message, _ := restful.Call("https://v2.jokeapi.dev/joke/Programming?blacklistFlags=nsfw,religious,political,racist,sexist,explicit&type=${type}", &options)
	fmt.Println(message)
}

```
## Embedded JavaScript Engine
Starting from version 0.4, Restful introduces an integrated ES2015/ES6 JavaScript runtime for scenarios where basic XPath and Transforms fall short. The options.JS parameter accepts a JavaScript literal that enables the transformation of input data. The JSON obtained from the API call is accessible through the *data* variable within the JavaScript code block. The JavaScript string should conclude with the value to be exported, as illustrated in the example below with the result being the exported output.

```golang
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
```

Sequencing of transforms (output of each transform is fed into the next step as input):
1. JS 
2. XPath
3. Transformer