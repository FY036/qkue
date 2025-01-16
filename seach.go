package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

// 结构
type Config struct {
	APIKey string `yaml:"apikey"`
}
type RequestBody struct {
	Include     []string `json:"include"`
	Size        int      `json:"size"`
	IgnoreCache string   `json:"ignore_cache"`
	Query       string   `json:"query"`
	Start       int      `json:"start"`
	Shortcuts   []string `json:"shortcuts"`
	Latest      string   `json:"latest"`
}
type Pagination struct {
	Total int `json:"total"`
}
type Meta struct {
	Pagination Pagination `json:"pagination"`
}
type Response struct {
	Meta Meta `json:"meta"`
}
type UserInfo struct {
    Data struct {
        User struct {
            FullName string `json:"fullname"` // 用户的全名
        } `json:"user"`
        MonthRemainingCredit int `json:"month_remaining_credit"` // 剩余的月度积分
    } `json:"data"`
}

// 发送请求函数
func sendRequest(key string, size int, sentence string, aaa, bbb, ccc, ddd bool) (string, error) {
	var shortcuts []string
	latest := "False"
	if aaa {
		latest = "True"
	}
	if bbb {
		shortcuts = append(shortcuts, "63734bfa9c27d4249ca7261c")
	}
	if ccc {
		shortcuts = append(shortcuts, "635fcb52cc57190bd8826d09")
	}
	if ddd {
		shortcuts = append(shortcuts, "635fcbaacc57190bd8826d0b")
	}
	requestBody := RequestBody{
		Include:     []string{"ip", "port", "service.http.host", "service.name"},
		Size:        size,
		IgnoreCache: "False",
		Query:       sentence,
		Start:       0,
		Shortcuts:   shortcuts,
		Latest:      latest,
	}

	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return "", fmt.Errorf("\033[31m[-]\033[0m %v", err)
	}

	client := &http.Client{}
	req, err := http.NewRequest("POST", "https://quake.360.net/api/v3/search/quake_service", bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("\033[31m[-]\033[0m 创建请求出错: %v", err)
	}

	req.Header.Set("X-QuakeToken", key)
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("\033[31m[-]\033[0m 发送请求出错: %v", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("\033[31m[-]\033[0m 读取响应出错: %v", err)
	}

	return string(body), nil
}

// 打印积分
func GetUserInfo(key string) (string, int, error) {
    client := &http.Client{}
    req, _ := http.NewRequest("GET", "https://quake.360.net/api/v3/user/info", nil)
    req.Header.Set("Host", "quake.360.net")
    req.Header.Set("X-QuakeToken", key)
    resp, err := client.Do(req)
    if err != nil {
        return "", 0, fmt.Errorf("发送请求时出错: %v", err)
    }
    defer resp.Body.Close()
    body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        return "", 0, fmt.Errorf("读取响应体时出错: %v", err)
    }
    var userInfo UserInfo
    err = json.Unmarshal(body, &userInfo)
    if err != nil {
        return "", 0, fmt.Errorf("解析JSON时出错: %v", err)
    }
    return userInfo.Data.User.FullName, userInfo.Data.MonthRemainingCredit, nil
}

func main() {
	// 解析yaml文件
	aaa := flag.Bool("a", false, "")
	bbb := flag.Bool("b", false, "")
	ccc := flag.Bool("c", false, "")
	ddd := flag.Bool("d", false, "")
	flag.Parse()
	exePath, _ := os.Executable()
	exeDir := filepath.Dir(exePath)
	configFilePath := filepath.Join(exeDir, "quake.yaml")
	configFile, err := ioutil.ReadFile(configFilePath)
	if err != nil {
		fmt.Println("\n\033[31m[-]\033[0m 没找到quake.yaml文件,请放在qkue.exe的目录下.")
		return
	}
	var config Config
	err = yaml.Unmarshal(configFile, &config)
	if err != nil {
		fmt.Println("\033[31m[-]\033[0m 无法解析yaml:", err)
		return
	}
	name,credit, _:= GetUserInfo(config.APIKey)
	logo := "\033[32m" + `	  _              ______ 
         | |            |  ____|
   __ _  | | __  _   _  | |__   
  / _' | | |/ / | | | | |  __|  
 | (_| | |   <  | |_| | | |____ 
  \__, | |_|\_\  \__,_| |______|
     | |                        
     |_|               
	   ` + "\033[0m" + `    (qkuE)::url导出 ` + "\033[32m" + name + ":" + strconv.Itoa(credit) + "\033[0m" + `

-a(最新数据) | -b(过滤无效请求) | -c(排除蜜罐) | -d(排除CDN)
`
	fmt.Print(logo)
	// 发送请求并获取总数
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("\n\033[36m[!]\u001B[0m 请输入查询语句:\n")
	Sentence, _ := reader.ReadString('\n')
	Sentence = strings.TrimSpace(Sentence)

	response, err := sendRequest(config.APIKey, 1, Sentence, *aaa, *bbb, *ccc, *ddd)
	if err != nil {
		fmt.Println("\033[31m[-]\033[0m ", err)
		return
	}
	var responseData Response
	err = json.Unmarshal([]byte(response), &responseData)
	if err != nil {
		fmt.Println("\n\033[31m[-]\033[0m 无法解析响应,请更换apiKey重试.\n")
		return
	}

	total := responseData.Meta.Pagination.Total
	fmt.Printf("\n\033[32m[+]\u001B[0m 查询结果的总条数为: %d\n", total)
	fmt.Print("\n\033[36m[!]\u001B[0m 请输入导出的条数(上限1w): ")
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)
	sum, err := strconv.Atoi(input)
	if err != nil {
		sum = total // 默认值
	}
	// 发送请求并获取响应
	response, err = sendRequest(config.APIKey, sum, Sentence, *aaa, *bbb, *ccc, *ddd)
	if err != nil {
		fmt.Println("\033[31m[-]\033[0m ", err)
		return
	}
	urls, err := ExtractHTTPUrls(response)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Print("\n\033[32m[+]\u001B[0m 导出其中的URL:\n\n")
	file, _ := os.OpenFile("url_export.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	defer file.Close()
	for _, url := range urls {
		file.WriteString(url + "\n")
		fmt.Println(url)
	}
	fmt.Print("\n\033[32m[+]\u001B[0m 已保存到url_export.txt\n\n")
}
