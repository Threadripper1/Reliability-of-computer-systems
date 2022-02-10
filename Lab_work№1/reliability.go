package reliability

import (
	"sort"
)

type calculator struct {
	selection []float64
}

type interval struct {
	from, to float64
}

func (i *interval) Len() float64 {
	return i.to - i.from
}

func (i *interval) isAccept(N float64) bool {
	return N >= i.from && N <= i.to
}

func NewReliabilityCalculator(selection []float64) *calculator {
	sort.Float64s(selection)
	return &calculator{
		selection: selection,
	}
}

func (c *calculator) Tcp() float64 {
	var sum float64
	for _, i := range c.selection {
		sum += i
	}
	return sum / float64(len(c.selection))
}

func (c *calculator) Max() float64 {
	return c.selection[len(c.selection)-1]
}

func (c *calculator) SplitOnIntervals(N int) []interval {
	intervalLen := c.Max() / float64(N)
	var from, to float64 = 0, intervalLen
	intervals := make([]interval, N)
	for i := 0; i < N; i++ {
		intervals[i] = interval{from, to}
		from, to = to, to+intervalLen
	}
	return intervals
}

func (c *calculator) FindStaticalDensitiesOnIntervals(intervals []interval) []float64 {
	densities := make([]float64, len(intervals))
	if len(intervals) < 1 {
		return densities
	}
	intervalLen := intervals[0].Len()
	var acceptCount int
	var from int
	for k, interval := range intervals {
		for i := from; i < len(c.selection); i++ {
			if !interval.isAccept(c.selection[i]) {
				from = i
				break
			}
			acceptCount++
		}
		densities[k] = float64(acceptCount) / (float64(len(c.selection)) * intervalLen)
		acceptCount = 0
	}
	return densities
}

// FindMTBF finds mean time between failures
func (c *calculator) FindMTBF(staticalDensities []float64, intervalSize float64) []float64 {
	mtbf := make([]float64, len(staticalDensities)+1)
	mtbf[0] = 1
	var prob float64 = 0
	for i, density := range staticalDensities {
		prob += density * intervalSize
		mtbf[i+1] = 1.0 - prob
	}
	return mtbf
}

func (c *calculator) FindStaticalMTBF(interval []interval, mtbf []float64, gamma float64) float64 {
	var p1, p2 float64
	offset := 0
	for i, p := range mtbf {
		if p < gamma {
			break
		}
		p1 = p
		offset = i
	}
	if offset != len(mtbf)-2 {
		p2 = mtbf[offset+1]
	}

	Ty := interval[offset].to - interval[offset].Len()*((p2-gamma)/(p2-p1))
	return Ty
}

func (c *calculator) FindReliableProbability(intervals []interval, densities []float64, hours float64) float64 {
	offset := findOffset(intervals, hours)
	sum := 0.0
	for i := 0; i < offset; i++ {
		sum += densities[i] * intervals[0].Len()
	}
	sum += densities[offset] * (hours - intervals[offset].from)
	return 1 - sum
}

func findOffset(intervals []interval, hours float64) int {
	offset := 0
	for offset = 0; offset < len(intervals); offset++ {
		if intervals[offset].to > hours {
			return offset
		}
	}
	return offset
}

func (c *calculator) FindFailureIntensity(intervals []interval, densities []float64, hours float64) float64 {
	offset := findOffset(intervals, hours)
	p := c.FindReliableProbability(intervals, densities, hours)
	return densities[offset] / p
}
