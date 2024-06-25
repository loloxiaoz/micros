# micros

微服务框架

## todo

- [x] api文档生成 [gin](https://github.com/swaggo/gin-swagger)

```shell
```

- [x] docker-compose
- [x] config 库
- [ ] 可观测性
- [ ] orm
- [ ] 数据库
- [ ] gitbook文档生成
- [ ] 云原生
- [ ] 配置文件开关

## 启动命令

```shell
swag init -g ./cmd/server/main.go  -o api
make build
docker run -d -p 8090:8090   micros:unknown
```
