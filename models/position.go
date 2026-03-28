package models

type Position string

const (
	PositionStriker     Position = "striker"
	PositionMidfielder  Position = "midfielder"
	PositionDefender    Position = "defender"
	PositionGoalkeeper  Position = "goalkeeper"
)

var ValidPositions = map[Position]struct{}{
	PositionStriker:    {},
	PositionMidfielder: {},
	PositionDefender:   {},
	PositionGoalkeeper: {},
}

func IsValidPosition(p Position) bool {
	_, ok := ValidPositions[p]
	return ok
}
