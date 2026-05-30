package event

import (
	"encoding/json"
)

var _ SpecificEvent = (*StandardCommandFailedEvent)(nil)

type StandardCommandFailedEvent struct {
	errorMessage string
	*StandardEvent
}

func NewStandardCommandFailedEvent(
	originAction SpecificEvent,
	err error,
	generator IDGenerator,
) (*StandardCommandFailedEvent, error) {
	e := &StandardCommandFailedEvent{
		errorMessage: err.Error(),
	}
	payload, err := e.Marshal()
	if err != nil {
		return nil, err
	}

	e.StandardEvent = NewStandardEvent(originAction.AggregateID(), originAction.Topic()+".failed", payload, generator)
	e.AddHeaders(originAction.Headers())
	return e, nil
}

func (e *StandardCommandFailedEvent) ErrorMessage() string {
	return e.errorMessage
}

func (e *StandardCommandFailedEvent) Marshal() ([]byte, error) {
	type Alias struct {
		ErrorMessage string `json:"errorMessage"`
	}
	return json.Marshal(&Alias{
		ErrorMessage: e.errorMessage,
	})
}

func (e *StandardCommandFailedEvent) Unmarshal(data []byte) error {
	type Alias struct {
		ErrorMessage string
	}
	var tmp Alias
	if err := json.Unmarshal(data, &tmp); err != nil {
		return err
	}
	e.errorMessage = tmp.ErrorMessage
	return nil
}
