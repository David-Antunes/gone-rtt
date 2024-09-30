package main

import (
	"bytes"
	"encoding/json"
	"log"
	"net"
	"os"
	"strings"
	"time"

	api "github.com/David-Antunes/gone-proxy/api"
)

var rttLog = log.New(os.Stdout, "RTT INFO: ", log.Ltime)

func main() {
	//Configure IP and Broadcast Addr
	listenAddr, err := net.ResolveUDPAddr("udp4", ":8000")

	if err != nil {
		panic(err)
	}

	rttLog.Println(listenAddr)

	port, err := net.ListenUDP("udp4", listenAddr)

	if err != nil {
		panic(err)
	}

	ief, err := net.InterfaceByName("eth0")

	if err != nil {
		panic(err)
	}

	addrs, err := ief.Addrs()

	if err != nil {
		panic(err)
	}

	ip := addrs[0].(*net.IPNet).IP.To4()

	rttLog.Println("IP address:", ip)

	splitAddr := strings.Split(addrs[0].String(), ".")

	if len(splitAddr) != 4 {
		panic("something went wrong with Ip address")
	}

	broadcastIp := splitAddr[0] + "." + splitAddr[1] + "." + splitAddr[2] + ".255:8000"

	rttLog.Println("Broadcast Address:", broadcastIp)

	conn, err := net.Dial("udp4", broadcastIp)

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
