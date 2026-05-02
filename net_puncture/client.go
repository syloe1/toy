package main

import (
	"bufio"
	"log"
	"net"
	"os"
)

func main() {
	// 连接服务器
	conn, err := net.Dial("tcp", "公网IP:9090")
	if err != nil {
		log.Fatal("连接失败：", err)
	}
	defer conn.Close()
	log.Println("✅ 已连接到中转服务器：公网IP:9090")
	// 读取服务器消息（后台协程）
	go func() {
		buf := make([]byte, 1024)
		for {
			n, err := conn.Read(buf)
			if err != nil {
				log.Println("服务器断开连接")
				return
			}
			log.Printf("📩 服务器返回：%s", string(buf[:n]))
		}
	}()

	// 读取控制台输入 → 发送给服务器
	scanner := bufio.NewScanner(os.Stdin)
	log.Println("请输入消息，回车发送：")
	for scanner.Scan() {
		msg := scanner.Text()
		_, _ = conn.Write([]byte(msg + "\n"))
	}
}
