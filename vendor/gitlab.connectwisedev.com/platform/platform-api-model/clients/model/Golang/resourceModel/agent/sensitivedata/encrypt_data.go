package sensitivedata
const (
	//GenerateKeyPairAction is the action to be performed upon GenerateKeyPair mailbox message
	GenerateKeyPairAction string = "generateKeyPair"
)

//Encryption is the struct defining the data to be encrypted for specific endpointID
type Encryption struct {
	EndpointID string `json:"endpointId,omitempty"`
	Data       string `json:"data,omitempty"`
}
