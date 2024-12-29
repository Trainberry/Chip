package server

import (
	"io"

	"github.com/soypat/seqs/httpx"
	"chip/internal/controls"
)

func setLights(respWriter io.Writer, resp *httpx.ResponseHeader, req *httpx.RequestHeader) {
	if string(req.Method()) == "POST" {
		controls.SetLedState(true)
	} else {
		controls.SetLedState(false)
	}
	resp.SetStatusCode(204)
	respWriter.Write(resp.Header())
}