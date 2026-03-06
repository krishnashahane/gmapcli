package maps

import "fmt"

var ErrNoAPIKey = fmt.Errorf("googlemapscli: api key is not set")

type InputError struct {
	Param   string
	Reason  string
}

func (e InputError) Error() string {
	return fmt.Sprintf("googlemapscli: bad %s: %s", e.Param, e.Reason)
}

type RemoteError struct {
	Code    int
	Payload string
}

func (e *RemoteError) Error() string {
	if e.Payload == "" {
		return fmt.Sprintf("googlemapscli: server returned %d", e.Code)
	}
	return fmt.Sprintf("googlemapscli: server returned %d: %s", e.Code, e.Payload)
}
