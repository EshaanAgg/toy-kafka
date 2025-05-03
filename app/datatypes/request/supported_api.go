package request

type SupportedAPI struct {
	MinVersion int16
	MaxVersion int16
	NewFn      func(*RequestHeader) (Request, error)
}

var RequestKeyMap = map[int16]SupportedAPI{
	18: {
		MinVersion: 4,
		MaxVersion: 4,
		NewFn:      NewAPIVersionV4Request,
	},
}
