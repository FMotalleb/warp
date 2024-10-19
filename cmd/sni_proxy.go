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
package cmd

import (
	"fmt"
	"net"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/FMotalleb/warp/config"
	transporter "github.com/FMotalleb/warp/sni/transport"
)

// sniProxyCmd represents the sniProxy command
var (
	sniProxyCmd = &cobra.Command{
		Use:   "sni",
		Short: "Proxy to target host based on SNI (Server Name Indication)",
		// Long:  ``,
		Run: func(_ *cobra.Command, args []string) {
			if sniParams == nil {
				return
			}
			sniParams.GlobalConfig = *globalParams
			logrus.Infof("Forwarding requests received from `%s`:`%d`, to target host:%d based on SNI", sniParams.ListenAddr, sniParams.ListenPort, sniParams.RemotePort)
			listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", sniParams.ListenAddr, sniParams.ListenPort))
			if err != nil {
				logrus.Fatalf("Fatal error during the creation of listener: %s", err)
				return
			}
			defer func() {
				err := listener.Close()
				if err != nil {
					logrus.Warnf("Failed to close listener: %v", err)
				}
			}()
			transporter.Listen(listener, sniParams)
		},
	}
	sniParams *config.SNIConfig = &config.SNIConfig{}
)

func init() {
	rootCmd.AddCommand(sniProxyCmd)
	sniProxyCmd.Flags().StringVarP(&sniParams.ListenAddr, listenAddrFlag, "l", "127.0.0.1", "Listen Address")
	sniProxyCmd.Flags().Uint16VarP(&sniParams.ListenPort, listenPortFlag, "o", 443, "Listen Port")
	sniProxyCmd.Flags().Uint16VarP(&sniParams.RemotePort, remotePortFlag, "p", 443, "Forward any request received from listen address to this port")
}
