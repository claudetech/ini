package ini

type tokenType int

const (
	newLineTokType tokenType = iota
	spaceTokType
	sepTokType
	commentTokType
	symbolTokType
	otherTokType
)

func (t tokenType) ToString() string {
	switch t {
	case newLineTokType:
		return "newline"
	case spaceTokType:
		return "space"
	case sepTokType:
		return "separator"
	case commentTokType:
		return "comment"
	case symbolTokType:
		return "symbol"
	default:
		return "normal char"
	}
}

type token interface {
	getType() tokenType
}

type newLineToken struct{}

func (t *newLineToken) getType() tokenType {
	return newLineTokType
}

type spaceToken struct {
	value string
}

func (t *spaceToken) getType() tokenType {
	return spaceTokType
}

type sepToken struct{}

func (t *sepToken) getType() tokenType {
	return sepTokType
}

type commentToken struct{}

func (t *commentToken) getType() tokenType {
	return commentTokType
}

type symbolToken struct {
	symbol string
}

func (t *symbolToken) getType() tokenType {
	return symbolTokType
}

type otherToken struct {
	value string
}

func (t *otherToken) getType() tokenType {
	return otherTokType
}
