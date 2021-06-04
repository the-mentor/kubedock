package cmd

import (
	"flag"
	"os"
	"path/filepath"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/klog"

	"github.com/joyrex2001/kubedock/internal"
	"github.com/joyrex2001/kubedock/internal/config"
)

var cfgFile string

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Start the kubedock api server",
	Run: func(cmd *cobra.Command, args []string) {
		flag.Set("v", viper.GetString("verbosity"))
		internal.Main()
	},
}

func init() {
	rootCmd.AddCommand(serverCmd)

	klog.InitFlags(nil)
	// pflag.CommandLine.AddGoFlagSet(flag.CommandLine)

	serverCmd.PersistentFlags().String("listen-addr", ":8999", "Webserver listen address")
	serverCmd.PersistentFlags().String("unix-socket", "", "Unix socket to listen to (instead of port)")
	serverCmd.PersistentFlags().Bool("tls-enable", false, "Enable TLS on api server")
	serverCmd.PersistentFlags().String("tls-key-file", "", "TLS keyfile")
	serverCmd.PersistentFlags().String("tls-cert-file", "", "TLS certificate file")
	serverCmd.PersistentFlags().StringP("namespace", "n", getContextNamespace(), "Namespace in which containers should be orchestrated")
	serverCmd.PersistentFlags().String("initimage", config.Image, "Image to use as initcontainer for volume setup")
	serverCmd.PersistentFlags().BoolP("inspector", "i", false, "Enable image inspect to fetch container port config from a registry")
	serverCmd.PersistentFlags().DurationP("timeout", "t", 1*time.Minute, "Container creating timeout")
	serverCmd.PersistentFlags().DurationP("reapmax", "r", 15*time.Minute, "Reap all resources older than this time")
	serverCmd.PersistentFlags().Bool("lock", false, "Lock namespace for this instance")
	serverCmd.PersistentFlags().Duration("lock-timeout", 15*time.Minute, "Max time trying to acquire namespace lock")
	serverCmd.PersistentFlags().StringP("verbosity", "v", "1", "Log verbosity level")
	serverCmd.PersistentFlags().BoolP("prune-start", "P", false, "Prune all existing kubedock resources before starting")

	viper.BindPFlag("server.listen-addr", serverCmd.PersistentFlags().Lookup("listen-addr"))
	viper.BindPFlag("server.socket", serverCmd.PersistentFlags().Lookup("unix-socket"))
	viper.BindPFlag("server.tls-enable", serverCmd.PersistentFlags().Lookup("tls-enable"))
	viper.BindPFlag("server.tls-cert-file", serverCmd.PersistentFlags().Lookup("tls-cert-file"))
	viper.BindPFlag("server.tls-key-file", serverCmd.PersistentFlags().Lookup("tls-key-file"))
	viper.BindPFlag("kubernetes.namespace", serverCmd.PersistentFlags().Lookup("namespace"))
	viper.BindPFlag("kubernetes.initimage", serverCmd.PersistentFlags().Lookup("initimage"))
	viper.BindPFlag("kubernetes.timeout", serverCmd.PersistentFlags().Lookup("timeout"))
	viper.BindPFlag("registry.inspector", serverCmd.PersistentFlags().Lookup("inspector"))
	viper.BindPFlag("reaper.reapmax", serverCmd.PersistentFlags().Lookup("reapmax"))
	viper.BindPFlag("lock.enabled", serverCmd.PersistentFlags().Lookup("lock"))
	viper.BindPFlag("lock.timeout", serverCmd.PersistentFlags().Lookup("lock-timeout"))
	viper.BindPFlag("verbosity", serverCmd.PersistentFlags().Lookup("verbosity"))
	viper.BindPFlag("prune-start", serverCmd.PersistentFlags().Lookup("prune-start"))

	viper.BindEnv("server.listen-addr", "SERVER_LISTEN_ADDR")
	viper.BindEnv("server.tls-enable", "SERVER_TLS_ENABLE")
	viper.BindEnv("server.tls-cert-file", "SERVER_TLS_CERT_FILE")
	viper.BindEnv("server.tls-key-file", "SERVER_TLS_KEY_FILE")
	viper.BindEnv("kubernetes.namespace", "NAMESPACE")
	viper.BindEnv("kubernetes.initimage", "INIT_IMAGE")
	viper.BindEnv("kubernetes.timeout", "TIME_OUT")
	viper.BindEnv("reaper.reapmax", "REAPER_REAPMAX")

	// kubeconfig
	if home := homeDir(); home != "" {
		serverCmd.PersistentFlags().String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		serverCmd.PersistentFlags().String("kubeconfig", "", "absolute path to the kubeconfig file")
	}
	viper.BindPFlag("kubernetes.kubeconfig", serverCmd.PersistentFlags().Lookup("kubeconfig"))
}

// homeDir returns the current home directory of the user.
func homeDir() string {
	if h := os.Getenv("HOME"); h != "" {
		return h
	}
	return os.Getenv("USERPROFILE") // windows
}

// getContextNamespace will return the namespace that is set in the current
// kubeconfig context, and returns 'default' if none is set.
func getContextNamespace() string {
	res := "default"
	rul := clientcmd.NewDefaultClientConfigLoadingRules()
	if rul == nil {
		return res
	}
	cfg, err := rul.Load()
	if err != nil {
		return res
	}
	ctx := cfg.Contexts[cfg.CurrentContext]
	if ctx != nil && ctx.Namespace != "" {
		res = ctx.Namespace
	}
	return res
}