package logger

import (
	"context"
	"errors"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

// GormLogger 操作对象，实现 gormlogger.Interface
type GormLogger struct {
	ZapLogger     *zap.Logger
	SlowThreshold time.Duration
}

func NewGormLogger() GormLogger {
	return GormLogger{
		ZapLogger:     Logger,                 // 使用全局的 logger.Logger 对象
		SlowThreshold: 200 * time.Millisecond, // 慢查询阈值，单位为千分之一秒
	}
}

func (g GormLogger) LogMode(level logger.LogLevel) logger.Interface {
	return NewGormLogger()
}

func (g GormLogger) Info(ctx context.Context, s string, i ...interface{}) {

	g.logger().Sugar().Infof(s, i...)
}

func (g GormLogger) Warn(ctx context.Context, s string, i ...interface{}) {
	g.logger().Sugar().Warnf(s, i...)
}

func (g GormLogger) Error(ctx context.Context, s string, i ...interface{}) {
	g.logger().Sugar().Errorf(s, i...)
}

func (g GormLogger) Trace(ctx context.Context, begin time.Time, fc func() (sql string, rowsAffected int64), err error) {

	dura := time.Since(begin)
	sql, rows := fc()

	logFields := []zap.Field{
		zap.String("zap", sql),
		zap.Int("time", int(dura.Milliseconds())),
		zap.Int("rows", int(rows)),
	}
	if err != nil {
		// 记录未找到的错误使用 warning 等级
		if errors.Is(err, gorm.ErrRecordNotFound) {
			g.logger().Warn("Database ErrRecordNotFound", logFields...)
		} else {
			// 其他错误使用 error 等级
			logFields = append(logFields, zap.Error(err))
			g.logger().Error("Database Error", logFields...)
		}
	}

	// 慢查询日志
	if g.SlowThreshold != 0 && dura > g.SlowThreshold {
		g.logger().Warn("Database Slow Log", logFields...)
	}

	// 记录所有 SQL 请求
	g.logger().Debug("Database Query", logFields...)
}

// 目的跳过无用的调用
func (g GormLogger) logger() *zap.Logger {
	// 跳过 gorm 内置的调用
	var (
		gormPackage    = filepath.Join("gorm.io", "gorm")
		zapgormPackage = filepath.Join("moul.io", "zapgorm2")
	)
	for i := 2; i < 20; i++ {
		_, file, _, ok := runtime.Caller(i)
		switch {
		case !ok:
		case strings.HasSuffix(file, "_test.go"):
		case strings.Contains(file, gormPackage):
		case strings.Contains(file, zapgormPackage):
		case strings.Contains(file, "mysql@"):
		default:
			return g.ZapLogger.WithOptions(zap.AddCallerSkip(i - 2))

		}
	}
	return g.ZapLogger
}
