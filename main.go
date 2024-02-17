package main

import (
	"os"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
)

var (
	ip4     layers.IPv4
	eth     layers.Ethernet
	options gopacket.SerializeOptions
)

type StaticArpTable struct {
	IP               string
	Mac              string
	NetworkInterface string
}

func main() {
	subCmd()

}

func subCmd() {
	args := os.Args

	if len(args) == 1 {
		return
	}

	if args[1] == "list" {
		devList()
	}
	if args[1] == "cat" {
		catch()
	}

	os.Exit(0)
}
