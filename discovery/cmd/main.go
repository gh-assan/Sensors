package main

import (
	"gh-assan/rmsq/discovery"
)

func main() {
	sensorDiscovery := discovery.NewSensorDiscovery()

	sensorDiscovery.Start()

	select {}
}
