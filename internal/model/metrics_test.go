package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMetricTypes(t *testing.T) {
	assert.Equal(t, "gauge", Gauge)
	assert.Equal(t, "counter", Counter)
}
