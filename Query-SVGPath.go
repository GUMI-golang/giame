package giame

import (
	"io"

	"github.com/go-gl/mathgl/mgl32"
	"github.com/pkg/errors"
	"bytes"
	"strconv"
	"fmt"
	"unicode"
	"github.com/GUMI-golang/gumi/gcore"
	"github.com/GUMI-golang/giame/tools/svg"
)

func svgpath(q InnerQuery, r io.Reader) (err error) {
	rd := svg.NewSVGPathParser(r)
	// command parser
	//
	var stack = newsvgpathStack()
	//
	var (
		bts []byte
		from int
	)
	for bts, from, err =rd.Raw(); err == nil; bts, from, err =rd.Raw() {
		var command svgpathCommand
		if len(bts) >1{
			command = svgpathCommand{
				Fn:bts[0],
				Args:bts[1:],
				From:from,
				To:from + len(bts),
			}
		}else {
			command = svgpathCommand{
				Fn:bts[0],
				From:from,
				To:from + len(bts),
			}
		}
		err = dosvgpath(q, command, stack)
		if err != nil {
			break
		}
	}
	if err == io.EOF{
		return nil
	}
	return err
}

type svgpathStack struct {
	bear *mgl32.Mat2
	last *mgl32.Vec2
	// for smooth quadto
	qcontrol0 *mgl32.Vec2
	keepq bool
	// for smooth cubeto
	ccontrol0 *mgl32.Vec2
	keepc bool
}
func newsvgpathStack() *svgpathStack {
	temp := mgl32.Rotate2D(0)
	return &svgpathStack{
		last: new(mgl32.Vec2),
		bear: &temp,

	}
}
func (s *svgpathStack ) do() {
	if s.keepq{
		s.keepq = false
	}else {
		s.qcontrol0 = nil
	}
	if s.keepc{
		s.keepc = false
	}else {
		s.ccontrol0 = nil
	}
}
type svgpathCommand struct {
	Fn byte
	Args []byte
	From int
	To int
}
func (s svgpathCommand) debugmsg() string {
	return fmt.Sprintf("[%d:%d] %s (Fn : %s, Args : %s)", s.From, s.To, string(append([]byte{s.Fn}, s.Args...)), string(s.Fn), string(s.Args))
}


var (
	ErrorSVGPath = errors.New("SVG")
	ErrorSVGPathParseRead = errors.Wrap(ErrorSVGPath, "Reading fail")
	ErrorSVGPathSyntax = errors.Wrap(ErrorSVGPath, "SyntaxError")
	ErrorSVGPathUnsupported = errors.Wrap(ErrorSVGPath, "Currently Unsupported")
)

