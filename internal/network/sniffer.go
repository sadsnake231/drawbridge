package network

import (
	"log"
	"log/slog"
	"time"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
)

func StartSniffing(device string, sm *StateManager) {
	handle, err := pcap.OpenLive(device, 1024, false, 10*time.Second)
	if err != nil {
		log.Fatal(err)
	}

	defer handle.Close()

	err = handle.SetBPFFilter("tcp[tcpflags] & (tcp-syn) != 0")
	if err != nil {
		log.Fatal(err)
	}

	packetSource := gopacket.NewPacketSource(handle, handle.LinkType())

	for packet := range packetSource.Packets() {
		processPacket(packet, sm)
	}
}

func processPacket(packet gopacket.Packet, sm *StateManager) {
	ipLayer := packet.Layer(layers.LayerTypeIPv4)
	if ipLayer == nil {
		return
	}

	ip, _ := ipLayer.(*layers.IPv4)

	tcpLayer := packet.Layer(layers.LayerTypeTCP)
	if tcpLayer == nil {
		return
	}

	tcp, _ := tcpLayer.(*layers.TCP)

	slog.Debug("knock received", "ip", ip.SrcIP.String(), "port", tcp.DstPort)

	sm.HandlePacket(ip.SrcIP.String(), int(tcp.DstPort))
}
