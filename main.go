package main

import (
	"fmt"
	"log"
	"os"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
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

	fmt.Println("----------Listen Packet----------")

	handle, err := pcap.OpenLive("ens18", 65536, true, pcap.BlockForever)
	if err != nil {
		log.Fatal("ERROR opening pcap:", err)
	}

	// if err := handle.SetBPFFilter("tcp and port 54321"); err != nil {
	// 	log.Fatal(err)
	// }

	defer handle.Close()

	src := gopacket.NewPacketSource(handle, handle.LinkType())

	for packet := range src.Packets() {
		fmt.Println(packet)
		packet.Layer(layers.LayerTypeIPv4)
	}

}

func subCmd() {
	args := os.Args

	if len(args) == 1 {
		return
	}

	if args[1] == "list" {
		devList()
	}
	os.Exit(0)
}
