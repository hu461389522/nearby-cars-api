# 附近车辆查询系统（Go + Redis + MySQL）

一个基于 **Go + Gin + Redis + MySQL** 的车辆调度后台原型，支持**附近车辆查询**、**动态添加/删除车辆**，数据**持久化存储**，并提供 **HTTP API** 接口。

## 技术栈

- Go 1.27+
- Gin Web 框架
- Redis（GEO 地理位置查询）
- MySQL（数据持久化）

## 功能列表

- `GET /nearby?lng=X&lat=Y`：查询 20 公里内的空闲车辆，按距离排序
- `POST /car`：添加新车（JSON 格式）
- `DELETE /car/{plate}`：删除车辆

## 快速运行

```bash
go mod tidy
go run main.go
```

## 测试示例

```bash
# 添加车辆
curl -X POST http://localhost:8080/car -H "Content-Type: application/json" -d "{\"plate\":\"TEST001\",\"lng\":116.48,\"lat\":39.92}"

# 查询附近车辆
curl http://localhost:8080/nearby?lng=116.48&lat=39.92

# 删除车辆
curl -X DELETE http://localhost:8080/car/TEST001
```