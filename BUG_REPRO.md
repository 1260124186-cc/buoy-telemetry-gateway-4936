# 修复前故障复现（Docker）

## 项目与标准命令

项目为 Go 1.26 的浮标遥测 HTTP 服务。标准命令为 `go build ./...` 和 `go test ./...`。

## 环境构建与编译

在 arm64 主机执行 `docker build -f benzhi.Dockerfile -t buoy-telemetry-gateway:bug .` 成功，容器内 `go version` 为 `go1.26.6 linux/arm64`，`go build ./...` 通过。

## 故障触发步骤

```bash
go test ./internal/service -run TestBuoyIdentityUsesOneCanonicalKey -count=20
```

## 实际错误输出

```text
--- FAIL: TestBuoyIdentityUsesOneCanonicalKey (0.00s)
    buoy_identity_test.go:23: calibration did not match canonical buoy: 4
FAIL
FAIL    buoy-telemetry-gateway/internal/service
```

## 期望行为

同一浮标的带空格或大小写变体应复用已有校准和历史窗口。
