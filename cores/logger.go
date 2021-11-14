package cores

import (
	"apioak-admin/app/packages"
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"io"
	"strings"
	"time"
)

func InitLogger(conf *ConfigGlobal) error {
	confLogger := conf.Logger

	infoLevel := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl < zapcore.WarnLevel
	})

	errorLevel := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl >= zapcore.WarnLevel
	})

	infoWriter := GetWriter(conf, confLogger.LogFileInfo)
	errorWriter := GetWriter(conf, confLogger.LogFileError)

	core := zapcore.NewTee(
		zapcore.NewCore(GetEncoder(), zapcore.AddSync(infoWriter), infoLevel),
		zapcore.NewCore(GetEncoder(), zapcore.AddSync(errorWriter), errorLevel),
	)

	log := zap.New(core, zap.AddCaller())
	packages.SetLogger(log.Sugar())

	return nil
}

func GetWriter(conf *ConfigGlobal, filename string) io.Writer {
	var confLogger = conf.Logger

	logReserve := time.Duration(confLogger.LogReserve)
	logReserve = time.Hour * 24 * logReserve

	hook, err := rotatelogs.New(
		strings.Replace(confLogger.LogPath+"/"+filename, ".log", "", -1)+"-%Y%m%d.log",
		rotatelogs.WithMaxAge(logReserve),
	)

	if err != nil {
		panic(err)
	}
	return hook
}

func GetEncoder() zapcore.Encoder {
	return zapcore.NewJSONEncoder(zapcore.EncoderConfig{
		MessageKey:  "message",
		LevelKey:    "level",
		EncodeLevel: zapcore.CapitalLevelEncoder,
		TimeKey:     "time",
		EncodeTime: func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
			enc.AppendString(t.Format("2006-01-02 15:04:05"))
		},
		CallerKey:    "file",
		EncodeCaller: zapcore.ShortCallerEncoder,
		EncodeDuration: func(d time.Duration, enc zapcore.PrimitiveArrayEncoder) {
			enc.AppendInt64(int64(d) / 1000000)
		},
	})
}
