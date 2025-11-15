package websocket

import (
	"encoding/json"
	myErrors "gochat/internal/shared/errors"
)

func EncodePong() ([]byte, error) {
	return json.Marshal(&Message{
		Type: PongType,
	})
}

func EncodeResponse(resp *Response) ([]byte, error) {
	payload, err := json.Marshal(resp)
	if err != nil {
		return nil, err
	}
	return json.Marshal(&Message{
		Type:    ResponseType,
		Payload: payload,
	})
}

func DecodeMessage(b []byte) (*Message, error) {
	m := new(Message)
	if err := json.Unmarshal(b, m); err != nil {
		return nil, myErrors.ErrDecode
	}
	return m, nil
}

func DecodeRequest(payload json.RawMessage) (*Request, error) {
	r := new(Request)
	if err := json.Unmarshal(payload, r); err != nil {
		return nil, myErrors.ErrDecode
	}
	return r, nil
}
