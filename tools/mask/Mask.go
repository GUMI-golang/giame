package mask

// Mask3 must larger than 1x1 matrix
// Mask3 must n x m size matrix
// Mask3 unit sum must be
// Mask3 width and height must be odd number
type Mask3 struct {
	Data [3][3]float64
	Name string
}
func (s Mask3) DataF32() ([3][3]float32) {
	return [3][3]float32{
		{float32(s.Data[0][0]), float32(s.Data[0][1]), float32(s.Data[0][2])},
		{float32(s.Data[1][0]), float32(s.Data[1][1]), float32(s.Data[1][2])},
		{float32(s.Data[2][0]), float32(s.Data[2][1]), float32(s.Data[2][2])},
	}
}
func (s Mask3) Size() (w, h int) {
	return 3,3
}
func (s Mask3) Strings() (string) {
	return s.Name
}