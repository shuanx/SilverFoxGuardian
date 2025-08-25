package backend

import (
	"testing"
)

// 测试 sendRequest 函数
func TestGetAllExternalConnections(t *testing.T) {
	// 替换为实际的 API 请求参数
	cookie := "_uetvid=b7e70890466011f0ae4c816546390a48; Hm_lvt_0e8161ac4f393ecc79c975079cc5fcc8=1749604166; UEDC_LOGIN_LANGUAGE=cn; oldPath=https://ti.sangfor.com.cn/sandbox-dashboard; vxLoginOldPath=https://ti.sangfor.com.cn/sandbox-dashboard; csrftoken=a7oFB1UjLQODjyIWH3v3VHA4LguSDHUmIfXDVm8JN6PCztehmMng3Wo0Tp1qGiZ0; sessionid=3l666pd9vzxtiudug6ofh51om9fkgpv1; loginSuccess=true"
	xCsrftoken := "a7oFB1UjLQODjyIWH3v3VHA4LguSDHUmIfXDVm8JN6PCztehmMng3Wo0Tp1qGiZ0"

	// 调用函数获取返回的 JSON 数据 channel
	_, err := AnalyzeExternalConnectionsStream(cookie, xCsrftoken)
	if err != nil {
		t.Fatalf("Error: %v", err)
	}
}
