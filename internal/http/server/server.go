package server

import (
	"bufio"
	"io"
	"time"

	"github.com/soypat/seqs/httpx"
	"github.com/soypat/seqs/stacks"
	"chip/internal/setup"
)

func StartServer() {
	listener, err := stacks.NewTCPListener(setup.PortStack, stacks.TCPListenerConfig{
		MaxConnections: 2,
		ConnTxBufSize:  2030,
		ConnRxBufSize:  2030,
	})
	if err != nil {
		setup.Logger.Error("Error while creating listener: " + err.Error())
		panic("listener create:" + err.Error())
	}

	if err = listener.StartListening(80); err != nil {
		setup.Logger.Error("Error while starting to listen: " + err.Error())
		panic("listener start:" + err.Error())
	}

	// Reuse the same buffers for each connection to avoid heap allocations.
	var req httpx.RequestHeader
	var resp httpx.ResponseHeader
	buf := bufio.NewReaderSize(nil, 1024)

	for {
		conn, err := listener.Accept()
		if err != nil {
			continue
		}

		if err = conn.SetDeadline(time.Now().Add(3 * time.Second));  err != nil {
			continue
		}

		buf.Reset(conn)
		if err = req.Read(buf); err != nil {
			continue
		}

		resp.Reset()
		handler(conn, &resp, &req)
		conn.Close()
	}
}


func handler(respWriter io.Writer, resp *httpx.ResponseHeader, req *httpx.RequestHeader) {
	uri := string(req.RequestURI())
	resp.SetConnectionClose()
	switch uri {
	case "/state":
		getState(respWriter, resp, req)
	case "/ping":
		ping(respWriter, resp, req)
	case "/lights":
		setLights(respWriter, resp, req)
	case "/speed": 
		setSpeed(respWriter, resp, req)
	default:
		resp.SetStatusCode(404)
		respWriter.Write(resp.Header())
	}
}

func boolToString(value bool) string {
	if value {
		return "true"
	} else {
		return "false"
	}
}