package simulator

import (
    "math/rand"
    "time"
)

// DecisionMaker 决策器
type DecisionMaker struct {
    config *SimulationConfig
}

// NewDecisionMaker 创建新的决策器
func NewDecisionMaker(config *SimulationConfig) *DecisionMaker {
    rand.Seed(time.Now().UnixNano())
    return &DecisionMaker{
        config: config,
    }
}

// MakeDecision 为主要决策函数
func (dm *DecisionMaker) MakeDecision(bird *Bird, obstacles []Obstacle, otherBirds []Bird) *DecisionResult {
    result := &DecisionResult{
        BirdID: bird.ID,
    }

    // 优先级1：障碍物避障
    if avoidance, ok := AvoidObstacles(bird, obstacles); ok {
        result.Velocity = avoidance
        result.Priority = 1
        result.Reason = "obstacle_avoidance"
        return result
    }

    // 优先级2：与其他鸟交互
    if interaction, ok := InteractWithBirds(bird, otherBirds, dm.config.PerfectDistance); ok {
        result.Velocity = interaction
        result.Priority = 2
        result.Reason = "bird_interaction"
        return result
    }

    // 优先级3：随机移动
    random := RandomMove(bird)
    result.Velocity = random
    result.Priority = 3
    result.Reason = "random_move"

    return result
}

// UpdateBird 更新鸟的状态
func (dm *DecisionMaker) UpdateBird(bird *Bird, decision *DecisionResult) {
    // 合成当前速度与决策向量
    newVelocity := CompositeVelocities(bird.Velocity, decision.Velocity)
    
    // 应用速度和转向限制
    newVelocity = LimitVelocity(bird.Velocity, newVelocity, dm.config.MaxSpeed, dm.config.MaxTurnAngle)
    
    bird.Velocity = newVelocity

    // 更新位置
    rad := bird.Velocity.Angle * DegToRad
    bird.Position.X += bird.Velocity.Speed * math.Cos(rad) * dm.config.TimeStep
    bird.Position.Y += bird.Velocity.Speed * math.Sin(rad) * dm.config.TimeStep
}