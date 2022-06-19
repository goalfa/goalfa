package main

import (
	"fmt"
	"github.com/datafony/alfa/assets"
	"github.com/datafony/alfa/cmd/alfa/sdk"
	"github.com/spf13/cobra"
	"log"
	"net/http"
	"os"
)

var rootCmd = &cobra.Command{
	Use:   "buck",
	Short: "Hugo is a very fast static site generator",
	Long: `A Fast and Flexible Static Site Generator built with
                love by spf13 and friends in Go.
                Complete documentation is available at http://hugo.spf13.com`,
	Run: func(cmd *cobra.Command, args []string) {
		// Do Stuff Here
		fmt.Println("Hello world!")
		fs := http.FileServer(assets.Assets)
		http.Handle("/", fs)
		log.Println("Listening on :3000...")
		err := http.ListenAndServe(":3000", nil)
		if err != nil {
			log.Fatal(err)
		}
	},
}

func main() {
	rootCmd.AddCommand(
		sdk.Command,
	)
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
