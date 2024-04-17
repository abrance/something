package log

//
//import (
//	"go.uber.org/zap"
//	"go.uber.org/zap/zapcore"
//)
//
//func GetLogger() *zap.Logger {
//	config := DefaultZapLoggerConfig
//	config.Encoding = "console"
//	config.EncoderConfig.EncodeTime = zapcore.RFC3339TimeEncoder
//	lg, err := config.Build()
//	if err != nil {
//		cobrautl.ExitWithError(cobrautl.ExitBadArgs, err)
//	}
//	return lg
//}
