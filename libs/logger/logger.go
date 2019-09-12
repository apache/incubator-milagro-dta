// Licensed to the Apache Software Foundation (ASF) under one
// or more contributor license agreements.  See the NOTICE file
// distributed with this work for additional information
// regarding copyright ownership.  The ASF licenses this file
// to you under the Apache License, Version 2.0 (the
// "License"); you may not use this file except in compliance
// with the License.  You may obtain a copy of the License at
//
//   http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.

/*
Package logger - configurable logging middleware
*/
package logger

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/go-kit/kit/log"
)

const (
	timeFormat = "2006-01-02T15:04:05.000"

	logLevelNone = iota
	logLevelError
	logLevelInfo
	logLevelDebug
)

// Logger implements go-kit logger interface
// Adding additional functionality
type Logger struct {
	logger   log.Logger
	logLevel int
}

// NewLogger initializes a new Logger instance
// Supported log types are: text, json, none
func NewLogger(logType, logLevel string, options ...Option) (*Logger, error) {
	logger := logTypeFromString(logType)
	if logger == nil {
		return nil, errors.New("unuspported log format")
	}

	l := &Logger{
		logger:   logger,
		logLevel: logLevelFromString(logLevel),
	}

	for _, o := range options {
		if err := o(l); err != nil {
			return nil, err
		}
	}

	return l, nil
}

func logTypeFromString(logType string) log.Logger {
	var logger log.Logger
	switch logType {
	case "text":
		logger = NewLogTextLogger(os.Stdout)
	case "fmt":
		logger = log.NewLogfmtLogger(os.Stdout)
	case "json":
		logger = log.NewJSONLogger(os.Stdout)
	case "none":
		logger = log.NewNopLogger()
	default:
		return nil
	}

	timestampFormat := log.TimestampFormat(
		func() time.Time { return time.Now().UTC() },
		timeFormat,
	)

	logger = log.With(logger, "ts", timestampFormat)

	return logger
}

func logLevelFromString(logLevel string) int {
	switch logLevel {
	case "error":
		return logLevelError
	case "info":
		return logLevelInfo
	case "debug":
		return logLevelDebug
	}

	return logLevelNone
}

// Option for adding additional configurations of Logger instance
type Option func(l *Logger) error

// With adds additional values to the logger
func With(k, v interface{}) Option {
	return func(l *Logger) error {
		l.logger = log.With(l.logger, k, v)
		return nil
	}
}

// Log satisfies Logger interface
func (l *Logger) Log(kvals ...interface{}) error {
	return l.logger.Log(kvals...)
}

// Error log
func (l *Logger) Error(msg string, vals ...interface{}) {
	if l.logLevel < logLevelError {
		return
	}

	_ = l.Log("t", "err", "msg", fmt.Sprintf(msg, vals...))
}

// Info log
func (l *Logger) Info(msg string, vals ...interface{}) {
	if l.logLevel < logLevelInfo {
		return
	}
	_ = l.Log("t", "inf", "msg", fmt.Sprintf(msg, vals...))
}

// Request log
func (l *Logger) Request(reqID, method, path string) {
	if l.logLevel < logLevelInfo {
		return
	}
	_ = l.Log("t", "req", "reqID", reqID, "method", method, "path", path)
}

// Response log
func (l *Logger) Response(reqID, method, path string, statusCode int, statusText string, reqTime time.Duration) {
	if l.logLevel < logLevelInfo {
		return
	}

	_ = l.Log("t", "res", "reqID", reqID, "method", method, "path", path, "statusCode", statusCode, "statusText", statusText, "time", reqTime.String())
}

// Debug log
func (l *Logger) Debug(msg string, vals ...interface{}) {
	if l.logLevel < logLevelDebug {
		return
	}
	_ = l.Log("t", "dbg", "msg", fmt.Sprintf(msg, vals...))
}
