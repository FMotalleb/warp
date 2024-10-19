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
	"os"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"

	"github.com/FMotalleb/warp/config"
	"github.com/FMotalleb/warp/raw/transporter"
)

// rootCmd represents the base command when called without any subcommands
var (
	rootCmd = &cobra.Command{
		Use:   "warp",
		Short: "A brief description of your application",
		Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
		// Uncomment the following line if your bare application
		// has an action associated with it:
		Run: func(cmd *cobra.Command, args []string) {
			if rawParams == nil {
				return
			}
			rawParams.GlobalConfig = *globalParams

			logrus.Infof("forwarding requests received from `%s`:`%d`, to `%s`:`%d`", rawParams.ListenAddr, rawParams.ListenPort, rawParams.RemoteAddr, rawParams.RemotePort)
			listenAddr := fmt.Sprintf("%s:%d", rawParams.ListenAddr, rawParams.ListenPort)

			listener, err := net.Listen(rawParams.ListenProto, listenAddr)
			if err != nil {
				logrus.Fatalln(err)
			}
			defer func() {
				err := listener.Close()
				logrus.Warnf("failed to close listener: %v", err)
			}()
			for i := globalParams.Threads; i > 0; i-- {
				go transporter.Listen(listener, rawParams)
			}
			make(chan interface{}) <- 0
		},
	}

	rawParams    *config.RawConfig    = &config.RawConfig{}
	globalParams *config.GlobalConfig = &config.GlobalConfig{}
)

func getString(flags *pflag.FlagSet, name flagName) string {
	result, err := flags.GetString(name)
	if err != nil {
		logrus.Fatalln(err, ": ", name)
	}
	if result == "" {
		logrus.Fatalf("%s cannot be empty", name)
	}
	return result
}

func getUint16(flags *pflag.FlagSet, name flagName) uint16 {
	result, err := flags.GetUint16(name)
	if err != nil {
		logrus.Fatalln(err)
	}
	if result == 0 {
		logrus.Fatalf("%s cannot be 0", name)
	}
	return result
}

func getDuration(flags *pflag.FlagSet, name flagName) time.Duration {
	result, err := flags.GetDuration(name)
	if err != nil {
		logrus.Fatalln(err)
	}
	if result == 0 {
		logrus.Fatalf("%s cannot be 0", name)
	}
	return result
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().StringVarP(&rawParams.ListenAddr, listenAddrFlag, "l", "127.0.0.1", "Listen Address")
	rootCmd.Flags().Uint16VarP(&rawParams.ListenPort, listenPortFlag, "o", 8080, "Listen Port")
	rootCmd.Flags().StringVar(&rawParams.ListenProto, listenProtoFlag, "tcp", "Listen Protocol")

	rootCmd.Flags().StringVarP(&rawParams.RemoteAddr, remoteAddrFlag, "r", "", "Forward any request received from listen address to this address")
	rootCmd.Flags().Uint16VarP(&rawParams.RemotePort, remotePortFlag, "p", 0, "Forward any request received from listen address to this port")
	rootCmd.Flags().StringVar(&rawParams.RemoteProto, remoteProtoFlag, "tcp", "Remote protocol")

	rootCmd.PersistentFlags().Uint16Var(&globalParams.Threads, threadsFlag, 50, "Thread(Goroutine) count")
	rootCmd.PersistentFlags().DurationVarP(&globalParams.Timeout, timeoutFlag, "t", 0, "Connection Timeout")
	rootCmd.PersistentFlags().BoolVar(&globalParams.Intercept, interceptFlag, false, "Printout Transferring data")
	rootCmd.PersistentFlags().BoolVar(&globalParams.Base64Intercept, base64InterceptFlag, false, "Printout Transferring data base64 encoded")
}

type flagName = string

var (
	listenAddrFlag      flagName = "listen-address"
	listenPortFlag      flagName = "listen-port"
	listenProtoFlag     flagName = "listen-protocol"
	remoteAddrFlag      flagName = "remote-address"
	remotePortFlag      flagName = "remote-port"
	remoteProtoFlag     flagName = "remote-protocol"
	threadsFlag         flagName = "threads"
	timeoutFlag         flagName = "timeout"
	interceptFlag       flagName = "intercept"
	base64InterceptFlag flagName = "b64-intercept"
)
