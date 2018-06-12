package giame

import (
	"fmt"
	"github.com/GUMI-golang/gumi/gcore"
	"golang.org/x/image/font/gofont/goregular"
	"golang.org/x/image/font"
	"golang.org/x/image/math/fixed"
	"github.com/go-gl/mathgl/mgl32"
	"image"
	"github.com/golang/freetype/truetype"
	"container/list"
)


var DefaultVFont *VectorFont= NewVectorFont(gcore.MustValue(truetype.Parse(goregular.TTF)).(*truetype.Font), 22, font.HintingFull)
type FontName struct {
	Family, Name string
	FullName     string
}

func (s FontName) String() string {
	return fmt.Sprintf("VectorFont(%s, family : %s, name : %s)", s.FullName, s.Family, s.Name)
}


type VectorFont struct {
	f    *truetype.Font
	n    FontName
	size fixed.Int26_6
	hint font.Hinting
	//
	cachemap map[truetype.Index] *list.Element
	cache *list.List
}
type cachehelp struct {
	i truetype.Index
	b *truetype.GlyphBuf
}


func NewVectorFont(f *truetype.Font, size int, hint font.Hinting) *VectorFont {
	return &VectorFont{
		n: FontName{
			Name:     f.Name(truetype.NameIDFontSubfamily),
			Family:   f.Name(truetype.NameIDFontFamily),
			FullName: f.Name(truetype.NameIDFontFullName),
		},
		f:    f,
		size: fixed.I(size),
		hint: hint,
		//
		cachemap: make(map[truetype.Index] *list.Element),
		cache: list.New(),
	}
}
func (s *VectorFont) Name() FontName {
	return s.n
}
func (s *VectorFont) Size() int {
	return s.size.Round()
}
func (s *VectorFont) SetSize(size int) {
	s.size = fixed.I(size)
}
func (s *VectorFont) Hint() font.Hinting {
	return s.hint
}
func (s *VectorFont) SetHint(hint font.Hinting) {
	s.hint = hint
}
//
func (s *VectorFont) Text(cq InnerQuery, text string, point mgl32.Vec2, align gcore.Align)  {
	s.drawText(
		cq,
		text,
		alignHelp(align, point, s.MeasureText(text)),
	)
	return
}
func (s *VectorFont) TextInRect(cq InnerQuery, text string, rect image.Rectangle, align gcore.Align) {
	v, h := gcore.SplitAlign(align)

	var pt mgl32.Vec2
	switch v {
	case gcore.AlignTop:
		pt[1] = float32(rect.Min.Y)
	case gcore.AlignVertical:
		pt[1] = float32(rect.Min.Y) + float32(rect.Dy())/ 2
	case gcore.AlignBottom:
		pt[1] = float32(rect.Min.Y) + float32(rect.Dy())

	}
	switch h {
	case gcore.AlignLeft:
		pt[0] = float32(rect.Min.X)
	case gcore.AlignHorizontal:
		pt[0] = float32(rect.Min.X) + float32(rect.Dx())/ 2
	case gcore.AlignRight:
		pt[0] = float32(rect.Min.X) + float32(rect.Dx())
	}
	s.Text(cq, text, pt, align)
}
func (s *VectorFont) MeasureText(text string) (res mgl32.Vec2) {
	var previdx truetype.Index = 0
	for i, r := range []rune(text) {
		idx := s.f.Index(r)
		if i > 0 {
			res[0] += Fint32ToFloat32(s.f.Kern( s.size, previdx, idx))
		}

		hmat := s.f.HMetric(s.size, idx)
		res[0] += Fint32ToFloat32(hmat.AdvanceWidth)
		previdx = idx
	}
	res[1] = s.MeasureHeight()
	return res
}
func (s *VectorFont) MeasureHeight() (float32) {
	return Fint32ToFloat32(s.size)
}
//
func (s *VectorFont) drawText(cq InnerQuery, text string, point mgl32.Vec2) {
	var prevIdx truetype.Index = 0
	for i, r := range []rune(text) {
		idx := s.f.Index(r)
		b, err := s.load(idx)
		if err != nil {
			continue
		}
		if i > 0 {
			point[0] += Fint32ToFloat32(s.f.Kern(s.size, prevIdx, idx))
		}
		temp := point
		temp[1] += Fint32ToFloat32(b.Bounds.Min.Y)
		fontRaster(cq, b, temp)
		point[0] += Fint32ToFloat32(b.AdvanceWidth)
		prevIdx = idx
	}
}
const cachesize = 256
func (s *VectorFont) load(i truetype.Index) (*truetype.GlyphBuf, error) {
	// Current
	if c, ok := s.cachemap[i];ok{
		s.cache.MoveToFront(c)
		return c.Value.(cachehelp).b, nil
	}
	buf := &truetype.GlyphBuf{}
	err := buf.Load(s.f, s.size, i, s.hint)
	if err != nil {
		return nil, err
	}
	//
	if s.cache.Len() > cachesize{
		delete(s.cachemap, s.cache.Back().Value.(cachehelp).i)
		b := s.cache.Back()
		b.Value = cachehelp{
			i:i,
			b:buf,
		}
		s.cachemap[i] = b
	}else {
		s.cachemap[i] = s.cache.PushBack(cachehelp{
			i:i,
			b:buf,
		})
	}
	return buf, nil
}

