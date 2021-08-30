package main

import (
	"gh-assan/rmsq/receiver"
	"gh-assan/rmsq/sender"
)

func main() {

	go receiver.Receive()
	go sender.Send()

	select {} // block forever

}
