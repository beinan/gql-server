package graphql

import "encoding/json"

type IntOption interface {
	Get() int
	IsSet() bool
	GetOrElse(int) int
}

type StringOption interface {
	Get() string
	IsSet() bool
	GetOrElse(string) string
}

type intOptionStruct struct {
	value int
	isSet bool
}

func NewIntOption(v int) IntOption {
	return intOptionStruct{
		value: v,
		isSet: true,
	}
}

var EmptyIntOption = intOptionStruct{
	isSet: false,
}

func (o intOptionStruct) Get() int {
	if !o.isSet {
		panic("IntOption value not set.")
	}
	return o.value
}

func (o intOptionStruct) GetOrElse(v int) int {
	if o.isSet {
		return o.value
	}
	return v
}

func (o intOptionStruct) IsSet() bool {
	return o.isSet
}

type stringOptionStruct struct {
	value string
	isSet bool
}

func NewStringOption(v string) StringOption {
	return stringOptionStruct{
		value: v,
		isSet: true,
	}
}

var EmptyStringOption = stringOptionStruct{
	isSet: false,
}

func (o stringOptionStruct) Get() string {
	if !o.isSet {
		panic("StringOption value not set.")
	}
	return o.value
}

func (o stringOptionStruct) GetOrElse(v string) string {
	if o.isSet {
		return o.value
	}
	return v
}

func (o stringOptionStruct) IsSet() bool {
	return o.isSet
}

func (o stringOptionStruct) MarshalJSON() ([]byte, error) {
	if o.isSet {
		return json.Marshal(o.value)
	}
	return json.Marshal(nil)
}
