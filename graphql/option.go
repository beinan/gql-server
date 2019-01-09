package graphql

import "encoding/json"

type IntOption struct {
	Value int
	IsSet bool
}

type StringOption struct {
	Value string
	IsSet bool
}

type FloatOption struct {
	Value float32
	IsSet bool
}

type IDOption struct {
	Value string
	IsSet bool
}

type BooleanOption struct {
	Value bool
	IsSet bool
}

func (o StringOption) MarshalJSON() ([]byte, error) {
	if o.IsSet {
		return json.Marshal(o.Value)
	}
	return json.Marshal(nil)
}

func (o IntOption) MarshalJSON() ([]byte, error) {
	if o.IsSet {
		return json.Marshal(o.Value)
	}
	return json.Marshal(nil)
}

func (o BooleanOption) MarshalJSON() ([]byte, error) {
	if o.IsSet {
		return json.Marshal(o.Value)
	}
	return json.Marshal(nil)
}

func (o FloatOption) MarshalJSON() ([]byte, error) {
	if o.IsSet {
		return json.Marshal(o.Value)
	}
	return json.Marshal(nil)
}

func (o IDOption) MarshalJSON() ([]byte, error) {
	if o.IsSet {
		return json.Marshal(o.Value)
	}
	return json.Marshal(nil)
}
