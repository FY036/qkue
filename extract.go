package main

import (
	"encoding/json"
	"fmt"
	"strings"
)

// 定义提取内容
type Service struct {
	Name string `json:"name"`
	HTTP struct {
		Host string `json:"host"`
	} `json:"http"`
}
type Data struct {
	Port    int     `json:"port"`
	Service Service `json:"service"`
}
type SimplifiedResponse struct {
	Data    []Data `json:"data"`
	Message string `json:"message"`
}

// 提取函数接受响应体返回urls
func ExtractHTTPUrls(jsonData string) ([]string, error) {
	var response SimplifiedResponse
	err := json.Unmarshal([]byte(jsonData), &response)
	if err != nil {
		return nil, fmt.Errorf("\n\033[31m[-]\033[0m 解析JSON出错: %v\n", response.Message)
	}
	var urls []string
	for _, data := range response.Data {
		if strings.Contains(data.Service.Name, "http") {
			port := data.Port
			host := data.Service.HTTP.Host

			if host != "" {
				var url string
				if port == 443 {
					url = fmt.Sprintf("https://%s", host)
				} else if port == 80 {
					url = fmt.Sprintf("http://%s", host)
				} else {
					url = fmt.Sprintf("http://%s:%d", host, port)
				}
				urls = append(urls, url)
			}
		}
	}

	return urls, nil
}
