package yaml

import (
	"fmt"

	"gopkg.in/yaml.v2"
)

// From converts the json message to a map string interface
func From(msg interface{}) (interface{}, error) {
	var mii map[interface{}]interface{}
	err := yaml.Unmarshal([]byte(msg.(fmt.Stringer).String()), &mii)
	return mii, err
}
