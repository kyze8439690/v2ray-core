package all

import (
	// The following are necessary as they register handlers in their init functions.

	// Required features. Can't remove unless there is replacements.
	_ "v2ray.com/core/app/dispatcher"
	_ "v2ray.com/core/app/proxyman/inbound"
	_ "v2ray.com/core/app/proxyman/outbound"

	// Other optional features.
	_ "v2ray.com/core/app/log"
	_ "v2ray.com/core/app/policy"
	_ "v2ray.com/core/app/router"
	_ "v2ray.com/core/app/stats"

	// Inbound and outbound proxies.
	_ "v2ray.com/core/proxy/socks"
	_ "v2ray.com/core/proxy/vmess/outbound"

	// Transports
	_ "v2ray.com/core/transport/internet/http"
	_ "v2ray.com/core/transport/internet/tcp"
	_ "v2ray.com/core/transport/internet/tls"
	_ "v2ray.com/core/transport/internet/udp"
	_ "v2ray.com/core/transport/internet/websocket"

	// Transport headers
	_ "v2ray.com/core/transport/internet/headers/http"
	_ "v2ray.com/core/transport/internet/headers/tls"

	// JSON config support. Choose only one from the two below.
	// The following line loads JSON from v2ctl
	// _ "v2ray.com/core/main/json"
	// The following line loads JSON internally
	_ "v2ray.com/core/main/jsonem"

	// Load config from file or http(s)
	_ "v2ray.com/core/main/confloader/external"
)
