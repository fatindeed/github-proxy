package cmd

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/fatindeed/github-proxy/services"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/txn2/txeh"
)

func init() {
	rootCmd.PersistentFlags().StringVar(&certFile, "cert-file", "github.com.pem", "Certificate for the proxy")
	rootCmd.PersistentFlags().StringVar(&keyFile, "key-file", "github.com-key.pem", "Private key for the proxy")
	rootCmd.PersistentFlags().StringVar(&proxyTarget, "proxy", "https://mirror.ghproxy.com", "Proxy target")
}

var (
	// Version is the version of the CLI injected in compilation time
	Version     = "dev"
	certFile    string
	keyFile     string
	proxyTarget string
	rootCmd     = &cobra.Command{
		Use:     "github-proxy",
		Version: fmt.Sprintf("1.0.0, build %s", Version),
		Short:   "Start a GitHub proxy",
		RunE: func(cmd *cobra.Command, args []string) error {
			hostArr := []string{"github.com", "raw.githubusercontent.com", "gist.github.com", "gist.githubusercontent.com"}

			hosts, err := txeh.NewHostsDefault()
			if err != nil {
				return fmt.Errorf("new hosts error: %w", err)
			}
			hosts.AddHosts("127.0.0.1", hostArr)
			err = hosts.Save()
			if err != nil {
				return fmt.Errorf("hosts save error: %w", err)
			}
			logrus.Debugf("host saved")

			srv := &http.Server{Handler: services.NewReverseProxyHandler(hostArr, proxyTarget)}
			done := make(chan bool)
			go func() {
				sigs := make(chan os.Signal, 1)
				signal.Notify(sigs, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
				sig := <-sigs
				logrus.Debugf("signal received: %v", sig)

				hosts.RemoveHosts(hostArr)
				err = hosts.Save()
				if err != nil {
					logrus.Errorf("hosts save error: %v", err)
				}
				logrus.Debugf("host reverted")

				ctx, cancel := context.WithTimeout(cmd.Context(), 5*time.Second)
				defer cancel()
				err = srv.Shutdown(ctx)
				if err != nil {
					logrus.Errorf("server shutdown error: %v", err)
				}
				logrus.Debugf("server shutdown")

				done <- true
			}()

			logrus.Infof("starting github proxy")
			err = srv.ListenAndServeTLS(certFile, keyFile)
			if err != nil && err != http.ErrServerClosed {
				return fmt.Errorf("start server error: %w", err)
			}
			<-done
			logrus.Infof("github proxy exited")
			return nil
		},
		SilenceUsage: true,
	}
)

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
