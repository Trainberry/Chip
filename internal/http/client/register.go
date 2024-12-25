package client

import (
	"net/netip"
	"time"
	"strconv"
	"github.com/soypat/seqs"
	"github.com/soypat/seqs/stacks"
	"chip/internal/setup"
	"chip/internal/configuration"
)

// RegisterTrain creates a TCP connection to central server and tells him that a new train is alive
func RegisterTrain() {
	// Parse address from configuration
	svAddr, err := netip.ParseAddrPort(configuration.ServerAddr)
	if err != nil {
		setup.Logger.Error("Unable to parse server address: " + err.Error())
		panic("parsing server address:" + err.Error())
	}

	// Set-up port and create connection
	clientAddr := netip.AddrPortFrom(setup.PortStack.Addr(), uint16(64632))
	conn, err := stacks.NewTCPConn(setup.PortStack, stacks.TCPConnConfig{})

	if err != nil {
		setup.Logger.Error("Error while creating TCP connection: " + err.Error())
		panic("conn create:" + err.Error())
	}

	// Forge request by hand
	activePayload := []byte("{\"name\": \""+configuration.TrainName+"\", \"model\": \""+configuration.TrainModel+"\", \"ip\": \""+configuration.TrainIP+"\"}")
	httpRequest := []byte("POST /register HTTP/1.1\r\nHost: "+configuration.ServerAddr+"\r\nContent-Type: application/json; charset=utf-8\r\nContent-Length: "+strconv.Itoa(len(activePayload))+"\r\n\r\n")

	setup.Logger.Info("Sending information to central server...")

	// Make sure to timeout the connection if it takes too long.
	conn.SetDeadline(time.Now().Add(time.Second * 10))

	// Open connection
	if err = conn.OpenDialTCP(clientAddr.Port(), configuration.ServerMAC, svAddr, seqs.Value(64632));  err != nil {
		setup.Logger.Error("Unable to connect to central server: " + err.Error())
		panic("opening TCP: " + err.Error())
	}
	for conn.State() != seqs.StateEstablished {
		time.Sleep(100 * time.Millisecond)
	}

	// Send the request.
	if _, err = conn.Write(append(httpRequest, activePayload...));  err != nil {
		setup.Logger.Error("Unable to write request: " + err.Error())
		panic("writing request: " + err.Error())
	}

	// Read response (server should reply with "OK")
	rxBuf := make([]byte, 256)
	if n, err := conn.Read(rxBuf); n == 0 && err != nil {
		setup.Logger.Error("Unable to read response: " + err.Error())
		panic("reading response: " + err.Error())
	} else if n == 0 {
		setup.Logger.Error("Server did not responded in time")
		panic("no response")
	}

	setup.Logger.Info("Successfully registered to central server")
	conn.Close()
}