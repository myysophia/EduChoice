package response

type FrequentResponse struct {
	Code      string `json:"code"`
	Message   string `json:"message"`
	Data      string `json:"data"`
	Location  string `json:"location"`
	Encrydata string `json:"encrydata"`
}
