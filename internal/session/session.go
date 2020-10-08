package session

// Allowed protocols.
const (
	ProtocolRCON   = "rcon"
	ProtocolTELNET = "telnet"
)

// DefaultProtocol contains the default protocol for connecting to a
// remote server.
const DefaultProtocol = ProtocolRCON

// Session contains details for making a request on a remote server.
type Session struct {
	Address  string `json:"address" yaml:"address"`
	Password string `json:"password" yaml:"password"`
	Log      string `json:"log" yaml:"log"`
	Type     string `json:"type" yaml:"type"`
}
