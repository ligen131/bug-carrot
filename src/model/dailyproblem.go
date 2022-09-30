package model

import (
	"bug-carrot/param"
	"errors"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const dailyproblemCollectionName = "dailyproblem"

//userCollection returns the mongo.collection of users
func (m *model) dailyproblemCollection() *mongo.Collection {
	return m.database.Collection(dailyproblemCollectionName)
}

type DailyproblemInterface interface {
	AddDailyproblem(dailyproblem param.Dailyproblem) error
	GetDailyproblemByDate(id string) (param.Dailyproblem, error)
	UpdateDailyproblem(date string, div int, isProblem bool, link string) (param.Dailyproblem, error)
}

func (m *model) AddDailyproblem(dailyproblem param.Dailyproblem) error {
	filter := bson.M{"date": dailyproblem.Date}
	update := bson.M{"$setOnInsert": dailyproblem}

	boolTrue := true
	opt := options.UpdateOptions{
		Upsert: &boolTrue,
	}

	res, err := m.dailyproblemCollection().UpdateOne(m.context, filter, update, &opt)
	if err != nil {
		return err
	}
	if res.UpsertedCount == 0 {
		return errors.New("not create")
	}

	return nil
}

func (m *model) GetDailyproblemByDate(date string) (param.Dailyproblem, error) {
	filter := bson.M{"date": date}

	var dailyproblem param.Dailyproblem
	err := m.dailyproblemCollection().FindOne(m.context, filter).Decode(&dailyproblem)
	if err != nil {
		return param.Dailyproblem{}, err
	}

	return dailyproblem, nil
}

func (m *model) UpdateDailyproblem(date string, div int, isProblem bool, link string) (param.Dailyproblem, error) {
	filter := bson.M{"date": date}

	replacedDailyproblem, err := m.GetDailyproblemByDate(date)

	dailyproblem := replacedDailyproblem
	dailyproblem.Date = date
	if div == 1 {
		if isProblem {
			dailyproblem.Div1Problem = link
		} else {
			dailyproblem.Div1Solution = link
		}
	} else {
		if isProblem {
			dailyproblem.Div2Problem = link
		} else {
			dailyproblem.Div2Solution = link
		}
	}
	if err != nil {
		return dailyproblem, m.AddDailyproblem(dailyproblem)
	}
	err = m.dailyproblemCollection().FindOneAndReplace(m.context, filter, dailyproblem).Decode(&replacedDailyproblem)
	if err != nil {
		return param.Dailyproblem{}, err
	}

	return replacedDailyproblem, err
}
