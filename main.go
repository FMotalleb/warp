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
	"time"

	"github.com/sirupsen/logrus"

	"github.com/FMotalleb/warp/cmd"
)

var remoteAddr string

func main() {
	cmd.Execute()
	if cmd.Params == nil {
		return
	}
	logrus.Infof("forwarding requests received from `%s`, to `%s`", cmd.Params.Listen, cmd.Params.Remote)
	listenAddr := fmt.Sprintf("%s:%s", cmd.Params.Listen.Hostname(), cmd.Params.Listen.Port())
	remoteAddr = fmt.Sprintf("%s:%s", cmd.Params.Remote.Hostname(), cmd.Params.Remote.Port())
	listener, err := net.Listen(cmd.Params.Listen.Scheme, listenAddr)
	if err != nil {
		logrus.Fatalln(err)
	}
	defer listener.Close()
	for i := cmd.Params.Threads; i > 0; i-- {
		go listen(listener)
	}
	for {
		time.Sleep(time.Second)
	}
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
	remoteConnectionForwarded, err := net.Dial(cmd.Params.Remote.Scheme, remoteAddr)
	if err != nil {
		logrus.Warn(err)
		return
	}
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		io.Copy(local, remoteConnectionForwarded)
		wg.Done()
	}()
	wg.Add(1)
	go func() {
		io.Copy(remoteConnectionForwarded, local)
		wg.Done()
	}()
	wg.Wait()
}
