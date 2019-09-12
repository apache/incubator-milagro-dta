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

package logger

import (
	"fmt"
	"io"
	"log"
	"strings"

	kitlog "github.com/go-kit/kit/log"
)

type logTextLogger struct {
	log *log.Logger
}

// NewLogTextLogger returns a logger that encodes keyvals to the Writer in
// plain text format.
// The passed Writer must be safe for concurrent use by multiple goroutines if
// the returned Logger will be used concurrently.
func NewLogTextLogger(w io.Writer) kitlog.Logger {
	return &logTextLogger{log.New(w, "", 0)}
}

func (l logTextLogger) Log(keyvals ...interface{}) error {
	var out string

	for i, v := range keyvals {
		if i%2 == 0 {
			continue
		}

		k := keyvals[i-1]

		switch k {
		case "ts", "msg", "method", "path", "statusCode", "statusText", "reqID":
			out += fmt.Sprintf("%v ", v)
		case "t":
			s, ok := v.(string)
			if ok {
				out += fmt.Sprintf("[%s] ", strings.ToUpper(s))
			} else {
				out += fmt.Sprintf("%v ", v)
			}
		case "time":
			out += fmt.Sprintf("(%v) ", v)
		default:
			out += fmt.Sprintf("%v: %v ", k, v)
		}
	}

	l.log.Println(strings.TrimSpace(out))
	return nil
}
