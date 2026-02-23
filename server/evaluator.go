package server

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/kofuk/spaghettini/server/backend"
	"github.com/kofuk/spaghettini/server/types"
)

type Evaluator struct {
	backend       backend.Backend
	printResponse bool
}

func (t *Evaluator) Execute(r *types.Request) (*types.Response, error) {
	resp, err := t.backend.Handle(r)
	if err != nil {
		return nil, err
	}
	if t.printResponse {
		fmt.Fprintln(os.Stderr, string(resp))
	}

	parsedResp, err := http.ReadResponse(bufio.NewReader(bytes.NewBuffer(resp)), nil)
	if err != nil {
		return nil, err
	}
	defer parsedResp.Body.Close()

	body, err := io.ReadAll(parsedResp.Body)
	if err != nil {
		return nil, err
	}

	response := &types.Response{
		StatusCode: parsedResp.StatusCode,
		StatusText: parsedResp.Status,
		Header:     parsedResp.Header,
		Body:       body,
		Trailer:    parsedResp.Trailer,
	}

	return response, nil
}
