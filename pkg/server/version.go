package server

import (
	"fmt"
	"os"
	"runtime/debug"

	"github.com/dxhbiz/go-ntrip-proxy/pkg/version"
	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version of ntrip-proxy",
	Run: func(cmd *cobra.Command, args []string) {
		goVersion := "Unknown"
		build, ok := debug.ReadBuildInfo()
		if ok {
			goVersion = build.GoVersion
		}

		ver := fmt.Sprintf("ntrip-proxy: %s\r\ngo: %s\r\ngit: %s\r\ngit repo: %s\r\n", version.RELEASE, goVersion, version.COMMIT, version.REPO)
		fmt.Print(ver)
		os.Exit(0)
	},
}
