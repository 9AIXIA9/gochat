package zap

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func Initialize(config *LoggerConfig) error {
	encoder := getEncoder()
	var l = new(zapcore.Level)
	err := l.UnmarshalText([]byte(config.Level))
	if err != nil {
		return err
	}

	lg := zap.New(zapcore.NewCore(encoder, zapcore.Lock(os.Stdout), l), zap.AddCaller())

	//替换全局 logger对象
	zap.ReplaceGlobals(lg)
	return nil
}

// 设置编码器配置
func getEncoder() zapcore.Encoder {
	encoderConfig := zap.NewProductionEncoderConfig()             //创建新的编码器配置
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder         //设置时间格式为ISO 08601
	encoderConfig.TimeKey = "time"                                //时间输出字段名为time
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder       //设置日志级别以大写形式输出“ERROR”
	encoderConfig.EncodeDuration = zapcore.SecondsDurationEncoder //设置持续时间以秒为单位输出
	encoderConfig.EncodeCaller = zapcore.ShortCallerEncoder       //设置调用者信息以简短形式输出（文件名和行号）
	return zapcore.NewJSONEncoder(encoderConfig)
}
