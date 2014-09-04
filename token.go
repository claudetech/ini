package ini

type token interface{}

type newLineToken struct{}

type spaceToken struct {
	value string
}

type sepToken struct{}

type commentToken struct{}

type symbolToken struct {
	value string
}

type otherToken struct {
	value string
}
