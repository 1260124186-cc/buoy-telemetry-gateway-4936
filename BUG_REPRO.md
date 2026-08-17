# 修复前故障复现（Docker）

## 项目与标准命令

项目为 Go 1.26 的浮标遥测 HTTP 服务。标准命令为 `go build ./...` 和 `go test ./...`。

## 环境构建与编译

在 arm64 主机执行 `docker build -f benzhi.Dockerfile -t buoy-telemetry-gateway:bug .` 成功，容器内 `go version` 为 `go1.26.6 linux/arm64`，`go build ./...` 通过。

## 故障触发步骤

在容器内执行：

```bash
go test ./internal/httpapi -run TestCanceledIngestDoesNotCommitReading -count=20
```

## 实际错误输出

```text
--- FAIL: TestCanceledIngestDoesNotCommitReading (0.00s)
    cancellation_test.go:25: canceled request was committed: {"buoy_id":"east-2","sensor":"salinity","value":31,"unit":"ppt","sequence":1}
FAIL
FAIL    buoy-telemetry-gateway/internal/httpapi
```

## 期望行为

请求上下文取消后接口不得返回创建成功，也不得在最近读数或事件中留下该次上报。
