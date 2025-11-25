package persistence

import (
	"context"
	"time"

	"go.uber.org/zap"
	gormlogger "gorm.io/gorm/logger"
)

// ZapGormLogger is a small adapter that implements gorm/logger.Interface using zap.
type ZapGormLogger struct {
	logger        *zap.Logger
	slowThreshold time.Duration
	logLevel      gormlogger.LogLevel
}

// NewZapGormLogger creates a new adapter with the provided zap logger.
func NewZapGormLogger(z *zap.Logger, slow time.Duration, lvl gormlogger.LogLevel) gormlogger.Interface {
	return &ZapGormLogger{logger: z, slowThreshold: slow, logLevel: lvl}
}

func (z *ZapGormLogger) LogMode(level gormlogger.LogLevel) gormlogger.Interface {
	nz := *z
	nz.logLevel = level
	return &nz
}

func (z *ZapGormLogger) Info(ctx context.Context, msg string, data ...interface{}) {
	if z.logLevel < gormlogger.Info {
		return
	}
	z.logger.Sugar().Infof(msg, data...)
}

func (z *ZapGormLogger) Warn(ctx context.Context, msg string, data ...interface{}) {
	if z.logLevel < gormlogger.Warn {
		return
	}
	z.logger.Sugar().Warnf(msg, data...)
}

func (z *ZapGormLogger) Error(ctx context.Context, msg string, data ...interface{}) {
	if z.logLevel < gormlogger.Error {
		return
	}
	z.logger.Sugar().Errorf(msg, data...)
}

func (z *ZapGormLogger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	if z.logLevel <= gormlogger.Silent {
		return
	}
	elapsed := time.Since(begin)
	sql, rows := fc()

	// error case
	if err != nil && z.logLevel >= gormlogger.Error {
		z.logger.Error("gorm.query.error",
			zap.Error(err),
			zap.Duration("elapsed", elapsed),
			zap.Int64("rows", rows),
			zap.String("sql", sql),
		)
		return
	}

	// slow query
	if z.slowThreshold != 0 && elapsed > z.slowThreshold && z.logLevel >= gormlogger.Warn {
		z.logger.Warn("gorm.slow_query",
			zap.Duration("elapsed", elapsed),
			zap.Int64("rows", rows),
			zap.String("sql", sql),
		)
		return
	}

	// regular info-level query logging
	if z.logLevel >= gormlogger.Info {
		z.logger.Debug("gorm.query",
			zap.Duration("elapsed", elapsed),
			zap.Int64("rows", rows),
			zap.String("sql", sql),
		)
	}
}
