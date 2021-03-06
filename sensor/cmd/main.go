package main

import (
	"flag"

	"gh-assan/rmsq/sensor"
)

func main() {

	var name = flag.String("name", "sensor", "name of the sensor")
	var freq = flag.Int("freq", 1, "update frequency in cycles/sec")
	var max = flag.Float64("max", 5., "maximum value for generated readings")
	var min = flag.Float64("min", 1., "minimum value for generated readings")
	var stepSize = flag.Float64("step", 0.1, "maximum allowable change per measurement")

	flag.Parse()

	sensor := sensor.NewSensor(*name, *freq, *max, *min, *stepSize)

	go sensor.SendReading()

	go sensor.HandleDiscovery()

	select {}
}
