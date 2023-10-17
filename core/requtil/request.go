package requtil

type (
	Request struct {
		Url     string            `json:"url"`
		Body    string            `json:"body,omitempty"`
		Method  string            `json:"method,omitempty"`
		Headers map[string]string `json:"headers,omitempty"`
		Cookies map[string]string `json:"cookies,omitempty"`
		Params  map[string]string `json:"params,omitempty"`
	}

	Response struct {
		StatusCode  int
		respHeader  map[string]string
		respContent string
	}
)

func NewRequest(url string, method string, body string) *Response {
	return &Response{}
}
