# ColdChain 冷链物流温控与批次追溯平台

ColdChain 面向生鲜与医药冷链场景：仓库与运输车辆装配温度探点，探点按批次
周期上报温度并受配额约束，平台按时间窗口聚合温控曲线，按已发布阈值规则判定
超温并冻结批次，支持解冻、出库、签收等批次流转，全程操作留审计，可按代际
回放温控窗口。全部数据以 JSON 文件持久化在本机数据目录，无外部数据库依赖。

## 运行

```bash
go build -mod=vendor -o coldchain ./cmd/coldchain
./coldchain -addr :8080 -data ./data -seed
```

打开 http://localhost:8080/ 查看探点控制台，/ui/batches、/ui/temperature、
/ui/audit 分别查看批次、温控读数与操作审计页面；/api/status 返回平台汇总，
其余 /api/* 提供命名空间、探点、批次、温控、规则、告警、配额、追溯与审计的
JSON 接口。

## 测试

```bash
go test -mod=vendor ./...
go vet -mod=vendor ./...
```

## Docker

```bash
bash build_benzhi_docker.sh coldchain linux/amd64
docker run --rm -p 8080:8080 coldchain bash -c 'go run ./cmd/coldchain -addr :8080 -data /tmp/data -seed'
```
