package main

import (
	"flag"
	"fmt"
	"lab1/config"
	"lab1/reliability"
)

func arrayToStr(varName string, array []float64) string {
	var out string
	for i, elem := range array {
		out += fmt.Sprintf("\n\t%s%d = %f", varName, i, elem)
	}
	return out
}

func main() {
	c := flag.String("config", "config.yaml", "path to config")
	cfg, err := config.NewAppConfig(*c)
	if err != nil {
		panic("unable to parse config")
	}
	fmt.Println("incoming config:", cfg.String())

	calculator := reliability.NewReliabilityCalculator(cfg.Selection)
	fmt.Printf("max time: %f\n", calculator.Max())
	fmt.Printf("Tcp: %f\n", calculator.Tcp())

	intervals := calculator.SplitOnIntervals(cfg.IntervalSize)
	fmt.Printf("split on %d intervals: %f\n", cfg.IntervalSize, intervals)
	if len(intervals) < 1 {
		fmt.Println("no intervals")
		return
	}

	densities := calculator.FindStaticalDensitiesOnIntervals(intervals)
	fmt.Println("find statical density for above intervals: ", arrayToStr("f", densities))

	P := calculator.FindMTBF(densities, intervals[0].Len())
	fmt.Println("\nmean time between failures: ", arrayToStr("P", P))

	Ty := calculator.FindStaticalMTBF(intervals, P, cfg.Gamma)
	fmt.Println("\nTy: ", Ty)

	p1 := calculator.FindReliableProbability(intervals, densities, cfg.Hours[0])
	fmt.Printf("P(%.3fh): %f\n", cfg.Hours[0], p1)

	intensity := calculator.FindFailureIntensity(intervals, densities, cfg.Hours[1])
	fmt.Printf("λ(%.3fh): %f\n", cfg.Hours[1], intensity)
}
