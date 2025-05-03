package request

type SupportedAPI struct {
	MinVersion int16
	MaxVersion int16
	NewFn      func(*RequestHeader) (Request, error)
}

var RequestKeyMap = map[int16]SupportedAPI{
	1: {
		MinVersion: 16,
		MaxVersion: 16,
		NewFn:      NewFetchV16Request,
	},
	18: {
		MinVersion: 4,
		MaxVersion: 4,
		NewFn:      NewAPIVersionsV4Request,
	},
}
