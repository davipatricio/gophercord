package payloads

type BasePayload struct {
	Op uint8       `json:"op"`
	D  interface{} `json:"d"`
	S  uint32      `json:"s"`
	T  string      `json:"t"`
}
