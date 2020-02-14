// +build linux

package udp

import (
	"v2ray.com/core/common/net"
)

func ReadUDPMsg(conn *net.UDPConn, payload []byte, oob []byte) (int, int, int, *net.UDPAddr, error) {
	return conn.ReadMsgUDP(payload, oob)
}
