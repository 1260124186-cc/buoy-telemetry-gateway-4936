# 修复前故障复现（Docker）

## 项目与标准命令

项目为 Go 1.26 的浮标遥测 HTTP 服务。标准命令为 `go build ./...` 和 `go test ./...`。

## 环境构建与编译

在 arm64 主机执行 `docker build -f benzhi.Dockerfile -t buoy-telemetry-gateway:bug .` 成功，容器内 `go version` 为 `go1.26.6 linux/arm64`，`go build ./...` 通过。

## 故障触发步骤

在容器内执行：

```bash
go test ./internal/service -run TestHistoryWindowIsStableAndSummaryKeepsRange -count=20
```

## 实际错误输出

```text
--- FAIL: TestHistoryWindowIsStableAndSummaryKeepsRange (0.00s)
    history_window_test.go:25: unstable newest-first history: first=[sequence:1 sequence:2 sequence:3] second=[sequence:1 sequence:2 sequence:3]
FAIL
FAIL    buoy-telemetry-gateway/internal/service
```

## 期望行为

重复查询最近窗口始终返回相同的最新到最旧顺序，且汇总范围与均值不受查询操作影响。
