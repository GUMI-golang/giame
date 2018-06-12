package main

import (
	"github.com/GUMI-golang/giame/tools/svg"
	"strings"
	"fmt"
	"io"
)

func main() {
	//var a = []int{0,1,2,3,4}
	//fmt.Println(a[:5])
	//fmt.Println(a[6:])
	p := svg.NewSVGPathParser(strings.NewReader("M0,0 L16,64 L32,0 Z"))
	var (
		//v []byte
		//f int
		e error
		pd svg.PathData
	)
	//for v, f, e = p.Raw(); e == nil; v, f, e = p.Raw() {
	//	fmt.Println(string(v), ":" , f, "-", len(v))
	//}
	for pd, e = p.Next(); e == nil; pd, e = p.Next() {
		fmt.Println(pd)
	}
	if e != nil && e != io.EOF{
		panic(e)
	}
}