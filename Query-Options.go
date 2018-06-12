package giame


type StrokeJoin uint8
const (
	StrokeJoinBevel StrokeJoin = iota
	StrokeJoinRound StrokeJoin = iota
	StrokeJoinMiter StrokeJoin = iota
)

type StrokeCap uint8
const (
	StrokeCapButt   StrokeCap = iota
	StrokeCapRound  StrokeCap = iota
	StrokeCapSqaure StrokeCap = iota
)
