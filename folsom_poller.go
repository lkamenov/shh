package shh

import (
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/heroku/slog"
)

type FolsomMemory struct {
	Total         uint64 `json:"total"`
	Processes     uint64 `json:"processes"`
	ProcessesUsed uint64 `json:"processes_used"`
	System        uint64 `json:"system"`
	Atom          uint64 `json:"atom"`
	AtomUsed      uint64 `json:"atom_used"`
	Binary        uint64 `json:"binary"`
	Code          uint64 `json:"code"`
	Ets           uint64 `json:"ets"`
}

type FolsomStatistics struct {
	ContextSwitches   uint64                  `json:"context_switches"`
	GarbageCollection FolsomGarbageCollection `json:"garbage_collection"`
	Io                FolsomIo                `json:"io"`
	Reductions        FolsomReductions        `json:"reductions"`
	RunQueue          uint64                  `json:"run_queue"`
	Runtime           FolsomRuntime           `json:"runtime"`
	WallClock         FolsomWallClock         `json:"wall_clock"`
}

type FolsomGarbageCollection struct {
	NumOfGcs       uint64 `json:"number_of_gcs"`
	WordsReclaimed uint64 `json:"words_reclaimed"`
}

type FolsomIo struct {
	Input  uint64 `json:"input"`
	Output uint64 `json:"output"`
}

type FolsomReductions struct {
	Total     uint64 `json:"total_reductions"`
	SinceLast uint64 `json:"reductions_since_last_call"`
}

type FolsomRuntime struct {
	Total     uint64 `json:"total_run_time"`
	SinceLast uint64 `json:"time_since_last_call"`
}

type FolsomWallClock struct {
	Total     uint64 `json:"total_wall_clock_time"`
	SinceLast uint64 `json:"wall_clock_time_since_last_call"`
}

type FolsomType struct {
	Type string `json:"type"`
}

type FolsomValue struct {
	Name  string
	Type  string
	Value json.Number `json:"value"`
}

type FolsomPoller struct {
	measurements chan<- Measurement
	baseUrl      string
	client       *http.Client
}

func NewFolsomPoller(measurements chan<- Measurement, config Config) FolsomPoller {
	var url string

	if config.FolsomBaseUrl != nil {
		url = config.FolsomBaseUrl.String()
	}

	client := &http.Client{
		Transport: &http.Transport{
			ResponseHeaderTimeout: config.NetworkTimeout,
			Dial: func(network, address string) (net.Conn, error) {
				return net.DialTimeout(network, address, config.NetworkTimeout)
			},
		},
	}

	return FolsomPoller{
		measurements: measurements,
		baseUrl:      url,
		client:       client,
	}
}

func (poller FolsomPoller) Poll(tick time.Time) {
	if poller.baseUrl == "" {
		return
	}

	ctx := slog.Context{"poller": poller.Name(), "fn": "Poll", "tick": tick}

	poller.doMemoryPoll(ctx, tick)
	poller.doStatisticsPoll(ctx, tick)
	poller.doMetricsPoll(ctx, tick)
}

func (poller FolsomPoller) doMemoryPoll(ctx slog.Context, tick time.Time) {
	memory := FolsomMemory{}

	if err := poller.decodeReq("/_memory", &memory); err != nil {
		LogError(ctx, err, "while performing request for this tick")
		return
	}

	poller.measurements <- GaugeMeasurement{tick, poller.Name(), []string{"mem", "total"}, memory.Total, Bytes}
	poller.measurements <- GaugeMeasurement{tick, poller.Name(), []string{"mem", "procs", "total"}, memory.Processes, Bytes}
	poller.measurements <- GaugeMeasurement{tick, poller.Name(), []string{"mem", "procs", "used"}, memory.ProcessesUsed, Bytes}
	poller.measurements <- GaugeMeasurement{tick, poller.Name(), []string{"mem", "system"}, memory.System, Bytes}
	poller.measurements <- GaugeMeasurement{tick, poller.Name(), []string{"mem", "atom", "total"}, memory.Atom, Bytes}
	poller.measurements <- GaugeMeasurement{tick, poller.Name(), []string{"mem", "atom", "used"}, memory.AtomUsed, Bytes}
	poller.measurements <- GaugeMeasurement{tick, poller.Name(), []string{"mem", "binary"}, memory.Binary, Bytes}
	poller.measurements <- GaugeMeasurement{tick, poller.Name(), []string{"mem", "code"}, memory.Code, Bytes}
	poller.measurements <- GaugeMeasurement{tick, poller.Name(), []string{"mem", "ets"}, memory.Ets, Bytes}
}

