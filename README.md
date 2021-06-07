# tcping
tcping是一个检查域名指定端口是否可达的工具  
项目主要参考了jlyo/tcping
## 依赖
go 1.14版本  
支持linux和windows10

## 使用方法
    count表示ping次数
    hostname 域名
    port  端口号
    go run tcping.go -host $(hostname) -port $(port) -c $(count)

