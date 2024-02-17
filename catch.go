package main

import (
	"fmt"
	"log"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
)

func catch() {
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
