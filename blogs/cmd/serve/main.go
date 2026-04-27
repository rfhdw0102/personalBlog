package main

import (
	"blogs/internal/app"
	"log"
)

func main() {
	application := app.NewApp()
	if err := application.Initialize(); err != nil {
		log.Fatalf("应用初始化失败: %v", err)
	}
	application.Run()
}
