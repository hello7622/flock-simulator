package simulator

import (
	"math"
)

// BoidsSimulator 实现Boids算法
type BoidsSimulator struct {
	config *SimulationConfig
}

// NewBoidsSimulator 创建新的Boids模拟器
func NewBoidsSimulator(config *SimulationConfig) *BoidsSimulator {
	return &BoidsSimulator{
		config: config,
	}
}

// Distance 计算两点之间的距离
func (bs *BoidsSimulator) Distance(a, b Point) float64 {
	dx := a.X - b.X
	dy := a.Y - b.Y
	return math.Sqrt(dx*dx + dy*dy)
}

// Limit 限制向量大小
func (bs *BoidsSimulator) Limit(v Velocity, max float64) Velocity {
	magnitude := math.Sqrt(v.DX*v.DX + v.DY*v.DY)
	if magnitude > max {
		scale := max / magnitude
		v.DX *= scale
		v.DY *= scale
	}
	return v
}

// Separation 分离规则：避免与相邻个体过于拥挤
func (bs *BoidsSimulator) Separation(bird *Bird, birds []Bird) Velocity {
	var steer Velocity
	count := 0

	for _, other := range birds {
		if bird.ID == other.ID {
			continue
		}

		distance := bs.Distance(bird.Position, other.Position)
		if distance > 0 && distance < bs.config.SeparationDistance {
			// 计算远离其他鸟的向量
			diff := Velocity{
				DX: bird.Position.X - other.Position.X,
				DY: bird.Position.Y - other.Position.Y,
			}
			// 距离越近，排斥力越大
			diff.DX /= distance
			diff.DY /= distance
			steer.DX += diff.DX
			steer.DY += diff.DY
			count++
		}
	}

	if count > 0 {
		steer.DX /= float64(count)
		steer.DY /= float64(count)
		steer = bs.Limit(steer, bs.config.MaxForce)
	}

	return steer
}

// Alignment 对齐规则：与相邻个体保持相同方向
func (bs *BoidsSimulator) Alignment(bird *Bird, birds []Bird) Velocity {
	var avg Velocity
	count := 0

	for _, other := range birds {
		if bird.ID == other.ID {
			continue
		}

		distance := bs.Distance(bird.Position, other.Position)
		if distance > 0 && distance < bs.config.AlignmentDistance {
			avg.DX += other.Velocity.DX
			avg.DY += other.Velocity.DY
			count++
		}
	}

	if count > 0 {
		avg.DX /= float64(count)
		avg.DY /= float64(count)
		avg = bs.Limit(avg, bs.config.MaxForce)
	}

	return avg
}

// Cohesion 凝聚规则：向相邻个体的平均位置移动
func (bs *BoidsSimulator) Cohesion(bird *Bird, birds []Bird) Velocity {
	var center Point
	count := 0

	for _, other := range birds {
		if bird.ID == other.ID {
			continue
		}

		distance := bs.Distance(bird.Position, other.Position)
		if distance > 0 && distance < bs.config.CohesionDistance {
			center.X += other.Position.X
			center.Y += other.Position.Y
			count++
		}
	}

	if count > 0 {
		center.X /= float64(count)
		center.Y /= float64(count)

		// 计算指向中心的向量
		steer := Velocity{
			DX: center.X - bird.Position.X,
			DY: center.Y - bird.Position.Y,
		}
		steer = bs.Limit(steer, bs.config.MaxForce)
		return steer
	}

	return Velocity{}
}

// AvoidObstacles 躲避障碍物
func (bs *BoidsSimulator) AvoidObstacles(bird *Bird, obstacles []Obstacle) Velocity {
	var steer Velocity

	for _, obstacle := range obstacles {
		distance := bs.Distance(bird.Position, obstacle.Position)
		avoidDistance := obstacle.Radius + bird.Radius + 10 // 安全距离

		if distance < avoidDistance {
			// 计算远离障碍物的向量
			away := Velocity{
				DX: bird.Position.X - obstacle.Position.X,
				DY: bird.Position.Y - obstacle.Position.Y,
			}
			// 距离越近，躲避力越大
			force := 1.0 - (distance / avoidDistance)
			away.DX *= force
			away.DY *= force
			steer.DX += away.DX
			steer.DY += away.DY
		}
	}

	if steer.DX != 0 || steer.DY != 0 {
		steer = bs.Limit(steer, bs.config.MaxForce)
	}

	return steer
}

