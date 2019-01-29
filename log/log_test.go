/*
Copyright (C)  2018 Yahoo Japan Corporation Athenz team.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package log

import (
	"bytes"
	"fmt"
	"strings"
	"testing"

	"github.com/kpango/glg"
)

func TestNew(t *testing.T) {
	type args struct {
		requestID  string
		printValue string
	}
	tests := []struct {
		name  string
		args  args
		want  Logger
		wantW string
	}{
		{
			name: "log instance New test.",
			args: args{
				requestID:  "requestID-25",
				printValue: "test",
			},
			want: &logger{
				log: glg.New().
					SetPrefix("prefix-29").
					SetLevelWriter(glg.PRINT, bytes.NewBuffer(nil)).
					SetLevelMode(glg.PRINT, glg.WRITER),
			},
			wantW: "[requestID-25]:	test\n",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			w := bytes.NewBuffer(nil)

			got := New(w, tt.args.requestID)
			got.Printf(tt.args.printValue)

			if gotW := w.String(); !strings.HasSuffix(gotW, tt.wantW) {
				t.Errorf("New() = [%v], want [%v]", gotW, tt.wantW)
			}
		})
	}
}

func Test_logger_Printf(t *testing.T) {
	type fieldsArgs struct {
		prefix string
		buffer *bytes.Buffer
	}
	type args struct {
		format string
		args   []interface{}
	}
	tests := []struct {
		name       string
		fieldsArgs fieldsArgs
		args       args
		wantW      string
	}{
		{
			name: "log printf write nil test.",
			fieldsArgs: fieldsArgs{
				prefix: "prefix-72",
				buffer: bytes.NewBuffer(nil),
			},
			args: args{
				format: "%v", // format value empty case is not test target.
				args: func() []interface{} {
					return nil
				}(),
			},
			wantW: "[prefix-72]:	%!v(MISSING)\n",
		},
		{
			name: "log printf write test.",
			fieldsArgs: fieldsArgs{
				prefix: "prefix-86",
				buffer: bytes.NewBuffer(nil),
			},
			args: args{
				format: "%v",
				args: func() []interface{} {
					return []interface{}{
						"args-93",
					}
				}(),
			},
			wantW: "[prefix-86]:	args-93\n",
		},
		{
			name: "log printf write test.",
			fieldsArgs: fieldsArgs{
				prefix: "prefix-102",
				buffer: bytes.NewBuffer(nil),
			},
			args: args{
				format: "%v%v",
				args: func() []interface{} {
					return []interface{}{
						"args-109",
						"args-110",
					}
				}(),
			},
			wantW: "[prefix-102]:	args-109args-110\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := &logger{
				log: glg.New().
					SetPrefix(tt.fieldsArgs.prefix).
					SetLevelWriter(glg.PRINT, tt.fieldsArgs.buffer).
					SetLevelMode(glg.PRINT, glg.WRITER),
			}

			l.Printf(tt.args.format, tt.args.args...)

			if gotW := tt.fieldsArgs.buffer.String(); !strings.HasSuffix(gotW, tt.wantW) {
				t.Errorf("New() = [%v], want [%v]", gotW, tt.wantW)
			}
		})
	}
}

func Test_logger_Println(t *testing.T) {
	type fields struct {
		log *glg.Glg
	}
	type args struct {
		args []interface{}
	}
	type test struct {
		name         string
		fields       fields
		args         args
		checkFunc    func() error
		checkRecover func(interface{}) error
	}
	tests := []test{
		func() test {
			b := bytes.NewBuffer(nil)
			return test{
				name: "log println write test.",
				fields: fields{
					log: glg.New().
						SetPrefix("prefix-125").
						SetLevelWriter(glg.PRINT, b).
						SetLevelMode(glg.PRINT, glg.WRITER),
				},
				args: args{
					args: func() []interface{} {
						return []interface{}{
							"args-131",
						}
					}(),
				},
				checkFunc: func() error {
					got := b.String()
					want := "[prefix-125]:	args-131\n"

					if !strings.HasSuffix(got, want) {
						return fmt.Errorf("New() = [%v], want [%v]", got, want)
					}
					return nil
				},
			}
		}(),
		func() test {
			b := bytes.NewBuffer(nil)
			return test{
				name: "log println write test.",
				fields: fields{
					log: glg.New().
						SetPrefix("prefix-167").
						SetLevelWriter(glg.PRINT, b).
						SetLevelMode(glg.PRINT, glg.WRITER),
				},
				args: args{
					args: func() []interface{} {
						return []interface{}{
							"args-173",
							"args-174",
						}
					}(),
				},
				checkFunc: func() error {
					got := b.String()
					want := "[prefix-167]:	args-173 args-174\n"

					if !strings.HasSuffix(got, want) {
						return fmt.Errorf("New() = [%v], want [%v]", got, want)
					}
					return nil
				},
			}
		}(),
		func() test {
			b := &WriterMock{
				WriteFunc: func(p []byte) (n int, err error) {
					return 0, fmt.Errorf("Error")
				},
			}
			return test{
				name: "log println write error.",
				fields: fields{
					log: glg.New().
						SetLevelWriter(glg.PRINT, b).
						SetLevelMode(glg.PRINT, glg.WRITER),
				},
				args: args{
					args: []interface{}{
						"args-173",
					},
				},
				checkRecover: func(i interface{}) error {
					want := "Error"
					if i != want {
						return fmt.Errorf("logger error, recovered: %v, want: %v", i, want)
					}
					return nil
				},
			}
		}(),
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				if r := recover(); r != nil {
					if err := tt.checkRecover(r); err != nil {
						t.Error(err)
						return
					}
				}
			}()
			l := &logger{
				log: tt.fields.log,
			}
			glg.ReplaceExitFunc(func(i int) {})
			l.Println(tt.args.args...)

			if err := tt.checkFunc(); err != nil {
				t.Error(err)
			}
		})
	}
}
