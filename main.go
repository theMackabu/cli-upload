package main

import (
	"fmt"
	"os"

	"upload/cli"
	"upload/config"
)

func main() {
	filename, private, showHelp, err := cli.ParseArgs(os.Args)
	if err != nil {
		cli.PrintError("%s", err)
		fmt.Println()
		cli.PrintUsage()
		os.Exit(1)
	}

	if showHelp {
		cli.PrintUsage()
		return
	}

	config := config.NewConfig()
	uploader := cli.NewFileUploader(config)

	fmt.Printf("Uploading %s%s%s...\n", cli.ColorCyan, filename, cli.ColorReset)

	response, err := uploader.UploadFile(filename, private)
	if err != nil {
		cli.PrintError("%s", err)
		os.Exit(1)
	}

	cli.PrintSuccess(response, private)
}
