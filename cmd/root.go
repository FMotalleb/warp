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
	"net/url"
	"os"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
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
			listen, err := cmd.Flags().GetString("listen")
			if err != nil {
				logrus.Fatalln(err)
			}
			if listen == "" {
				logrus.Fatalln("listen address cannot be empty")
			}
			remote, err := cmd.Flags().GetString("remote")
			if err != nil {
				logrus.Fatalln(err)
			}
			if remote == "" {
				logrus.Fatalln("remote address cannot be empty")
			}

			listenUri, err := url.Parse(listen)
			if err != nil {
				logrus.Fatalln(err)
			}
			remoteUri, err := url.Parse(remote)
			if err != nil {
				logrus.Fatalln(err)
			}
			threads, err := cmd.Flags().GetUint16("threads")
			if err != nil {
				logrus.Fatalln(err)
			}
			if threads == 0 {
				logrus.Fatalln("threads value cannot be 0")
			}
			Params = &Config{
				listenUri,
				remoteUri,
				threads,
			}
		},
	}

	Params *Config
)

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().StringP("listen", "l", "tcp://127.0.0.1:8080", "Listen Address")
	rootCmd.Flags().StringP("remote", "r", "", "Forward any request received from listen address to this address")
	rootCmd.Flags().Uint16P("threads", "t", 50, "Thread(Goroutine) count")
}
