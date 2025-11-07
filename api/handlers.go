package api

import (
    "encoding/json"
    "net/http"
    "strconv"
    "sync"

    "github.com/your-username/flock-simulator/simulator"
)

var (
    currentSimulator *simulator.FlockSimulator
    simulatorMutex   sync.RWMutex
)

// CreateSimulationHandler 创建新的模拟
func CreateSimulationHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        return
    }

    var req CreateSimulationRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, "Invalid request body", http.StatusBadRequest)
        return
    }

    state := &simulator.SimulationState{
        Birds:     req.Birds,
        Obstacles: req.Obstacles,
        Step:      0,
    }

    config := req.Config
    if config == nil {
        // 使用默认配置
        _, defaultConfig := simulator.CreateDefaultSimulation()
        config = defaultConfig
    }

    simulatorMutex.Lock()
    currentSimulator = simulator.NewFlockSimulator(state, config)
    simulatorMutex.Unlock()

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(SimulationResponse{
        State: state,
        Step:  0,
    })
}

// StepSimulationHandler 执行模拟步进
func StepSimulationHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        return
    }

    var req StepRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, "Invalid request body", http.StatusBadRequest)
        return
    }

    steps := req.Steps
    if steps <= 0 {
        steps = 1
    }

    simulatorMutex.Lock()
    if currentSimulator == nil {
        simulatorMutex.Unlock()
        http.Error(w, "No active simulation", http.StatusBadRequest)
        return
    }

    var finalState *simulator.SimulationState
    for i := 0; i < steps; i++ {
        finalState = currentSimulator.Step()
    }
    simulatorMutex.Unlock()

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(SimulationResponse{
        State: finalState,
        Step:  finalState.Step,
    })
}

// GetStateHandler 获取当前模拟状态
func GetStateHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodGet {
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        return
    }

    simulatorMutex.RLock()
    defer simulatorMutex.RUnlock()

    if currentSimulator == nil {
        http.Error(w, "No active simulation", http.StatusBadRequest)
        return
    }

    state := currentSimulator.GetState()

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(SimulationResponse{
        State: state,
        Step:  state.Step,
    })
}

// ResetSimulationHandler 重置模拟
func ResetSimulationHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        return
    }

    state, config := simulator.CreateDefaultSimulation()

    simulatorMutex.Lock()
    currentSimulator = simulator.NewFlockSimulator(state, config)
    simulatorMutex.Unlock()

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(SuccessResponse{
        Message: "Simulation reset to default state",
    })
}