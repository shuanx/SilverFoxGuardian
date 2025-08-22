package backend

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"net"
	"strings"

	"github.com/shirou/gopsutil/v3/process"
)

type ProcessCheckResult struct {
	ProcessName []string `json:"processName"`
	Describe    string   `json:"describe"`
	IsExist     string   `json:"isExist"`
	PID         string   `json:"pid"`
	Connections []string `json:"connections"`
}

type RemoteConfig struct {
	RemoteControl []struct {
		ProcessName []string `json:"processName"`
		Describe    string   `json:"describe"`
	} `json:"remote_control"`
}

//go:embed configs/remote.json
var embeddedConfig string // 嵌入的配置文件内容

// LoadRemoteConfig 从嵌入的文件中加载配置
func LoadRemoteConfig() (*RemoteConfig, error) {
	var config RemoteConfig

	// 使用 json.Unmarshal 将嵌入的 JSON 数据解析为结构体
	err := json.Unmarshal([]byte(embeddedConfig), &config)
	if err != nil {
		return nil, fmt.Errorf("failed to decode remote config: %w", err)
	}

	return &config, nil
}

// isPrivateIP 判断是否为内网 IP
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

// filterLocalConnections 过滤掉本地连接、内网 IP、黑名单 IP 和 ":0:LISTEN" 格式的连接
func filterLocalConnections(conns []string) []string {
	var filteredConns []string

	// 定义黑名单 IP
	blacklistedIPs := map[string]struct{}{
		"127.0.0.1": {}, // 本地回环地址
		"0.0.0.0":   {}, // 未绑定地址
		"1.0.0.1":   {}, // 需要排除的地址
	}

	for _, conn := range conns {
		// 提取 IP 地址部分
		ip := strings.Split(conn, ":")[0]

		// 如果 IP 在黑名单中，跳过
		if _, isBlacklisted := blacklistedIPs[ip]; isBlacklisted {
			continue
		}

		// 如果是 ":0:LISTEN" 格式，跳过
		if strings.Contains(conn, ":0:LISTEN") {
			continue
		}

		// 如果是内网 IP，标记为 "内网 IP"
		if isPrivateIP(ip) {
			filteredConns = append(filteredConns, fmt.Sprintf("%s [内网 IP]", conn))
			continue
		}

		// 如果不是内网 IP，直接保留
		filteredConns = append(filteredConns, conn)
	}
	return filteredConns
}

// filterExternalIPs 从 IP 列表中筛选出外网 IP
func filterExternalIPs(ipList []string) []string {
	var externalIPs []string
	for _, ip := range ipList {
		if !isPrivateIP(ip) { // 如果不是内网 IP，则是外网 IP
			externalIPs = append(externalIPs, ip)
		}
	}
	return externalIPs
}

func CheckProcesses() ([]ProcessCheckResult, error) {
	config, err := LoadRemoteConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to load remote config: %w", err)
	}

	var results []ProcessCheckResult

	processes, err := process.Processes()
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve processes: %w", err)
	}

	for _, item := range config.RemoteControl {
		var pids []string
		var connsSet = make(map[string]struct{}) // 用于去重
		var ipSet = make(map[string]struct{})    // 收集所有外联 IP

		for _, procName := range item.ProcessName {
			lowerProcName := strings.ToLower(procName)
			for _, proc := range processes {
				currentProcName, _ := proc.Name()
				if strings.ToLower(currentProcName) == lowerProcName {
					// 收集 PID
					pids = append(pids, fmt.Sprintf("%d", proc.Pid))

					// 获取外联 IP
					conns, err := getProcessConnections(proc)
					if err != nil {
						continue
					}

					// 过滤掉本地连接、内网 IP 和特定格式的连接
					conns = filterLocalConnections(conns)

					// 去重处理
					for _, conn := range conns {
						connsSet[conn] = struct{}{}

						// 提取 IP 地址部分
						ip := strings.Split(conn, ":")[0]
						ipSet[ip] = struct{}{}
					}
				}
			}
		}

		// 将所有 IP 转换为切片
		var ipList []string
		for ip := range ipSet {
			ipList = append(ipList, ip)
		}

		// 筛选出外网 IP
		externalIPs := filterExternalIPs(ipList)

		// 如果存在外网 IP，则调用 BatchIPLocation 进行探测
		var ipLocations map[string]IPLocation
		if len(externalIPs) > 0 {
			ipLocations, err = BatchIPLocation(externalIPs)
			if err != nil {
				fmt.Printf("Failed to get batch IP locations: %v\n", err)
			}
		} else {
			ipLocations = make(map[string]IPLocation) // 外网 IP 为空时，返回空结果
		}

		// 将去重后的连接结果转换为切片，并补充归属地信息
		var uniqueConns []string
		for conn := range connsSet {
			ip := strings.Split(conn, ":")[0]
			location, found := ipLocations[ip]
			if found {
				// 格式化连接信息，加入归属地
				connWithLocation := fmt.Sprintf("%s [%s-%s]", conn, location.Country, location.City)
				uniqueConns = append(uniqueConns, connWithLocation)
			} else {
				uniqueConns = append(uniqueConns, conn) // 如果未找到归属地信息，保留原始信息
			}
		}

		// 结果汇总
		results = append(results, ProcessCheckResult{
			ProcessName: item.ProcessName,
			Describe:    item.Describe,
			IsExist:     ifExist(len(pids) > 0),
			PID:         strings.Join(pids, ", "), // 汇总所有 PIDs
			Connections: uniqueConns,              // 汇总所有外联 IP
		})
	}

	fmt.Println(results)
	return results, nil
}

func getProcessConnections(proc *process.Process) ([]string, error) {
	conns, err := proc.Connections()
	if err != nil {
		return nil, err
	}

	var connInfos []string
	for _, conn := range conns {
		if conn.Type == 1 { // 仅处理 TCP 类型的连接
			connStr := fmt.Sprintf("%s:%d:%s", conn.Raddr.IP, conn.Raddr.Port, conn.Status)
			connInfos = append(connInfos, connStr)
		}
	}

	return connInfos, nil
}

func ifExist(exist bool) string {
	if exist {
		return "True"
	}
	return "No"
}
