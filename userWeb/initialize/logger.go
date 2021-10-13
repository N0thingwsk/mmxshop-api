package initialize

import "go.uber.org/zap"

func InitLogger() {
	logger, err := zap.NewDevelopment()
	if err != nil {
		panic("Zap logger error" + err.Error())
	}
	zap.ReplaceGlobals(logger)
}
