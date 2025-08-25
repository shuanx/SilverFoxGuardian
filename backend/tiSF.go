package backend

import (
	"encoding/json"
	"fmt"
	"log"
	"os/exec"
	"regexp"
	"strings"
)

// 使用 curl 发送 HTTPS 请求测试 IP 是否恶意
func sendHttpsRequestWithCurl(batch []string, cookie string, xCsrftoken string) (map[string]interface{}, error) {
	// API 地址
	url := "https://ti.sangfor.com.cn:443/api/v1/threats/iocs"

	// 将批量数据转为 JSON 格式
	batchStr, err := json.Marshal(batch)
	if err != nil {
		return nil, fmt.Errorf("error marshaling batch: %v", err)
	}

	// 构建 curl 命令
	curlCommand := []string{
		"curl",
		"-s",         // 静默模式，避免进度条等信息干扰输出
		"-X", "POST", // 使用 POST 方法
		url,                                                     // API 地址
		"-H", "Content-Type: application/x-www-form-urlencoded", // 设置 Content-Type
		"-H", "Accept: application/json, text/plain, */*", // 设置 Accept 头
		"-H", fmt.Sprintf("User-Agent: %s", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/90.0.4430.212 Safari/537.36"), // 设置 User-Agent
		"-H", fmt.Sprintf("X-Csrftoken: %s", xCsrftoken), // 设置 CSRF Token
		"-H", fmt.Sprintf("Cookie: %s", cookie), // 设置 Cookie
		"--data-urlencode", fmt.Sprintf("batch=%s", string(batchStr)), // 设置请求体
		"-k", // 忽略 SSL 证书验证
	}

	// 打印调试信息
	log.Printf("Executing curl command: %s\n", strings.Join(curlCommand, " "))

	// 执行 curl 命令
	cmd := exec.Command(curlCommand[0], curlCommand[1:]...)

	// 捕获标准输出和标准错误
	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Printf("Error executing curl command: %v", err)
		log.Printf("Curl Output: %s\n", string(output))
		return nil, fmt.Errorf("执行 curl 命令失败: %v", err)
	}

	// 解析命令输出（日志记录原始输出内容）
	outputStr := string(output)
	log.Printf("Raw Curl Response: %s\n", outputStr)

	// 使用正则表达式提取合法的 JSON 数据（找到第一个 `{` 开头的 JSON）
	re := regexp.MustCompile(`\{.*\}`)
	jsonStr := re.FindString(outputStr)
	if jsonStr == "" {
		log.Printf("未找到合法的 JSON 数据，响应内容: %s", outputStr)
		return nil, fmt.Errorf("未找到合法的 JSON 数据")
	}

	// 检查和解析 JSON 响应
	var response map[string]interface{}
	err = json.Unmarshal([]byte(jsonStr), &response)
	if err != nil {
		log.Printf("JSON 解析失败: %v\n原始响应: %s", err, jsonStr)
		return nil, fmt.Errorf("JSON 解析失败: %v", err)
	}

	// 打印解析结果
	log.Printf("JSON 解析结果: %+v\n", response)
	return response, nil
}
