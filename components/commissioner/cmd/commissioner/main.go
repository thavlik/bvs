package main

import (
	"fmt"
	"os"
	"strings"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var atom zap.AtomicLevel = zap.NewAtomicLevel()
var log *zap.Logger

func setLogLevel() {
	logLevel, ok := os.LookupEnv("LOG_LEVEL")
	if ok {
		logLevel = strings.ToLower(logLevel)
		switch logLevel {
		case "info":
			atom.SetLevel(zap.InfoLevel)
		case "debug":
			atom.SetLevel(zap.DebugLevel)
		case "warn":
			atom.SetLevel(zap.WarnLevel)
		case "error":
			atom.SetLevel(zap.ErrorLevel)
		case "dpanic":
			atom.SetLevel(zap.DPanicLevel)
		case "panic":
			atom.SetLevel(zap.PanicLevel)
		case "fatal":
			atom.SetLevel(zap.FatalLevel)
		default:
			panic(fmt.Sprintf("unrecognized log level \"%s\"", logLevel))
		}
	}
}

func init() {
	setLogLevel()
	encoderCfg := zap.NewProductionEncoderConfig()
	log = zap.New(zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderCfg),
		zapcore.Lock(os.Stdout),
		atom,
	))
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		log.Error("main", zap.String("err", err.Error()))
		os.Exit(1)
	}
}
