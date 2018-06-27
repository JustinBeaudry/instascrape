package cli

import (
	"github.com/spf13/cobra"
	"fmt"
	"os"
)

var rootCmd = &cobra.Command{
	Use: "instascrape",
	Short: "scrape instagram photos/captions from users",
	Long: "scrape images/videos/captions from instagram by user",
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
