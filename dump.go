package dump

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"reflect"
	"strings"
	"sync"
	"unsafe"
)

var (
	NestedPadding = "  "
	LogWriter     = io.Writer(os.Stderr)
	LogChunkSize  = 1024
	logMutex      = sync.Mutex{}
	logBuffer     = bytes.Buffer{}
)

func New(writer io.Writer) *Dumper {
	return &Dumper{buffer: noBuf{Writer: writer}}
}

type Dumper struct {
	prefix string
	buffer
}

type tLog func(v ...any) tLog

func Log(v ...any) tLog {
	var gen func(*Dumper, ...any) tLog
	gen = func(d *Dumper, v ...any) tLog {
		logMutex.Lock()
		defer logMutex.Unlock()

		d.Dumpln(v...)
		return func(v ...any) tLog {
			return gen(d.WithPrefix(NestedPadding), v...)
		}
	}

	return gen(&Dumper{buffer: &buf{
		Writer:    LogWriter,
		chunkSize: LogChunkSize,
		buf:       logBuffer,
	}}, v...)
}

func Dump(v ...any) string {
	builder := &strings.Builder{}
	New(builder).Dump(v...)
	return builder.String()
}

func (d Dumper) WithPrefix(prefix string) *Dumper {
	d.prefix += prefix
	return &d
}

func (d *Dumper) WithBuffer(chunkSize int) *Dumper {
	d.buffer = &buf{
		Writer:    d.buffer,
		chunkSize: chunkSize,
	}
	return d
}

func (d *Dumper) Dump(v ...any) {
	if len(v) > 0 {
		d.multiDump(v...)
		d.buffer.flush()
	}
}

func (d *Dumper) Dumpln(v ...any) {
	if len(v) > 0 {
		d.multiDump(v...)
	}
	_ = d.buffer.add([]byte{'\n'})
	d.buffer.flush()
}

func (d *Dumper) multiDump(v ...any) {
	d.append(d.prefix)
	d.dump(d.prefix, reflect.ValueOf(v[0]), false)
	for _, v := range v[1:] {
		d.append(" ")
		d.dump(d.prefix, reflect.ValueOf(v), false)
	}
}

func (d *Dumper) dump(prefix string, v reflect.Value, isNested bool) {
	switch v.Kind() {
	case reflect.String:
		d.addMultilineString(prefix, v.String(), isNested)

	case reflect.Invalid:
		d.append("nil")

	case reflect.Pointer, reflect.Interface:
		if v.IsNil() {
			d.append("nil")
		} else {
			d.dump(prefix, v.Elem(), isNested)
		}

	case reflect.Slice, reflect.Array:
		if v.Len() == 0 {
			d.append("[]")
		} else {
			d.appendln("[")
			nextPrefix := prefix + NestedPadding
			for i := 0; i < v.Len(); i++ {
				d.append(nextPrefix)
				d.dump(nextPrefix, v.Index(i), false)
				d.appendln(",")
			}
			d.append(prefix, "]")
		}

	case reflect.Map:
		if v.Len() == 0 {
			d.append("{}")
		} else {
			d.appendln("{")
			nextPrefix := prefix + NestedPadding
			for iter := v.MapRange(); iter.Next(); {
				d.append(nextPrefix)
				d.dump(nextPrefix, iter.Key(), true)
				d.append(": ")
				d.dump(nextPrefix, iter.Value(), true)
				d.appendln(",")
			}
			d.append(prefix, "}")
		}

	case reflect.Struct:
		if v.NumField() == 0 {
			d.append("{}")
		} else {
			d.appendln("{")
			nextPrefix := prefix + NestedPadding
			for i := 0; i < v.NumField(); i++ {
				fieldName := v.Type().Field(i).Name
				d.append(nextPrefix, fieldName, ": ")
				d.dump(nextPrefix, v.FieldByName(fieldName), true)
				d.appendln("")
			}
			d.append(prefix, "}")
		}

	case reflect.Func:
		d.append(fmt.Sprintf("func()=%v", v))

	default:
		d.append(fmt.Sprint(v))
	}
}

func (d *Dumper) append(strs ...string) {
	for _, str := range strs {
		_ = d.buffer.add(unsafe.Slice(unsafe.StringData(str), len(str)))
	}
}

func (d *Dumper) appendln(strs ...string) {
	d.append(strs...)
	d.append("\n")
}

func (d *Dumper) addMultilineString(prefix, str string, needNewLine bool) {
	if i := strings.IndexByte(str, '\n'); i >= 0 {
		padding := ""
		if needNewLine {
			d.append("|\n", prefix)
			padding = " "
		}

		d.append(padding, "`")
		for ; i >= 0; i = strings.IndexByte(str, '\n') {
			d.append(str[:i], "\n", prefix, padding, " ")
			str = str[i+1:]
		}
	} else {
		d.append("`")
	}
	d.append(str, "`")
}

type buffer interface {
	io.Writer
	add([]byte) error
	flush() error
}

type noBuf struct{ io.Writer }

func (b noBuf) add(bytes []byte) (err error) {
	_, err = b.Writer.Write(bytes)
	return
}

func (b noBuf) flush() error { return nil }

type buf struct {
	io.Writer
	buf       bytes.Buffer
	chunkSize int
}

func (b *buf) add(bytes []byte) (err error) {
	if _, err = b.buf.Write(bytes); err == nil && b.buf.Len() >= b.chunkSize {
		err = b.flush()
	}
	return
}

func (b *buf) flush() (err error) {
	_, err = b.buf.WriteTo(b.Writer)
	return
}
