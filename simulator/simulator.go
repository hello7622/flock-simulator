package simulator

import (
	"math/rand"
	"strconv"
	"time"
)

// FlockSimulator 鸟群模拟器
type FlockSimulator struct {
	state  *SimulationState
	config *SimulationConfig
	boids  *BoidsSimulator
}

// NewFlockSimulator 创建新的鸟群模拟器
func NewFlockSimulator(config *SimulationConfig) *FlockSimulator {
	if config == nil {
		config = DefaultConfig()
	}

	// 初始化随机种子
	rand.Seed(time.Now().UnixNano())

	initialState := &SimulationState{
		Birds:     []Bird{},
		Obstacles: []Obstacle{},
		Attractor: Attractor{Active: false},
		Step:      0,
		Running:   true,
	}

	return &FlockSimulator{
		state:  initialState,
		config: config,
		boids:  NewBoidsSimulator(config),
	}
}

// DefaultConfig 返回默认配置
func DefaultConfig() *SimulationConfig {
	return &SimulationConfig{
		SeparationDistance: 25,
		AlignmentDistance:  50,
		CohesionDistance:   50,
		MaxSpeed:           4,
		MaxForce:           0.3, // 增加最大力
		SeparationWeight:   1.5,
		AlignmentWeight:    1.0,
		CohesionWeight:     1.0,
		AvoidanceWeight:    2.0,
		AttractionWeight:   0.7, // 显著增加引导点权重
	}
}

// Step 执行单步模拟
func (fs *FlockSimulator) Step() *SimulationState {
	if !fs.state.Running {
		return fs.state
	}

	// 更新所有鸟的状态
	for i := range fs.state.Birds {
		bird := &fs.state.Birds[i]
		fs.boids.UpdateBird(bird, fs.state.Birds, fs.state.Obstacles, fs.state.Attractor)
	}

	// 检查碰撞
	fs.state.Birds = fs.boids.CheckCollisions(fs.state.Birds, fs.state.Obstacles)

	fs.state.Step++
	return fs.state
}

// AddBird 添加鸟
func (fs *FlockSimulator) AddBird(position Point) {
	bird := Bird{
		ID:       generateID(),
		Position: position,
		Velocity: Velocity{
			DX: (rand.Float64() - 0.5) * 2,
			DY: (rand.Float64() - 0.5) * 2,
		},
		Radius: 3,
	}
	fs.state.Birds = append(fs.state.Birds, bird)
}

// AddObstacle 添加障碍物
func (fs *FlockSimulator) AddObstacle(position Point, radius float64) {
	obstacle := Obstacle{
		ID:       generateID(),
		Position: position,
		Radius:   radius,
	}
	fs.state.Obstacles = append(fs.state.Obstacles, obstacle)
}

// SetAttractor 设置引导点
func (fs *FlockSimulator) SetAttractor(position Point, active bool) {
	fs.state.Attractor = Attractor{
		Position: position,
		Active:   active,
	}
}

// ToggleRunning 切换运行状态
func (fs *FlockSimulator) ToggleRunning() {
	fs.state.Running = !fs.state.Running
}

// GetState 获取当前状态
func (fs *FlockSimulator) GetState() *SimulationState {
	return fs.state
}

// Reset 重置模拟器
func (fs *FlockSimulator) Reset() {
	fs.state = &SimulationState{
		Birds:     []Bird{},
		Obstacles: []Obstacle{},
		Attractor: Attractor{Active: false},
		Step:      0,
		Running:   true,
	}
}

func generateID() string {
	return strconv.FormatInt(rand.Int63(), 36)
}
