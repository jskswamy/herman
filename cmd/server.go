// Copyright © 2019 Krishnaswamy Subramanian
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/hashicorp/mdns"
	"github.com/jskswamy/herman/api"
	"github.com/jskswamy/herman/cmd/cli"
	"github.com/jskswamy/herman/pkg/db"
	"github.com/spf13/cobra"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// serverCmd represents the server command
var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Bind herman server and broadcast its information using mdns",
	Run: func(cmd *cobra.Command, args []string) {
		router := gin.Default()
		hostname, err := os.Hostname()
		cli.DieIf(err)

		err = db.Open(dbPath)
		cli.DieIf(err)

		txtRcd := []string{
			"na=the ultimate broadcaster",
			fmt.Sprintf("ve=%s", BuildVersion()),
		}
		mdnsService, err := mdns.NewMDNSService(hostname, "_herman._tcp", "", "", bindPort, nil, txtRcd)
		cli.DieIf(err)

		mdnsServer, err := mdns.NewServer(&mdns.Config{Zone: mdnsService})
		cli.DieIf(err)

		cleanup := func() {
			cli.Info("stopping advertising")
			err := mdnsServer.Shutdown()
			cli.DieIf(err)

			err = db.Close()
			cli.DieIf(err)
		}

		api.Bind(router)
		srv := &http.Server{
			Addr:    fmt.Sprintf("%s:%d", bindAddress, bindPort),
			Handler: router,
		}

		go func() {
			// service connections
			if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
				cli.DieIf(err)
			}
		}()

		// Wait for interrupt signal to gracefully shutdown the mdnsServer with
		// a timeout of 5 seconds.
		quit := make(chan os.Signal)
		// kill (no param) default send syscanll.SIGTERM
		// kill -2 is syscall.SIGINT
		// kill -9 is syscall. SIGKILL but can"t be catch, so don't need add it
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		<-quit
		cli.Warn("\nshutdown Server ...")

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := srv.Shutdown(ctx); err != nil {
			cli.Success("server Shutdown:", err)
		}
		// catching ctx.Done(). timeout of 5 seconds.
		select {
		case <-ctx.Done():
			cleanup()
		}
		cli.Success("done")
	},
}

var bindAddress string
var bindPort int
var dbPath string

func init() {
	rootCmd.AddCommand(serverCmd)
	serverCmd.Flags().StringVarP(&bindAddress, "bind-address", "", "0.0.0.0", "specify the advertise address to use")
	serverCmd.Flags().IntVarP(&bindPort, "bind-port", "p", 5624, "specify the advertise port to use")
	serverCmd.Flags().StringVarP(&dbPath, "db-path", "", "herman.db", "specify the database path where db will be stored")
}
