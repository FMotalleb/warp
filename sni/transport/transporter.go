package transporter

import (
	"fmt"
	"io"
	"net"

	"github.com/sirupsen/logrus"

	"github.com/FMotalleb/warp/config"
	"github.com/FMotalleb/warp/sni/parser"
)

// Listen and accept connections and forward them to the target host based on SNI
func Listen(listener net.Listener, sniParams *config.SNIConfig) {
	for {
		conn, err := listener.Accept()
		if err != nil {
			logrus.Warnf("Failed to accept connection: %v", err)
			continue
		}
		for i := uint16(0); i < sniParams.Threads; i++ {
			go handleSNI(conn, sniParams)
		}
	}
}

func handleSNI(conn net.Conn, rawParams *config.SNIConfig) {
	buff := make([]byte, 9)
	_, err := conn.Read(buff)
	if err != nil {
		logrus.Warnf("Failed to read TCP handshake: %v", err)
		return
	}
	followupSize := int(buff[7])*255 + int(buff[8]) + 1
	finalBuffer := make([]byte, len(buff)+followupSize)
	copy(finalBuffer, buff)
	_, err = conn.Read(finalBuffer[len(buff):])
	if err != nil {
		logrus.Warnf("Failed to read the full TCP handshake: %v", err)
		return
	}
	host, err := parser.GetHostname(finalBuffer)
	logrus.Infof("%v\n%s\n", finalBuffer, finalBuffer)
	if err != nil {
		logrus.Warnf("Failed to resolve sni target: %v", err)
		return
	}

	targetAddr := fmt.Sprintf("%s:%d", host, rawParams.RemotePort)
	target, err := net.Dial("tcp", targetAddr)
	if err != nil {
		logrus.Warnf("Failed to connect to target (%s): %v", targetAddr, err)
	}

	_, err = target.Write(finalBuffer)
	if err != nil {
		logrus.Warnf("Failed to push TCP handshake to target(%s): %v", targetAddr, err)
		return
	}
	go func() {
		_, err = io.Copy(conn, target)
		if err != nil {
			logrus.Warnf("Failed to copy data from target(%s) to client(%s): %v", target.RemoteAddr(), conn.RemoteAddr(), err)
			return
		}
	}()
	go func() {
		_, err = io.Copy(target, conn)
		if err != nil {
			logrus.Warnf("Failed to copy data from client(%s) to target(%s): %v", conn.RemoteAddr(), target.RemoteAddr(), err)
			return
		}
	}()
}
