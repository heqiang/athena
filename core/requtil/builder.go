package requtil

type urlBuilder struct {
	baseUrl string
	scheme  string
	host    string
	paths   map[string]string
	params  map[string]string
}
