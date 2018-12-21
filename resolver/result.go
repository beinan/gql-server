package resolver

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/beinan/gql-server/concurrent/future"
)

type Result struct {
	Alias       string
	FutureValue future.Future
}

type Results []Result

//result array to json object
func (this Results) MarshalJSON() ([]byte, error) {
	buffer := bytes.NewBufferString("{")
	length := len(this)
	for i, value := range this {
		jsonValue, err := json.Marshal(value.FutureValue)
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
