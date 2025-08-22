package main

import (
	"SilverFoxGuardian/backend"
	"context"
	"embed"
	"encoding/json"
	"fmt"
	"io/fs"
	"net/http"
	"os"
	"os/signal"

	webview "github.com/webview/webview_go"
)

//go:embed frontend/dist/*
var embeddedFiles embed.FS

func main() {
	// 创建一个 Webview 窗口
	debug := true
	w := webview.New(debug)
	defer w.Destroy()

	w.SetTitle("SilverFox Guardian")
	w.SetSize(800, 600, webview.HintNone)

	// 创建一个 HTTP 服务器实例
	server := &http.Server{Addr: "localhost:8080"}

	// 提供嵌入式静态文件服务
	go func() {
		// 提取嵌入的文件系统（子目录 frontend/dist）
		fsys, err := fs.Sub(embeddedFiles, "frontend/dist")
		if err != nil {
			fmt.Println("无法加载嵌入的文件系统:", err)
			return
		}

		// 使用嵌入的文件系统提供静态文件服务
		http.Handle("/", http.FileServer(http.FS(fsys)))

		fmt.Println("静态文件服务器启动在 http://localhost:8080")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fmt.Println("静态文件服务器启动失败:", err)
		}
	}()

	// 绑定 WebView 中的 JavaScript 函数到 Go 函数
	w.Bind("checkProcesses", func() string {
		fmt.Println("用户点击了 'Check Processes' 按钮")
		// 调用后端逻辑
		results, err := backend.CheckProcesses()
		if err != nil {
			return fmt.Sprintf("Error: %s", err)
		}

		// 将结果转换为 JSON 格式
		jsonData, err := json.MarshalIndent(results, "", "  ")
		if err != nil {
			return fmt.Sprintf("Error: %s", err)
		}

		// 打印输出到控制台（用于调试）
		fmt.Println("检查结果:", string(jsonData))

		return string(jsonData)
	})

	// 加载嵌入式 HTML 文件
	w.Navigate("http://localhost:8080/index.html")

	// 信号监听器，用于优雅关闭服务器
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	// 启动一个 goroutine 监听 WebView 的关闭事件
	go func() {
		<-ctx.Done() // 等待信号
		fmt.Println("正在关闭服务器...")
		// 优雅关闭 HTTP 服务器
		if err := server.Shutdown(context.Background()); err != nil {
			fmt.Println("关闭服务器时发生错误:", err)
		} else {
			fmt.Println("服务器已成功关闭")
		}
	}()

	// 启动 Webview
	w.Run()

	// 触发关闭信号，停止监听
	stop()
}
