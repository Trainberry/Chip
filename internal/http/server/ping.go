package server

import (
	"io"

	"github.com/soypat/seqs/httpx"
)

func ping(respWriter io.Writer, resp *httpx.ResponseHeader, req *httpx.RequestHeader) {
	resp.SetStatusCode(204)
	respWriter.Write(resp.Header())
}