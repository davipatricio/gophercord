package payloads

type HelloPayloadData struct {
	BasePayload
	HeartbeatInterval int32 `json:"heartbeat_interval"`
}

func CreateHelloPayload(data BasePayload) HelloPayloadData {
	return HelloPayloadData{
		BasePayload: BasePayload{
			Op: data.Op,
			D:  data.D,
			S:  data.S,
			T:  data.T,
		},
		HeartbeatInterval: int32(data.D.(map[string]interface{})["heartbeat_interval"].(float64)),
	}
}
