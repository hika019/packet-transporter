package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"strings"

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

type Device struct {
	macAddr  net.HardwareAddr
	ipv4Addr net.IP
}

func main() {
	myDevice := getInterface("ens18")

	var dstIp = "192.168.20.101"

	eth := &layers.Ethernet{
		SrcMAC:       myDevice.macAddr,
		DstMAC:       parseMac("FF:FF:FF:FF:FF:FF"),
		EthernetType: layers.EthernetTypeARP,
	}
	arpreq := &layers.ARP{
		AddrType:          layers.LinkTypeEthernet,
		Protocol:          layers.EthernetTypeIPv4,
		HwAddressSize:     6,
		ProtAddressSize:   4,
		Operation:         layers.ARPRequest,
		SourceHwAddress:   myDevice.macAddr,
		SourceProtAddress: myDevice.ipv4Addr.To4(),
		DstHwAddress:      []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00},
		DstProtAddress:    net.ParseIP(dstIp).To4(),
	}
	buff := gopacket.NewSerializeBuffer()

	err := gopacket.SerializeLayers(buff, gopacket.SerializeOptions{FixLengths: true}, eth, arpreq)
	if err != nil {
		log.Fatal(err)
	}

	handle, err := pcap.OpenLive("ens18", 65536, true, pcap.BlockForever)
	if err != nil {
		log.Fatal(err)
	}

	defer handle.Close()

	err = handle.WritePacketData(buff.Bytes())
	if err != nil {
		log.Fatal(err)
	}

	packetSource := gopacket.NewPacketSource(handle, handle.LinkType())
	for packet := range packetSource.Packets() {
		ethernetLayer := packet.Layer(layers.LayerTypeEthernet)
		ethernetPacket := ethernetLayer.(*layers.Ethernet)
		if ethernetPacket.EthernetType.LayerType() == layers.LayerTypeARP {
			fmt.Printf("packet is %+v\n", packet)
			break
		}
	}

}

func parseMac(macAddr string) net.HardwareAddr {
	parsedMac, _ := net.ParseMAC(macAddr)
	return parsedMac
}

func getInterface(ifname string) Device {
	netifs, err := net.Interfaces()
	if err != nil {
		log.Fatal(err)
	}

	for _, netif := range netifs {
		if netif.Name == ifname {
			addrs, _ := netif.Addrs()
			for _, addr := range addrs {
				if !strings.Contains(addr.String(), ":") && strings.Contains(addr.String(), ".") {
					ip, _, _ := net.ParseCIDR(addr.String())
					return Device{
						macAddr:  netif.HardwareAddr,
						ipv4Addr: ip,
					}
				}
			}
		}
	}
	return Device{}
}

func subCmd() {
	args := os.Args

	dev := getInterface("ens18")
	fmt.Println(dev.ipv4Addr.String(), dev.macAddr.String())

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
