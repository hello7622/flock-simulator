package main

import (
	"log"
	"net/http"

	"github.com/hello7622/flock-simulator/api"
)

func main() {
	// 设置路由
	http.HandleFunc("/api/simulation/create", api.CreateSimulationHandler)
	http.HandleFunc("/api/simulation/step", api.StepSimulationHandler)
	http.HandleFunc("/api/simulation/state", api.GetStateHandler)
	http.HandleFunc("/api/simulation/reset", api.ResetSimulationHandler)

	// 静态文件服务（用于前端）
	http.Handle("/", http.FileServer(http.Dir("./static")))

	log.Println("Flock Simulator API server starting on :8080")
	log.Println("Available endpoints:")
	log.Println("  POST   /api/simulation/create - Create new simulation")
	log.Println("  POST   /api/simulation/step   - Advance simulation")
	log.Println("  GET    /api/simulation/state  - Get current state")
	log.Println("  POST   /api/simulation/reset  - Reset simulation")

	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal("Server failed to start: ", err)
	}
}
