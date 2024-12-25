package server

import (
	"io"
    "strconv"

	"github.com/soypat/seqs/httpx"
	"chip/internal/controls"
)

func setSpeed(respWriter io.Writer, resp *httpx.ResponseHeader, req *httpx.RequestHeader) {
	speed := req.Peek("Speed")
	direction := req.Peek("Direction")

	val, err := strconv.Atoi(string(speed))
	if err != nil {
		return
	}

	controls.SetSpeed(uint32(val))

	if string(direction) == "1" {
		controls.SetDirection(true)
	} else {
		controls.SetDirection(false)
	}
	resp.SetStatusCode(204)
}