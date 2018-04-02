package util

import (
	"encoding/json"
	"fmt"
)

func DUMP(o interface{}) {
	fmt.Printf("%+v\n", o)
}

func DUMP_JSON(o interface{}) {
	data, err := json.Marshal(o)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%s\n", data)
}
