package gtf

import (
	"fmt"
	"os"
	// "path"
	"reflect"
	// "runtime"
	// "strconv"
	"strings"
)

// assert that the expression expr is true
func (t *Test) True(expr bool, args ...interface{}) {
	assert(expr, args, []interface{}{"The expression is not true, it is [%T:%[1]v]", expr})
}

// assert that the expression expr is false
func (t *Test) False(expr bool, args ...interface{}) {
	assert(!expr, args, []interface{}{"The expression is not false, it is [%T:%[1]v]", expr})
}

// assert that the expression expr is nil
func (t *Test) Nil(expr interface{}, args ...interface{}) {
	assert(isNil(expr), args, []interface{}{"Nil failed，actual value:[%T:%[1]v]", expr})
}

// assert that the expression expr is not nil
func (t *Test) NotNil(expr interface{}, args ...interface{}) {
	assert(!isNil(expr), args, []interface{}{"NotNil failed，actual value: [%T:%[1]v]", expr})
}

// assert that actual and expected is equal
func (t *Test) Equal(actual, expected interface{}, args ...interface{}) {
	assert(isEqual(actual, expected), args, []interface{}{"Equal failed, actual: [%T:%[1]v]; expected: [%T:%[2]v]", actual, expected})
}

// assert that actual and expected is not equal
func (t *Test) NotEqual(actual, expected interface{}, args ...interface{}) {
	assert(!isEqual(actual, expected), args, []interface{}{"NotEqual failed, actual: [%T:%[1]v]; expected: [%T:%[2]v]", actual, expected})
}

// assert that the expression expr is empty(nil,"",0,false)
func (t *Test) Empty(expr interface{}, args ...interface{}) {
	assert(isEmpty(expr), args, []interface{}{"The expression expr is not empty(nil,\"\",0,false), it is [%T:%[1]v]", expr})
}

// assert that the expression expr is empty(nil,"",0,false)
func (t *Test) NotEmpty(expr interface{}, args ...interface{}) {
	assert(!isEmpty(expr), args, []interface{}{"he expression expr is empty(nil,\"\",0,false), it is [%T:%[1]v]", expr})
}

// assert that the file exists
func (t *Test) FileExists(path string, args ...interface{}) {
	_, err := os.Stat(path)

	if err != nil && !os.IsExist(err) {
		assert(false, args, []interface{}{"Eorro Info：%v", err.Error()})
	}
}

// assert that the file does not exists
func (t *Test) FileNotExists(path string, args ...interface{}) {
	_, err := os.Stat(path)
	assert(os.IsNotExist(err), args, []interface{}{"Eorro Info：%v", err.Error()})
}

// assert that the container includes item or includes all the elements of the item
// for detail, please refer to IsContains()
func (t *Test) Contains(container, item interface{}, args ...interface{}) {
	assert(isContains(container, item), args, []interface{}{"container:[%v] does not contain item: [%v]", container, item})
}

// assert that the container does not include item or does not include all the elements of the item
func (t *Test) NotContains(container, item interface{}, args ...interface{}) {
	assert(!isContains(container, item), args, []interface{}{"container:[%v] contains item: [%v]", container, item})
}

// 断言函数会发生panic，否则输出错误信息。
func (t *Test) Panic(fn func(), args ...interface{}) {
	has, _ := hasPanic(fn)
	assert(has, args, []interface{}{"并未发生panic"})
}

// 断言函数会发生panic，且panic信息中包含指定的字符串内容，否则输出错误信息。
func (t *Test) PanicString(fn func(), str string, args ...interface{}) {
	if has, msg := hasPanic(fn); has {
		index := strings.Index(fmt.Sprint(msg), str)
		assert(index >= 0, args, []interface{}{"并未发生panic"})
	}
}

// 断言函数会发生panic，且panic返回的类型与typ的类型相同。
func (t *Test) PanicType(fn func(), typ interface{}, args ...interface{}) {
	has, msg := hasPanic(fn)
	if !has {
		return
	}

	t1 := reflect.TypeOf(msg)
	t2 := reflect.TypeOf(typ)
	assert(t1 == t2, args, []interface{}{"PanicType失败，v1[%v]的类型与v2[%v]的类型不相同", t1, t2})

}

// 断言函数会发生panic，否则输出错误信息。
func (t *Test) NotPanic(fn func(), args ...interface{}) {
	has, msg := hasPanic(fn)
	assert(!has, args, []interface{}{"发生了panic，其信息为[%v]", msg})
}