// AttractToPoint 被引导点吸引
func (bs *BoidsSimulator) AttractToPoint(bird *Bird, attractor Attractor) Velocity {
	if !attractor.Active {
		return Velocity{}
	}

	// 计算指向引导点的向量
	steer := Velocity{
		DX: attractor.Position.X - bird.Position.X,
		DY: attractor.Position.Y - bird.Position.Y,
	}

	// 修复：增强吸引力，让鸟更容易被引导
	distance := bs.Distance(bird.Position, attractor.Position)
	if distance > 0 {
		// 距离越远，吸引力越强（但有限制）
		force := math.Min(1.0, 100.0/distance) // 增加基础吸引力
		steer.DX *= force
		steer.DY *= force
	}

	// 限制吸引力大小
	steer = bs.Limit(steer, bs.config.MaxForce*2) // 增加最大吸引力

	return steer
}

// UpdateBird 更新鸟的状态
func (bs *BoidsSimulator) UpdateBird(bird *Bird, birds []Bird, obstacles []Obstacle, attractor Attractor) {
	// 应用所有规则
	separation := bs.Separation(bird, birds)
	alignment := bs.Alignment(bird, birds)
	cohesion := bs.Cohesion(bird, birds)
	avoidance := bs.AvoidObstacles(bird, obstacles)
	attraction := bs.AttractToPoint(bird, attractor)

	// 加权合成
	separation.DX *= bs.config.SeparationWeight
	separation.DY *= bs.config.SeparationWeight
	alignment.DX *= bs.config.AlignmentWeight
	alignment.DY *= bs.config.AlignmentWeight
	cohesion.DX *= bs.config.CohesionWeight
	cohesion.DY *= bs.config.CohesionWeight
	avoidance.DX *= bs.config.AvoidanceWeight
	avoidance.DY *= bs.config.AvoidanceWeight
	attraction.DX *= bs.config.AttractionWeight
	attraction.DY *= bs.config.AttractionWeight

	// 更新速度
	bird.Velocity.DX += separation.DX + alignment.DX + cohesion.DX + avoidance.DX + attraction.DX
	bird.Velocity.DY += separation.DY + alignment.DY + cohesion.DY + avoidance.DY + attraction.DY

	// 限制最大速度
	bird.Velocity = bs.Limit(bird.Velocity, bs.config.MaxSpeed)

	// 更新位置
	bird.Position.X += bird.Velocity.DX
	bird.Position.Y += bird.Velocity.DY

	// 修复：边界环绕（使用画布实际尺寸）
	if bird.Position.X < 0 {
		bird.Position.X = 1200 // 与前端画布宽度一致
	} else if bird.Position.X > 1200 {
		bird.Position.X = 0
	}
	if bird.Position.Y < 0 {
		bird.Position.Y = 800 // 与前端画布高度一致
	} else if bird.Position.Y > 800 {
		bird.Position.Y = 0
	}
}

// CheckCollisions 检查碰撞并移除死亡的鸟
func (bs *BoidsSimulator) CheckCollisions(birds []Bird, obstacles []Obstacle) []Bird {
	var aliveBirds []Bird

	for _, bird := range birds {
		collided := false
		for _, obstacle := range obstacles {
			distance := bs.Distance(bird.Position, obstacle.Position)
			if distance < obstacle.Radius+bird.Radius {
				collided = true
				break
			}
		}
		if !collided {
			aliveBirds = append(aliveBirds, bird)
		}
	}

	return aliveBirds
}
