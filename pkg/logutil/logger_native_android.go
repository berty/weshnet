//go:build android

package logutil

/*
#cgo LDFLAGS: -llog
#include <android/log.h>
#include <stdlib.h>
*/
import "C"

import (
	"fmt"
	"unsafe"

	"go.uber.org/zap/zapcore"
)

func NativeLog(logLevel zapcore.Level, namespace string, message string) {
	var level C.int = C.ANDROID_LOG_INFO

	switch logLevel {
	case zapcore.DebugLevel:
		level = C.ANDROID_LOG_DEBUG
	case zapcore.InfoLevel:
		level = C.ANDROID_LOG_INFO
	case zapcore.WarnLevel:
		level = C.ANDROID_LOG_WARN
	case zapcore.ErrorLevel:
		level = C.ANDROID_LOG_ERROR
	}

	cNamespace := C.CString(namespace)
	defer C.free(unsafe.Pointer(cNamespace))
	cMessage := C.CString(fmt.Sprintf("[%s] %s", logLevel.CapitalString(), message))
	defer C.free(unsafe.Pointer(cMessage))

	C.__android_log_write(level, cNamespace, cMessage)
}
