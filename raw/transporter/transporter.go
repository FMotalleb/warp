package transporter

import (
	"fmt"
	"io"
	"net"
	"sync"

	"github.com/sirupsen/logrus"

	"github.com/FMotalleb/warp/config"
	"github.com/FMotalleb/warp/interceptor"
)

func Listen(listener net.Listener, cfg *config.RawConfig) {
	for {
		conn, err := listener.Accept()
		if err != nil {
			logrus.Warn(err)
			continue
		}
		go handlePortForward(conn, cfg)
	}
}

func handlePortForward(local net.Conn, cfg *config.RawConfig) {
	var remoteConnectionForwarded net.Conn
	var err error

	if cfg.Timeout != 0 {
		remoteConnectionForwarded, err = net.DialTimeout(cfg.RemoteProto, fmt.Sprintf("%s:%d", cfg.RemoteAddr, cfg.RemotePort), cfg.Timeout)
	} else {
		remoteConnectionForwarded, err = net.Dial(cfg.RemoteProto, fmt.Sprintf("%s:%d", cfg.RemoteAddr, cfg.RemotePort))
	}
	if err != nil {
		logrus.Warn(err)
		return
	}
	defer func() {
		err := remoteConnectionForwarded.Close()
		if err != nil {
			logrus.Warn(err)
		}
	}()

	wg := sync.WaitGroup{}
	wg.Add(2)
	go func() {
		out := &interceptor.Interceptor{
			Prefix:        fmt.Sprintf("%s <", local.RemoteAddr()),
			HasStdOutput:  cfg.Intercept || cfg.Base64Intercept,
			Base64Encoded: cfg.Base64Intercept,
			Output:        local,
		}
		_, err := io.Copy(out, remoteConnectionForwarded)
		if err != nil {
			logrus.Warn(err)
		}
		wg.Done()
	}()
	go func() {
		out := &interceptor.Interceptor{
			Prefix:        fmt.Sprintf("%s >", local.RemoteAddr()),
			HasStdOutput:  cfg.Intercept || cfg.Base64Intercept,
			Base64Encoded: cfg.Base64Intercept,
			Output:        remoteConnectionForwarded,
		}
		_, err := io.Copy(out, local)
		if err != nil {
			logrus.Warn(err)
		}
		wg.Done()
	}()
	wg.Wait()
}
