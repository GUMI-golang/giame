package svg

import "fmt"

type SVGError struct {
	Where [2]int
	What string
	Why string
}

func (s *SVGError) Error() string {
	return fmt.Sprintf("SVGError : %s, [%d:%d] - %s", s.What, s.Where[0], s.Where[1], s.Why)
}

