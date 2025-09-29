package decoder

import (
	"encoding/json"

	"github.com/BurntSushi/toml"
	"gopkg.in/yaml.v3"
)

type IDecoder interface {
	Decode(data []byte, v interface{}) error
}

type TOMLDecoder struct{}

func (d TOMLDecoder) Decode(data []byte, v interface{}) error {
	return toml.Unmarshal(data, v)
}

type YAMLDecoder struct{}

func (d YAMLDecoder) Decode(data []byte, v interface{}) error {
	return yaml.Unmarshal(data, v)
}

type JSONDecoder struct{}

func (d JSONDecoder) Decode(data []byte, v interface{}) error {
	return json.Unmarshal(data, v)
}
