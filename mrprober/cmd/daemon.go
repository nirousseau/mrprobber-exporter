package cmd

import (
	"context"
	"fmt"
	"github.com/VictoriaMetrics/metrics"
	"github.com/spf13/cobra"
	"log"
	"mrprober/conf"
	"mrprober/engine"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const DAEMON_SHUTDOWN_TIMEOUT = 5 * time.Second

func newDaemonCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "daemon",
		Short:   "Start as a daemon/service and execute rules on demand.",
		Aliases: []string{"serve"},
		Args:    cobra.ExactArgs(0),
		Run: func(cmd *cobra.Command, args []string) {

			// Read web server configuration
			config := conf.SafeConfiguration.Get()

			// Expose the registered metrics at `/metrics` path.
			http.HandleFunc(config.Global.Web.MetricsPath, func(w http.ResponseWriter, req *http.Request) {
				metrics.WritePrometheus(w, true)
			})

			// Attach raw handler
			http.HandleFunc("/raw", func(w http.ResponseWriter, r *http.Request) {

				// Execute and print results
				for r := range engine.OneShotRun() {
					_, err := fmt.Fprintln(w, r)
					if err != nil {
						log.Fatal(err)
					}
				}
			})

			// Catch OS proc signals
			stop := make(chan os.Signal, 1)
			signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

			// Start background polling
			quit := make(chan int)
			engine.StartActivePolling(quit)

			// handle on http service
			h := &http.Server{Addr: config.Global.Web.ListenAddress, Handler: nil}

			// Bootstrapping web server
			go func() {
				log.Println(fmt.Sprintf("Listening on %s", config.Global.Web.ListenAddress))
				if err := h.ListenAndServe(); err != nil {
					log.Fatal(err)
				}
			}()
			// Block until signal is catched
			<-stop

			// notify end of active polling
			quit <- 0

			ctx, cancel := context.WithTimeout(context.Background(), DAEMON_SHUTDOWN_TIMEOUT)
			defer cancel()

			log.Println("Shutting down...")
			if err := h.Shutdown(ctx); err != nil {
				log.Fatalf("Error: %v\n", err)
			}
		},
	}

	return cmd
}
