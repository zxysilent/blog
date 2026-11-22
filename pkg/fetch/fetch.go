package fetch

import (
	"bytes"
	"encoding/json"
	"net/http"
	"time"
)

// http 请求用
// 总计超时时间限制5s
var client = http.Client{Timeout: 5 * time.Second}

// Get HTTP 工具类, GET 并解析返回的报文，如果有错误，返回 error
// url 统一资源定位符（uniform resource locator;URL)
// reply 返回值
func Get(url string, reply interface{}) error {
	resp, err := client.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return json.NewDecoder(resp.Body).Decode(reply)
}

func Post(url string, data interface{}, reply interface{}) error {
	body, _ := json.Marshal(data)
	resp, err := client.Post(url, "application/json", bytes.NewBuffer(body))
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return json.NewDecoder(resp.Body).Decode(reply)
}
