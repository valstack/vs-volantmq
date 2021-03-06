package connection

import (
	"sync"

	"github.com/VolantMQ/vlapi/mqttp"
)

type onRelease func(o, n mqttp.IFace)

type ackQueue struct {
	messages  sync.Map
	onRelease onRelease
}

func (a *ackQueue) store(pkt mqttp.IFace, replace bool) bool {
	id, _ := pkt.ID()

	if _, ok := a.messages.Load(id); ok && !replace {
		return false
	}

	a.messages.Store(id, pkt)

	return true
}

func (a *ackQueue) release(pkt mqttp.IFace) bool {
	id, _ := pkt.ID()

	if value, ok := a.messages.Load(id); ok {
		if orig, k := value.(mqttp.IFace); k && a.onRelease != nil {
			a.onRelease(orig, pkt)
		}
		a.messages.Delete(id)

		return true
	}

	return false
}
