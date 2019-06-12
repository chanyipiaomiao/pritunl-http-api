package common

import (
	"net"
	"time"
)

// IsPortOpen 检测端口是否已经打开状态
func IsPortOpen(ip, port string) (bool, error) {
	var (
		conn net.Conn
		err  error
	)

	if conn, err = net.DialTimeout("tcp", ip+":"+port, 2*time.Second); err != nil {
		return false, err
	}

	defer conn.Close()

	return true, nil
}
