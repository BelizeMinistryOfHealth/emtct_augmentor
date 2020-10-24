package http

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
)

func parseBody(body io.ReadCloser, o interface{}) error {
	b, err := ioutil.ReadAll(body)
	defer body.Close()
	if err != nil {
		return fmt.Errorf("could not parse the body posted: %v", err)
	}
	err = json.Unmarshal(b, &o)
	if err != nil {
		return err
	}
	return nil
}
