package config

import (
	"errors"
	"os"
	"regexp"
	"strings"
//   "fmt"
)

var ptn *regexp.Regexp
var parse_err error
var conf = make(map[string]string)

func init() {
	ptn, parse_err = regexp.Compile(`[#]*[A-Za-z0-9\_\-\.\t ]+[=].*`)
	if parse_err != nil {
		panic(parse_err)
	}
}

func Parse(filename string) (err error) {
	fi, err := os.Stat(filename)
	if err != nil {
		return
	}

	f, err := os.Open(filename)
	if err != nil {
		return
	}

	buff := make([]byte, fi.Size())
	f.Read(buff)
	f.Close()

	str := ptn.FindAllString(string(buff), -1)

	for i := 0; i < len(str); i++ {
		switch {
		case strings.Index(str[i], "#") == 0:
			// 说明是注释行，什么也不做
		case strings.Index(str[i], "=") > 0:
			kvs := strings.Split(str[i], "=")
			key, val := strings.TrimSpace(kvs[0]), strings.Trim(kvs[1], `"`)
			conf[key] = val
		default:
			err = errors.New("parse config files failed! please check this line:" + str[i])
			return
		}
	}
	return
}

//调试用，返回conf对象的指针，一般情况下最好不要使用
func GetAll() *map[string]string {
	return &conf
}

func Get(key string) (value string) {
	value = conf[key]
	return
}
