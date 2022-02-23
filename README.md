# go-hcl
hcl is a replacement for log which wraps hc-log 
it does not support Fatal or Panic function

```go
ppackage main

import "github.com/vogtp/go-hcl"

func main() {
	hcl.Print("Hello HCL")
	hcl.Printf("Hello %s", "HCL")
	hcl.Error("I am getting bored")
}

```

Output: 

```
2022-02-23T15:48:19+01:00 [INFO]  hcltests: Hello HCL
2022-02-23T15:48:19+01:00 [INFO]  hcltests: Hello HCL
2022-02-23T15:48:19+01:00 [ERROR] hcltests: I am getting bored
```