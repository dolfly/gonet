// Copyright 2016 fatedier, fatedier@gmail.com
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package log

import (
	"log"
	"os"
)

var Log *log.Logger

func init() {
	Log = log.New(os.Stdout, "gonet", log.LstdFlags)
}

// wrap log
func Error(format string, v ...interface{}) {
	Log.Printf(format, v...)
}

func Warn(format string, v ...interface{}) {
	Log.Printf(format, v...)
}

func Info(format string, v ...interface{}) {
	Log.Printf(format, v...)
}

func Debug(format string, v ...interface{}) {
	Log.Printf(format, v...)
}

func Trace(format string, v ...interface{}) {
	Log.Printf(format, v...)
}

// Logger
type Logger interface {
	AddLogPrefix(string)
	GetPrefixStr() string
	GetAllPrefix() []string
	ClearLogPrefix()
	Error(string, ...interface{})
	Warn(string, ...interface{})
	Info(string, ...interface{})
	Debug(string, ...interface{})
	Trace(string, ...interface{})
}

type PrefixLogger struct {
	prefix    string
	allPrefix []string
}

func NewPrefixLogger(prefix string) *PrefixLogger {
	logger := &PrefixLogger{
		allPrefix: make([]string, 0),
	}
	logger.AddLogPrefix(prefix)
	return logger
}

func (pl *PrefixLogger) AddLogPrefix(prefix string) {
	if len(prefix) == 0 {
		return
	}

	pl.prefix += "[" + prefix + "] "
	pl.allPrefix = append(pl.allPrefix, prefix)
}

func (pl *PrefixLogger) GetPrefixStr() string {
	return pl.prefix
}

func (pl *PrefixLogger) GetAllPrefix() []string {
	return pl.allPrefix
}

func (pl *PrefixLogger) ClearLogPrefix() {
	pl.prefix = ""
	pl.allPrefix = make([]string, 0)
}

func (pl *PrefixLogger) Error(format string, v ...interface{}) {
	Log.Printf(pl.prefix+format, v...)
}

func (pl *PrefixLogger) Warn(format string, v ...interface{}) {
	Log.Printf(pl.prefix+format, v...)
}

func (pl *PrefixLogger) Info(format string, v ...interface{}) {
	Log.Printf(pl.prefix+format, v...)
}

func (pl *PrefixLogger) Debug(format string, v ...interface{}) {
	Log.Printf(pl.prefix+format, v...)
}

func (pl *PrefixLogger) Trace(format string, v ...interface{}) {
	Log.Printf(pl.prefix+format, v...)
}
