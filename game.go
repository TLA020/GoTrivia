package GoTrivia

import (
	"time"
)

type Game struct {
	Id int64
	Question *Question
	StartTime time.Time
	EndTime time.Time
}

type Question struct {
	question string
	answer string
	hint string
}

