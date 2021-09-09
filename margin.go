package reportengine

type Margin struct {
	top    float64
	bottom float64
	left   float64
	right  float64
}

func NewMargin(margin float64) Margin {
	return Margin{margin, margin, margin, margin}
}
func NewVerticalMargin(margin float64) Margin {
	return Margin{top: margin, bottom: margin}
}
