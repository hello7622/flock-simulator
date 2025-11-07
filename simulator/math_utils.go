package simulator

import "math"

const (
    TwoPi = 2 * math.Pi
    RadToDeg = 180.0 / math.Pi
    DegToRad = math.Pi / 180.0
)

// NormalizeAngle 将角度标准化到 [0, 360) 范围
func NormalizeAngle(angle float64) float64 {
    angle = math.Mod(angle, 360.0)
    if angle < 0 {
        angle += 360.0
    }
    return angle
}

// AngleDifference 计算两个角度之间的最小差值，范围 [-180, 180]
func AngleDifference(a, b float64) float64 {
    diff := a - b
    for diff > 180 {
        diff -= 360
    }
    for diff < -180 {
        diff += 360
    }
    return diff
}

// Distance 计算两点之间的距离
func Distance(a, b Point) float64 {
    dx := a.X - b.X
    dy := a.Y - b.Y
    return math.Sqrt(dx*dx + dy*dy)
}

// IsInFOV 判断目标点是否在鸟的扇形视野内
func IsInFOV(bird *Bird, target Point) bool {
    // 计算距离
    dist := Distance(bird.Position, target)
    if dist > bird.DetectionRadius {
        return false
    }

    // 计算目标相对于鸟的角度
    dx := target.X - bird.Position.X
    dy := target.Y - bird.Position.Y
    targetAngle := math.Atan2(dy, dx) * RadToDeg
    targetAngle = NormalizeAngle(targetAngle)

    // 计算角度差异
    angleDiff := AngleDifference(targetAngle, bird.Velocity.Angle)

    return math.Abs(angleDiff) <= bird.FOVAngle
}

// CompositeVelocities 合成两个速度向量
func CompositeVelocities(v1, v2 Velocity) Velocity {
    // 将极坐标转换为直角坐标
    angle1Rad := v1.Angle * DegToRad
    angle2Rad := v2.Angle * DegToRad

    x1 := v1.Speed * math.Cos(angle1Rad)
    y1 := v1.Speed * math.Sin(angle1Rad)
    x2 := v2.Speed * math.Cos(angle2Rad)
    y2 := v2.Speed * math.Sin(angle2Rad)

    // 合成向量
    resultX := x1 + x2
    resultY := y1 + y2

    // 转换回极坐标
    resultSpeed := math.Sqrt(resultX*resultX + resultY*resultY)
    
    var resultAngle float64
    if resultSpeed < 1e-10 {
        resultAngle = 0
    } else {
        resultAngle = math.Atan2(resultY, resultX) * RadToDeg
        resultAngle = NormalizeAngle(resultAngle)
    }

    return Velocity{
        Speed: resultSpeed,
        Angle: resultAngle,
    }
}

// LimitVelocity 限制速度大小和转向角度
func LimitVelocity(current, desired Velocity, maxSpeed, maxTurnAngle float64) Velocity {
    // 限制速度大小
    if desired.Speed > maxSpeed {
        desired.Speed = maxSpeed
    }

    // 限制转向角度
    angleDiff := AngleDifference(desired.Angle, current.Angle)
    if math.Abs(angleDiff) > maxTurnAngle {
        if angleDiff > 0 {
            desired.Angle = NormalizeAngle(current.Angle + maxTurnAngle)
        } else {
            desired.Angle = NormalizeAngle(current.Angle - maxTurnAngle)
        }
    }

    return desired
}