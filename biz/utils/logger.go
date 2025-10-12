package utils

import (
	"fmt"
	"runtime"
	"strings"

	"github.com/sirupsen/logrus"
)

var Logger *logrus.Logger

// InitLogger 初始化日志器
func InitLogger() {
	Logger = logrus.New()
	
	// 设置日志格式
	Logger.SetFormatter(&logrus.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: "2006/01/02 15:04:05",
		ForceColors:     true,
	})
	
	// 设置日志级别
	Logger.SetLevel(logrus.DebugLevel)
}

// LogError 记录错误日志，包含调用栈信息
func LogError(err error, message string) {
	if Logger == nil {
		InitLogger()
	}
	
	// 获取调用者信息
	pc, file, line, ok := runtime.Caller(1)
	if !ok {
		if err != nil {
			Logger.WithError(err).Error(message)
		} else {
			Logger.Error(message)
		}
		return
	}
	
	// 获取函数名
	fn := runtime.FuncForPC(pc)
	funcName := "unknown"
	if fn != nil {
		funcName = fn.Name()
		// 简化函数名，只保留最后一部分
		if idx := strings.LastIndex(funcName, "."); idx != -1 {
			funcName = funcName[idx+1:]
		}
	}
	
	// 简化文件路径，只保留相对路径
	if idx := strings.LastIndex(file, "/orbia_api/"); idx != -1 {
		file = file[idx+1:]
	}
	
	fields := logrus.Fields{
		"file":     fmt.Sprintf("%s:%d", file, line),
		"function": funcName,
	}
	
	if err != nil {
		fields["error"] = err.Error()
	}
	
	Logger.WithFields(fields).Error(message)
}

// LogPanic 记录panic日志，包含详细的调用栈信息
func LogPanic(recovered interface{}, message string) {
	if Logger == nil {
		InitLogger()
	}
	
	// 获取调用者信息
	pc, file, line, ok := runtime.Caller(1)
	if !ok {
		Logger.WithField("panic", recovered).Error(message)
		return
	}
	
	// 获取函数名
	fn := runtime.FuncForPC(pc)
	funcName := "unknown"
	if fn != nil {
		funcName = fn.Name()
		// 简化函数名，只保留最后一部分
		if idx := strings.LastIndex(funcName, "."); idx != -1 {
			funcName = funcName[idx+1:]
		}
	}
	
	// 简化文件路径，只保留相对路径
	if idx := strings.LastIndex(file, "/orbia_api/"); idx != -1 {
		file = file[idx+1:]
	}
	
	// 获取完整的调用栈
	buf := make([]byte, 4096)
	n := runtime.Stack(buf, false)
	stack := string(buf[:n])
	
	Logger.WithFields(logrus.Fields{
		"file":     fmt.Sprintf("%s:%d", file, line),
		"function": funcName,
		"panic":    recovered,
		"stack":    stack,
	}).Error(message)
}

// LogInfo 记录信息日志
func LogInfo(message string, fields ...logrus.Fields) {
	if Logger == nil {
		InitLogger()
	}
	
	if len(fields) > 0 {
		Logger.WithFields(fields[0]).Info(message)
	} else {
		Logger.Info(message)
	}
}

// LogDebug 记录调试日志
func LogDebug(message string, fields ...logrus.Fields) {
	if Logger == nil {
		InitLogger()
	}
	
	if len(fields) > 0 {
		Logger.WithFields(fields[0]).Debug(message)
	} else {
		Logger.Debug(message)
	}
}

// LogWarn 记录警告日志
func LogWarn(message string, fields ...logrus.Fields) {
	if Logger == nil {
		InitLogger()
	}
	
	if len(fields) > 0 {
		Logger.WithFields(fields[0]).Warn(message)
	} else {
		Logger.Warn(message)
	}
}