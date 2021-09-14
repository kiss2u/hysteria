package transport

import (
	"net"
)

type Transport interface {
	QUICResolveUDPAddr(address string) (*net.UDPAddr, error)
	QUICListenUDP(laddr *net.UDPAddr) (*net.UDPConn, error)

	LocalResolveIPAddr(address string) (*net.IPAddr, error)
	LocalResolveTCPAddr(address string) (*net.TCPAddr, error)
	LocalResolveUDPAddr(address string) (*net.UDPAddr, error)
	LocalDial(network, address string) (net.Conn, error)
	LocalDialTCP(laddr, raddr *net.TCPAddr) (*net.TCPConn, error)
	LocalListenTCP(laddr *net.TCPAddr) (*net.TCPListener, error)
	LocalListenUDP(laddr *net.UDPAddr) (*net.UDPConn, error)
}
