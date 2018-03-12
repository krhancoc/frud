package plug

import (
	"reflect"
)

func CheckUnimplimented(obj interface{}, i interface{}) []string {
	var unimplemented []string
	t := reflect.TypeOf(i).Elem()
	for i := 0; i < t.NumMethod(); i++ {
		f := t.Method(i).Name
		check := reflect.TypeOf(obj)
		_, ok := check.MethodByName(f)
		if !ok {
			unimplemented = append(unimplemented, f)
		}
	}
	return unimplemented
}

// func (plug *Plug) CheckUnimplimented(i interface{}) []string {

// 	var unimplemented []string
// 	t := reflect.TypeOf(i).Elem()
// 	for i := 0; i < t.NumMethod(); i++ {
// 		f := t.Method(i).Name
// 		check := reflect.TypeOf(*plug.Main)
// 		_, ok := check.MethodByName(f)
// 		if !ok {
// 			unimplemented = append(unimplemented, f)
// 		}
// 	}
// 	return unimplemented
// }
