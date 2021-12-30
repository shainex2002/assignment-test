package cmd

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"net/http"
	"net/http/httputil"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "simple-dump-server",
	Short: "A http server that dumps all information of requests it gets",
	Run: func(cmd *cobra.Command, args []string) {
		logrus.SetLevel(logrus.DebugLevel)

		//Serve endpoint
		go serve()

		sigterm := make(chan os.Signal, 1)
		signal.Notify(sigterm, syscall.SIGTERM)
		signal.Notify(sigterm, syscall.SIGINT)
		<-sigterm

	},
}

var (
	httpCount = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "harness_canary_skeleton_hhtp_count",
			Help: "Number of http requests to each endpoint",
		},
		[]string{"path"},
	)
)

//HTTP Endpoints
func serve() {
	logrus.Infof("Hello World")

	r := mux.NewRouter()
	//Liveness Endpoint
	r.HandleFunc("/healthz", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	//Readiness Endpoint
	r.HandleFunc("/readyz", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	//Metrics Endpoint
	r.Handle("/metrics", promhttp.Handler())

	//Endpoints
	r.Handle("/test200", dumpRequest(http.StatusOK, "test200"))
	r.Handle("/test400", dumpRequest(http.StatusBadRequest, "test400"))
	r.Handle("/test401", dumpRequest(http.StatusUnauthorized, "test401"))
	r.Handle("/test404", dumpRequest(http.StatusNotFound, "test404"))

	srv := &http.Server{
		Handler:      r,
		Addr:         ":8080",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	logrus.Infof("Starting Server on port 8080....")
	logrus.Fatal(srv.ListenAndServe())
}

func dumpRequest(statusCode int, path string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		// Save a copy of this request for debugging.
		requestDump, err := httputil.DumpRequest(req, true)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(string(requestDump))

		logrus.Infof("Sending response back")
		w.WriteHeader(statusCode)
		httpCount.WithLabelValues(path).Inc()
		return
	})
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		logrus.Fatal(err)
	}
}

func init() {
}
