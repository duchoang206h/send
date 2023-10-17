package cmd

import (
	"context"
	"os"
	"sync"

	"github.com/duchoang206h/send-cli/internal/send"
	"github.com/spf13/cobra"
)

var (
	pathFileFlag      string
	pathDirectoryFlag string
)
// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "",
	Short: "A brief description of your application",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
	Run:  func(_ *cobra.Command, _ []string) { 
		ctx := context.Background()
		rsChan := make(chan string)
		errChan :=make(chan error)
		var wait sync.WaitGroup
		wait.Add(1)
		go func (rs chan string, err chan error, wg *sync.WaitGroup)  {
			defer wg.Done()
			send.PrintResult(rs, err)
		}(rsChan, errChan, &wait)
		var (
			result string
			err error
		)
		switch true {
			case pathFileFlag != "":
				result, err = send.SendFile(ctx, pathFileFlag);
			case pathDirectoryFlag != "":
				result, err = send.SendDirectory(ctx, pathDirectoryFlag);
		}
		if err != nil {
			errChan <- err
		}else {
			rsChan <-result
		}
		wait.Wait()
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.send.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().StringVarP(&pathFileFlag ,"file", "f", "", "Path to file")
	rootCmd.Flags().StringVarP(&pathDirectoryFlag,"directory", "d", "", "Path to directory")
}


