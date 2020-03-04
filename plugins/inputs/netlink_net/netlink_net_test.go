package netlink_net

import (
	"math"
	"testing"

	"github.com/influxdata/telegraf/testutil"
)

func TestPkTest(t *testing.T) {
	s := &Netlink{
		Amplitude: 10.0,
	}

	for i := 0.0; i < 10.0; i++ {

		var acc testutil.Accumulator

		//sine := math.Sin((i*math.Pi)/5.0) * s.Amplitude
		//cosine := math.Cos((i*math.Pi)/5.0) * s.Amplitude

		s.Gather(&acc)

		fields := make(map[string]interface{})
		fields["sine"] = 5
		fields["cosine"] = 33

		//acc.AssertContainsFields(t, "netlink", fields)
	}
}
