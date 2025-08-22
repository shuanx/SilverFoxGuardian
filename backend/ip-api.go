package backend

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type IPLocation struct {
	Query      string  `json:"query"`
	Country    string  `json:"country"`
	RegionName string  `json:"regionName"`
	City       string  `json:"city"`
	Org        string  `json:"org"`
	Lat        float64 `json:"lat"`
	Lon        float64 `json:"lon"`
}

// BatchIPLocation 批量查询 IP 的归属地信息
func BatchIPLocation(ips []string) (map[string]IPLocation, error) {
	fmt.Println(ips)
	// 构造请求 URL
	apiURL := "http://ip-api.com/batch"

	// 将 IP 列表转换为 JSON
	requestBody, err := json.Marshal(ips)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal IP list: %w", err)
	}

	// 构建 HTTP 请求
	req, err := http.NewRequest("POST", apiURL, bytes.NewBuffer(requestBody))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	// 发起请求
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	// 检查响应状态码
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("batch query failed with status code: %d", resp.StatusCode)
	}

	// 解析响应
	var locations []IPLocation
	if err := json.NewDecoder(resp.Body).Decode(&locations); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	// 将结果转换为 map
	results := make(map[string]IPLocation)
	for _, loc := range locations {
		results[loc.Query] = loc
	}

	return results, nil
}
