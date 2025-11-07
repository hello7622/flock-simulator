# 鸟群模拟系统 - Go 后端 API

## 项目信息

- **技术栈**: Go 1.23 + 标准库 HTTP Server
- **架构**: 模块化设计，RESTful API

## 项目结构

```plaintext
flock-simulator/
├── go.mod
├── main.go                 # HTTP 服务器入口
├── simulator/              # 核心模拟引擎
│   ├── types.go           # 数据结构定义
│   ├── math_utils.go      # 数学计算工具
│   ├── behaviors.go       # 行为决策逻辑
│   ├── decision.go        # 决策器
│   └── simulator.go       # 模拟器核心
└── api/                   # HTTP API 层
    ├── handlers.go        # 请求处理器
    └── types.go           # API 数据类型
```

## 核心设计思路

### 1. 优先级行为系统

- **优先级1**: 障碍物避障 - 紧急安全行为
- **优先级2**: 鸟群交互 - 群体协调行为  
- **优先级3**: 随机移动 - 默认探索行为

### 2. 物理模型

- 基于极坐标的速度向量系统
- 扇形视野检测机制
- 向量合成物理运动

### 3. 状态管理

- 线程安全的模拟状态管理
- 时间步进式模拟推进
- 完整的重置和初始化支持

## API 接口文档

### 1. 创建模拟实例

**端点**: `POST /api/simulation/create`

**请求体**:

```json
{
  "birds": [
    {
      "id": "bird_1",
      "position": {"x": 0, "y": 0},
      "velocity": {"speed": 2.0, "angle": 0},
      "detection_radius": 20.0,
      "default_speed": 2.0,
      "fov_angle": 45.0
    }
  ],
  "obstacles": [
    {
      "id": "obstacle_1", 
      "position": {"x": 15, "y": 15},
      "radius": 3.0
    }
  ],
  "config": {
    "time_step": 1.0,
    "perfect_distance": 8.0,
    "max_speed": 5.0,
    "max_turn_angle": 30.0
  }
}
```

### 2. 执行模拟步进

**端点**: `POST /api/simulation/step`

**请求体**:

```json
{
  "steps": 1
}
```

### 3. 获取当前状态

**端点**: `GET /api/simulation/state`

### 4. 重置模拟

**端点**: `POST /api/simulation/reset`

## 数据模型

### Bird (鸟类实体)

```go
type Bird struct {
    ID             string   `json:"id"`
    Position       Point    `json:"position"`        // 当前位置
    Velocity       Velocity `json:"velocity"`        // 速度向量
    DetectionRadius float64 `json:"detection_radius"` // 检测半径
    DefaultSpeed   float64 `json:"default_speed"`    // 默认速度
    FOVAngle       float64 `json:"fov_angle"`        // 视野角度
}
```

### SimulationState (模拟状态)

```go
type SimulationState struct {
    Birds     []Bird    `json:"birds"`
    Obstacles []Obstacle `json:"obstacles"`
    Step      int       `json:"step"`       // 当前步数
    Timestamp int64     `json:"timestamp"`  // 时间戳
}
```

## 启动方式

```bash
go run main.go
```

服务启动在 `http://localhost:8080`

## 扩展性说明

- 支持动态添加/移除鸟类和障碍物
- 可配置的模拟参数
- 线程安全的并发访问
- 易于集成可视化前端

此文档为网页端开发提供了完整的后端 API 参考，前端可通过标准的 HTTP RESTful 接口与模拟系统进行交互。
