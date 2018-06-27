package cli

import (
	"github.com/spf13/cobra"
	"instascrape/services"
	"log"
	"encoding/json"
	"time"
	"io/ioutil"
	"path/filepath"
)

var verbose = false
var force = false
var outputPath = ""
var cacheDir = ""

func init() {
	rootCmd.AddCommand(scrapeCmd)
	scrapeCmd.Flags().StringVarP(&outputPath, "output", "o", "", "output path")
	scrapeCmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "verbose output")
	scrapeCmd.Flags().BoolVarP(&force, 	"force", "f", false, "forces scrape to ignore cache")
	scrapeCmd.Flags().StringVarP(&cacheDir, "cachedir", "c", "/tmp/instascrape_cache", "directory for scraping cache")
}

var scrapeCmd = &cobra.Command{
	Use: "scrape [slug]",
	Short: "execute a scrape session",
	Long: "run a scrape session as a single-run, a daemon, or as a cron",
	Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if verbose {
			log.Print("starting scrape of User: " + args[0])
		}
		start := time.Now()
		raw, scrapeErr := services.Scrape(args[0], force, cacheDir)
		if verbose {
			log.Print("scrape complete in " + time.Since(start).String() + "!")
		}
		if scrapeErr != nil {
			log.Fatal(scrapeErr)
		}
		if raw == nil {
			log.Fatal("no feed information from user: " + args[0])
		}
		if verbose {
			log.Print("parsing feed information")
		}
		feed := services.Transform(raw)
		f, jsonErr := json.MarshalIndent(feed, "", " ")
		if jsonErr != nil {
			log.Fatal(jsonErr)
		}
		if verbose {
			log.Print("parse complete!")
		}
		if outputPath != "" {
			if filepath.Ext(outputPath) != ".json" {
				log.Fatal("--write output path must include the .json extension")
			}
			if verbose {
				log.Print("writing to disk!")
			}
			byteFeed := []byte(string(f))
			writeErr := ioutil.WriteFile(outputPath, byteFeed, 0644)
			if writeErr != nil {
				log.Fatal(writeErr)
			}
			if verbose {
				log.Print("disk write complete")
			}
		} else {
			log.Print(string(f))
		}
 	},
}

