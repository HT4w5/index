package index

import "log"

type Logger interface {
	Debugf(format string, a ...any)
	Infof(format string, a ...any)
	Warnf(format string, a ...any)
	Errorf(format string, a ...any)
}

type SimpleLogger struct {
}

func (l *SimpleLogger) Debugf(format string, a ...any) {
	log.Printf("[DEBUG]"+format, a...)
}

func (l *SimpleLogger) Infof(format string, a ...any) {
	log.Printf("[INFO]"+format, a...)
}

func (l *SimpleLogger) Warnf(format string, a ...any) {
	log.Printf("[WARN]"+format, a...)
}

func (l *SimpleLogger) Errorf(format string, a ...any) {
	log.Printf("[ERROR]"+format, a...)
}

type DiscardLogger struct {
}

func (l *DiscardLogger) Debugf(format string, a ...any) {
}

func (l *DiscardLogger) Infof(format string, a ...any) {
}

func (l *DiscardLogger) Warnf(format string, a ...any) {
}

func (l *DiscardLogger) Errorf(format string, a ...any) {
}
