package simulator

import (
	"math"
	"math/rand"
)

// AvoidObstacles 优先级1：障碍物避障
func AvoidObstacles(bird *Bird, obstacles []Obstacle) (Velocity, bool) {
	for _, obstacle := range obstacles {
		if IsInFOV(bird, obstacle.Position) {
			dist := Distance(bird.Position, obstacle.Position)
			safetyDistance := obstacle.Radius + 5.0 // 安全距离

			if dist < safetyDistance {
				// 计算避障方向
				rotateDirection := 1.0
				if bird.Velocity.Speed < 0 {
					rotateDirection = -1.0
				}

				avoidAngle := NormalizeAngle(bird.Velocity.Angle + 90.0*rotateDirection)

				avoidance := Velocity{
					Speed: 0, // 停止前进
					Angle: avoidAngle,
				}

				return avoidance, true
			}
		}
	}

	return Velocity{}, false
}

// InteractWithBirds 优先级2：与其他鸟的交互
func InteractWithBirds(bird *Bird, otherBirds []Bird, perfectDist float64) (Velocity, bool) {
	composite := Velocity{}
	count := 0

	for _, other := range otherBirds {
		if bird.ID == other.ID {
			continue // 跳过自己
		}

		if IsInFOV(bird, other.Position) {
			dist := Distance(bird.Position, other.Position)

			// 计算到其他鸟的角度
			dx := other.Position.X - bird.Position.X
			dy := other.Position.Y - bird.Position.Y
			angleToBird := math.Atan2(dy, dx) * RadToDeg
			angleToBird = NormalizeAngle(angleToBird)

			var interaction Velocity

			if dist < perfectDist {
				// 太近，排斥
				interaction.Speed = -1.0
				interaction.Angle = NormalizeAngle(angleToBird + 180)
			} else if dist > perfectDist {
				// 太远，吸引
				interaction.Speed = 1.0
				interaction.Angle = angleToBird
			} else {
				// 完美距离，无相互作用
				continue
			}

			// 检查是否同向
			angleDiff := AngleDifference(other.Velocity.Angle, bird.Velocity.Angle)
			if math.Abs(angleDiff) < 45.0 {
				interaction.Speed *= 2.0 // 同向时作用更强
			}

			composite = CompositeVelocities(composite, interaction)
			count++
		}
	}

	if count > 0 {
		composite.Speed /= float64(count) // 平均化
		return composite, true
	}

	return Velocity{}, false
}

// RandomMove 优先级3：随机移动
func RandomMove(bird *Bird) Velocity {
	// 小幅度随机旋转 [-1, 1] 度
	randomRotation := (float64(rand.Intn(2000)) / 1000.0) - 1.0

	return Velocity{
		Speed: bird.DefaultSpeed,
		Angle: NormalizeAngle(bird.Velocity.Angle + randomRotation),
	}
}
