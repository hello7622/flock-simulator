package api

import "github.com/your-username/flock-simulator/simulator"

// API 请求和响应类型

type CreateSimulationRequest struct {
    Birds     []simulator.Bird    `json:"birds"`
    Obstacles []simulator.Obstacle `json:"obstacles"`
    Config    *simulator.SimulationConfig `json:"config"`
}

type SimulationResponse struct {
    State *simulator.SimulationState `json:"state"`
    Step  int                       `json:"step"`
}

type StepRequest struct {
    Steps int `json:"steps"` // 执行的步数
}

type ErrorResponse struct {
    Error string `json:"error"`
}

type SuccessResponse struct {
    Message string `json:"message"`
}