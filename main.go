package main

import (
	"bytes"
	"log"
	"net"
	"os"
	"time"

	"encoding/json"

	"github.com/David-Antunes/gone-proxy/api"
	"github.com/seancfoley/ipaddress-go/ipaddr"
)

var rttLog = log.New(os.Stdout, "RTT INFO: ", log.Ltime)
var DEFAULT_PORT = ":8000"
var DEFAULT_IFACE = "eth0"

func main() {
	var err error
	var ief *net.Interface
	if ief, err = net.InterfaceByName(DEFAULT_IFACE); err != nil {
		panic(err)
	}

	var addrs []net.Addr
	if addrs, err = ief.Addrs(); err != nil {
		panic(err)
	}

	var ip net.IP
	var ipNet *net.IPNet

	if ip, ipNet, err = net.ParseCIDR(addrs[0].String()); err != nil {
		panic(err)
	}

	var bcast string

	if adr, err := ipaddr.NewIPAddressFromNetIPNet(ipNet); err != nil {
		panic(err)
	} else {
		b, _ := adr.ToIPv4().ToBroadcastAddress()
		bcast = b.GetNetIPAddr().IP.String() + DEFAULT_PORT
	}

	rttLog.Println("IP address:", ip)

	rttLog.Println("Broadcast Address:", bcast)

	listenAddr, err := net.ResolveUDPAddr("udp4", DEFAULT_PORT)

	if err != nil {
		panic(err)
	}

	rttLog.Println(listenAddr)

	port, err := net.ListenUDP("udp4", listenAddr)

	if err != nil {
		panic(err)
	}

	conn, err := net.Dial("udp4", bcast)

	if err != nil {
		panic(err)
	}

	var size int
	var addr net.Addr
	var ipSender net.IP

	for {
		buf := make([]byte, 2048)
		size, addr, err = port.ReadFrom(buf)

		if err != nil {
			panic(err)
		}

		buf = buf[:size]
		ipSender = addr.(*net.UDPAddr).IP.To4()
		if ip.Equal(ipSender) {
			continue
		}
		resp := &api.RTTRequest{}
		d := json.NewDecoder(bytes.NewReader(buf))
		err = d.Decode(resp)
		if err != nil {
			panic(err)
		}

		resp.ReceiveTime = time.Now()
		resp.TransmitTime = time.Now()

		req, err := json.Marshal(resp)

		if err != nil {
			panic(err)
		}

		conn.Write(req)
		rttLog.Println("StartTime:", resp.StartTime, "ReceiveTime:", resp.ReceiveTime, "Difference:", resp.ReceiveTime.Sub(resp.StartTime))
	}
}
