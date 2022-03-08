package config

import "time"

// Allowed protocols.
const (
	ProtocolRCON    = "rcon"
	ProtocolTELNET  = "telnet"
	ProtocolWebRCON = "web"
)

// DefaultProtocol contains the default protocol for connecting to a
// remote server.
const DefaultProtocol = ProtocolRCON

// DefaultTimeout contains the default dial and execute timeout.
const DefaultTimeout = 10 * time.Second

// Session contains details for making a request on a remote server.
type Session struct {
	Address  string `json:"address" yaml:"address"`
	Password string `json:"password" yaml:"password"`
	// Log is the name of the file to which requests will be logged.
	// If not specified, no logging will be performed.
	Log        string        `json:"log" yaml:"log"`
	Type       string        `json:"type" yaml:"type"`
	SkipErrors bool          `json:"skip_errors" yaml:"skip_errors"`
	Timeout    time.Duration `json:"timeout" yaml:"timeout"`
}
