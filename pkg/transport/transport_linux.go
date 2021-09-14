package transport

import (
	"net"
	"syscall"
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
	conn, err := net.ListenUDP("udp", laddr)
	if err != nil {
		return nil, err
	}
	if err := t.setMark(conn, t.QUICMark); err != nil {
		_ = conn.Close()
		return nil, err
	}
	return conn, nil
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
	conn, err := dialer.Dial(network, address)
	if err != nil {
		return nil, err
	}
	if tc, ok := conn.(*net.TCPConn); ok {
		if err := t.setMark(tc, t.LocalMark); err != nil {
			_ = conn.Close()
			return nil, err
		}
	} else if uc, ok := conn.(*net.UDPConn); ok {
		if err := t.setMark(uc, t.LocalMark); err != nil {
			_ = conn.Close()
			return nil, err
		}
	}
	return conn, nil
}

func (t *transport) LocalDialTCP(laddr, raddr *net.TCPAddr) (*net.TCPConn, error) {
	dialer := &net.Dialer{Timeout: t.LocalTimeout, LocalAddr: laddr}
	conn, err := dialer.Dial("tcp", raddr.String())
	if err != nil {
		return nil, err
	}
	if err := t.setMark(conn.(*net.TCPConn), t.LocalMark); err != nil {
		_ = conn.Close()
		return nil, err
	}
	return conn.(*net.TCPConn), nil
}

func (t *transport) LocalListenTCP(laddr *net.TCPAddr) (*net.TCPListener, error) {
	l, err := net.ListenTCP("tcp", laddr)
	if err != nil {
		return nil, err
	}
	if err := t.setMark(l, t.LocalMark); err != nil {
		_ = l.Close()
		return nil, err
	}
	return l, nil
}

func (t *transport) LocalListenUDP(laddr *net.UDPAddr) (*net.UDPConn, error) {
	l, err := net.ListenUDP("udp", laddr)
	if err != nil {
		return nil, err
	}
	if err := t.setMark(l, t.LocalMark); err != nil {
		_ = l.Close()
		return nil, err
	}
	return l, nil
}

func (t *transport) setMark(conn syscall.Conn, mark uint32) error {
	rc, err := conn.SyscallConn()
	if err != nil {
		return err
	}
	return rc.Control(func(fd uintptr) {
		_ = syscall.SetsockoptInt(int(fd), syscall.SOL_SOCKET, syscall.SO_MARK, int(mark))
	})
}
