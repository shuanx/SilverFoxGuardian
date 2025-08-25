package backend

import (
	"fmt"
	"testing"
)

// 测试 sendRequest 函数
func TestSendRequest(t *testing.T) {
	// 测试数据
	batch := []string{"149.88.67.69", "149.88.67.68"} // 输入的 batch 参数
	cookie := "cookie"
	xCsrftoken := "csrftoken"

	// 调用函数
	response, err := sendHttpsRequestWithCurl(batch, cookie, xCsrftoken)
	if err != nil {
		fmt.Printf("Request failed: %v\n", err)
		return
	}

	// 打印最终响应
	fmt.Printf("Final Response: %+v\n", response)
}
