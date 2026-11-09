package tty

type Rect struct {
	X      int
	Y      int
	Width  int
	Height int
}

func (rect Rect) IsEmpty() bool {
	return rect.Width <= 0 || rect.Height <= 0
}

func (rect Rect) Contains(x int, y int) bool {
	return x >= rect.X && x < rect.X+rect.Width && y >= rect.Y && y < rect.Y+rect.Height
}

func (rect Rect) Intersect(other Rect) Rect {

	x1 := max(rect.X, other.X)
	y1 := max(rect.Y, other.Y)
	x2 := min(rect.X+rect.Width, other.X+other.Width)
	y2 := min(rect.Y+rect.Height, other.Y+other.Height)

	if x2 <= x1 || y2 <= y1 {
		return Rect{}
	} else {
		return Rect{
			X:      x1,
			Y:      y1,
			Width:  x2 - x1,
			Height: y2 - y1,
		}
	}

}
