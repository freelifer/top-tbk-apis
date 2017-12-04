package tbk

import (
	"bytes"
	"crypto/hmac"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"sort"
	"strings"
	"time"
)

const (
	SIGN_METHOD_MD5  string = "md5"
	SIGN_METHOD_HMAC string = "hmac"
)

func executeRequest(req ITbkRequest, resp interface{}) ([]byte, error) {
	// 1. 创建http.Client
	client := &http.Client{}

	reqUrl := req.GetReqUrl()
	log.Println(reqUrl)

	d := req.GetBody()
	log.Println("-------------6666", string(d))
	request, err := http.NewRequest("POST", reqUrl, strings.NewReader(d))
	if err != nil {
		log.Println("-------------5555", string(d))
		return nil, err
	}

	// 下面这句必须要，否则Post参数不对
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	log.Println(request)
	response, err := client.Do(request)
	if err != nil {
		log.Println("------------444444444", err)
		return nil, err
	}
	respBody := response.Body
	defer respBody.Close()

	data, err := ioutil.ReadAll(respBody)
	if err != nil {
		log.Println("------------333333", err)
		return nil, err
	}

	var errResp TbkErrResponse
	err = json.Unmarshal(data, &errResp)
	if err != nil {
		log.Println("------------222222", err)
		return data, err
	}
	errstr := errResp.GetErr()
	if errstr != "" {
		log.Println("------------1111111", errstr)
		return data, errors.New(errstr)
	}

	log.Println("------------", errstr)
	return data, json.Unmarshal(data, resp)
}

func currentTime() string {
	t := time.Now()
	return t.Format("2006-01-02 15:04:05")
}

func signTopRequest(reqSysMap *TbkMap, reqAppMap *TbkMap, app_secret string, sign_method string) string {
	table := make(map[string]string)
	tableLen := 0
	if reqSysMap != nil {
		tableLen += len(reqSysMap.Values)
		for k, v := range reqSysMap.Values {
			table[k] = v
		}

	}
	if reqAppMap != nil {
		tableLen += len(reqAppMap.Values)
		for k, v := range reqAppMap.Values {
			table[k] = v
		}
	}

	// 第一步：检查参数是否已经排序
	keys := make([]string, len(table))
	i := 0
	for k, _ := range table {
		keys[i] = k
		i++
	}
	sort.Strings(keys)

	// 第二步：把所有参数名和参数值串在一起
	var buffer bytes.Buffer
	if strings.Compare(sign_method, SIGN_METHOD_MD5) == 0 {
		buffer.WriteString(app_secret)
	}

	log.Print(keys)
	for _, k := range keys {
		value := table[k]
		if areNotEmpty(k, value) {
			buffer.WriteString(k)
			buffer.WriteString(value)
		}
	}

	log.Println(buffer.String())
	var sign string
	// 第三步：使用MD5/HMAC加密
	if strings.Compare(sign_method, SIGN_METHOD_HMAC) == 0 {
		sign = Hmac(app_secret, buffer.String())
	} else {
		buffer.WriteString(app_secret)
		sign = Md5(buffer.String())
	}
	return sign
}

func Hmac(key, data string) string {
	hmac := hmac.New(md5.New, []byte(key))
	hmac.Write([]byte(data))
	return strings.ToUpper(hex.EncodeToString(hmac.Sum([]byte(""))))
}

func Md5(data string) string {
	md5 := md5.New()
	md5.Write([]byte(data))
	md5Data := md5.Sum([]byte(""))
	return hex.EncodeToString(md5Data)
}

func Md5_2(data string) string {
	md5 := md5.New()
	md5.Write([]byte(data))
	md5Data := md5.Sum([]byte(""))
	return strings.ToUpper(fmt.Sprintf("%x", md5Data))
}

func createUrl(urlPath string, reqSysMap *TbkMap) string {
	var buffer bytes.Buffer
	buffer.WriteString(urlPath + "?")
	i := 0
	tableLen := len(reqSysMap.Values)
	for k, v := range reqSysMap.Values {
		i++
		buffer.WriteString(k)
		buffer.WriteString("=")
		buffer.WriteString(url.QueryEscape(v))
		if tableLen != i {
			buffer.WriteString("&")
		}
	}
	return buffer.String()
}

func getBody(reqAppMap *TbkMap) string {
	if reqAppMap == nil {
		return ""
	}
	var buf bytes.Buffer
	keys := make([]string, 0, len(reqAppMap.Values))
	for k := range reqAppMap.Values {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		vs := reqAppMap.Values[k]
		prefix := url.QueryEscape(k) + "="
		if buf.Len() > 0 {
			buf.WriteByte('&')
		}
		buf.WriteString(prefix)
		buf.WriteString(url.QueryEscape(vs))
	}
	return buf.String()
}
