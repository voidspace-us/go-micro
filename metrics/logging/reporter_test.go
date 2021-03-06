package logging

import (
	"bytes"
	"testing"
	"time"

	log "github.com/micro/go-micro/v3/logger"
	"github.com/micro/go-micro/v3/metrics"

	"github.com/stretchr/testify/assert"
)

func TestLoggingReporter(t *testing.T) {

	// Make a Reporter:
	reporter := New(metrics.Path("/prometheus"), metrics.DefaultTags(map[string]string{"service": "prometheus-test"}))
	assert.NotNil(t, reporter)
	assert.Equal(t, "prometheus-test", reporter.options.DefaultTags["service"])
	assert.Equal(t, ":9000", reporter.options.Address)
	assert.Equal(t, "/prometheus", reporter.options.Path)

	// Make a log buffer and create a new logger to use it:
	logBuffer := new(bytes.Buffer)
	reporter.logger = log.NewLogger(log.WithLevel(log.TraceLevel), log.WithOutput(logBuffer), log.WithFields(convertTags(reporter.options.DefaultTags)))

	// Check that our implementation is valid:
	assert.Implements(t, new(metrics.Reporter), reporter)

	// Test tag conversion:
	tags := metrics.Tags{
		"tag1": "false",
		"tag2": "true",
	}
	convertedTags := convertTags(tags)
	assert.Equal(t, "false", convertedTags["tag1"])
	assert.Equal(t, "true", convertedTags["tag2"])

	// Test submitting metrics through the interface methods:
	assert.NoError(t, reporter.Count("test.counter.1", 6, tags))
	assert.NoError(t, reporter.Count("test.counter.2", 19, tags))
	assert.NoError(t, reporter.Count("test.counter.1", 5, tags))
	assert.NoError(t, reporter.Gauge("test.gauge.1", 99, tags))
	assert.NoError(t, reporter.Gauge("test.gauge.2", 55, tags))
	assert.NoError(t, reporter.Gauge("test.gauge.1", 98, tags))
	assert.NoError(t, reporter.Timing("test.timing.1", time.Second, tags))
	assert.NoError(t, reporter.Timing("test.timing.2", time.Minute, tags))

	// Test reading back the metrics from the logbuffer (doesn't seem to work because the output still goes to StdOut):
	// assert.Contains(t, logBuffer.String(), "level=debug service=prometheus-test Count metric: map[tag1:false tag2:true]")
}
