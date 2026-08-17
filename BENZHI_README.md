# buoy-telemetry-gateway__003 Docker 交付说明

## 项目概览
- Buoy Telemetry Gateway is a small Go HTTP service for coastal engineering teams that receive measurements from autonomous buoys. It validates incoming readings, applies per-sensor 
- Go module: `buoy-telemetry-gateway`

## 标准命令

```bash
go build ./...
go test ./...
```

## 实际启动入口

```bash
go run ./cmd/server
```

## Docker 构建

```bash
./build_benzhi_docker.sh buoy-telemetry-gateway__003-benzhi linux/amd64
docker run --rm -it buoy-telemetry-gateway__003-benzhi bash
```

## 环境

- 基础镜像: `golang:1.26`
- 依赖在镜像构建阶段预下载，容器内可直接执行 Go 构建和测试命令。
- 代码目录: `/app`
- 源码中检测到的服务端口: `31`, `8080`
