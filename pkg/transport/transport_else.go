// +build !linux

package transport

import (
	"net"
	"time"
)

func NewTransport(localTimeout int, localMark, quicMark uint32) Transport {
	if localTimeout <= 0 {
		localTimeout = 8
	}
	return &transport{
		LocalTimeout: time.Duration(localTimeout) * time.Second,
		LocalMark:    localMark,
		QUICMark:     quicMark,
	}
}

type transport struct {
	LocalTimeout time.Duration
	LocalMark    uint32
	QUICMark     uint32
}

func (t *transport) QUICResolveUDPAddr(address string) (*net.UDPAddr, error) {
	return net.ResolveUDPAddr("udp", address)
}

func (t *transport) QUICListenUDP(laddr *net.UDPAddr) (*net.UDPConn, error) {
	return net.ListenUDP("udp", laddr)
}

func (t *transport) LocalResolveIPAddr(address string) (*net.IPAddr, error) {
	return net.ResolveIPAddr("ip", address)
}

func (t *transport) LocalResolveTCPAddr(address string) (*net.TCPAddr, error) {
	return net.ResolveTCPAddr("tcp", address)
}

func (t *transport) LocalResolveUDPAddr(address string) (*net.UDPAddr, error) {
	return net.ResolveUDPAddr("udp", address)
}

func (t *transport) LocalDial(network, address string) (net.Conn, error) {
	dialer := &net.Dialer{Timeout: t.LocalTimeout}
	return dialer.Dial(network, address)
}

func (t *transport) LocalDialTCP(laddr, raddr *net.TCPAddr) (*net.TCPConn, error) {
	dialer := &net.Dialer{Timeout: t.LocalTimeout, LocalAddr: laddr}
	conn, err := dialer.Dial("tcp", raddr.String())
	if err != nil {
		return nil, err
	}
	return conn.(*net.TCPConn), nil
}

func (t *transport) LocalListenTCP(laddr *net.TCPAddr) (*net.TCPListener, error) {
	return net.ListenTCP("tcp", laddr)
}

func (t *transport) LocalListenUDP(laddr *net.UDPAddr) (*net.UDPConn, error) {
	return net.ListenUDP("udp", laddr)
}
