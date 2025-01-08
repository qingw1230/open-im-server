package utils

import (
	"net"

	"github.com/qingw1230/studyim/pkg/common/config"
)

var ServerIP = ""

func init() {
	if config.Config.ServerIP != "" {
		ServerIP = config.Config.ServerIP
		return
	}

	conn, err := net.Dial("udp", "114.114.114.114:80")
	if err != nil {
		panic(err.Error())
	}
	defer conn.Close()
	localAddr := conn.LocalAddr().(*net.UDPAddr)
	ServerIP = localAddr.IP.String()
}
