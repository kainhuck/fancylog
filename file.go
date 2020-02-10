package fancylog

import (
	"fmt"
	"os"
	"path"
	"time"
)

// 向文件中写入日志
var (
	// MaxSize 通道缓存
	MaxSize = 50000
)

// FileLogger 日志结构体
type FileLogger struct {
	Level       LogLevel
	filePath    string // 日志文件保存路径
	fileName    string // 日志文件名
	fileObj     *os.File
	errFileObj  *os.File
	maxFileSize int64 // 最大文件大小
	logChan     chan *logMsg
}

type logMsg struct {
	Level       LogLevel
	LevelString string
	msg         string
	funcName    string
	fileName    string
	timeStamp   string
	line        int
}

// NewFileLogger 构造函数
func NewFileLogger(levelString, fp, fn string, maxSize int64) *FileLogger {
	logLevel, err := parseLogLevel(levelString)
	if err != nil {
		panic(err)
	}
	fl := &FileLogger{
		Level:       logLevel,
		filePath:    fp,
		fileName:    fn,
		maxFileSize: maxSize,
		logChan:     make(chan *logMsg, MaxSize),
	}
	err = fl.initFile() // 按照文件路径和文件名将文件打开
	if err != nil {
		panic(err)
	}
	return fl
}

// 打开日志文件
func (f *FileLogger) initFile() error {
	fullFileName := path.Join(f.filePath, f.fileName)
	fileObj, err := os.OpenFile(fullFileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Printf("open log file failed, err: %v", err)
		return err
	}
	errFileObj, err := os.OpenFile(fullFileName+".err", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Printf("open err log file failed, err: %v", err)
		return err
	}
	//日志文件都已经打开
	f.fileObj = fileObj
	f.errFileObj = errFileObj

	// 开启后台goroutine写日志
	go f.writeLogBackground()
	return nil
}

// 检查日志文件大小是否需要切割
func (f *FileLogger) checkSize(file *os.File) bool {
	fileInfo, err := file.Stat()
	if err != nil {
		fmt.Printf("Open file failed, err:%v\n", err)
		return false
	}
	// 如果当前文件大小大于等于日志文件的最大值就返回true
	return fileInfo.Size() >= f.maxFileSize
}

// 检查是否需要记录日志
func (f *FileLogger) enable(level LogLevel) bool {
	return level >= f.Level
}

// 切割文件
func (f *FileLogger) splitFile(file *os.File) (*os.File, error) {
	// 需要切割
	// 1. 备份一下 rename xxx.log -> xxx.log.bak202001010000000
	nowStr := time.Now().Format("20060102150405000")
	fileInfo, err := file.Stat()
	if err != nil {
		fmt.Printf("open log file failed, err:%v\n", err)
		return nil, err
	}
	logName := path.Join(f.filePath, fileInfo.Name())      // 拿到当前日志文件的完整路径
	newLogName := fmt.Sprintf("%s.bak%s", logName, nowStr) // 拼接一个日志文件备份名字

	file.Close() // 关闭当前文件

	os.Rename(logName, newLogName)
	// 2. 打开一个新的文件
	fileObj, err := os.OpenFile(logName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Printf("open new log file failed, err:%v\n", err)
		return nil, err
	}
	// 3. 将新打开的日志文件赋值给 f.fileObj
	return fileObj, nil
}

func (f *FileLogger) writeLogBackground() {
	for {
		if f.checkSize(f.fileObj) {
			newFile, err := f.splitFile(f.fileObj)
			if err != nil {
				fmt.Printf("check size failed,err:%v\n", err)
				return
			}
			f.fileObj = newFile
		}

		select {
		case tempMsg := <-f.logChan:
			logInfo := fmt.Sprintf("[%s] [%s] [%s:%s:%d] %s\n", tempMsg.timeStamp, tempMsg.LevelString, tempMsg.fileName, tempMsg.funcName, tempMsg.line, tempMsg.msg)
			fmt.Fprintf(f.fileObj, logInfo)
			if tempMsg.Level >= ERROR {
				// 如果要记录的日志等级大于ERROR，需要记录到专门的日志文件
				if f.checkSize(f.errFileObj) {
					newFile, err := f.splitFile(f.errFileObj)
					if err != nil {
						fmt.Printf("check size failed,err:%v\n", err)
						return
					}
					f.errFileObj = newFile
				}
				fmt.Fprintf(f.errFileObj, logInfo)
			}
		default:
			// 取不到日志休息500ms
			time.Sleep(500 * time.Millisecond)
		}
	}
}

// 记录日志
func (f *FileLogger) log(lv LogLevel, format string, a ...interface{}) {
	if f.enable(lv) {
		msg := fmt.Sprintf(format, a...)
		levelString, err := getLevelString(lv)
		if err != nil {
			panic(err)
		}
		now := time.Now()
		funcName, fileName, lineNo := getInfo(3)
		// 先把日志发送到通道
		tempMsg := &logMsg{
			Level:       lv,
			LevelString: levelString,
			msg:         msg,
			funcName:    funcName,
			fileName:    fileName,
			timeStamp:   now.Format("20060102150405000"),
			line:        lineNo,
		}
		select {
		case f.logChan <- tempMsg:
		default:
			// 丢掉，防止阻塞
		}

	}
}

// Debug 调试
func (f *FileLogger) Debug(format string, a ...interface{}) {
	f.log(DEBUG, format, a...)
}

// Trace 追踪信息
func (f *FileLogger) Trace(format string, a ...interface{}) {
	f.log(TRACE, format, a...)
}

// Info 信息
func (f *FileLogger) Info(format string, a ...interface{}) {
	f.log(INFO, format, a...)
}

// Warning 警告
func (f *FileLogger) Warning(format string, a ...interface{}) {
	f.log(WARNING, format, a...)
}

// Error 错误
func (f *FileLogger) Error(format string, a ...interface{}) {
	f.log(ERROR, format, a...)
}

// Fatal 致命错误
func (f *FileLogger) Fatal(format string, a ...interface{}) {
	f.log(FATAL, format, a...)
}
