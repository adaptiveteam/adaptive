package mapping

import (
	"github.com/jinzhu/gorm"
	"strconv"
	logger2 "github.com/adaptiveteam/adaptive/adaptive-utils-go/logger"
	"github.com/adaptiveteam/adaptive/lambdas/reporting-transformed-model-streaming-lambda/model"
)

func intToBoolean(i int) (op bool) {
	if i == 1 {
		op = true
	}
	return
}

func stringToFloat(s string) (op float64) {
	op, _ = strconv.ParseFloat(s, 32)
	return
}

func stringToInt(s string) int {
	i, _ := strconv.Atoi(s)
	return i
}

// DBEntity is an interface of a type class assiciated with each entity. 
type DBEntity = interface {
	TableName() string
	ParseUnsafe(js []byte, logger logger2.AdaptiveLogger) interface{}
	HandleStreamEntityUnsafe(e2 model.StreamEntity, conn *gorm.DB, logger logger2.AdaptiveLogger)
}