func fontRaster(cq InnerQuery, buf *truetype.GlyphBuf, point mgl32.Vec2) {
	var start int
	for _, end := range buf.Ends {
		fontContour(cq, buf.Points[start:end], point)
		start = end
	}
}
func Fint32ToFloat32(i fixed.Int26_6) float32 {
	return float32(i) / float32(0x40)
}
func Float32ToFint32(f float32) fixed.Int26_6 {
	return fixed.Int26_6(f * 0x40)
}
func fontContour(cq InnerQuery, points []truetype.Point, point mgl32.Vec2) {
	if len(points) == 0 {
		return
	}
	var first mgl32.Vec2
	var ifirst, ilast = 0, len(points)
	if points[0].Flags&0x01 != 0 {
		ifirst = 1
		first = mgl32.Vec2{
			point[0] + Fint32ToFloat32(points[0].X),
			point[1] - Fint32ToFloat32(points[0].Y),
		}
	} else {
		last := mgl32.Vec2{
			point[0] + Fint32ToFloat32(points[ilast-1].X),
			point[1] - Fint32ToFloat32(points[ilast-1].Y),
		}
		if points[ilast-1].Flags&0x01 != 0 {
			first = last
			ilast = ilast - 1
		} else {
			first = mgl32.Vec2{
				(first.X() + last.X()) / 2,
				(first.Y() + last.Y()) / 2,
			}
		}
	}
	//==================================
	// drawloop
	// start point
	cq.MoveTo(first)
	var q0, q0on = first, true
	for i := ifirst; i < ilast; i++ {

		p := points[i]

		var q, qon = mgl32.Vec2{
			point.X() + Fint32ToFloat32(p.X),
			point.Y() - Fint32ToFloat32(p.Y),
		}, p.Flags&0x01 != 0
		if qon {
			if q0on {
				cq.LineTo(q)
			} else {
				cq.QuadTo(q0, q)
			}
		} else {
			if !q0on {
				cq.QuadTo(q0, q0.Add(q).Mul(0.5))
			}
		}
		q0, q0on = q, qon
	}
	if q0on {
	} else {
		cq.QuadTo(q0, first)
	}
	cq.CloseTo()

}
func alignHelp(align gcore.Align, point, size mgl32.Vec2) (res mgl32.Vec2) {
	v, h := gcore.SplitAlign(align)
	switch v {
	case gcore.AlignTop:
		res[1] = point[1] + size[1]
	case gcore.AlignVertical:
		res[1] = point[1] + size[1] / 2
	case gcore.AlignBottom:
		res[1] = point[1]
	}
	switch h {
	case gcore.AlignLeft:
		res[0] = point[0]
	case gcore.AlignHorizontal:
		res[0] = point[0] - size[0] / 2
	case gcore.AlignRight:
		res[0] = point[0] - size[0]
	}
	return
}