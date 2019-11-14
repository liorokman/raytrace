package tuple

type Color struct {
	Tuple
}

func NewColor(r, g, b float64) Color {
	return Color{
		Tuple{[4]float64{r, g, b, 0.0}},
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
		Tuple{[4]float64{c.Red() * r.Red(), c.Green() * r.Green(), c.Blue() * r.Blue(), 0.0}},
	}
}
