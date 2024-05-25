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
	"os"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
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
			listenProto := getString(cmd.Flags(), listenProtoFlag)
			listen := getString(cmd.Flags(), listenAddrFlag)
			listenPort := getUint16(cmd.Flags(), listenPortFlag)
			remoteProto := getString(cmd.Flags(), remoteProtoFlag)
			remote := getString(cmd.Flags(), remoteAddrFlag)
			remotePort := getUint16(cmd.Flags(), remotePortFlag)
			threads := getUint16(cmd.Flags(), threadsFlag)

			Params = &Config{
				ListenProto: listenProto,
				ListenAddr:  listen,
				ListenPort:  listenPort,
				RemoteProto: remoteProto,
				RemoteAddr:  remote,
				RemotePort:  remotePort,
				Threads:     threads,
				Timeout:     time.Minute,
			}
		},
	}

	Params *Config
)

func getString(flags *pflag.FlagSet, name FlagName) string {
	result, err := flags.GetString(name)
	if err != nil {
		logrus.Fatalln(err, ": ", name)
	}
	if result == "" {
		logrus.Fatalf("%s cannot be empty", name)
	}
	return result
}

func getUint16(flags *pflag.FlagSet, name FlagName) uint16 {
	result, err := flags.GetUint16(name)
	if err != nil {
		logrus.Fatalln(err)
	}
	if result == 0 {
		logrus.Fatalf("%s cannot be 0", name)
	}
	return result
}

func getDuration(flags *pflag.FlagSet, name FlagName) time.Duration {
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
	rootCmd.Flags().StringP(listenAddrFlag, "l", "127.0.0.1", "Listen Address")
	rootCmd.Flags().Uint16P(listenPortFlag, "p", 8080, "Listen Port")
	rootCmd.Flags().String(listenProtoFlag, "tcp", "Listen Protocol")

	rootCmd.Flags().String(remoteProtoFlag, "tcp", "Remote protocol")
	rootCmd.Flags().StringP(remoteAddrFlag, "r", "", "Forward any request received from listen address to this address")
	rootCmd.Flags().Uint16(remotePortFlag, 0, "Forward any request received from listen address to this port")

	rootCmd.Flags().Uint16(threadsFlag, 50, "Thread(Goroutine) count")
	rootCmd.Flags().DurationP(timeoutFlag, "t", time.Minute, "Connection Timeout")
}

type FlagName = string

var (
	listenAddrFlag  FlagName = "listen-address"
	listenPortFlag  FlagName = "listen-port"
	listenProtoFlag FlagName = "listen-protocol"
	remoteAddrFlag  FlagName = "remote-address"
	remotePortFlag  FlagName = "remote-port"
	remoteProtoFlag FlagName = "remote-protocol"
	threadsFlag     FlagName = "threads"
	timeoutFlag     FlagName = "timeout"
)
