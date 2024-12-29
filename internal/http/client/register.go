package client

import (
	"net/netip"
	"time"
	"strconv"
	"github.com/soypat/seqs"
	"github.com/soypat/seqs/stacks"
	"chip/internal/setup"
	"math/rand"
	"chip/internal/configuration"
)

// RegisterTrain creates a TCP connection to central server and tells him that a new train is alive
func RegisterTrain() {
	start := time.Now()

	var err error
	var conn *stacks.TCPConn
	var svAddr netip.AddrPort
	
	// Parse address from configuration
	for svAddr, err = netip.ParseAddrPort(configuration.ServerAddr) ; err != nil ; {
		setup.Logger.Error("Unable to parse server address: " + err.Error() + ". Will try again in 5 seconds")
		time.Sleep(time.Second * 5)
	}

	// Set-up port and create connection
	rng := rand.New(rand.NewSource(int64(time.Now().Sub(start))))

	clientAddr := netip.AddrPortFrom(setup.PortStack.Addr(), uint16(rng.Intn(65535-4096)+4096))

	for conn, err = stacks.NewTCPConn(setup.PortStack, stacks.TCPConnConfig{}) ; err != nil ; {
		setup.Logger.Error("Error while creating TCP connection: " + err.Error() + ". Will try again in 5 seconds")
		time.Sleep(time.Second * 5)
	}

	// Forge request by hand
	activePayload := []byte("{\"name\": \""+configuration.TrainName+"\", \"model\": \""+configuration.TrainModel+"\", \"ip\": \""+configuration.TrainIP+"\"}")
	httpRequest := []byte("POST /register HTTP/1.1\r\nHost: "+configuration.ServerAddr+"\r\nContent-Type: application/json; charset=utf-8\r\nContent-Length: "+strconv.Itoa(len(activePayload))+"\r\n\r\n")

	setup.Logger.Info("Sending information to central server...")

	// Make sure to timeout the connection if it takes too long.
	conn.SetDeadline(time.Now().Add(time.Second * 10))

	// Open connection

	if err = conn.OpenDialTCP(clientAddr.Port(), configuration.ServerMAC, svAddr, seqs.Value(uint16(rng.Intn(65535-4096)+4096))) ; err != nil {
		setup.Logger.Error("Unable to connect to central server: " + err.Error() + ". Will try again in 5 seconds")
		time.Sleep(time.Second * 5)
		conn.Close()
		RegisterTrain()
		return
	}

	count := 0
	for conn.State() != seqs.StateEstablished && count < 50 { 
		time.Sleep(100 * time.Millisecond)
		count += 1
	}
	
	if conn.State() != seqs.StateEstablished {
		setup.Logger.Error("Connection is still unestablished after 5 seconds waiting ; restarting TCP transaction now.")
		conn.SetDeadline(time.Now().Add(time.Second * 10))
		conn.Close()
		RegisterTrain()
		return
	}

	// Send the request.
	for _, err = conn.Write(append(httpRequest, activePayload...));  err != nil ;  {
		setup.Logger.Error("Unable to write request: " + err.Error() + ". Will try again in 5 seconds")
		time.Sleep(time.Second * 5)
	}

	// Read response (server should reply with "OK")
	rxBuf := make([]byte, 256)
	if n, err := conn.Read(rxBuf); n == 0 && err != nil {
		setup.Logger.Error("Unable to read response: " + err.Error() + ". Will try again in 5 seconds")
		conn.Close()
		time.Sleep(time.Second * 5)
		RegisterTrain()
		return
	} else if n == 0 {
		setup.Logger.Error("Server did not responded in time. Will try again in 5 seconds")
		conn.Close()
		time.Sleep(time.Second * 5)
		RegisterTrain()
		return
	}

	setup.Logger.Info("Successfully registered to central server")
	conn.Close()
}