func (poller FolsomPoller) doStatisticsPoll(ctx slog.Context, tick time.Time) {
	stats := FolsomStatistics{}
	if err := poller.decodeReq("/_statistics", &stats); err != nil {
		LogError(ctx, err, "while performing request for this tick")
		return
	}

	poller.measurements <- CounterMeasurement{tick, poller.Name(), []string{"stats", "context-switches"}, stats.ContextSwitches, ContextSwitches}
	poller.measurements <- CounterMeasurement{tick, poller.Name(), []string{"stats", "gc", "num"}, stats.GarbageCollection.NumOfGcs, Empty}
	poller.measurements <- CounterMeasurement{tick, poller.Name(), []string{"stats", "gc", "reclaimed"}, stats.GarbageCollection.WordsReclaimed, Words}
	poller.measurements <- CounterMeasurement{tick, poller.Name(), []string{"stats", "io", "input"}, stats.Io.Input, Bytes}
	poller.measurements <- CounterMeasurement{tick, poller.Name(), []string{"stats", "io", "output"}, stats.Io.Output, Bytes}
	poller.measurements <- GaugeMeasurement{tick, poller.Name(), []string{"stats", "reductions"}, stats.Reductions.SinceLast, Reductions}
	poller.measurements <- GaugeMeasurement{tick, poller.Name(), []string{"stats", "run-queue"}, stats.RunQueue, Processes}
	poller.measurements <- GaugeMeasurement{tick, poller.Name(), []string{"stats", "runtime"}, stats.Runtime.SinceLast, MilliSeconds}
	poller.measurements <- GaugeMeasurement{tick, poller.Name(), []string{"stats", "wall-clock"}, stats.WallClock.SinceLast, MilliSeconds}
}

func (poller FolsomPoller) doMetricsPoll(ctx slog.Context, tick time.Time) {
	metrics := make(map[string]FolsomType)
	if err := poller.decodeReq("/_metrics?info=true", &metrics); err != nil {
		LogError(ctx, err, "while performing request for this tick")
		return
	}

	for key, ft := range metrics {
		v := FolsomValue{Name: key, Type: ft.Type}
		if err := poller.decodeReq("/_metrics/"+v.Name, &v); err != nil {
			LogError(ctx, err, "while performing request for "+v.Name+"this tick")
			return
		}

		if m, err := poller.genMeasurement(tick, v); err != nil {
			LogError(ctx, err, "while performing request for "+v.Name+"this tick")
			return
		} else {
			poller.measurements <- m
		}
	}
}

func (poller FolsomPoller) genMeasurement(tick time.Time, v FolsomValue) (Measurement, error) {
	var err error

	switch v.Type {
	case "counter":
		var val int64
		if val, err = v.Value.Int64(); err == nil {
			return CounterMeasurement{tick, poller.Name(), []string{v.Name}, uint64(val), Empty}, nil
		}
	case "gauge":
		if strings.Contains(v.Value.String(), ".") {
			var val float64
			if val, err = v.Value.Float64(); err == nil {
				return FloatGaugeMeasurement{tick, poller.Name(), []string{v.Name}, val, Empty}, nil
			}
		} else {
			var val int64
			if val, err = v.Value.Int64(); err == nil {
				return GaugeMeasurement{tick, poller.Name(), []string{v.Name}, uint64(val), Empty}, nil
			}
		}
	default:
		err = errors.New("Unsupported metric type: " + v.Type)
	}

	return nil, err
}

func (poller FolsomPoller) decodeReq(path string, v interface{}) error {
	req, err := http.NewRequest("GET", poller.baseUrl+path, nil)
	if err != nil {
		return err
	}

	resp, rerr := poller.client.Do(req)
	if rerr != nil {
		return rerr
	} else if resp.StatusCode >= 300 {
		resp.Body.Close()
		return fmt.Errorf("Response returned a %d", resp.StatusCode)
	}

	defer resp.Body.Close()

	decoder := json.NewDecoder(resp.Body)
	if derr := decoder.Decode(v); derr != nil {
		return derr
	}

	return nil
}

func (poller FolsomPoller) Name() string {
	return "folsom"
}

func (poller FolsomPoller) Exit() {}