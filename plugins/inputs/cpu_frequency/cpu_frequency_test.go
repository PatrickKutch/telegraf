package cpu_frequency

import (
	"math"
	"testing"

	"github.com/influxdata/telegraf/testutil"
)

func TestCPU_Frequency(t *testing.T) {
	s := &CPU_Frequency{
		Amplitude: 10.0,
	}

	for i := 0.0; i < 10.0; i++ {

		var acc testutil.Accumulator

		s.Gather(&acc)

		fields := make(map[string]interface{})

		acc.AssertContainsFields(t, "cpu0", fields)
	}
}
