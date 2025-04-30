package painter

import (
	"image"
	"image/color"

	"golang.org/x/exp/shiny/screen"
)

// Operation змінює вхідну текстуру.
type Operation interface {
	// Do виконує зміну операції, повертаючи true, якщо текстура вважається готовою для відображення.
	Do(t screen.Texture) (ready bool)
}

// OperationList групує список операції в одну.
type OperationList []Operation

func (ol OperationList) Do(t screen.Texture) (ready bool) {
	for _, o := range ol {
		ready = o.Do(t) || ready
	}
	return
}

// UpdateOp операція, яка не змінює текстуру, але сигналізує, що текстуру потрібно розглядати як готову.
var UpdateOp = updateOp{}

type updateOp struct{}

func (op updateOp) Do(t screen.Texture) bool { return true }

// OperationFunc використовується для перетворення функції оновлення текстури в Operation.
type OperationFunc func(t screen.Texture) bool

func (f OperationFunc) Do(t screen.Texture) bool {
	return f(t)
}

func NewWhiteOp(loop *Loop) Operation {
	return OperationFunc(func(t screen.Texture) bool {
		loop.bgColor = color.White
		return true
	})
}

// NewGreenOp - операція, що фарбує фон у зелений колір.
func NewGreenOp(loop *Loop) Operation {
	return OperationFunc(func(t screen.Texture) bool {
		loop.bgColor = color.RGBA{G: 0xFF, A: 0xFF}
		return true
	})
}

// NewBgRectOp - операція для малювання чорного прямокутника на фоні.
func NewBgRectOp(loop *Loop, x1, y1, x2, y2 float64) Operation {
	return OperationFunc(func(t screen.Texture) bool {
		width := t.Bounds().Dx()
		height := t.Bounds().Dy()

		r := image.Rect(
			int(float64(width)*x1),
			int(float64(height)*y1),
			int(float64(width)*x2),
			int(float64(height)*y2),
		)
		loop.rect = &r
		return true
	})
}

// NewFigureOp - операція для додавання нової фігури на фон.
func NewFigureOp(loop *Loop, x, y float64) Operation {
	return OperationFunc(func(t screen.Texture) bool {
		width := t.Bounds().Dx()
		height := t.Bounds().Dy()

		point := image.Point{
			X: int(float64(width) * x),
			Y: int(float64(height) * y),
		}
		loop.figures = append(loop.figures, point)
		return true
	})
}

// NewMoveOp - операція для переміщення всіх фігур на нові координати.
func NewMoveOp(loop *Loop, x, y float64) Operation {
	return OperationFunc(func(t screen.Texture) bool {
		width := t.Bounds().Dx()
		height := t.Bounds().Dy()

		newPoint := image.Point{
			X: int(float64(width) * x),
			Y: int(float64(height) * y),
		}
		for i := range loop.figures {
			loop.figures[i] = newPoint
		}
		return true
	})
}

// NewResetOp - операція для скидання текстури до початкового стану.
func NewResetOp(loop *Loop) Operation {
	return OperationFunc(func(t screen.Texture) bool {
		loop.bgColor = color.Black
		loop.rect = nil
		loop.figures = nil
		t.Fill(t.Bounds(), color.Black, screen.Src)
		return true
	})
}

// NewUpdateOp - операція для оновлення текстури на основі поточного стану.
func NewUpdateOp(loop *Loop) Operation {
	return OperationFunc(func(t screen.Texture) bool {
		t.Fill(t.Bounds(), loop.bgColor, screen.Src)

		if loop.rect != nil {
			t.Fill(*loop.rect, color.Black, screen.Src)
		}

		for _, point := range loop.figures {
			drawFigure(t, point)
		}

		return true
	})
}

func drawFigure(t screen.Texture, center image.Point) {
	width := t.Bounds().Dx()
	height := t.Bounds().Dy()

	stemWidth := width * 150 / 800
	stemHeight := height * 50 / 800
	headWidth := width * 50 / 800
	headHeight := height * 200 / 800

	yellow := color.RGBA{R: 0xFF, G: 0xFF, A: 0xFF}

	stemRect := image.Rect(
		center.X-stemWidth/2,
		center.Y-stemHeight/2,
		center.X+stemWidth/2,
		center.Y+stemHeight/2,
	)
	t.Fill(stemRect, yellow, screen.Src)

	headRect := image.Rect(
		center.X+stemWidth/2-headWidth/2,
		center.Y-headHeight/2,
		center.X+stemWidth/2+headWidth/2,
		center.Y+headHeight/2,
	)
	t.Fill(headRect, yellow, screen.Src)
}
