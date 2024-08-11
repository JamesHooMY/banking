package v1

type ErrResponse struct {
	Data any    `json:"data"`
	Msg  string `json:"msg"`
}
