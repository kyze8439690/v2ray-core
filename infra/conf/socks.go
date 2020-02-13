package conf

import (
	"github.com/golang/protobuf/proto"
	"v2ray.com/core/common/protocol"
	"v2ray.com/core/proxy/socks"
)

const (
	AuthMethodNoAuth = "noauth"
)

type SocksServerConfig struct {
	AuthMethod string   `json:"auth"`
	UDP        bool     `json:"udp"`
	Host       *Address `json:"ip"`
}

func (v *SocksServerConfig) Build() (proto.Message, error) {
	config := new(socks.ServerConfig)
	switch v.AuthMethod {
	case AuthMethodNoAuth:
		config.AuthType = socks.AuthType_NO_AUTH
	default:
		//newError("unknown socks auth method: ", v.AuthMethod, ". Default to noauth.").AtWarning().WriteToLog()
		config.AuthType = socks.AuthType_NO_AUTH
	}

	config.UdpEnabled = v.UDP
	if v.Host != nil {
		config.Address = v.Host.Build()
	}

	return config, nil
}

type SocksRemoteConfig struct {
	Address *Address `json:"address"`
	Port    uint16   `json:"port"`
}
type SocksClientConfig struct {
	Servers []*SocksRemoteConfig `json:"servers"`
}

func (v *SocksClientConfig) Build() (proto.Message, error) {
	config := new(socks.ClientConfig)
	config.Server = make([]*protocol.ServerEndpoint, len(v.Servers))
	for idx, serverConfig := range v.Servers {
		server := &protocol.ServerEndpoint{
			Address: serverConfig.Address.Build(),
			Port:    uint32(serverConfig.Port),
		}
		config.Server[idx] = server
	}
	return config, nil
}
