package main

import (
	"fmt"
	"github.com/shirou/gopsutil/v3/net"
	"net"
)

// 判断是否为内网 IP
func isPrivateIP(ip string) bool {
	parsedIP := net.ParseIP(ip)
	if parsedIP == nil {
		return false
	}

	// 检查是否属于内网 IP 范围
	if parsedIP.IsPrivate() || parsedIP.IsLoopback() || ip == "0.0.0.0" {
		return true
	}
	return false
}

// 获取所有外联 IP
func getAllExternalConnections() ([]string, error) {
	connections, err := net.Connections("inet") // 获取所有网络连接 (IPv4 和 IPv6)
	if err != nil {
		return nil, fmt.Errorf("failed to get network connections: %w", err)
	}

	externalIPs := make(map[string]struct{}) // 用于去重
	for _, conn := range connections {
		// 只处理 ESTABLISHED 状态的连接
		if conn.Status == "ESTABLISHED" {
			ip := conn.Raddr.IP // 获取远程地址的 IP
			if !isPrivateIP(ip) && ip != "" {
				externalIPs[ip] = struct{}{} // 只记录外网 IP
			}
		}
	}

	// 转换为切片
	var results []string
	for ip := range externalIPs {
		results = append(results, ip)
	}
	return results, nil
}

// 主函数
func main() {
	externalIPs, err := getAllExternalConnections()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	// 打印所有外联 IP
	fmt.Println("External Connections:")
	for _, ip := range externalIPs {
		fmt.Println(ip)
	}
}
