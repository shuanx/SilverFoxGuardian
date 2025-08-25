package backend

import (
	"encoding/json"
	"fmt"
	psnet "github.com/shirou/gopsutil/v3/net" // 获取网络连接信息
	"github.com/shirou/gopsutil/v3/process"   // 获取进程信息
	"net"
)

// 判断是否为内网 IP 或链路本地地址
func isPrivateOrLinkLocalIP(ip string) bool {
	parsedIP := net.ParseIP(ip)
	if parsedIP == nil {
		return false
	}

	// 检查是否属于内网 IP、回环地址或链路本地地址
	if parsedIP.IsPrivate() || parsedIP.IsLoopback() || ip == "0.0.0.0" {
		return true
	}

	// 链路本地地址 (IPv6 fe80::/10)
	if parsedIP.IsLinkLocalUnicast() {
		return true
	}

	return false
}

// 获取所有外联 IP 及其对应的进程、PID 和端口
func getAllExternalConnections() ([]map[string]string, error) {
	connections, err := psnet.Connections("inet") // 获取所有网络连接 (IPv4 和 IPv6)
	if err != nil {
		return nil, fmt.Errorf("failed to get network connections: %w", err)
	}

	results := []map[string]string{} // 用于存储最终结果
	seen := map[string]struct{}{}    // 用于去重
	for _, conn := range connections {
		// 只处理 ESTABLISHED 状态的连接
		ip := conn.Raddr.IP // 获取远程地址的 IP

		// 跳过内网 IP、链路本地地址和空 IP
		if isPrivateOrLinkLocalIP(ip) || ip == "" || ip == "::" {
			continue
		}

		// 根据 PID 获取进程信息
		pid := conn.Pid
		proc, err := process.NewProcess(pid)
		if err != nil {
			continue // 如果无法获取进程信息，跳过
		}

		name, err := proc.Name()
		if err != nil {
			name = "Unknown" // 如果无法获取进程名，设置为 "Unknown"
		}

		// 远程端口
		port := conn.Raddr.Port

		// 构建唯一键
		uniqueKey := fmt.Sprintf("%s:%d:%s:%d", name, pid, ip, port)
		if _, exists := seen[uniqueKey]; exists {
			continue // 如果已存在于去重集合中，跳过
		}

		// 标记为已处理
		seen[uniqueKey] = struct{}{}

		// 保存结果
		results = append(results, map[string]string{
			"ProcessName": name,
			"PID":         fmt.Sprintf("%d", pid),
			"IP":          ip,
			"Port":        fmt.Sprintf("%d", port),
			"status":      conn.Status,
		})
	}

	return results, nil
}

// 分块输出日志，避免日志被截断
func logChunked(data string, chunkSize int) {
	for i := 0; i < len(data); i += chunkSize {
		end := i + chunkSize
		if end > len(data) {
			end = len(data)
		}
		fmt.Println(data[i:end])
	}
}

// 执行外联 IP 的批量测试，并逐条返回 JSON 格式的结果
func AnalyzeExternalConnectionsStream(cookie string, xCsrftoken string) ([]map[string]interface{}, error) {
	// 获取外联 IP 和对应的进程信息
	connections, err := getAllExternalConnections()
	if err != nil {
		return nil, fmt.Errorf("failed to get external connections: %w", err)
	}

	// 提取所有 IP
	var ips []string
	for _, conn := range connections {
		ips = append(ips, conn["IP"])
	}

	// 打印待测试的外联 IP
	fmt.Printf("测试的外联 IP 为： %v\n", ips)

	// 调用批量测试函数，假设返回的数据结构类似于您提供的 `results`
	results, err := sendHttpsRequestWithCurl(ips, cookie, xCsrftoken)
	if err != nil {
		return nil, fmt.Errorf("failed to analyze IPs: %w", err)
	}

	// 定义最终返回的结果
	finalResults := []map[string]interface{}{}

	// 遍历 connections，匹配 results 中的 IP
	for _, conn := range connections {
		ip := conn["IP"]

		// 初始化 Tags、安全状态、location
		tags := []string{}
		safetyType := "未知" // 默认为未知
		location := ""     // 初始化 location

		// 在 results["data"] 中查找匹配的 IP
		if data, ok := results["data"].([]interface{}); ok {
			for _, item := range data {
				if itemMap, ok := item.(map[string]interface{}); ok {
					if resultIP, exists := itemMap["ioc"].(string); exists && resultIP == ip {
						// 如果匹配成功，提取安全状态
						if safety, exists := itemMap["safety_type"].(string); exists {
							switch safety {
							case "0":
								safetyType = "恶意"
							case "1":
								safetyType = "安全"
							case "2":
								safetyType = "未知"
							}
						}

						// 提取标签信息
						if labels, exists := itemMap["label"].([]interface{}); exists {
							for _, label := range labels {
								if labelMap, ok := label.(map[string]interface{}); ok {
									if tagValue, exists := labelMap["value"].(string); exists {
										tags = append(tags, tagValue)
									}
								}
							}
						}

						// 提取 location 信息
						if loc, exists := itemMap["location"].(map[string]interface{}); exists {
							country := ""
							province := ""
							if c, ok := loc["country"].(string); ok {
								country = c
							}
							if p, ok := loc["province"].(string); ok {
								province = p
							}
							// 拼接 country 和 province
							location = fmt.Sprintf("%s %s", country, province)
						}
						break // 匹配到后退出循环
					}
				}
			}
		}

		// 添加到最终结果
		finalResults = append(finalResults, map[string]interface{}{
			"process":  conn["ProcessName"],
			"pid":      conn["PID"],
			"remote":   conn["IP"] + ":" + conn["Port"],
			"tiResult": safetyType, // 映射后的安全类型
			"Tags":     tags,
			"geo":      location, // 新增 location 信息
			"state":    conn["status"],
		})
	}

	// 打印最终结果（分块输出，避免日志被截断）
	finalResultsJSON, _ := json.MarshalIndent(finalResults, "", "  ")
	logChunked(string(finalResultsJSON), 1000)

	return finalResults, nil
}
