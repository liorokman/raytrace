package tuple

var (
	Black = NewColor(0, 0, 0)
	White = NewColor(1, 1, 1)
	Red   = NewColor(1, 0, 0)
	Green = NewColor(0, 1, 0)
	Blue  = NewColor(0, 0, 1)
)

type Color struct {
	Tuple
}

func NewColor(r, g, b float64) Color {
	return Color{
		Tuple{r, g, b, 0.0},
	}
}

func (c Color) Red() float64 {
	return c.X()
}

func (c Color) Green() float64 {
	return c.Y()
}

func (c Color) Blue() float64 {
	return c.Z()
}

func (c Color) Add(r Color) Color {
	return Color{
		c.Tuple.Add(r.Tuple),
	}
}

func (c Color) Subtract(r Color) Color {
	return Color{
		c.Tuple.Subtract(r.Tuple),
	}
}

func (c Color) Mult(r float64) Color {
	return Color{
		c.Tuple.Mult(r),
	}
}

// Technically the Shur product or Hadamard product
func (c Color) MultColor(r Color) Color {
	return Color{
		Tuple{c.Red() * r.Red(), c.Green() * r.Green(), c.Blue() * r.Blue(), 0.0},
	}
}

func (c Color) Equals(r Color) bool {
	return c.Tuple.Equals(r.Tuple)
}
