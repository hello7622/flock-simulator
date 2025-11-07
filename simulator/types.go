package simulator

// Point 表示二维空间中的点
type Point struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}

// Velocity 表示速度向量
type Velocity struct {
	DX float64 `json:"dx"`
	DY float64 `json:"dy"`
}

// Bird 表示单个飞鸟
type Bird struct {
	ID       string   `json:"id"`
	Position Point    `json:"position"`
	Velocity Velocity `json:"velocity"`
	Radius   float64  `json:"radius"`
}

// Obstacle 表示障碍物
type Obstacle struct {
	ID       string  `json:"id"`
	Position Point   `json:"position"`
	Radius   float64 `json:"radius"`
}

// Attractor 表示全局引导点
type Attractor struct {
	Position Point `json:"position"`
	Active   bool  `json:"active"`
}

// SimulationConfig 模拟配置
type SimulationConfig struct {
	SeparationDistance float64 `json:"separation_distance"`
	AlignmentDistance  float64 `json:"alignment_distance"`
	CohesionDistance   float64 `json:"cohesion_distance"`
	MaxSpeed           float64 `json:"max_speed"`
	MaxForce           float64 `json:"max_force"`
	SeparationWeight   float64 `json:"separation_weight"`
	AlignmentWeight    float64 `json:"alignment_weight"`
	CohesionWeight     float64 `json:"cohesion_weight"`
	AvoidanceWeight    float64 `json:"avoidance_weight"`
	AttractionWeight   float64 `json:"attraction_weight"`
}

// SimulationState 模拟状态
type SimulationState struct {
	Birds     []Bird     `json:"birds"`
	Obstacles []Obstacle `json:"obstacles"`
	Attractor Attractor  `json:"attractor"`
	Step      int        `json:"step"`
	Running   bool       `json:"running"`
}
