package main

import (
	"SilverFoxGuardian/backend"
	"context"
	"embed"
	"encoding/json"
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	webview "github.com/webview/webview_go"
)

//go:embed frontend/dist/*
var embeddedFiles embed.FS

//go:embed frontend/statics/icon.ico
var iconData []byte

func main() {
	// 初始化日志
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)

	// 创建一个 WebView 窗口
	debug := true
	w := webview.New(debug)
	defer w.Destroy()

	w.SetTitle("SilverFox Guardian")
	w.SetSize(900, 700, webview.HintNone)

	// 创建一个 HTTP 服务器实例
	server := &http.Server{
		Addr:         "localhost:8080",
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  30 * time.Second,
	}

	// 提取嵌入式文件系统的子目录
	fsys, err := fs.Sub(embeddedFiles, "frontend/dist")
	if err != nil {
		log.Fatalf("无法加载嵌入式文件系统: %s", err)
	}

	// 提供嵌入式静态文件服务
	http.Handle("/", http.FileServer(http.FS(fsys)))

	// 启动 HTTP 服务器
	go func() {
		log.Println("静态文件服务器启动在 http://localhost:8080")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("静态文件服务器启动失败: %s", err)
		}
	}()

	// 绑定 WebView 中的 JavaScript 函数到 Go 函数
	w.Bind("checkProcesses", func() string {
		resultChan := make(chan string, 1)

		// 异步调用后端逻辑
		go func() {
			results, err := backend.CheckProcesses()
			if err != nil {
				resultChan <- fmt.Sprintf("Error: %s", err)
				return
			}

			jsonData, err := json.MarshalIndent(results, "", "  ")
			if err != nil {
				resultChan <- fmt.Sprintf("Error: %s", err)
				return
			}

			resultChan <- string(jsonData)
		}()

		// 等待结果返回
		return <-resultChan
	})

	// 绑定 WebView 中的 JavaScript 函数到 Go 函数
	w.Bind("analyzeExternalConnections", func() string {
		resultChan := make(chan string, 1)

		// 异步调用后端逻辑
		go func() {
			cookie := "cookie"
			xCsrftoken := "csrftoken"

			// 调用函数获取返回的 JSON 数据 channel
			results, err := backend.AnalyzeExternalConnectionsStream(cookie, xCsrftoken)
			if err != nil {
				resultChan <- fmt.Sprintf("Error: %s", err)
				return
			}

			jsonData, err := json.MarshalIndent(results, "", "  ")
			if err != nil {
				resultChan <- fmt.Sprintf("Error: %s", err)
				return
			}

			resultChan <- string(jsonData)
		}()

		// 等待结果返回
		return <-resultChan
	})

	// 加载嵌入式 HTML 文件
	w.Navigate("http://localhost:8080/index.html")

	// 信号监听器，用于优雅关闭服务器
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	go func() {
		<-ctx.Done()
		log.Println("正在关闭服务器...")

		// 优雅关闭 HTTP 服务器
		if err := server.Shutdown(context.Background()); err != nil {
			log.Printf("关闭服务器时发生错误: %s", err)
		} else {
			log.Println("服务器已成功关闭")
		}
	}()

	// 启动 WebView
	w.Run()
}
