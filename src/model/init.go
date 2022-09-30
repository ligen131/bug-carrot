package model

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
)

type Model interface {
	Close()
	HomeworkInterface
	FoodInterface
	ScheduleInterface
	KeyWordInterface
	DailyproblemInterface
}

type model struct {
	database *mongo.Database
	context  context.Context
	cancel   context.CancelFunc
}

//GetModel returns a mongo model which will be used in controller
//to help call for the functions that work on mongoDB
func GetModel() Model {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)

	mongoModel := &model{
		database: getMongoDataBase(ctx),
		context:  ctx,
		cancel:   cancel,
	}

	return mongoModel
}
