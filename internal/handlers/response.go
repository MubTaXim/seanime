package handlers

type SeaResponse[R any] struct {
	Error string `json:"error,omitempty"`
	Data  R      `json:"data,omitempty"`
}

func NewDataResponse[R any](data R) SeaResponse[R] {
	res := SeaResponse[R]{
		Data: data,
	}
	return res
}

func NewErrorResponse(err error) SeaResponse[any] {
	if err == nil {
		return SeaResponse[any]{
			Error: "Unknown error",
		}
	}
	res := SeaResponse[any]{
		Error: err.Error(),
	}
	return res
}
