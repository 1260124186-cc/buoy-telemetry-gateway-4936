# 修复前故障复现（Docker）

## 项目与标准命令

项目为 Go 1.26 的浮标遥测 HTTP 服务。标准命令为 `go build ./...` 和 `go test ./...`。

## 环境构建与编译

在 arm64 主机执行 `docker build -f benzhi.Dockerfile -t buoy-telemetry-gateway:bug .` 成功，容器内 `go version` 为 `go1.26.6 linux/arm64`，`go build ./...` 通过。

## 故障触发步骤

在容器内执行：

```bash
go test ./internal/httpapi -run TestReadingWithoutOptionalLabelsIsAccepted -count=20
```

## 实际错误输出

```text
--- FAIL: TestReadingWithoutOptionalLabelsIsAccepted (0.00s)
panic: assignment to entry in nil map [recovered, repanicked]

buoy-telemetry-gateway/internal/httpapi.(*Handler).reading(...)
    internal/httpapi/handler.go:47
FAIL    buoy-telemetry-gateway/internal/httpapi
```

## 期望行为

客户端省略可选 labels 时读数仍被正常接收，后续保存和查询不会 panic 或污染其他读数。
