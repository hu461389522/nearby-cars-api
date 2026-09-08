# 🚗 附近车辆查询 API

一个基于 **Go + Gin + Redis + MySQL** 的车辆调度后台原型，支持附近车辆查询、动态注册与删除，数据持久化存储。

---

## 🌐 在线体验

> **无需下载，直接访问以下链接即可看到 API 返回结果：**
http://139.199.18.67:8080/nearby?lng=116.48&lat=39.92
（当前数据库为空，所以返回 `data: null`，这是正常的。你可以通过 POST 请求添加车辆后再查询。）

---

## ✨ 功能特性

| 功能 | 方法 | 接口 | 说明 |
|------|------|------|------|
| 查询附近车辆 | GET | `/nearby?lng=X&lat=Y` | 返回 20 公里内的车辆，按距离升序 |
| 添加车辆 | POST | `/car` | 同时写入 MySQL 和 Redis 缓存 |
| 删除车辆 | DELETE | `/car/{车牌号}` | 同时从 MySQL 和 Redis 移除 |

---

## 🛠️ 技术栈

| 组件 | 用途 |
|------|------|
| Go 1.23 | 后端开发语言 |
| Gin | Web 框架，提供 RESTful API |
| Redis GEO | 高性能地理位置查询与距离计算 |
| MySQL | 车辆数据持久化存储 |
| GitHub | 代码托管与版本管理 |

---

## 🏗️ 系统架构
┌─────────┐ ┌─────────┐ ┌─────────┐
│ 客户端 │ ──→ │ Gin │ ──→ │ Redis │（缓存，毫秒级查询）
│ (浏览器) │ ←── │ (Go) │ ←── │ (GEO) │
└─────────┘ └─────────┘ └─────────┘
│
▼
┌─────────┐
│ MySQL │（持久化存储，数据不丢失）
└─────────┘

- 服务启动时，自动将 MySQL 中的车辆数据**全量加载**到 Redis 缓存中。
- 查询请求**直接从 Redis 读取**，保证高并发下的实时性。
- 增删操作**同时更新 MySQL 和 Redis**，保持数据一致性。

---

## 📦 快速开始（本地运行）

### 环境要求

- Go 1.23+
- MySQL 8.0
- Redis 6.0+

### 1. 克隆项目

```bash
git clone https://github.com/hu461389522/nearby-cars-api.git
cd nearby-cars-api

2. 创建数据库
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
dsn := "root:你的密码@tcp(127.0.0.1:3306)/car_db?charset=utf8mb4&parseTime=True"

4. 运行
bash
go mod tidy
go run main.go


📡 API 接口文档
1. 查询附近车辆
请求

text
GET /nearby?lng={经度}&lat={纬度}
示例

bash
curl "http://localhost:8080/nearby?lng=116.48&lat=39.92"
响应
{
  "code": 200,
  "message": "success",
  "data": [
    {"name": "京A88888", "dist_km": 2.1547},
    {"name": "京B66666", "dist_km": 5.7097}
  ]
}

2. 添加车辆
请求
POST /car
Content-Type: application/json

示例

bash
curl -X POST "http://localhost:8080/car" \
  -H "Content-Type: application/json" \
  -d '{"plate":"京A88888","lng":116.397,"lat":39.908}'
响应

json
{
  "code": 200,
  "message": "✅ 车辆 京A88888 入库并同步缓存成功",
  "data": {"plate":"京A88888","lng":116.397,"lat":39.908}
}

3. 删除车辆
请求

text
DELETE /car/{车牌号}
示例

bash
curl -X DELETE "http://localhost:8080/car/京A88888"
响应

json
{
  "code": 200,
  "message": "✅ 车辆 京A88888 已从数据库和缓存中删除"
}

📁 项目结构
text
nearby-cars-api/
├── main.go          # 主程序（含所有接口）
├── go.mod           # Go 模块依赖
├── README.md        # 项目说明
└── .gitignore       # Git 忽略文件
🚀 部署
项目已部署至腾讯云轻量应用服务器，外网可访问：

text
http://139.199.18.67:8080/nearby?lng=116.48&lat=39.92
📝 未来计划
□ 添加 JWT 鉴权，保护接口安全
□ 支持车辆状态（空闲/繁忙）筛选
□ 接入高德地图，可视化展示车辆位置
□ 增加分页查询，支持大数据量场景
👤 关于作者
一个正在转行 Go 后端的开发者，用这个项目作为学习和求职的作品。

📄 License
MIT License







