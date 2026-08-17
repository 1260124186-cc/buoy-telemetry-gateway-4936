# 修复前故障复现（Docker）

## 项目与标准命令

项目为 Go 1.26 的浮标遥测 HTTP 服务。标准命令为 `go build ./...` 和 `go test ./...`。

## 环境构建与编译

在 arm64 主机执行 `docker build -f benzhi.Dockerfile -t buoy-telemetry-gateway:bug .` 成功，容器内 `go version` 为 `go1.26.6 linux/arm64`，`go build ./...` 通过。

## 故障触发步骤

在容器内执行：

```bash
go test ./internal/httpapi -run TestErrorContractPreservesClassification -count=20
```

## 实际错误输出

```text
--- FAIL: TestErrorContractPreservesClassification (0.00s)
    error_contract_test.go:20: invalid calibration status = 500, body={"error":"invalid calibration: scale must not be zero"}
FAIL
FAIL    buoy-telemetry-gateway/internal/httpapi
```

## 期望行为

不合法的校准请求返回客户端错误；缺少可选校准配置时仍可接收原始读数。
