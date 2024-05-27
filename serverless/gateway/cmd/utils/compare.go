package serveless_utils

import "github.com/tidwall/gjson"

func IntegerEqual(data string, variable string, result int) bool {
	res := gjson.Get(data, variable)
	return res.Int() == int64(result)
}
func IntegerNotEqual(data string, variable string, result int) bool {
	res := gjson.Get(data, variable)
	return res.Int() != int64(result)
}
func IntegerGreaterThan(data string, variable string, result int) bool {
	res := gjson.Get(data, variable)
	return res.Int() > int64(result)
}
func IntegerLessThan(data string, variable string, result int) bool {
	res := gjson.Get(data, variable)
	return res.Int() < int64(result)
}
func BooleanEqual(data string, variable string, result bool) bool {
	res := gjson.Get(data, variable)
	return res.Bool() == result
}
func BooleanNotEqual(data string, variable string, result bool) bool {
	res := gjson.Get(data, variable)
	return res.Bool() != result
}
func StringEqual(data string, variable string, result string) bool {
	res := gjson.Get(data, variable)
	return res.String() == result
}
func StringNotEqual(data string, variable string, result string) bool {
	res := gjson.Get(data, variable)
	return res.String() != result
}
func FloatEqual(data string, variable string, result float64) bool {
	res := gjson.Get(data, variable)
	return res.Float() == result
}
func FloatNotEqual(data string, variable string, result float64) bool {
	res := gjson.Get(data, variable)
	return res.Float() != result
}
func FloatGreaterThan(data string, variable string, result float64) bool {
	res := gjson.Get(data, variable)
	return res.Float() > result
}
func FloatLessThan(data string, variable string, result float64) bool {
	res := gjson.Get(data, variable)
	return res.Float() < result
}
