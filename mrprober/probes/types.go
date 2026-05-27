package probes

import (
	"fmt"
	"github.com/VictoriaMetrics/metrics"
	"mrprober/conf"
	"slices"
	"sort"
	"strings"
)

// Enum of supported probes
const (
	NetProbeT  string = "net"
	FakeProbeT        = "fake"
)

// ErrorCode Enum of known error codes
type ErrorCode int

const (
	Success ErrorCode = iota
	Failure
)

// ----------------------------------------------------------------------------

type Result struct {
	ProbeID string
	// MetricName is used as a domain identifier for the metric (https://prometheus.io/docs/practices/naming/)
	MetricName string
	// MetricLabels is a list of labels that are part of the metric name
	MetricLabels map[string]string
	ReturnCode   ErrorCode
	// Standard output
	Msg string
}

// Update register if needed to the prometheus registry and update the Gauge value.
func (pr Result) Update() {

	// Sort labels by key otherwise we may duplicate entries in the registry
	labels := make([]string, 0, len(pr.MetricLabels))
	for k := range pr.MetricLabels {
		labels = append(labels, k)
	}
	sort.Strings(labels)

	// flatten labels as string
	var labelsFlatten strings.Builder
	for _, k := range labels {
		_, _ = labelsFlatten.WriteString(fmt.Sprintf(`%s="%s",`, k, pr.MetricLabels[k]))
	}

	metricLabels := strings.TrimSuffix(labelsFlatten.String(), ",")
	finalMetricName := fmt.Sprintf(`mrprober_%s{%s}`, pr.MetricName, metricLabels)

	if slices.Contains(metrics.ListMetricNames(), finalMetricName) {
		// for gauges, we have to unregister the metric before updating the value
		metrics.UnregisterMetric(finalMetricName)
	}

	// Register as new gauge
	// A gauge is a metric that represents a single numerical value that can arbitrarily go up and down.
	metrics.NewGauge(finalMetricName, func() float64 {
		return float64(pr.ReturnCode)
	})
}

func (pr Result) String() string {

	var msg string
	switch pr.ReturnCode {
	case Success:
		msg = "pass"
	default:
		msg = "fail"
	}

	return fmt.Sprintf(`%s (%d): %s`, msg, pr.ReturnCode, pr.Msg)
}

type Probe interface {
	Exec() Result
}

// ----------------------------------------------------------------------------

// New is a factory method that creates a probe from a rule
func New(r *conf.Rule) (Probe, error) {

	var p Probe
	var err error

	switch t := r.Probe; t {
	case FakeProbeT:
		p, err = NewFake(r.Name, r.Args)
	case NetProbeT:
		p, err = NewNet(r.Name, r.Args)
	default:
		err = fmt.Errorf("probe `%s` has no implementation and therefore is not supported", t)
	}

	if err != nil {
		return nil, err
	}
	return p, nil
}