func dosvgpath(cq InnerQuery, c svgpathCommand, stack *svgpathStack) error {
	stack.do()
	//
	switch c.Fn {
	case 'M':
		return dosvgpathM(cq, c, stack)
	case 'm':
		return dosvgpathm(cq, c, stack)
	case 'L':
		return dosvgpathL(cq, c, stack)
	case 'l':
		return dosvgpathl(cq, c, stack)
	case 'Q':
		return dosvgpathQ(cq, c, stack)
	case 'q':
		return dosvgpathq(cq, c, stack)
	case 'T':
		return dosvgpathT(cq, c, stack)
	case 't':
		return dosvgpatht(cq, c, stack)
	case 'C':
		return dosvgpathC(cq, c, stack)
	case 'c':
		return dosvgpathc(cq, c, stack)
	case 'S':
		return dosvgpathS(cq, c, stack)
	case 's':
		return dosvgpaths(cq, c, stack)

	//
	case 'B':
		temp, err := strconv.ParseFloat(string(c.Args), 32)
		if err != nil {
			return err
		}
		*stack.bear = mgl32.Rotate2D(float32(gcore.ToRadian(temp)))
	case 'b':
		temp, err := strconv.ParseFloat(string(c.Args), 32)
		if err != nil {
			return err
		}
		*stack.bear = stack.bear.Mul2(mgl32.Rotate2D(float32(gcore.ToRadian(temp))))
	case 'A':
		// TODO
		//args, err := argsForAa(c.Args)
		//if err != nil {
		//	return err
		//}
		//for _, arg := range args {
		//	center := stack.last.Add(arg.to).Mul(0.5)
		//	p := rot90.Mul2(mgl32.Rotate2D(arg.xAxisRotation)).Mul2x1(arg.to.Sub(*stack.last)).Normalize()
		//	r :=
		//	//
		//
		//	//
		//	*stack.last = arg.to
		//}
	case 'a':
		// TODO
	case 'Z':
		fallthrough
	case 'z':
		cq.CloseTo()
		return nil
	}

	return errors.WithMessage(ErrorSVGPathUnsupported, c.debugmsg())
}
// MoveTo
func dosvgpathM(cq InnerQuery, c svgpathCommand, stack *svgpathStack) error {
	vts, err := vectors(c.Args)
	if err != nil {
		return errors.WithMessage(ErrorSVGPathSyntax, err.Error() + " : " + c.debugmsg())
	}
	if len(vts) < 1{
		return errors.WithMessage(ErrorSVGPathSyntax, "Vecter length must be more than 1 : " + c.debugmsg())
	}
	for _, v := range vts {
		cq.MoveTo(v)
		*stack.last = v
	}
	return nil
}
func dosvgpathm(cq InnerQuery, c svgpathCommand, stack *svgpathStack) error {
	vts, err := vectors(c.Args)
	if err != nil {
		return errors.WithMessage(ErrorSVGPathSyntax, err.Error() + " : " + c.debugmsg())
	}
	if len(vts) < 1{
		return errors.WithMessage(ErrorSVGPathSyntax, "Vecter length must be more than 1 : " + c.debugmsg())
	}
	for _, v := range vts {
		p := stack.last.Add(stack.bear.Mul2x1(v))
		cq.MoveTo(p)
		*stack.last = p
	}
	return nil
}
// LineTo
func dosvgpathL(cq InnerQuery, c svgpathCommand, stack *svgpathStack) error {
	vts, err := vectors(c.Args)
	if err != nil {
		return errors.WithMessage(ErrorSVGPathSyntax, err.Error() + " : " + c.debugmsg())
	}
	if len(vts) < 1{
		return errors.WithMessage(ErrorSVGPathSyntax, "Vecter length must be more than 1 : " + c.debugmsg())
	}
	for _, v := range vts {
		cq.LineTo(v)
		*stack.last = v
	}
	return nil
}
func dosvgpathl(cq InnerQuery, c svgpathCommand, stack *svgpathStack) error {
	vts, err := vectors(c.Args)
	if err != nil {
		return errors.WithMessage(ErrorSVGPathSyntax, err.Error() + " : " + c.debugmsg())
	}
	if len(vts) < 1{
		return errors.WithMessage(ErrorSVGPathSyntax, "Vecter length must be more than 1 : " + c.debugmsg())
	}
	for _, v := range vts {
		p := stack.last.Add(stack.bear.Mul2x1(v))
		cq.LineTo(p)
		*stack.last = p
	}
	return nil
}
// QuadTo
func dosvgpathQ(cq InnerQuery, c svgpathCommand, stack *svgpathStack) error {
	vts, err := vectors(c.Args)
	if err != nil {
		return errors.WithMessage(ErrorSVGPathSyntax, err.Error() + " : " + c.debugmsg())
	}
	if len(vts) < 2{
		return errors.WithMessage(ErrorSVGPathSyntax, "Vecter length must be more than 2 : " + c.debugmsg())
	}
	if len(vts) % 2 == 0{
		return errors.WithMessage(ErrorSVGPathSyntax, "Vecter length must be a multiple of 2 : " + c.debugmsg())
	}
	for i := 0; i < len(vts); i += 2 {
		cq.QuadTo(vts[i], vts[i + 1])
	}
	*stack.last = vts[len(vts) - 1]
	*stack.qcontrol0 = vts[len(vts) - 2]
	stack.keepq = true
	return nil
}
func dosvgpathq(cq InnerQuery, c svgpathCommand, stack *svgpathStack) error {
	vts, err := vectors(c.Args)
	if err != nil {
		return errors.WithMessage(ErrorSVGPathSyntax, err.Error() + " : " + c.debugmsg())
	}
	if len(vts) < 2{
		return errors.WithMessage(ErrorSVGPathSyntax, "Vecter length must be more than 2 : " + c.debugmsg())
	}
	if len(vts) % 2 == 0{
		return errors.WithMessage(ErrorSVGPathSyntax, "Vecter length must be a multiple of 2 : " + c.debugmsg())
	}
	for i := 0; i < len(vts); i += 2 {
		p0 := stack.last.Add(stack.bear.Mul2x1(vts[0]))
		to := stack.last.Add(stack.bear.Mul2x1(vts[1]))
		cq.QuadTo(p0, to)
		*stack.last = to
		*stack.qcontrol0 = p0
		stack.keepq = true
	}
	return nil
}
// smooth QuadTo
func dosvgpathT(cq InnerQuery, c svgpathCommand, stack *svgpathStack) error {
	vts, err := vectors(c.Args)
	if err != nil {
		return errors.WithMessage(ErrorSVGPathSyntax, err.Error() + " : " + c.debugmsg())
	}
	if len(vts) < 1{
		return errors.WithMessage(ErrorSVGPathSyntax, "Vecter length must be more than 1 : " + c.debugmsg())
	}
	// first case
	var p0, to mgl32.Vec2
	if stack.qcontrol0 ==nil{
		p0 = *stack.last
	}else {
		p0 = mirrorByPoint(*stack.qcontrol0, *stack.last)
	}
	to = vts[0]
	//
	cq.QuadTo(p0, to)
	*stack.last = to
	*stack.qcontrol0 = p0
	//
	for i := 1; i < len(vts); i += 1 {
		var p0, to mgl32.Vec2
		p0 = mirrorByPoint(*stack.qcontrol0, vts[i - 1])
		to = vts[i]
		//
		cq.QuadTo(p0, to)
		*stack.last = to
		*stack.qcontrol0 = p0
	}
	stack.keepq = true
	return nil
}
func dosvgpatht(cq InnerQuery, c svgpathCommand, stack *svgpathStack) error {
	vts, err := vectors(c.Args)
	if err != nil {
		return errors.WithMessage(ErrorSVGPathSyntax, err.Error() + " : " + c.debugmsg())
	}
	if len(vts) < 1{
		return errors.WithMessage(ErrorSVGPathSyntax, "Vecter length must be more than 1 : " + c.debugmsg())
	}
	// first case
	var p0, to mgl32.Vec2
	if stack.qcontrol0 ==nil{
		p0 = *stack.last
	}else {
		p0 = mirrorByPoint(*stack.qcontrol0, *stack.last)
	}
	prev := stack.last.Add(stack.bear.Mul2x1(vts[0]))
	to = prev
	//
	cq.QuadTo(p0, to)
	*stack.last = to
	*stack.qcontrol0 = p0
	//
	for i := 1; i < len(vts); i += 1 {
		var p0, to mgl32.Vec2
		p0 = mirrorByPoint(*stack.qcontrol0, prev)
		to = stack.last.Add(stack.bear.Mul2x1(vts[i]))
		//
		cq.QuadTo(p0, to)
		*stack.last = to
		*stack.qcontrol0 = p0
		prev = to
	}
	stack.keepq = true
	return nil
}
// CubeTo
func dosvgpathC(cq InnerQuery, c svgpathCommand, stack *svgpathStack) error {
	vts, err := vectors(c.Args)
	if err != nil {
		return errors.WithMessage(ErrorSVGPathSyntax, err.Error() + " : " + c.debugmsg())
	}
	if len(vts) < 3{
		return errors.WithMessage(ErrorSVGPathSyntax, "Vecter length must be more than 3 : " + c.debugmsg())
	}
	if len(vts) % 3 == 0{
		return errors.WithMessage(ErrorSVGPathSyntax, "Vecter length must be a multiple of 3 : " + c.debugmsg())
	}
	for i := 0; i < len(vts); i += 3 {
		cq.CubeTo(vts[i], vts[i + 1], vts[i + 2])
	}
	*stack.ccontrol0 = vts[len(vts) - 2]
	*stack.last = vts[len(vts) - 1]
	stack.keepc = true
	return nil
}
func dosvgpathc(cq InnerQuery, c svgpathCommand, stack *svgpathStack) error {
	vts, err := vectors(c.Args)
	if err != nil {
		return errors.WithMessage(ErrorSVGPathSyntax, err.Error() + " : " + c.debugmsg())
	}
	if len(vts) < 3{
		return errors.WithMessage(ErrorSVGPathSyntax, "Vecter length must be more than 3 : " + c.debugmsg())
	}
	if len(vts) % 3 == 0{
		return errors.WithMessage(ErrorSVGPathSyntax, "Vecter length must be a multiple of 3 : " + c.debugmsg())
	}
	for i := 0; i < len(vts); i += 3 {
		p0 := stack.last.Add(stack.bear.Mul2x1(vts[i]))
		p1 := stack.last.Add(stack.bear.Mul2x1(vts[i + 1]))
		to := stack.last.Add(stack.bear.Mul2x1(vts[i + 2]))
		cq.CubeTo(p0, p1, to)

		*stack.ccontrol0 = p1
		*stack.last = to
		stack.keepc = true
	}
	return nil
}
// smooth CubeTo
func dosvgpathS(cq InnerQuery, c svgpathCommand, stack *svgpathStack) error {
	vts, err := vectors(c.Args)
	if err != nil {
		return errors.WithMessage(ErrorSVGPathSyntax, err.Error() + " : " + c.debugmsg())
	}
	if len(vts) < 2{
		return errors.WithMessage(ErrorSVGPathSyntax, "Vecter length must be more than 2 : " + c.debugmsg())
	}
	if len(vts) % 2 == 0{
		return errors.WithMessage(ErrorSVGPathSyntax, "Vecter length must be a multiple of 2 : " + c.debugmsg())
	}
	// first case
	var p0, p1, to mgl32.Vec2
	if stack.ccontrol0 ==nil{
		p0 = *stack.last
	}else {
		p0 = mirrorByPoint(*stack.ccontrol0, *stack.last)
	}
	p1 = vts[0]
	to = vts[1]
	//
	cq.CubeTo(p0, p1, to)
	*stack.last = to
	*stack.ccontrol0 = p1
	//
	for i := 2; i < len(vts); i += 2 {
		var p0, to mgl32.Vec2
		p0 = mirrorByPoint(*stack.ccontrol0, vts[i-2])
		p1 = vts[i]
		to = vts[i + 1]
		//
		cq.CubeTo(p0, p1, to)
		*stack.last = to
		*stack.ccontrol0 = p1
	}
	stack.keepc = true
	return nil
}
func dosvgpaths(cq InnerQuery, c svgpathCommand, stack *svgpathStack) error {
	vts, err := vectors(c.Args)
	if err != nil {
		return errors.WithMessage(ErrorSVGPathSyntax, err.Error() + " : " + c.debugmsg())
	}
	if len(vts) < 2{
		return errors.WithMessage(ErrorSVGPathSyntax, "Vecter length must be more than 2 : " + c.debugmsg())
	}
	if len(vts) % 2 == 0{
		return errors.WithMessage(ErrorSVGPathSyntax, "Vecter length must be a multiple of 2 : " + c.debugmsg())
	}
	// first case
	var p0, p1, to mgl32.Vec2
	if stack.ccontrol0 ==nil{
		p0 = *stack.last
	}else {
		p0 = mirrorByPoint(*stack.ccontrol0, *stack.last)
	}
	p1 = stack.last.Add(stack.bear.Mul2x1(vts[0]))
	to = stack.last.Add(stack.bear.Mul2x1(vts[1]))
	//
	cq.CubeTo(p0, p1, to)
	*stack.last = to
	*stack.ccontrol0 = p1
	//
	for i := 2; i < len(vts); i += 2 {
		var p0, to mgl32.Vec2
		p0 = mirrorByPoint(*stack.ccontrol0, stack.last.Add(stack.bear.Mul2x1(vts[i - 2])))
		p1 = stack.last.Add(stack.bear.Mul2x1(vts[i]))
		to = stack.last.Add(stack.bear.Mul2x1(vts[i + 1]))
		//
		cq.CubeTo(p0, p1, to)
		*stack.last = to
		*stack.ccontrol0 = p1
	}
	stack.keepc = true
	return nil
}

