package param

import "time"

type Homework struct {
	Subject    string    `bson:"subject"`
	Context    string    `bson:"context"`
	HandInTime time.Time `bson:"hand_in_time"`
}
