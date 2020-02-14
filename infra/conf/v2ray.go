package conf

import (
	"encoding/json"

	"v2ray.com/core"
	"v2ray.com/core/app/dispatcher"
	"v2ray.com/core/app/proxyman"
	"v2ray.com/core/common/serial"
)

var (
	inboundConfigLoader = NewJSONConfigLoader(ConfigCreatorCache{
		"socks": func() interface{} { return new(SocksServerConfig) },
	}, "protocol", "settings")

	outboundConfigLoader = NewJSONConfigLoader(ConfigCreatorCache{
		"vmess": func() interface{} { return new(VMessOutboundConfig) },
	}, "protocol", "settings")
)

type InboundDetourConfig struct {
	Protocol  string           `json:"protocol"`
	PortRange *PortRange       `json:"port"`
	ListenOn  *Address         `json:"listen"`
	Settings  *json.RawMessage `json:"settings"`
	Tag       string           `json:"tag"`
}

// Build implements Buildable.
func (c *InboundDetourConfig) Build() (*core.InboundHandlerConfig, error) {
	receiverSettings := &proxyman.ReceiverConfig{}

	if c.PortRange == nil {
		return nil, newError("port range not specified in InboundDetour.")
	}
	receiverSettings.PortRange = c.PortRange.Build()

	if c.ListenOn != nil {
		if c.ListenOn.Family().IsDomain() {
			return nil, newError("unable to listen on domain address: ", c.ListenOn.Domain())
		}
		receiverSettings.Listen = c.ListenOn.Build()
	}

	settings := []byte("{}")
	if c.Settings != nil {
		settings = ([]byte)(*c.Settings)
	}
	rawConfig, err := inboundConfigLoader.LoadWithID(settings, c.Protocol)
	if err != nil {
		return nil, newError("failed to load inbound detour config.").Base(err)
	}
	ts, err := rawConfig.(Buildable).Build()
	if err != nil {
		return nil, err
	}

	return &core.InboundHandlerConfig{
		Tag:              c.Tag,
		ReceiverSettings: serial.ToTypedMessage(receiverSettings),
		ProxySettings:    serial.ToTypedMessage(ts),
	}, nil
}

type OutboundDetourConfig struct {
	Protocol      string           `json:"protocol"`
	Tag           string           `json:"tag"`
	Settings      *json.RawMessage `json:"settings"`
	StreamSetting *StreamConfig    `json:"streamSettings"`
}

// Build implements Buildable.
func (c *OutboundDetourConfig) Build() (*core.OutboundHandlerConfig, error) {
	senderSettings := &proxyman.SenderConfig{}

	if c.StreamSetting != nil {
		ss, err := c.StreamSetting.Build()
		if err != nil {
			return nil, err
		}
		senderSettings.StreamSettings = ss
	}

	settings := []byte("{}")
	if c.Settings != nil {
		settings = ([]byte)(*c.Settings)
	}
	rawConfig, err := outboundConfigLoader.LoadWithID(settings, c.Protocol)
	if err != nil {
		return nil, newError("failed to parse to outbound detour config.").Base(err)
	}
	ts, err := rawConfig.(Buildable).Build()
	if err != nil {
		return nil, err
	}

	return &core.OutboundHandlerConfig{
		SenderSettings: serial.ToTypedMessage(senderSettings),
		Tag:            c.Tag,
		ProxySettings:  serial.ToTypedMessage(ts),
	}, nil
}

type Config struct {
	LogConfig       *LogConfig             `json:"log"`
	RouterConfig    *RouterConfig          `json:"routing"`
	InboundConfigs  []InboundDetourConfig  `json:"inbounds"`
	OutboundConfigs []OutboundDetourConfig `json:"outbounds"`
	Policy          *PolicyConfig          `json:"policy"`
}

// Build implements Buildable.
func (c *Config) Build() (*core.Config, error) {
	config := &core.Config{
		App: []*serial.TypedMessage{
			serial.ToTypedMessage(&dispatcher.Config{}),
			serial.ToTypedMessage(&proxyman.InboundConfig{}),
			serial.ToTypedMessage(&proxyman.OutboundConfig{}),
		},
	}

	var logConfMsg *serial.TypedMessage
	if c.LogConfig != nil {
		logConfMsg = serial.ToTypedMessage(c.LogConfig.Build())
	} else {
		logConfMsg = serial.ToTypedMessage(DefaultLogConfig())
	}
	// let logger module be the first App to start,
	// so that other modules could print log during initiating
	config.App = append([]*serial.TypedMessage{logConfMsg}, config.App...)

	if c.RouterConfig != nil {
		routerConfig, err := c.RouterConfig.Build()
		if err != nil {
			return nil, err
		}
		config.App = append(config.App, serial.ToTypedMessage(routerConfig))
	}

	if c.Policy != nil {
		pc, err := c.Policy.Build()
		if err != nil {
			return nil, err
		}
		config.App = append(config.App, serial.ToTypedMessage(pc))
	}

	var inbounds []InboundDetourConfig

	if len(c.InboundConfigs) > 0 {
		inbounds = append(inbounds, c.InboundConfigs...)
	}

	for _, rawInboundConfig := range inbounds {
		ic, err := rawInboundConfig.Build()
		if err != nil {
			return nil, err
		}
		config.Inbound = append(config.Inbound, ic)
	}

	var outbounds []OutboundDetourConfig

	if len(c.OutboundConfigs) > 0 {
		outbounds = append(outbounds, c.OutboundConfigs...)
	}

	for _, rawOutboundConfig := range outbounds {
		oc, err := rawOutboundConfig.Build()
		if err != nil {
			return nil, err
		}
		config.Outbound = append(config.Outbound, oc)
	}

	return config, nil
}