func mirrorByPoint(a, mirror mgl32.Vec2) mgl32.Vec2 {
	return mirror.Add(mirror.Sub(a))
}


func vectors(bts []byte) (res []mgl32.Vec2, err error) {
	bts = bytes.TrimSpace(bts)
	var x *float32
	var from = 0
	for to, b := range bts {
		var temp float64
		if unicode.IsSpace(rune(b)) || b == ','{
			temp, err = strconv.ParseFloat(string(bts[from:to]), 32)
			if err != nil {
				return nil, err
			}
			if x == nil{
				temp32 := float32(temp)
				x = &temp32
			}else {
				res = append(res, mgl32.Vec2{*x, float32(temp)})
				x = nil
			}
			from = to + 1
		}
	}
	temp, err := strconv.ParseFloat(string(bts[from:]), 32)
	if err != nil {
		return nil, err
	}
	if x == nil{
		temp32 := float32(temp)
		x = &temp32
	}else {
		res = append(res, mgl32.Vec2{*x, float32(temp)})
		x = nil
	}
	if x != nil{
		return nil, errors.New("There is remain float")
	}
	return
}
func floats(bts []byte) (res []float32, err error) {
	bts = bytes.TrimSpace(bts)
	var from = 0
	for to, b := range bts {
		var temp float64
		if unicode.IsSpace(rune(b)) || b == ','{
			temp, err = strconv.ParseFloat(string(bts[from:to]), 32)
			if err != nil {
				return nil, err
			}
			res = append(res, float32(temp))
			from = to + 1
		}
	}
	temp, err := strconv.ParseFloat(string(bts[from:]), 32)
	if err != nil {
		return nil, err
	}
	res = append(res, float32(temp))
	return
}

