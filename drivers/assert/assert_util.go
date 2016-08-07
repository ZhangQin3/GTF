// Copyright 2014 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.
// https://github.com/issue9/assert
package assert

import (
	"bytes"
	"fmt"
	"gtf/drivers/log"
	"reflect"
	"strings"
	"time"
)

// 当expr条件不成立时，输出错误信息。
//
// expr 返回结果值为bool类型的表达式；
// msg1,msg2输出的错误信息，之所以提供两组信息，是方便在用户没有提供的情况下，
// 可以使用系统内部提供的信息，优先使用msg1中的信息，若不存在，则使用msg2的内容。
func assert(expr bool, msg1 []interface{}, msg2 []interface{}) {
	if !expr {
		log.Error(formatMessage(msg1, msg2))
	}
}

// 格式化错误提示信息。
//
// msg1中的所有参数将依次被传递给fmt.Sprintf()函数，
// 所以msg1[0]必须可以转换成string(如:string, []byte, []rune, fmt.Stringer)
//
// msg2参数格式与msg1完全相同，在msg1为空的情况下，会使用msg2的内容，
// 否则msg2不会启作用。
func formatMessage(msg1 []interface{}, msg2 []interface{}) string {
	if len(msg1) == 0 {
		msg1 = msg2
	}

	if len(msg1) == 0 {
		return "<未提供任何错误信息>"
	}

	format := ""
	switch v := msg1[0].(type) {
	case []byte:
		format = string(v)
	case []rune:
		format = string(v)
	case string:
		format = v
	case fmt.Stringer:
		format = v.String()
	default:
		return "<无法正确转换错误提示信息>"
	}

	return fmt.Sprintf(format, msg1[1:]...)
}

// 判断一个值是否为空(0, "", false, 空数组等)。
// []string{""}空数组里套一个空字符串，不会被判断为空。
func isEmpty(expr interface{}) bool {
	if expr == nil {
		return true
	}

	switch v := expr.(type) {
	case bool:
		return !v
	case int:
		return 0 == v
	case int8:
		return 0 == v
	case int16:
		return 0 == v
	case int32:
		return 0 == v
	case int64:
		return 0 == v
	case uint:
		return 0 == v
	case uint8:
		return 0 == v
	case uint16:
		return 0 == v
	case uint32:
		return 0 == v
	case uint64:
		return 0 == v
	case string:
		return len(v) == 0
	case float32:
		return 0 == v
	case float64:
		return 0 == v
	case time.Time:
		return v.IsZero()
	case *time.Time:
		return v.IsZero()
	}

	// 符合IsNil条件的，都为Empty
	if isNil(expr) {
		return true
	}

	// 长度为0的数组也是empty
	v := reflect.ValueOf(expr)
	switch v.Kind() {
	case reflect.Slice, reflect.Map, reflect.Chan:
		return 0 == v.Len()
	}

	return false
}

// 判断一个值是否为nil。
// 当特定类型的变量，已经声明，但还未赋值时，也将返回true
func isNil(expr interface{}) bool {
	if nil == expr {
		return true
	}

	v := reflect.ValueOf(expr)
	k := v.Kind()

	return (k == reflect.Chan ||
		k == reflect.Func ||
		k == reflect.Interface ||
		k == reflect.Map ||
		k == reflect.Ptr ||
		k == reflect.Slice) &&
		v.IsNil()
}

// 判断两个值是否相等。
//
// 除了通过reflect.DeepEqual()判断值是否相等之外，一些类似
// 可转换的数值也能正确判断，比如以下值也将会被判断为相等：
//  int8(5)                     == int(5)
//  []int{1,2}                  == []int8{1,2}
//  []int{1,2}                  == [2]int8{1,2}
//  []int{1,2}                  == []float32{1,2}
//  map[string]int{"1":"2":2}   == map[string]int8{"1":1,"2":2}
//
//  // map的键值不同，即使可相互转换也判断不相等。
//  map[int]int{1:1,2:2}        != map[int8]int{1:1,2:2}
func isEqual(v1, v2 interface{}) bool {
	if reflect.DeepEqual(v1, v2) {
		return true
	}

	vv1 := reflect.ValueOf(v1)
	vv2 := reflect.ValueOf(v2)

	// NOTE: 这里返回false，而不是true
	if !vv1.IsValid() || !vv2.IsValid() {
		return false
	}

	if vv1 == vv2 {
		return true
	}

	vv1Type := vv1.Type()
	vv2Type := vv2.Type()

	// 过滤掉已经在reflect.DeepEqual()进行处理的类型
	switch vv1Type.Kind() {
	case reflect.Struct, reflect.Ptr, reflect.Func, reflect.Interface:
		return false
	case reflect.Slice, reflect.Array:
		// vv2.Kind()与vv1的不相同
		if vv2.Kind() != reflect.Slice && vv2.Kind() != reflect.Array {
			// 虽然类型不同，但可以相互转换成vv1的，如：vv2是string，vv2是[]byte，
			if vv2Type.ConvertibleTo(vv1Type) {
				return isEqual(vv1.Interface(), vv2.Convert(vv1Type).Interface())
			}
			return false
		}

		// reflect.DeepEqual()未考虑类型不同但是类型可转换的情况，比如：
		// []int{8,9} == []int8{8,9}，此处重新对slice和array做比较处理。
		if vv1.Len() != vv2.Len() {
			return false
		}

		for i := 0; i < vv1.Len(); i++ {
			if !isEqual(vv1.Index(i).Interface(), vv2.Index(i).Interface()) {
				return false
			}
		}
		return true // for中所有的值比较都相等，返回true
	case reflect.Map:
		if vv2.Kind() != reflect.Map {
			return false
		}

		if vv1.IsNil() != vv2.IsNil() {
			return false
		}
		if vv1.Len() != vv2.Len() {
			return false
		}
		if vv1.Pointer() == vv2.Pointer() {
			return true
		}

		// 两个map的键名类型不同
		if vv2Type.Key().Kind() != vv1Type.Key().Kind() {
			return false
		}

		for _, index := range vv1.MapKeys() {
			vv2Index := vv2.MapIndex(index)
			if !vv2Index.IsValid() {
				return false
			}

			if !isEqual(vv1.MapIndex(index).Interface(), vv2Index.Interface()) {
				return false
			}
		}
		return true // for中所有的值比较都相等，返回true
	case reflect.String:
		if vv2.Kind() == reflect.String {
			return vv1.String() == vv2.String()
		}
		if vv2Type.ConvertibleTo(vv1Type) { // 考虑v1是string，v2是[]byte的情况
			return isEqual(vv1.Interface(), vv2.Convert(vv1Type).Interface())
		}

		return false
	}

	if vv1Type.ConvertibleTo(vv2Type) {
		return vv2.Interface() == vv1.Convert(vv2Type).Interface()
	} else if vv2Type.ConvertibleTo(vv1Type) {
		return vv1.Interface() == vv2.Convert(vv1Type).Interface()
	}

	return false
}

