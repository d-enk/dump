package dump

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"reflect"
	"strings"
)

var (
	NestedPadding = "  "
	LogWriter     = io.Writer(os.Stderr)
	ChunkSize     = 1024
)

func New(writer io.Writer) *Dumper {
	return &Dumper{writer: writer}
}

type Dumper struct {
	prefix string
	writer io.Writer
	buffer bytes.Buffer
}

type tLog func(v ...any) tLog

func Log(v ...any) tLog {
	var gen func(*Dumper, ...any) tLog
	gen = func(d *Dumper, v ...any) tLog {
		d.Dumpln(v...)
		return func(v ...any) tLog {
			return gen(d.WithPrefix(NestedPadding), v...)
		}
	}

	return gen(New(LogWriter), v...)
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

func (d *Dumper) Dump(v ...any) {
	if len(v) > 0 {
		d.multiDump(v...)
		d.flushBuffer()
	}
}

func (d *Dumper) Dumpln(v ...any) {
	if len(v) > 0 {
		d.multiDump(v...)
	}
	_ = d.buffer.WriteByte('\n')
	d.flushBuffer()
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

	default:
		d.append(fmt.Sprint(v))
	}
}

func (d *Dumper) flushBuffer() {
	_, _ = d.buffer.WriteTo(d.writer)
}

func (d *Dumper) append(strs ...string) {
	for _, str := range strs {
		_, _ = d.buffer.WriteString(str)

		if d.buffer.Len() >= ChunkSize {
			d.flushBuffer()
		}
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
