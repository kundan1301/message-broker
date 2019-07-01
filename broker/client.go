package broker

import (
	"net"
)

type info struct {
	clientID  string
	keepalive uint16
	localIP   string
	remoteIP  string
}

type Client struct {
	conn net.Conn
	info info
}
