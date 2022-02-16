package events

type Event interface {
	String() string
	Is(Event) bool
}

type GameEvent struct {
	id   int
	Name string
	Tick int
	Data map[string]string
}

func (ge GameEvent) String() string {
	return ge.Name
}

func (ge GameEvent) Is(e Event) bool {
	return ge.String() == e.String()
}
