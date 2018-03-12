package plug

import (
	"reflect"
)

// CheckUnimplimented will check if interface obj has all the functions asked for by i, it will then
// output a list of the functions not implimented.
// Please note that obj has to be a pointer to the struct in question and not the struct object itself.
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
