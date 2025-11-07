package api

import (
	"encoding/json"
	"net/http"
	"sync"

	"flock-simulator/simulator"
)

var (
	currentSimulator *simulator.FlockSimulator
	simulatorMutex   sync.RWMutex
)

// 初始化模拟器
func init() {
	simulatorMutex.Lock()
	defer simulatorMutex.Unlock()
	// 使用默认配置创建模拟器
	config := simulator.DefaultConfig()
	currentSimulator = simulator.NewFlockSimulator(config)
}

// GetStateHandler 获取当前状态
func GetStateHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	simulatorMutex.RLock()
	state := currentSimulator.GetState()
	simulatorMutex.RUnlock()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(SimulationResponse{State: state})
}

// StepHandler 执行单步模拟
func StepHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	simulatorMutex.Lock()
	state := currentSimulator.Step()
	simulatorMutex.Unlock()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(SimulationResponse{State: state})
}

// AddBirdHandler 添加鸟
func AddBirdHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req AddBirdRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	simulatorMutex.Lock()
	currentSimulator.AddBird(simulator.Point{X: req.X, Y: req.Y})
	state := currentSimulator.GetState()
	simulatorMutex.Unlock()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(SimulationResponse{State: state})
}

// AddObstacleHandler 添加障碍物
func AddObstacleHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req AddObstacleRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	simulatorMutex.Lock()
	currentSimulator.AddObstacle(simulator.Point{X: req.X, Y: req.Y}, req.Radius)
	state := currentSimulator.GetState()
	simulatorMutex.Unlock()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(SimulationResponse{State: state})
}

// SetAttractorHandler 设置引导点
func SetAttractorHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req SetAttractorRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	simulatorMutex.Lock()
	currentSimulator.SetAttractor(simulator.Point{X: req.X, Y: req.Y}, req.Active)
	state := currentSimulator.GetState()
	simulatorMutex.Unlock()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(SimulationResponse{State: state})
}

// ToggleRunningHandler 切换运行状态
func ToggleRunningHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	simulatorMutex.Lock()
	currentSimulator.ToggleRunning()
	state := currentSimulator.GetState()
	simulatorMutex.Unlock()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(SimulationResponse{State: state})
}

// ResetHandler 重置模拟
func ResetHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	simulatorMutex.Lock()
	currentSimulator.Reset()
	simulatorMutex.Unlock()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(SuccessResponse{
		Message: "Simulation reset",
		Success: true,
	})
}
