package fork

import "encoding/json"

type Marshaller interface {
	Decode(b []byte, v interface{}) error
	Encode(v interface{}) ([]byte, error)
}

type JsonMarshal struct {
}

func (j *JsonMarshal) Decode(b []byte, v interface{}) error {
	return json.Unmarshal(b, v)
}

func (j *JsonMarshal) Encode(v interface{}) ([]byte, error) {
	return json.Marshal(v)
}

var (
	DefaultMarshal = &JsonMarshal{}
)
