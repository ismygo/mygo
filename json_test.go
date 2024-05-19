package mygo

import "testing"

func TestMarshalJSON(t *testing.T) {
	data, err := JSON.MarshalJSON(map[string]string{"foo": "bar"})
	if nil != err {
		t.Fail()
	}

	if "{\"foo\":\"bar\"}" != (string(data)) {
		t.Fail()
	}
}

func Test(t *testing.T) {
	m := map[string]string{}
	err := JSON.UnmarshalJSON([]byte(`{"foo":"bar"}`), &m)
	if nil != err {
		t.Fail()
	}

	if "bar" != m["foo"] {
		t.Fail()
	}
}

func TestMarshalIndentJSON(t *testing.T) {
	data, err := JSON.MarshalIndentJSON(map[string]string{"foo": "bar"}, "", "  ")
	if nil != err {
		t.Fail()
	}

	logger.Debug(string(data))
	if "{\n  \"foo\": \"bar\"\n}" != (string(data)) {
		t.Fail()
	}
}
