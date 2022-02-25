# go-hcl
hcl is a replacement for log which wraps hc-log 
it does not support Fatal or Panic function

It redirects std lib log to itself

```go
package main

import "github.com/vogtp/go-hcl"

func main() {
	hcl.Print("Hello HCL")
	hcl.Printf("Hello %s", "HCL")
	hcl.Error("I am getting bored")
}

```

Output: 

```
2022-02-23T15:48:19+01:00 [INFO]  executable-name: Hello HCL
2022-02-23T15:48:19+01:00 [INFO]  executable-name: Hello HCL
2022-02-23T15:48:19+01:00 [ERROR] executable-name: I am getting bored
```
