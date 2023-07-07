package slice

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestJsonCompare(t *testing.T) {
	a := "{\"id\":1,\"change\":{\"name\":\"aa\",\"age\":12}}"

	b := "{\"id\":2,\"change\":{\"name\":\"bb\",\"age\":12}}"
	var (
		json1 map[string]interface{}
		json2 map[string]interface{}
	)
	_ = json.Unmarshal([]byte(a), &json1)
	_ = json.Unmarshal([]byte(b), &json2)
	res := JsonCompare(json1, json2)
	fmt.Println(res)
}
