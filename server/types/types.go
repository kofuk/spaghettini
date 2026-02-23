package types

type Response struct {
	StatusCode int
	StatusText string
	Header     map[string][]string
	Body       []byte
	Trailer    map[string][]string
}

type Request struct {
	Method  string
	Path    string
	Header  map[string][]string
	Body    []byte
	Trailer map[string][]string
}
