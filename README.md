# chatrooms
一个Go语言写的在线聊天室demo
可通过广播群发 消息
可以通过--name=私发用户名 消息 
可以通过--rename=自定义名字 修改自己的用户名
可以通过--who 查询所有在线用户名
## 运行服务器

```Go
go run server.go
```

## 运行第一个客户端并修改名称
```Go
go run client.go
--rename=ikun666
```
## 运行第二个客户端并修改名称
```Go
go run client.go
--rename=cxk
```
## 运行第三个客户端并修改名称
```Go
go run client.go
--rename=ikun
--who
hello,everyone
--name=cxk sing dance rap and play basketball
```
更多功能等待实现中