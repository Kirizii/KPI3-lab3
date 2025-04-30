package painter

import (
	"image"
	"image/color"
	"sync"
	"time"

	"golang.org/x/exp/shiny/screen"
)

// Receiver отримує текстуру, яка була підготовлена в результаті виконання команд у циклі подій.
type Receiver interface {
	Update(t screen.Texture)
}

// Loop реалізує цикл подій для формування текстури отриманої через виконання операцій отриманих з внутрішньої черги.
type Loop struct {
	Receiver Receiver
	bgColor  color.Color
	rect     *image.Rectangle
	figures  []image.Point

	next screen.Texture // текстура, яка зараз формується
	prev screen.Texture // текстура, яка була відправлення останнього разу у Receiver

	mq messageQueue

	stop    chan struct{}
	stopReq bool
}

var size = image.Pt(400, 400)

// Start запускає цикл подій. Цей метод потрібно запустити до того, як викликати на ньому будь-які інші методи.
func (l *Loop) Start(s screen.Screen) {
	l.next, _ = s.NewTexture(size)
	l.prev, _ = s.NewTexture(size)
	l.stop = make(chan struct{})

	go func() {
		for {
			select {
			case <-l.stop:
				return
			default:
				if !l.mq.empty() {
					op := l.mq.pull()
					if op == nil {
						continue
					}
					if update := op.Do(l.next); update {
						l.Receiver.Update(l.next)
						l.next, l.prev = l.prev, l.next
					}
				} else {
					time.Sleep(10 * time.Millisecond)
				}
			}
		}
	}()
}

// Post додає нову операцію у внутрішню чергу.
func (l *Loop) Post(op Operation) {
	l.mq.push(op)
}

// StopAndWait сигналізує про необхідність завершити цикл та блокується до моменту його повної зупинки.
func (l *Loop) StopAndWait() {
	close(l.stop)
}

type messageQueue struct {
	items []Operation
	mu    sync.Mutex
}

func (mq *messageQueue) push(op Operation) {
	mq.mu.Lock()
	defer mq.mu.Unlock()
	mq.items = append(mq.items, op)
}

func (mq *messageQueue) pull() Operation {
	mq.mu.Lock()
	defer mq.mu.Unlock()
	if len(mq.items) == 0 {
		return nil
	}
	op := mq.items[0]
	mq.items = mq.items[1:]
	return op
}

func (mq *messageQueue) empty() bool {
	mq.mu.Lock()
	defer mq.mu.Unlock()
	return len(mq.items) == 0
}
