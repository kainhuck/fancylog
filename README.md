# fancylog
fancylog for golang

# features
- 支持终端输出和文件输👍
- 终端彩色输出🎽
- 日志输出时间和代码定位🚙
- 文件切割📜
- 错误单独储存🚀

# How to use
```go
go get github.com/kainhuck/fancylog
```

```go
var Log fancylog.Logger
ConsoleLog := fancylog.NewConsoleLogger("debug")
FileLog := fancylog.NewFileLogger("debug", "./log", "fancylog.log", 1024*8)
Log = ConsoleLog
//Log = FileLog
Log.Info("hello %s", "fancylog")
```