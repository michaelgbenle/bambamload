package logger

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sync"
	"time"
)

type DailyFileWriter struct {
	mu         sync.Mutex
	baseDir    string
	prefix     string
	currentDay string
	file       *os.File
	writer     io.Writer
}

func NewDailyFileWriter(baseDir, prefix string) (*DailyFileWriter, error) {
	w := &DailyFileWriter{
		baseDir: baseDir,
		prefix:  prefix,
	}
	if err := w.rotate(); err != nil {
		return nil, err
	}
	go w.autoRotate()
	return w, nil
}

func (w *DailyFileWriter) Write(p []byte) (int, error) {
	w.mu.Lock()
	defer w.mu.Unlock()
	return w.writer.Write(p)
}

func (w *DailyFileWriter) rotate() error {
	w.mu.Lock()
	defer w.mu.Unlock()

	currentDay := time.Now().Format("2006-01-02")
	if w.currentDay == currentDay {
		return nil // no rotation needed
	}

	_ = os.MkdirAll(w.baseDir, os.ModePerm)

	if w.file != nil {
		_ = w.file.Close()
	}

	filename := filepath.Join(w.baseDir, fmt.Sprintf("%s-%s.log", w.prefix, currentDay))
	file, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return err
	}

	w.currentDay = currentDay
	w.file = file
	w.writer = file
	return nil
}

func (w *DailyFileWriter) autoRotate() {
	for {
		next := time.Now().Add(24 * time.Hour).Truncate(24 * time.Hour)
		time.Sleep(time.Until(next))
		_ = w.rotate()
	}
}
