package event

var _ SpecificEvent = (*StandardCommandSucceedEvent)(nil)

type StandardCommandSucceedEvent struct {
	*StandardEvent
}

func NewStandardCommandSucceedEvent(
	originAction SpecificEvent,
	generator IDGenerator,
) (*StandardCommandSucceedEvent, error) {
	e := &StandardCommandSucceedEvent{}
	payload, err := e.Marshal()
	if err != nil {
		return nil, err
	}

	e.StandardEvent = NewStandardEvent(originAction.AggregateID(), originAction.Topic()+".succeeded", payload, generator)
	e.AddHeaders(originAction.Headers())
	return e, nil
}

func (e *StandardCommandSucceedEvent) Marshal() ([]byte, error) {
	return []byte(""), nil
}

func (e *StandardCommandSucceedEvent) Unmarshal([]byte) error {
	return nil
}
