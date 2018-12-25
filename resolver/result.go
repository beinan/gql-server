package resolver

import (
	"bytes"
	"encoding/json"
	"fmt"
)

//GqlResult represents a single field (alias & value pair)
type GqlResult struct {
	Alias string
	Value GqlResultValue
}

//GqlResults represent a json object
type GqlResults []GqlResult

//MarshalJSON : result array to json object
func (results GqlResults) MarshalJSON() ([]byte, error) {
	buffer := bytes.NewBufferString("{")
	length := len(results)
	for i, value := range results {
		jsonValue, err := json.Marshal(value.Value)
		if err != nil {
			return nil, err
		}
		buffer.WriteString(fmt.Sprintf("\"%s\":%s", value.Alias, string(jsonValue)))
		if i < length-1 {
			buffer.WriteString(",")
		}
	}
	buffer.WriteString("}")
	return buffer.Bytes(), nil
}
