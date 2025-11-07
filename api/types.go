package api

import "flock-simulator/simulator"

// API请求和响应类型

type AddBirdRequest struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}

type AddObstacleRequest struct {
	X      float64 `json:"x"`
	Y      float64 `json:"y"`
	Radius float64 `json:"radius"`
}

type SetAttractorRequest struct {
	X      float64 `json:"x"`
	Y      float64 `json:"y"`
	Active bool    `json:"active"`
}

type SimulationResponse struct {
	State *simulator.SimulationState `json:"state"`
}

type SuccessResponse struct {
	Message string `json:"message"`
	Success bool   `json:"success"`
}
