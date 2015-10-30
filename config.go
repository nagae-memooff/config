package config

import (
	"errors"
	"os"
	"regexp"
	"strconv"
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
func GetModel() *map[string]string {
	return &conf
}

// 获取单个键值
func Get(key string) (value string) {
	value = conf[key]
	return
}

// 临时设置键值。本方法不会覆盖配置文件，因此重启后无法保存。
func Set(key string, value string) {
	conf[key] = value
	return
}

// 以可变参数的形式获取多个键值，并按顺序返回[]interface{}类型的切片。
func GetMulti(keys ...string) (values []interface{}) {
	for _, key := range keys {
		values = append(values, conf[key])
	}
	return
}

// 以可变参数的形式获取多个键值，并按顺序返回[]string类型的切片。
func GetMultiStrings(keys ...string) (values []string) {
	for _, key := range keys {
		values = append(values, conf[key])
	}
	return
}

//TODO:实现根据键值获取指定类型的值
func GetBool(key string) (value bool) {
	if conf[key] == "true" {
		value = true
	} else {
		value = false
	}
	return
}

func GetInt(key string) int {
	value, err := strconv.Atoi(conf[key])
	if err != nil {
		value = 0
	}
	return value
}

func GetFloat64(key string) float64 {
	value, err := strconv.ParseFloat(conf[key], 64)
	if err != nil {
		value = 0.0
	}
	return value
}

//func GetArray(key string) {
//}

func Default(key, value string) {
	if Get(key) == "" {
		Set(key, value)
	}
}
