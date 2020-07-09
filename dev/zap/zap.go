package zap

/*
 * @Description:
 * @Author: leisc
 * @Version: 1.0.0
 * @Date: 2020-06-17 10:43:30
 * @LastEditTime: 2020-06-18 10:27:58
 */

import (
	"os"
	"path/filepath"
	"sync"

	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Options struct {
	zap.Config
	LogFileDir    string `json:logFileDir`
	AppName       string `json:"appName"`
	ErrorFileName string `json:"errorFileName"`
	WarnFileName  string `json:"warnFileName"`
	InfoFileName  string `json:"infoFileName"`
	DebugFileName string `json:"debugFileName"`
	MaxSize       int    `json:"maxSize"` // megabytes
	MaxBackups    int    `json:"maxBackups"`
	MaxAge        int    `json:"maxAge"` // days
}

var (
	l                              *Logger
	sp                             = string(filepath.Separator)
	errWS, warnWS, infoWS, debugWS zapcore.WriteSyncer       // IO输出
	debugConsoleWS                 = zapcore.Lock(os.Stdout) // 控制台标准输出
	errorConsoleWS                 = zapcore.Lock(os.Stderr)
)

func init() {
	l = &Logger{
		Opts: &Options{},
	}
	initLogger()
}

type Logger struct {
	*zap.Logger
	sync.RWMutex
	Opts      *Options `json:"opts"`
	zapConfig zap.Config
	inited    bool
}

func initLogger() {

	l.Lock()
	defer l.Unlock()

	if l.inited {
		//l.Info("[initLogger] logger Inited")
		return
	}

	l.loadCfg()
	l.init()
	//l.Info("[initLogger] zap plugin initializing completed")
	l.inited = true
}

// GetLogger returns logger
func GetLogger() (ret *Logger) {
	return l
}

func (l *Logger) init() {

	l.setSyncers()
	var err error

	l.Logger, err = l.zapConfig.Build(l.cores())
	if err != nil {
		panic(err)
	}

	defer l.Logger.Sync()
}

func (l *Logger) loadCfg() {

	l.zapConfig = zap.Config{
		Level:             zap.NewAtomicLevelAt(zap.ErrorLevel),
		Development:       false,
		DisableCaller:     false,
		DisableStacktrace: false,
		Encoding:          "json",
		EncoderConfig: zapcore.EncoderConfig{
			TimeKey:        "t",
			LevelKey:       "level",
			NameKey:        "log",
			CallerKey:      "caller",
			MessageKey:     "msg",
			StacktraceKey:  "trace",
			LineEnding:     zapcore.DefaultLineEnding,
			EncodeLevel:    zapcore.LowercaseLevelEncoder,
			EncodeTime:     zapcore.ISO8601TimeEncoder,
			EncodeDuration: zapcore.StringDurationEncoder,
			EncodeCaller:   zapcore.ShortCallerEncoder,
			EncodeName:     zapcore.FullNameEncoder,
		},
		OutputPaths:      []string{"stdout"},
		ErrorOutputPaths: []string{"stderr"},
		InitialFields:    map[string]interface{}{"ringle123": "facelive"},
	}

	/*if l.Opts.Development {
		l.zapConfig = zap.NewDevelopmentConfig()
	} else {
		l.zapConfig = zap.NewProductionConfig()
	}*/

	// application log output path
	/*
		if l.Opts.OutputPaths == nil || len(l.Opts.OutputPaths) == 0 {
			l.zapConfig.OutputPaths = []string{"stdout"}
		}

		//  error of zap-self log
		if l.Opts.ErrorOutputPaths == nil || len(l.Opts.ErrorOutputPaths) == 0 {
			l.zapConfig.OutputPaths = []string{"stderr"}
		}

		// 默认输出到程序运行目录的logs子目录
		if l.Opts.LogFileDir == "" {
			l.Opts.LogFileDir, _ = filepath.Abs(filepath.Dir(filepath.Join(".")))
			l.Opts.LogFileDir += sp + "logs" + sp
		}

		if l.Opts.AppName == "" {
			l.Opts.AppName = "app"
		}

		if l.Opts.ErrorFileName == "" {
			l.Opts.ErrorFileName = "error.log"
		}

		if l.Opts.WarnFileName == "" {
			l.Opts.WarnFileName = "warn.log"
		}

		if l.Opts.InfoFileName == "" {
			l.Opts.InfoFileName = "info.log"
		}

		if l.Opts.DebugFileName == "" {
			l.Opts.DebugFileName = "debug.log"
		}

		if l.Opts.MaxSize == 0 {
			l.Opts.MaxSize = 50
		}
		if l.Opts.MaxBackups == 0 {
			l.Opts.MaxBackups = 3
		}
		if l.Opts.MaxAge == 0 {
			l.Opts.MaxAge = 30
		}
	*/
}

func (l *Logger) setSyncers() {

	f := func(fN string) zapcore.WriteSyncer {
		return zapcore.AddSync(&lumberjack.Logger{
			Filename:   "./logs/" + sp + "facelive" + "-" + fN,
			MaxSize:    100,
			MaxBackups: 1000,
			MaxAge:     180,
			Compress:   true,
			LocalTime:  true,
		})
	}

	errWS = f("error.log")
	warnWS = f("warn.log")
	infoWS = f("info.log")
	debugWS = f("debug.log")

	return
}

func (l *Logger) cores() zap.Option {

	fileEncoder := zapcore.NewJSONEncoder(l.zapConfig.EncoderConfig)
	consoleEncoder := zapcore.NewConsoleEncoder(l.zapConfig.EncoderConfig)

	errPriority := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl >= zapcore.ErrorLevel && zapcore.ErrorLevel >= l.zapConfig.Level.Level()
	})
	warnPriority := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl == zapcore.WarnLevel && zapcore.WarnLevel >= l.zapConfig.Level.Level()
	})
	infoPriority := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl == zapcore.InfoLevel && zapcore.InfoLevel >= l.zapConfig.Level.Level()
	})
	debugPriority := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl == zapcore.DebugLevel && zapcore.DebugLevel >= l.zapConfig.Level.Level()
	})

	cores := []zapcore.Core{
		// region 日志文件

		// error 及以上
		zapcore.NewCore(fileEncoder, errWS, errPriority),

		// warn
		zapcore.NewCore(fileEncoder, warnWS, warnPriority),

		// info
		zapcore.NewCore(fileEncoder, infoWS, infoPriority),

		// debug
		zapcore.NewCore(fileEncoder, debugWS, debugPriority),

		// endregion

		// region 控制台

		// 错误及以上
		zapcore.NewCore(consoleEncoder, errorConsoleWS, errPriority),

		// 警告
		zapcore.NewCore(consoleEncoder, debugConsoleWS, warnPriority),

		// info
		zapcore.NewCore(consoleEncoder, debugConsoleWS, infoPriority),

		// debug
		zapcore.NewCore(consoleEncoder, debugConsoleWS, debugPriority),

		// endregion
	}

	return zap.WrapCore(func(c zapcore.Core) zapcore.Core {
		return zapcore.NewTee(cores...)
	})
}
