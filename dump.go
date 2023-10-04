package dump

import (
	"fmt"
	"log"
	"os"
	"reflect"
	"strings"
)

type Dumper struct {
	prefix string
	strings.Builder
}

func (d *Dumper) WithPrefix(prefix string) *Dumper {
	d.prefix = prefix
	return d
}

func Dump(ss ...any) string {
	return (&Dumper{}).Dump(ss...).String()
}

func (d *Dumper) Dump(ss ...any) *Dumper {
	if len(ss) > 0 {
		d.WriteString(d.prefix)
		d.dump(d.prefix, reflect.ValueOf(ss[0]))
		for _, s := range ss[1:] {
			d.WriteByte(' ')
			d.dump(d.prefix, reflect.ValueOf(s))
		}
	}

	return d
}

func (d *Dumper) dump(prefix string, v reflect.Value) *Dumper {
	switch v.Kind() {
	case reflect.String:
		d.addStr(prefix, v.String())

	case reflect.Pointer, reflect.Interface:
		if v.IsNil() {
			d.Add("nil")
		} else {
			d.dump(prefix, v.Elem())
		}

	case reflect.Slice, reflect.Array:
		if v.Len() == 0 {
			d.Add("[]")
		} else {
			d.Addln("[")
			for i := 0; i < v.Len(); i++ {
				d.Add(prefix, " ")
				d.dump(prefix+" ", v.Index(i))
				d.Addln(",")
			}
			d.Add(prefix, "]")
		}

	case reflect.Map:
		if v.Len() == 0 {
			d.Add("{}")
		} else {
			d.Addln("{")
			for iter := v.MapRange(); iter.Next(); {
				d.Add(prefix, " ")
				d.dump(prefix+" ", iter.Key())
				d.Add(": ")
				d.dump(prefix+" ", iter.Value())
				d.Addln(",")
			}
			d.Add(prefix, "}")
		}

	case reflect.Struct:
		if v.NumField() > 0 {
			d.Addln("{")

			for i := 0; i < v.NumField(); i++ {
				fieldName := v.Type().Field(i).Name
				d.Add(prefix, "  ", fieldName, ": ")
				d.dump(prefix+"  ", v.FieldByName(fieldName))
				d.Addln()
			}
			d.Add(prefix, "}")
		} else {
			d.Add("{}")
		}

	default:
		d.Add(fmt.Sprint(v))
	}

	return d
}

func (d *Dumper) Add(strs ...string) {
	for _, str := range strs {
		d.WriteString(str)
	}
}

func (d *Dumper) Addln(strs ...string) {
	d.Add(strs...)
	d.WriteByte('\n')
}

func (d *Dumper) addStr(prefix, str string) {
	if i := strings.IndexByte(str, '\n'); i >= 0 {
		d.Add("\n", prefix, " `")
		for ; i >= 0; i = strings.IndexByte(str, '\n') {
			d.Add(str[:i], "\n", prefix, "  ")
			str = str[i+1:]
		}
	} else {
		d.WriteByte('`')
	}
	d.Add(str)
	d.WriteByte('`')
}

type tLog func(s ...any) tLog

var info = log.New(os.Stderr, "", log.Lmsgprefix).Println

func Log(s ...any) tLog {
	var gen func(Dumper, ...any) tLog
	gen = func(d Dumper, ss ...any) tLog {
		d.Dump(ss...)
		info(d.String())
		d.Reset()

		d.prefix += "    "
		return func(s ...any) tLog {
			return gen(d, s...)
		}
	}

	return gen(Dumper{}, s...)
}
