# REST API Manager 

Allows you to call any API and then transform the returned json data and return back.

Example use

```golang
func main() {
	// Get a message and print it.
	options := restful.Options{}
	options.Method = "GET"
	options.Headers = make(map[string]string)
	options.Headers["Content-Type"] = "application/json"
	options.Transformer = `[
		{
			"operation": "shift", 
			  "spec": {
				"data.phrase": "phrase"
			}
		},
		{
		"operation": "default",
		"spec": {
		  "event": "message"
		}
	  }
	  ]`
	message, _ := restful.Call("https://corporatebs-generator.sameerkumar.website/", &options)
	fmt.Println(message)
}
```

The above ex. makes the API call, formats the response JSON `{phrase: "Hello world"}` into `{event: "message", "data": "Hello world"}` and returns back.

JSON Transformation supports Jolt like transformations using [Kazaam](https://github.com/qntfy/kazaam)

Supports Basic authorization, and token/key based authentication.