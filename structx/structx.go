package structx

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"github.com/odycenter/std-library/crypto/md5"
	"net/url"
	"reflect"
	"regexp"
	"sort"
	"strings"
)

type EncodeValue url.Values

// Encode
// url.Encode
// 第一个参数是连接符默认 =
// 第二个参数是分隔符默认 &
// 输出结果：
// A=10&B=20.1&C=name
func (ev EncodeValue) Encode(symbol ...string) string {
	var joinSymbol, separatorSymbol string
	if len(symbol) == 1 {
		joinSymbol = symbol[0]
	} else {
		separatorSymbol = "&"
	}

	if len(symbol) == 2 {
		joinSymbol = symbol[0]
		separatorSymbol = symbol[1]
	} else {
		joinSymbol = "="
		separatorSymbol = "&"
	}

	if ev == nil {
		return ""
	}
	var buf strings.Builder
	keys := make([]string, 0, len(ev))
	for k := range ev {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		vs := ev[k]
		keyEscaped := url.QueryEscape(k)
		for _, v := range vs {
			if buf.Len() > 0 {
				buf.WriteString(separatorSymbol)
			}
			buf.WriteString(keyEscaped)
			buf.WriteString(joinSymbol)
			buf.WriteString(url.QueryEscape(v))
		}
	}
	return buf.String()
}

func Map(in any) (map[string]any, error) {

	result := make(map[string]any)

	v := reflect.ValueOf(in)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	// we only accept structs
	if v.Kind() != reflect.Struct {
		return nil, fmt.Errorf("only accepts structs but got %T", v)
	}

	typ := v.Type()
	for i := 0; i < v.NumField(); i++ {

		field := typ.Field(i)

		result[field.Name] = v.Field(i).Interface()

	}
	return result, nil
}

// Sign
// 对struct进行签名，如果未传要签名的字段。默认对所有字段进行签名
// 签名方式
//
//	默认对传递的 field 进行签名，如果未传，则对所有的field进行签名
//	获取field对应的orm:"column()"作为key，如果没有，则使用field name
func Sign(in any, salt string, fields ...string) string {
	v := reflect.ValueOf(in)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	if v.Kind() != reflect.Struct {
		return ""
	}

	var values = make(url.Values)
	getSignFields(v, values, fields...)
	values.Del("sign")

	signStr := fmt.Sprintf("%s:salt=%s", values.Encode(), salt)
	return md5.Sum([]byte(signStr)).Hex()
}

func getSignFields(v reflect.Value, values url.Values, fields ...string) {

	if len(fields) == 0 {
		// 获取所有的属性
		for i := 0; i < v.NumField(); i++ {
			field := v.Type().Field(i)
			fieldV := v.FieldByName(field.Name)

			if fieldV.Kind() == reflect.Struct {
				getSignFields(fieldV, values, fields...)
			} else {
				// 获取orm注解
				key := ormTag(field)
				values[key] = []string{fmt.Sprint(fieldV.Interface())}
			}
		}
	} else {
		// 获取指定的属性
		for _, f := range fields {
			if field, ok := v.Type().FieldByName(f); ok {
				fieldV := v.FieldByName(field.Name)

				if fieldV.Kind() == reflect.Struct {
					getSignFields(fieldV, values, fields...)
				} else {
					// 获取orm注解
					key := ormTag(field)
					values[key] = []string{fmt.Sprint(fieldV.Interface())}
				}
			}
		}
	}
}

func ormTag(field reflect.StructField) (ormField string) {
	tag := field.Tag.Get("orm")
	match := regexp.MustCompile(`column\(([a-zA-Z0-9_]+)\)`).FindStringSubmatch(tag)
	if len(match) > 1 {
		return match[1]
	} else {
		return field.Name
	}
}

func UrlValue(in any, emptyEncode ...bool) (out EncodeValue) {
	out = make(EncodeValue)

	v := reflect.ValueOf(in)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	if v.Kind() != reflect.Struct {
		return nil
	}

	if len(emptyEncode) > 0 && emptyEncode[0] == true {
		getFieldsAndValue(v, out, true)
	} else {
		getFieldsAndValue(v, out, false)
	}

	return
}

func getFieldsAndValue(v reflect.Value, values EncodeValue, emptyEncode bool) {
	for i := 0; i < v.NumField(); i++ {
		// 第归
		if v.Field(i).Kind() == reflect.Struct {
			getFieldsAndValue(v.Field(i), values, emptyEncode)
			continue
		}

		field := v.Type().Field(i)

		// 获取json注解
		var tag string
		jsonTag := field.Tag.Get("json")
		if len(jsonTag) > 0 {
			tag = strings.Split(jsonTag, ",")[0]
		}
		if tag == "-" {
			continue
		}
		if len(tag) == 0 {
			tag = field.Name
		}

		v := fmt.Sprintf("%v", v.Field(i).Interface())
		if emptyEncode || len(v) > 0 && v != "0" {
			values[tag] = []string{v}
		}
	}

	return
}

func Field(t any, fieldName string) (reflect.StructField, error) {
	s := reflect.TypeOf(t)
	if s.Kind() == reflect.Ptr {
		s = s.Elem()
	}

	if s.Kind() != reflect.Struct {
		return reflect.StructField{}, fmt.Errorf("need reflect.Struct")
	}

	for i := 0; i < s.NumField(); i++ {
		key := s.Field(i).Name

		if key == fieldName {
			return s.Field(i), nil
		}
	}

	return reflect.StructField{}, fmt.Errorf("not find struct field:%q", fieldName)
}

// Fields 获取 struct 的所有属性
func Fields(t reflect.Type) (keys []string) {
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	for i := 0; i < t.NumField(); i++ {
		key := t.Field(i).Name

		if t.Field(i).Type.Kind() == reflect.Struct {
			keys = append(Fields(t.Field(i).Type), keys...)
			continue
		} else {
			keys = append(keys, key)
		}
	}

	sort.Strings(keys)
	return
}

// DeepCopy 深层拷贝
func DeepCopy(dst, src any) error {
	var buf bytes.Buffer
	if err := gob.NewEncoder(&buf).Encode(src); err != nil {
		return err
	}
	return gob.NewDecoder(bytes.NewBuffer(buf.Bytes())).Decode(dst)
}