type arcArg struct {
	rad mgl32.Vec2
	xAxisRotation float32
	largeArcFlag, sweepFlag bool
	to mgl32.Vec2
}
func argsForAa(bts []byte) (res []arcArg , err error) {
	bts = bytes.Replace(bytes.TrimSpace(bts), []byte(","), []byte(" "), -1)
	temp := bytes.Split(bts, []byte(" "))
	if len(temp) % 7 != 0 && len(temp) < 7{
		return nil, errors.New("invalid")
	}
	for i := 0; i < len(temp); i += 7 {
		rx, err := strconv.ParseFloat(string(temp[i + 0]), 32)
		if err != nil {
			return nil, err
		}
		ry, err := strconv.ParseFloat(string(temp[i + 1]), 32)
		if err != nil {
			return nil, err
		}

		xar, err := strconv.ParseFloat(string(temp[i + 2]), 32)
		if err != nil {
			return nil, err
		}
		//
		laf, err := strconv.ParseInt(string(temp[i + 3]), 10, 32)
		if err != nil {
			return nil, err
		}
		sf, err := strconv.ParseInt(string(temp[i + 4]), 10, 32)
		if err != nil {
			return nil, err
		}

		x, err := strconv.ParseFloat(string(temp[i + 5]), 32)
		if err != nil {
			return nil, err
		}
		y, err := strconv.ParseFloat(string(temp[i + 6]), 32)
		if err != nil {
			return nil, err
		}
		res = append(res, arcArg{
			rad:mgl32.Vec2{float32(rx), float32(ry)},
			xAxisRotation: float32(gcore.ToRadian(xar)),
			largeArcFlag: laf == 1,
			sweepFlag: sf == 1,
			to:mgl32.Vec2{float32(x), float32(y)},
		})
	}
	return
}