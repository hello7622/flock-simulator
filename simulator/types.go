package simulator

// Point 表示二维空间中的点
type Point struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}

// Velocity 表示速度向量（极坐标形式）
type Velocity struct {
	Speed float64 `json:"speed"` // 速度大小
	Angle float64 `json:"angle"` // 角度，范围 [0, 360) 度
}

// Bird 表示单个鸟或无人机
type Bird struct {
	ID              string   `json:"id"`
	Position        Point    `json:"position"`
	Velocity        Velocity `json:"velocity"`
	DetectionRadius float64  `json:"detection_radius"`
	DefaultSpeed    float64  `json:"default_speed"`
	FOVAngle        float64  `json:"fov_angle"` // 视野角度，默认 45 度
}

// Obstacle 表示障碍物
type Obstacle struct {
	ID       string  `json:"id"`
	Position Point   `json:"position"`
	Radius   float64 `json:"radius"`
}

// SimulationConfig 模拟配置
type SimulationConfig struct {
	TimeStep        float64 `json:"time_step"`        // 时间步长
	PerfectDistance float64 `json:"perfect_distance"` // 理想间距
	MaxSpeed        float64 `json:"max_speed"`        // 最大速度限制
	MaxTurnAngle    float64 `json:"max_turn_angle"`   // 最大转向角度
}

// SimulationState 模拟状态
type SimulationState struct {
	Birds     []Bird     `json:"birds"`
	Obstacles []Obstacle `json:"obstacles"`
	Step      int        `json:"step"`
	Timestamp int64      `json:"timestamp"`
}

// DecisionResult 决策结果
type DecisionResult struct {
	BirdID   string   `json:"bird_id"`
	Velocity Velocity `json:"velocity"`
	Priority int      `json:"priority"` // 触发的优先级
	Reason   string   `json:"reason"`   // 决策原因
}
