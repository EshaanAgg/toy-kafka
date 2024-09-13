package protocol

const NO_ERROR_BODY = 0
const UNSUPPORTED_VERSION_ERROR_CODE = 35

const MOCK_THROTTLE_TIME = 0

type APIVersionListItem struct {
	APIVersion int16
	MinVersion int16
	MaxVersion int16
}

func (v *APIVersionListItem) Encode(r *Response) {
	r.WriteInt16(v.APIVersion)
	r.WriteInt16(v.MinVersion)
	r.WriteInt16(v.MaxVersion)
	r.WriteBytes(0x00) // Tagged Fields
}

type APIVersionBody struct {
	ErrorCode int16
	APIKeys   []APIVersionListItem
	Throttle  int32
}

func (b *APIVersionBody) Encode(r *Response) {
	r.WriteInt16(b.ErrorCode)

	r.WriteVarInt(len(b.APIKeys) + 1)
	for _, apiVersion := range b.APIKeys {
		apiVersion.Encode(r)
	}

	r.WriteInt32(b.Throttle)
	r.WriteBytes(0x00) // Tagged Fields
}

var API_VERSION_LIST = []APIVersionListItem{
	{
		// API Versions API
		APIVersion: 18,
		MinVersion: 0,
		MaxVersion: 4,
	},
	{
		// Fetch API
		APIVersion: 1,
		MinVersion: 0,
		MaxVersion: 16,
	},
}

func GetAPIVersionBody(r *Request) *APIVersionBody {
	var errorCode int16 = NO_ERROR_BODY
	if !r.Header.ValidateAPIVersion() {
		errorCode = UNSUPPORTED_VERSION_ERROR_CODE
	}

	return &APIVersionBody{
		ErrorCode: errorCode,
		APIKeys:   API_VERSION_LIST,
		Throttle:  MOCK_THROTTLE_TIME,
	}
}
