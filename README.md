# config
parse config from a config file.
a config file like this:
```
//server.cfg
url="http://10.56.0.94:8802/gdmon/index.asp"
key=a value
mysql_username=root
mysql_pwd=pwd
mysql_host=127.0.0.1
mysql_port=3306
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
	mysql_url := fmt.Sprintf("%s:%s@tcp(%s:%s)/esns_production?charset=utf8", config.GetMulti("mysql_username", "mysql_pwd", "mysql_host", "mysql_port")...)
}
```

and the output will like this:
```
key is:  a value
config is:  &map[url:http://10.56.0.94:8802/gdmon/index.asp key:a value mysql_username:root mysql_pwd:pwd mysql_host:127.0.0.1 mysql_port:3306]
root:pwd@tcp(127.0.0.1:3306)/esns_production?charset=utf8

```
