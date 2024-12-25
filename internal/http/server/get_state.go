package server

import (
	"io"
    "strconv"

	"github.com/soypat/seqs/httpx"
	"chip/internal/setup"
)

func getState(respWriter io.Writer, resp *httpx.ResponseHeader, req *httpx.RequestHeader) {
	data := "{\"speed\": " +strconv.FormatUint(uint64(setup.CurrentSpeed), 10) + ", \"isGoingForward\": " + boolToString(setup.CurrentDirection) + ", \"lightsActivated\": " + boolToString(setup.CurrentLedState) + "}"
	resp.SetContentLength(len([]byte(data)))
	resp.SetContentType("application/json")
	respWriter.Write(resp.Header())
	respWriter.Write([]byte(data))
}