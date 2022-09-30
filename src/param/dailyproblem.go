package param

type Dailyproblem struct {
	Date         string `bson:"date"` // format: 20060101 yyyymmdd
	Div1Problem  string `bson:"div1_problem,omitempty"`
	Div1Solution string `bson:"div1_solution,omitempty"`
	Div2Problem  string `bson:"div2_problem,omitempty"`
	Div2Solution string `bson:"div2_solution,omitempty"`
}
