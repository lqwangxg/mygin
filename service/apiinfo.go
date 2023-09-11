package service

type ApiInfo struct {
	Method string
	Params []string
	Result string
}

func NewApiDescInfo(method string, Params []string, result string) *ApiInfo {
	return &ApiInfo{method, Params, result}
}
