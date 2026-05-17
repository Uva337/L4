package logger

import (
	"fmt"
	"os"
	"time"
)

// обрабатывает логи в отдельной горутине
type AsyncLogger struct {
	logChan chan string
}

func NewAsyncLogger(bufferSize int) *AsyncLogger {
	l := &AsyncLogger{
		logChan: make(chan string, bufferSize),
	}
	// запускаем воркер, который будет разгребать канал
	go l.worker()
	return l
}

func (l *AsyncLogger) worker() {
	// единственная горутина, кто имеет право писать в stdout
	for msg := range l.logChan {
		fmt.Fprintln(os.Stdout, msg)
	}
}

// отправляет сообщение в канал логгера
func (l *AsyncLogger) Info(method, path, duration string) {
	msg := fmt.Sprintf("[%s] %s %s | ⏱ %s", time.Now().Format("2006-01-02 15:04:05"), method, path, duration)

	select {
	case l.logChan <- msg:
	default:
	}
}
