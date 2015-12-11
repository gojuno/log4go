/* term.go
 *
 * Copyright (c) 2010, Kyle Lemons <kyle@kylelemons.net> (creator).
 * All rights reserved.
 *
 * This software may be modified and distributed under the terms
 * of the New BSD license.  See the LICENSE file for details.
 */
package log4go

import (
	"fmt"
	"io"
	"os"
)

var stdout io.Writer = os.Stdout

// This is the standard writer that prints to standard output.
type ConsoleLogWriter struct {
	recordsChan chan *LogRecord
	format      string
}

// This creates a new ConsoleLogWriter
func NewConsoleLogWriter() *ConsoleLogWriter {
	clw := ConsoleLogWriter{
		format: "[%D %T] [%L] (%S) %M",
	}
	clw.recordsChan = make(chan *LogRecord, LogBufferLength)
	go clw.run(stdout)
	return &clw
}

func (w *ConsoleLogWriter) SetFormat(format string) {
	w.format = format
}

func (w *ConsoleLogWriter) run(out io.Writer) {
	for rec := range w.recordsChan {
		fmt.Fprint(out, FormatLogRecord(w.format, rec))
	}
}

// This is the ConsoleLogWriter's output method.  This will block if the output
// buffer is full.
func (w *ConsoleLogWriter) LogWrite(rec *LogRecord) {
	w.recordsChan <- rec
}

// Close stops the logger from sending messages to standard output.  Attempts to
// send log messages to this logger after a Close have undefined behavior.
func (w *ConsoleLogWriter) Close() {
	close(w.recordsChan)
}
