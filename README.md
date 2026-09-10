markdown
# 内部工具后台 API 系统

一个基于 Go + Gin + Redis + MySQL 的后台管理 API 服务，支持数据增删查、缓存同步、接口化操作，可作为内部工具、AI 应用、管理后台等场景的后端底座。

---

## 🌐 在线体验

API 接口已部署至公网，可直接访问测试（无需下载运行），预置了测试数据：

👉 http://139.199.18.67:8080/nearby?lng=116.42&lat=39.90

返回示例：
```json
{
  "code": 200,
  "data": [
    {"name": "测试A", "dist_km": 2.15},
    {"name": "测试B", "dist_km": 5.71},
    {"name": "测试C", "dist_km": 10.02}
  ],
  "message": "success"
}

## ✨ 功能特性

| 功能 | 方法 | 接口 | 说明 |
|------|------|------|------|
| 数据查询 | GET | `/nearby?lng=X&lat=Y` | 按条件筛选，返回 JSON 数据 |
| 数据添加 | POST | `/car` | 新增数据，同步写入 MySQL + Redis |
| 数据删除 | DELETE | `/car/{plate}` | 删除数据，同步清理缓存 |
| **AI 智能客服** | POST | `/chat` | 对接 DeepSeek 大模型，实现智能问答 |
| **聊天页面** | GET | `/` | 浏览器打开即可与 AI 对话 |

> **在线体验：**
> - 🤖 AI 聊天页面：http://139.199.18.67:8080/
> - 🚗 找车 API：http://139.199.18.67:8080/nearby?lng=116.42&lat=39.90

🛠️ 技术栈
组件	用途
Go 1.23	后端开发语言
Gin	Web 框架，提供 RESTful API
GORM	ORM 框架，简化 MySQL 操作
Redis	缓存加速，毫秒级查询响应
MySQL	数据持久化存储
🏗️ 系统架构（分层设计）
text
┌─────────┐      ┌─────────────┐      ┌─────────┐
│ 前端/   │ ──→ │  Controller │ ──→ │  Redis  │（缓存层）
│ API调用 │ ←── │  (Gin)      │ ←── │         │
└─────────┘      └─────────────┘      └─────────┘
                        │
                        ▼
                 ┌─────────────┐
                 │   Service   │（业务逻辑层）
                 └─────────────┘
                        │
                        ▼
                 ┌─────────────┐
                 │   Model     │（数据模型 + GORM）
                 └─────────────┘
                        │
                        ▼
                 ┌─────────────┐
                 │   MySQL     │（持久化存储）
                 └─────────────┘
Controller 层：处理 HTTP 请求，参数校验

Service 层：实现核心业务逻辑（MySQL + Redis 同步）

Model 层：定义数据实体，使用 GORM 操作数据库

服务启动时自动加载 MySQL 数据至 Redis 缓存

增删操作同步更新 MySQL + Redis，保证数据一致性

📦 快速开始
环境要求
Go 1.23+

MySQL 8.0

Redis 6.0+

1. 克隆项目
bash
git clone https://github.com/hu461389522/nearby-cars-api.git
cd nearby-cars-api
2. 创建数据库
sql
CREATE DATABASE IF NOT EXISTS car_db;
USE car_db;
CREATE TABLE IF NOT EXISTS cars (
    id INT AUTO_INCREMENT PRIMARY KEY,
    plate VARCHAR(20) NOT NULL UNIQUE,
    lng DECIMAL(10,7) NOT NULL,
    lat DECIMAL(10,7) NOT NULL,
    status TINYINT DEFAULT 1,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
3. 修改数据库连接
打开 main.go，修改 dsn 中的数据库密码：

go
dsn := "root:你的密码@tcp(127.0.0.1:3306)/car_db?charset=utf8mb4&parseTime=True"
4. 运行
bash
go mod tidy
go run main.go
📡 API 接口文档
1. 查询数据
请求

text
GET /nearby?lng={经度}&lat={纬度}
示例

bash
curl "http://localhost:8080/nearby?lng=116.42&lat=39.90"
响应

json
{
  "code": 200,
  "message": "success",
  "data": [
    {"name": "测试A", "dist_km": 2.15},
    {"name": "测试B", "dist_km": 5.71}
  ]
}
2. 添加数据
请求

text
POST /car
Content-Type: application/json
示例

bash
curl -X POST "http://localhost:8080/car" \
  -H "Content-Type: application/json" \
  -d '{"plate":"测试D","lng":116.40,"lat":39.90}'
响应

json
{
  "code": 200,
  "message": "添加成功",
  "data": {"plate":"测试D","lng":116.40,"lat":39.90}
}
3. 删除数据
请求

text
DELETE /car/{plate}
示例

bash
curl -X DELETE "http://localhost:8080/car/测试D"
响应

json
{
  "code": 200,
  "message": "删除成功"
}
📁 项目结构（分层架构）
text
nearby-cars-api/
├── main.go              # 入口文件
├── config/
│   └── db.go            # MySQL 连接配置
├── model/
│   └── car.go           # 数据模型定义
├── service/
│   └── car_service.go   # 业务逻辑层
├── controller/
│   └── car_controller.go # 接口处理层
├── router/
│   └── router.go        # 路由注册
├── go.mod
├── go.sum
└── README.md
🚀 部署
项目已部署至腾讯云轻量应用服务器，外网可访问：

text
http://139.199.18.67:8080/nearby?lng=116.42&lat=39.90
📝 关于本项目
本项目是一个后台管理 API 系统的示例，涵盖：

分层架构设计（Controller → Service → Model）

RESTful API 接口设计

MySQL 数据持久化 + GORM ORM 框架

Redis 缓存加速（GEO 地理位置查询）

服务启动时自动加载缓存

增删操作同步更新 MySQL + Redis

这套架构可直接复用到：

内部管理平台后端

AI 工具后台数据接口

智能客服系统数据层

企业数字化工具底座

👤 关于作者
电商运营转 Go 后端开发，熟悉业务需求分析、数据看板搭建，具备独立完成 API 系统从开发到部署上线的能力，持续学习 AI 工具落地方向。

📄 License
MIT License