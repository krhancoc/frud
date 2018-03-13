package main

type ModelOnly struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

func Create(init map[string]string) interface{} {

	return ModelOnly{
		Id:   init["Id"],
		Name: init["Name"],
	}
}

func Modify(init map[string]string) interface{} {

	return ModelOnly{
		Id:   init["Id"],
		Name: init["Name"],
	}
}

func Delete(init map[string]string) interface{} {

	return ModelOnly{
		Id:   init["Id"],
		Name: init["Name"],
	}
}

func Read(init map[string]string) interface{} {

	return ModelOnly{
		Id:   init["Id"],
		Name: init["Name"],
	}
}
