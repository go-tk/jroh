package fooapi

import (
	"encoding/json"
	"errors"
)

var _ json.Marshaler = (*MyStructString)(nil)

func (m *MyStructString) MarshalJSON() ([]byte, error) {
	if m.TheStringA == "taboo" {
		return nil, errors.New("bad word")
	}
	type T MyStructString
	return json.Marshal((*T)(m))
}
