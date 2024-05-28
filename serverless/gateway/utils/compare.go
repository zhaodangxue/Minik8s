package serveless_utils

import "github.com/tidwall/gjson"

func IntegerEqual(data string, variable string, result interface{}) bool {
	res := gjson.Get(data, variable)
	return res.Int() == result.(int64)
}
func IntegerNotEqual(data string, variable string, result interface{}) bool {
	res := gjson.Get(data, variable)
	return res.Int() != result.(int64)
}
func IntegerGreaterThan(data string, variable string, result interface{}) bool {
	res := gjson.Get(data, variable)
	return res.Int() > result.(int64)
}
func IntegerLessThan(data string, variable string, result interface{}) bool {
	res := gjson.Get(data, variable)
	return res.Int() < result.(int64)
}
func BooleanEqual(data string, variable string, result interface{}) bool {
	res := gjson.Get(data, variable)
	return res.Bool() == result.(bool)
}
func BooleanNotEqual(data string, variable string, result interface{}) bool {
	res := gjson.Get(data, variable)
	return res.Bool() != result.(bool)
}
func StringEqual(data string, variable string, result interface{}) bool {
	res := gjson.Get(data, variable)
	return res.String() == result.(string)
}
func StringNotEqual(data string, variable string, result interface{}) bool {
	res := gjson.Get(data, variable)
	return res.String() != result.(string)
}
func FloatEqual(data string, variable string, result interface{}) bool {
	res := gjson.Get(data, variable)
	return res.Float() == result.(float64)
}
func FloatNotEqual(data string, variable string, result interface{}) bool {
	res := gjson.Get(data, variable)
	return res.Float() != result
}
func FloatGreaterThan(data string, variable string, result interface{}) bool {
	res := gjson.Get(data, variable)
	return res.Float() > result.(float64)
}
func FloatLessThan(data string, variable string, result interface{}) bool {
	res := gjson.Get(data, variable)
	return res.Float() < result.(float64)
}
