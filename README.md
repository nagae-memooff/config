# config
parse config from a config file.
a config file like this:
```
//server.cfg
url="http://10.56.0.94:8802/gdmon/index.asp"
key=a value
```

parse the file like this:


```go
package main

import (
	"fmt"
  "github.com/nagae-memooff/config"
	logger "log"
)


func main() {
	err := config.Parse("server.cfg")
	if err != nil {
		logger.Fatal("load config failed." + err.Error())
	}

	fmt.Println("key is: ", config.Get("key"))
	fmt.Println("config is: ", config.GetAll())
}
```

and the output will like this:
```
key is:  a value
config is:  &map[url:http://10.56.0.94:8802/gdmon/index.asp key:a value]

```
