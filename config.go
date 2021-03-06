package config

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
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

func ParseToModel(filename string) (model map[string]string, err error) {
	model = make(map[string]string)

	fi, err := os.Stat(filename)
	if err != nil {
		return
	}

	f, err := os.Open(filename)
	defer f.Close()
	if err != nil {
		return
	}

	buff := make([]byte, fi.Size())
	f.Read(buff)

	str := ptn.FindAllString(string(buff), -1)

	for i := 0; i < len(str); i++ {
		switch {
		case strings.Index(str[i], "#") == 0:
			// 说明是注释行，什么也不做
		case strings.Index(str[i], "=") > 0:
			index := strings.Index(str[i], "=")
			key, val := strings.TrimSpace(str[i][:index]), strings.Trim(strings.TrimSpace(str[i][index+1:]), `"`)
			//       kvs := strings.Split(str[i], "=")
			//       key, val := strings.TrimSpace(kvs[0]), strings.Trim(kvs[1], `"`)
			model[key] = val
		default:
			err = errors.New("parse config files failed! please check this line:" + str[i])
			return
		}
	}

	encrypted_key := model["encrypted_key"]
	if encrypted_key != "" {
		for key, value := range model {
			if strings.Index(key, "_encrypted_") == 0 {
				// if key == "_encrypted_mysql_pwd" {
				//   decrypted_key := "mysql_pwd"
				decrypted_key := strings.TrimLeft(key, "_encrypted_")
				decrypted_value := Decrypt(value, encrypted_key, "a0fe7c7c98e09e8c")

				model[decrypted_key] = decrypted_value
			}
		}
	}

	return
}

func Parse(filename string) (err error) {
	fi, err := os.Stat(filename)
	if err != nil {
		return
	}

	f, err := os.Open(filename)
	defer f.Close()
	if err != nil {
		return
	}

	buff := make([]byte, fi.Size())
	f.Read(buff)

	str := ptn.FindAllString(string(buff), -1)

	for i := 0; i < len(str); i++ {
		switch {
		case strings.Index(str[i], "#") == 0:
			// 说明是注释行，什么也不做
		case strings.Index(str[i], "=") > 0:
			index := strings.Index(str[i], "=")
			key, val := strings.TrimSpace(str[i][:index]), strings.Trim(strings.TrimSpace(str[i][index+1:]), `"`)
			//       kvs := strings.Split(str[i], "=")
			//       key, val := strings.TrimSpace(kvs[0]), strings.Trim(kvs[1], `"`)
			conf[key] = val
		default:
			err = errors.New("parse config files failed! please check this line:" + str[i])
			return
		}
	}

	encrypted_key := conf["encrypted_key"]
	if encrypted_key != "" {
		for key, value := range conf {
			if strings.Index(key, "_encrypted_") == 0 {
				// if key == "_encrypted_mysql_pwd" {
				//   decrypted_key := "mysql_pwd"
				decrypted_key := strings.TrimLeft(key, "_encrypted_")
				decrypted_value := Decrypt(value, encrypted_key, "a0fe7c7c98e09e8c")

				conf[decrypted_key] = decrypted_value
			}
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

func Clean() {
	conf = make(map[string]string)
}

func Decrypt(encrypted, key, iv string) string {
	_key := []byte(md5_digest(md5_hexdigest(key)))

	block, err := aes.NewCipher(_key)
	if err != nil {
		return ""
	}

	_iv := []byte(iv)
	ciphertext, _ := hex.DecodeString(encrypted)

	if len(ciphertext)%aes.BlockSize != 0 {
		return ""
	}

	mode := cipher.NewCBCDecrypter(block, _iv)
	mode.CryptBlocks(ciphertext, ciphertext)

	return string(pkcs7_unpadding(ciphertext))
}

func Encrypt(data, key, iv string) string {
	_key := []byte(md5_digest(md5_hexdigest(key)))

	block, err := aes.NewCipher(_key)
	if err != nil {
		return ""
	}

	_iv := []byte(iv)
	padding_data := pkcs7_padding([]byte(data))

	encrypted_data := make([]byte, len(padding_data))
	mode := cipher.NewCBCEncrypter(block, _iv)
	mode.CryptBlocks(encrypted_data, padding_data)

	return hex.EncodeToString(encrypted_data)

}

func md5_hexdigest(str string) string {
	data := []byte(str)

	return fmt.Sprintf("%x", md5.Sum(data))
}

func md5_digest(str string) string {
	data := []byte(str)

	return fmt.Sprintf("%s", md5.Sum(data))
}

func pkcs7_unpadding(origData []byte) []byte {
	length := len(origData)

	unpadding := int(origData[length-1])
	if length <= unpadding {
		return make([]byte, 0)
	}

	return origData[:(length - unpadding)]
}

func pkcs7_padding(origData []byte) []byte {
	blockSize := aes.BlockSize

	padding := blockSize - len(origData)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(origData, padtext...)
}
