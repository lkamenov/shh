package pollers

import (
	"bytes"
	"github.com/freeformz/shh/mm"
	"github.com/freeformz/shh/utils"
	"io/ioutil"
	"time"
)

const (
	CONNTRACK_DATA = "/proc/sys/net/netfilter/nf_conntrack_count"
)

type Conntrack struct {
	measurements chan<- *mm.Measurement
}

func NewConntrackPoller(measurements chan<- *mm.Measurement) Conntrack {
	return Conntrack{measurements: measurements}
}

func (poller Conntrack) Poll(tick time.Time) {
	ctx := utils.Slog{"poller": poller.Name(), "fn": "Poll", "tick": tick}

	count, err := ioutil.ReadFile(CONNTRACK_DATA)
	if err != nil {
		ctx.Error(err, "reading "+CONNTRACK_DATA)
	}

	poller.measurements <- &mm.Measurement{tick, poller.Name(), []string{"count"}, utils.Atofloat64(string(bytes.TrimSpace(count)))}
}

func (poller Conntrack) Name() string {
	return "conntrack"
}

func (poller Conntrack) Exit() {}
