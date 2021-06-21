package vo

type Response struct {
	Request interface{} `json:"request"`
	Result  interface{} `json:"result"`
}