// 判断fn函数是否会发生panic
// 若发生了panic，将把msg一起返回。
func hasPanic(fn func()) (has bool, msg interface{}) {
	defer func() {
		if msg = recover(); msg != nil {
			has = true
		}
	}()
	fn()

	return
}

// 判断container是否包含了item的内容。若是指针，会判断指针指向的内容，
// 但是不支持多重指针。
//
// 若container是字符串(string、[]byte和[]rune，不包含fmt.Stringer接口)，
// 都将会以字符串的形式判断其是否包含item。
// 若container是个列表(array、slice、map)则判断其元素中是否包含item中的
// 的所有项，或是item本身就是container中的一个元素。
func isContains(container, item interface{}) bool {
	if container == nil { // nil不包含任何东西
		return false
	}

	cv := reflect.ValueOf(container)
	iv := reflect.ValueOf(item)

	if cv.Kind() == reflect.Ptr {
		cv = cv.Elem()
	}

	if iv.Kind() == reflect.Ptr {
		iv = iv.Elem()
	}

	if isEqual(container, item) {
		return true
	}

	// 判断是字符串的情况
	switch c := cv.Interface().(type) {
	case string:
		switch i := iv.Interface().(type) {
		case string:
			return strings.Contains(c, i)
		case []byte:
			return strings.Contains(c, string(i))
		case []rune:
			return strings.Contains(c, string(i))
		case byte:
			return bytes.IndexByte([]byte(c), i) != -1
		case rune:
			return bytes.IndexRune([]byte(c), i) != -1
		}
	case []byte:
		switch i := iv.Interface().(type) {
		case string:
			return bytes.Contains(c, []byte(i))
		case []byte:
			return bytes.Contains(c, i)
		case []rune:
			return strings.Contains(string(c), string(i))
		case byte:
			return bytes.IndexByte(c, i) != -1
		case rune:
			return bytes.IndexRune(c, i) != -1
		}
	case []rune:
		switch i := iv.Interface().(type) {
		case string:
			return strings.Contains(string(c), string(i))
		case []byte:
			return strings.Contains(string(c), string(i))
		case []rune:
			return strings.Contains(string(c), string(i))
		case byte:
			return strings.IndexByte(string(c), i) != -1
		case rune:
			return strings.IndexRune(string(c), i) != -1
		}
	}

	if (cv.Kind() == reflect.Slice) || (cv.Kind() == reflect.Array) {
		if !cv.IsValid() || cv.Len() == 0 { // 空的，就不算包含另一个，即使另一个也是空值。
			return false
		}

		if !iv.IsValid() {
			return false
		}

		// item是container的一个元素
		for i := 0; i < cv.Len(); i++ {
			if isEqual(cv.Index(i).Interface(), iv.Interface()) {
				return true
			}
		}

		// 开始判断item的元素是否与container中的元素相等。

		// 若item的长度为0，表示不包含
		if (iv.Kind() != reflect.Slice) || (iv.Len() == 0) {
			return false
		}

		// item的元素比container的元素多，必须在判断完item不是container中的一个元素之
		if iv.Len() > cv.Len() {
			return false
		}

		// 依次比较item的各个子元素是否都存在于container，且下标都相同
		ivIndex := 0
		for i := 0; i < cv.Len(); i++ {
			if isEqual(cv.Index(i).Interface(), iv.Index(ivIndex).Interface()) {
				if (ivIndex == 0) && (i+iv.Len() > cv.Len()) {
					return false
				}
				ivIndex++
				if ivIndex == iv.Len() { // 已经遍历完iv
					return true
				}
			} else if ivIndex > 0 {
				return false
			}
		}
		return false
	} // end cv.Kind == reflect.Slice and reflect.Array

	if cv.Kind() == reflect.Map {
		if cv.Len() == 0 {
			return false
		}

		if (iv.Kind() != reflect.Map) || (iv.Len() == 0) {
			return false
		}

		if iv.Len() > cv.Len() {
			return false
		}

		// 判断所有item的项都存在于container中
		for _, key := range iv.MapKeys() {
			cvItem := iv.MapIndex(key)
			if !cvItem.IsValid() { // container中不包含该值。
				return false
			}
			if !isEqual(cvItem.Interface(), iv.MapIndex(key).Interface()) {
				return false
			}
		}
		// for中的所有判断都成立，返回true
		return true
	}

	return false
}
