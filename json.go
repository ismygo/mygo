package mygo

import "encoding/json"

func (*GoJSON) UnmarshalJSON(data []byte, v interface{}) error {
	return json.Unmarshal(data, v)
}

func (*GoJSON) MarshalJSON(v interface{}) ([]byte, error) {
	return json.Marshal(v)
}

func (*GoJSON) MarshalIndentJSON(v interface{}, prefix, indent string) ([]byte, error) {
	return json.MarshalIndent(v, prefix, indent)
}
