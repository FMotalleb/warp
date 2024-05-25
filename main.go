/*
Copyright © 2024 Motalleb Fallahnezhad

This program is free software; you can redistribute it and/or
modify it under the terms of the GNU General Public License
as published by the Free Software Foundation; either version 2
of the License, or (at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program. If not, see <http://www.gnu.org/licenses/>.
*/
package main

import (
	"fmt"
	"io"
	"net"
	"sync"

	"github.com/sirupsen/logrus"

	"github.com/FMotalleb/warp/cmd"
)

var remoteAddr string

func main() {
	cmd.Execute()
	if cmd.Params == nil {
		return
	}
	logrus.Infof("forwarding requests received from `%s`:`%d`, to `%s`:`%d`", cmd.Params.ListenAddr, cmd.Params.ListenPort, cmd.Params.RemoteAddr, cmd.Params.RemotePort)
	listenAddr := fmt.Sprintf("%s:%d", cmd.Params.ListenAddr, cmd.Params.ListenPort)
	remoteAddr = fmt.Sprintf("%s:%d", cmd.Params.RemoteAddr, cmd.Params.RemotePort)
	listener, err := net.Listen(cmd.Params.ListenProto, listenAddr)
	if err != nil {
		logrus.Fatalln(err)
	}
	defer listener.Close()
	for i := cmd.Params.Threads; i > 0; i-- {
		go listen(listener)
	}
	make(chan interface{}) <- 0
}

func listen(listener net.Listener) {
	for {
		conn, err := listener.Accept()
		if err != nil {
			logrus.Warn(err)
			continue
		}
		go handlePortForward(conn)
	}
}

func handlePortForward(local net.Conn) {
	remoteConnectionForwarded, err := net.DialTimeout(cmd.Params.RemoteProto, remoteAddr, cmd.Params.Timeout)
	if err != nil {
		logrus.Warn(err)
		return
	}
	defer remoteConnectionForwarded.Close()

	wg := sync.WaitGroup{}
	wg.Add(2)
	go func() {
		_, err := io.Copy(local, remoteConnectionForwarded)
		if err != nil {
			logrus.Warn(err)
		}
		wg.Done()
	}()
	go func() {
		_, err := io.Copy(remoteConnectionForwarded, local)
		if err != nil {
			logrus.Warn(err)
		}
		wg.Done()
	}()
	wg.Wait()
}
