package painter

import (
	"image"
	"image/color"
	"image/draw"
	"reflect"
	"testing"
	"time"

	"golang.org/x/exp/shiny/screen"
)

func TestLoop_Post(t *testing.T) {
	var l Loop
	var tr testReceiver
	l.Receiver = &tr

	var testOps []string

	l.Start(mockScreen{})

	// Додаємо операції вручну через OperationFunc
	l.Post(OperationFunc(func(t screen.Texture) bool {
		testOps = append(testOps, "white")
		return true
	}))
	l.Post(OperationFunc(func(t screen.Texture) bool {
		testOps = append(testOps, "bgrect")
		return true
	}))
	l.Post(OperationFunc(func(t screen.Texture) bool {
		testOps = append(testOps, "figure")
		return true
	}))
	l.Post(OperationFunc(func(t screen.Texture) bool {
		testOps = append(testOps, "green")
		return true
	}))
	l.Post(OperationFunc(func(t screen.Texture) bool {
		testOps = append(testOps, "move")
		return true
	}))

	// Тепер додаємо операцію update в testOps
	l.Post(OperationFunc(func(t screen.Texture) bool {
		testOps = append(testOps, "update")
		return true
	}))

	l.Post(OperationFunc(func(t screen.Texture) bool {
		testOps = append(testOps, "reset")
		return true
	}))

	// Затримка для асинхронного виконання
	time.Sleep(100 * time.Millisecond)
	l.StopAndWait()

	// Очікувані операції
	expected := []string{"white", "bgrect", "figure", "green", "move", "update", "reset"}
	if !reflect.DeepEqual(testOps, expected) {
		t.Errorf("expected %v, got %v", expected, testOps)
	}
}

type testReceiver struct {
	lastTexture screen.Texture
}

func (tr *testReceiver) Update(t screen.Texture) {
	tr.lastTexture = t
}

type mockScreen struct{}

func (m mockScreen) NewBuffer(size image.Point) (screen.Buffer, error) {
	return nil, nil
}

func (m mockScreen) NewTexture(size image.Point) (screen.Texture, error) {
	return &mockTexture{}, nil
}

func (m mockScreen) NewWindow(opts *screen.NewWindowOptions) (screen.Window, error) {
	return nil, nil
}

type mockTexture struct {
	FilledRects []image.Rectangle
	Colors      []color.Color
}

func (m *mockTexture) Fill(dr image.Rectangle, src color.Color, op draw.Op) {
	m.FilledRects = append(m.FilledRects, dr)
	m.Colors = append(m.Colors, src)
}

func (m *mockTexture) Bounds() image.Rectangle {
	return image.Rect(0, 0, 400, 400)
}

func (m *mockTexture) Size() image.Point {
	return image.Point{X: 400, Y: 400}
}

func (m *mockTexture) Upload(dp image.Point, src screen.Buffer, sr image.Rectangle) {}

func (m *mockTexture) Release() {}
