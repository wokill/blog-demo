# mk-blog-svc

> 注意不要将代码提交到远程分支

## 初始化环境

[参考这里](https://git.dustess.com/mk-training/training-env-compose)

## 启动

```shell script
git clone ssh://git@git.dustess.com:60022/mk-training/training-env-compose.git
git checkout -b {{ username }}
cp config.template.json config.json
go run cmd/server/main.go
```
