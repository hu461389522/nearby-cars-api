# 内部工具后台 API 系统

一个基于 **Go + Gin + Redis + MySQL** 的后台管理 API 服务，支持数据增删查、缓存同步、接口化操作，可作为内部工具、AI 应用、管理后台等场景的后端底座。

---

## 🌐 在线体验

> **API 接口可直接访问测试（无需下载运行）：**

```
http://139.199.18.67:8080/nearby?lng=116.48&lat=39.92
```

（当前数据库为空，返回空数据，可通过 POST 添加测试数据。）

---

## ✨ 功能特性

| 功能 | 方法 | 接口 | 说明 |
|------|------|------|------|
| 数据查询 | GET | `/nearby?lng=X&lat=Y` | 按条件筛选数据，返回 JSON |
| 数据添加 | POST | `/car` | 新增数据，同步写入 MySQL + Redis |
| 数据删除 | DELETE | `/car/{id}` | 删除数据，同步清理缓存 |

> 这套 CRUD 接口模式可直接复用到内部管理系统、AI 工具后台、智能客服数据管理等场景。

---

## 🛠️ 技术栈

| 组件 | 用途 |
|------|------|
| Go 1.23 | 后端开发语言 |
| Gin | Web 框架，提供 RESTful API |
| Redis | 缓存加速，提升查询性能 |
| MySQL | 数据持久化存储 |
| GORM | ORM 框架，简化数据库操作 |
| Docker | 容器化部署（规划中） |

---

## 🏗️ 系统架构（适用于内部工具后台）

```
┌─────────┐      ┌─────────┐      ┌─────────┐
│ 前端/   │ ──→ │   Gin   │ ──→ │  Redis  │（缓存层）
│ AI应用  │ ←── │ (Go)    │ ←── │         │
└─────────┘      └─────────┘      └─────────┘
                      │
                      ▼
                 ┌─────────┐
                 │  MySQL  │（持久化存储）
                 └─────────┘
```

- 服务启动时自动加载 MySQL 数据至 Redis 缓存
- 查询请求优先走 Redis，降低数据库压力
- 增删操作同步更新 MySQL + Redis，保持数据一致性

---

## 📦 快速开始

### 环境要求

- Go 1.23+
- MySQL 8.0
- Redis 6.0+

### 1. 克隆项目

```bash
git clone https://github.com/hu461389522/nearby-cars-api.git
cd nearby-cars-api
```

### 2. 创建数据库

```sql
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
```

### 3. 修改数据库连接

打开 `main.go`，修改 `dsn` 中的数据库密码：

```go
dsn := "root:你的密码@tcp(127.0.0.1:3306)/car_db?charset=utf8mb4&parseTime=True"
```

### 4. 运行

```bash
go mod tidy
go run main.go
```

---

## 📡 API 接口文档

### 1. 查询数据

**请求**
```
GET /nearby?lng={经度}&lat={纬度}
```

**示例**
```bash
curl "http://localhost:8080/nearby?lng=116.48&lat=39.92"
```

**响应**
```json
{
  "code": 200,
  "message": "success",
  "data": [
    {"name": "数据A", "dist_km": 2.15},
    {"name": "数据B", "dist_km": 5.71}
  ]
}
```

---

### 2. 添加数据

**请求**
```
POST /car
Content-Type: application/json
```

**示例**
```bash
curl -X POST "http://localhost:8080/car" \
  -H "Content-Type: application/json" \
  -d '{"plate":"数据001","lng":116.397,"lat":39.908}'
```

**响应**
```json
{
  "code": 200,
  "message": "✅ 数据 数据001 入库成功",
  "data": {"plate":"数据001","lng":116.397,"lat":39.908}
}
```

---

### 3. 删除数据

**请求**
```
DELETE /car/{id}
```

**示例**
```bash
curl -X DELETE "http://localhost:8080/car/数据001"
```

**响应**
```json
{
  "code": 200,
  "message": "✅ 数据 数据001 已删除"
}
```

---

## 📁 项目结构

```
nearby-cars-api/
├── main.go          # 主程序（接口 + 数据库 + 缓存）
├── go.mod           # Go 模块依赖
├── README.md        # 项目说明
└── .gitignore       # Git 忽略文件
```

---

## 🚀 部署

项目已部署至腾讯云轻量应用服务器，外网可访问：

```
http://139.199.18.67:8080/nearby?lng=116.48&lat=39.92
```

---

## 📝 关于本项目

本项目是一个**后台管理 API 系统**的示例，涵盖：
- RESTful API 接口设计
- MySQL 数据持久化
- Redis 缓存加速
- 服务启动时自动加载缓存
- 增删操作同步更新 MySQL + Redis

这套架构可直接复用到：
- 内部管理平台后端
- AI 工具后台数据接口
- 智能客服系统数据层
- 企业数字化工具底座

---

## 👤 关于我

电商运营转 Go 后端开发，熟悉业务需求分析、数据看板搭建，具备独立完成 API 系统从开发到部署上线的能力，持续学习 AI 工具落地方向。

---

## 📄 License

MIT License
