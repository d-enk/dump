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
	d.WriteString(d.prefix)

	for _, s := range ss {
		d.dump(d.prefix, s)
		d.WriteByte(' ')
	}

	return d
}

func (d *Dumper) dump(prefix string, s any) *Dumper {
	switch v := s.(type) {
	case string, error:
		d.addStr(prefix, v)
		return d
	}

	switch v := reflect.ValueOf(s); v.Kind() {
	case reflect.Pointer:
		if v.IsNil() {
			d.dump(prefix, "<nil>")
		} else {
			d.dump(prefix, v.Elem().Interface())
		}

	case reflect.Slice, reflect.Array:
		if v.Len() == 0 {
			d.Add("[]")
		} else {
			d.Addln("[")
			for i := 0; i < v.Len(); i++ {
				d.Add(prefix, " ")
				d.dump(prefix+" ", v.Index(i).Interface())
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
				d.dump(prefix+" ", iter.Key().Interface())
				d.Add(": ")
				d.dump(prefix+" ", iter.Value().Interface())
				d.Addln(",")
			}
			d.Add(prefix, "}")
		}

	case reflect.Struct:

		var fields bool

		for i := 0; i < v.NumField(); i++ {
			if fields = v.Type().Field(i).IsExported(); fields {
				break
			}
		}

		if fields {
			d.Addln("{")
			for i := 0; i < v.NumField(); i++ {
				field := v.Field(i)

				fieldType := v.Type().Field(i)
				if fieldType.IsExported() {
					d.Add(prefix, " ", fieldType.Name, ": ")
					d.dump(prefix+"  ", field.Interface())
					d.Addln()
				} else {
					d.Addln(prefix, " ", fieldType.Name, ": <>")
				}
			}
			d.Add(prefix, "}")
		} else if v, ok := s.(interface{ String() string }); ok {
			d.addStr(prefix, v)
		} else {
			d.Add("{}")
		}

	case reflect.Interface:
		d.dump(prefix, v.Elem())

	default:
		d.Add(fmt.Sprint(s))
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

func (d *Dumper) addStr(prefix string, v any) {
	str := fmt.Sprint(v)
	stri := str

	if i := strings.IndexByte(stri, '\n'); i >= 0 {
		d.Add("\n", prefix, " `")
		for ; i >= 0; i = strings.IndexByte(stri, '\n') {
			d.Add(stri[:i], "\n", prefix, "  ")
			stri = stri[i+1:]
		}
	} else {
		d.WriteByte('`')
	}
	d.Add(stri)
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
