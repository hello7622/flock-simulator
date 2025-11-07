package simulator

import "time"

// FlockSimulator 鸟群模拟器
type FlockSimulator struct {
    state    *SimulationState
    config   *SimulationConfig
    decisionMaker *DecisionMaker
}

// NewFlockSimulator 创建新的鸟群模拟器
func NewFlockSimulator(initialState *SimulationState, config *SimulationConfig) *FlockSimulator {
    return &FlockSimulator{
        state:    initialState,
        config:   config,
        decisionMaker: NewDecisionMaker(config),
    }
}

// Step 执行单步模拟
func (fs *FlockSimulator) Step() *SimulationState {
    decisions := make([]*DecisionResult, 0, len(fs.state.Birds))
    
    // 第一阶段：所有鸟做出决策
    for i := range fs.state.Birds {
        bird := &fs.state.Birds[i]
        
        // 创建其他鸟的列表（排除自己）
        otherBirds := make([]Bird, 0, len(fs.state.Birds)-1)
        for j := range fs.state.Birds {
            if i != j {
                otherBirds = append(otherBirds, fs.state.Birds[j])
            }
        }
        
        decision := fs.decisionMaker.MakeDecision(bird, fs.state.Obstacles, otherBirds)
        decisions = append(decisions, decision)
    }
    
    // 第二阶段：更新所有鸟的状态
    for i, decision := range decisions {
        fs.decisionMaker.UpdateBird(&fs.state.Birds[i], decision)
    }
    
    // 更新模拟状态
    fs.state.Step++
    fs.state.Timestamp = time.Now().UnixMilli()
    
    return fs.state
}

// GetState 获取当前模拟状态
func (fs *FlockSimulator) GetState() *SimulationState {
    return fs.state
}

// Reset 重置模拟器到初始状态
func (fs *FlockSimulator) Reset(initialState *SimulationState) {
    fs.state = initialState
}

// CreateDefaultSimulation 创建默认模拟配置
func CreateDefaultSimulation() (*SimulationState, *SimulationConfig) {
    config := &SimulationConfig{
        TimeStep:       1.0,
        PerfectDistance: 8.0,
        MaxSpeed:       5.0,
        MaxTurnAngle:   30.0,
    }
    
    state := &SimulationState{
        Birds: []Bird{
            {
                ID:             "bird_1",
                Position:       Point{X: 0, Y: 0},
                Velocity:       Velocity{Speed: 2.0, Angle: 0},
                DetectionRadius: 20.0,
                DefaultSpeed:   2.0,
                FOVAngle:       45.0,
            },
            {
                ID:             "bird_2",
                Position:       Point{X: 10, Y: 10},
                Velocity:       Velocity{Speed: 2.0, Angle: 45},
                DetectionRadius: 20.0,
                DefaultSpeed:   2.0,
                FOVAngle:       45.0,
            },
            {
                ID:             "bird_3",
                Position:       Point{X: -10, Y: 5},
                Velocity:       Velocity{Speed: 2.0, Angle: 90},
                DetectionRadius: 20.0,
                DefaultSpeed:   2.0,
                FOVAngle:       45.0,
            },
        },
        Obstacles: []Obstacle{
            {
                ID:       "obstacle_1",
                Position: Point{X: 15, Y: 15},
                Radius:   3.0,
            },
            {
                ID:       "obstacle_2",
                Position: Point{X: 25, Y: 5},
                Radius:   2.0,
            },
        },
        Step:      0,
        Timestamp: time.Now().UnixMilli(),
    }
    
    return state, config
}