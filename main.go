package main

import (
	"log"
	"net"
	"net/http"
	"os"
	"strconv"

	"flock-simulator/api"

	"github.com/gorilla/mux"
)

func main() {
	router := mux.NewRouter()

	// API路由
	router.HandleFunc("/api/state", api.GetStateHandler).Methods("GET")
	router.HandleFunc("/api/step", api.StepHandler).Methods("POST")
	router.HandleFunc("/api/bird", api.AddBirdHandler).Methods("POST")
	router.HandleFunc("/api/obstacle", api.AddObstacleHandler).Methods("POST")
	router.HandleFunc("/api/attractor", api.SetAttractorHandler).Methods("POST")
	router.HandleFunc("/api/toggle", api.ToggleRunningHandler).Methods("POST")
	router.HandleFunc("/api/reset", api.ResetHandler).Methods("POST")

	// 静态文件服务
	router.PathPrefix("/").Handler(http.FileServer(http.Dir("./static/")))

	port := getPort()
	log.Printf("Flock Simulator starting on :%s", port)
	log.Printf("Open http://localhost:%s in your browser", port)

	if err := http.ListenAndServe(":"+port, router); err != nil {
		log.Fatal("Server failed to start: ", err)
	}
}

func getPort() string {
	port := "8080"
	if envPort := os.Getenv("PORT"); envPort != "" {
		port = envPort
	}

	// 检查端口是否被占用
	for i := 8080; i <= 8090; i++ {
		testPort := strconv.Itoa(i)
		if isPortAvailable(testPort) {
			if testPort != port {
				log.Printf("Port %s is occupied, using port %s instead", port, testPort)
			}
			return testPort
		}
	}

	log.Printf("Port %s is occupied, but no alternative found. Continuing anyway...", port)
	return port
}

func isPortAvailable(port string) bool {
	conn, err := net.Listen("tcp", ":"+port)
	if err != nil {
		return false
	}
	conn.Close()
	return true
}
