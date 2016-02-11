package main
import "eventsourcing"

type E struct {
	eventsourcing.BaseEvent
	v int
}

func main() {
	var es = eventsourcing.NewStore()
	es.Save(new(E))
}
