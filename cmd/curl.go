// Copyright © 2017 Volt Grid Pty Ltd <admin@voltgrid.com>
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"log"
	"log/syslog"
	"strings"
	"bytes"
	"syscall"
	"path"

	"github.com/spf13/cobra"
	// "github.com/davecgh/go-spew/spew"
)

// curlCmd represents the curl command
var curlCmd = &cobra.Command{
	Use:   "curl",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: runCurl,
	Args: cobra.ArbitraryArgs,
}

func init() {
	RootCmd.AddCommand(curlCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// curlCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// curlCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func runCurl(cmd *cobra.Command, args []string) {

	// Setup the logger
	var syslogTag = path.Base(os.Args[0])
	sysLog, err := syslog.Dial("unixgram", "/dev/log", syslog.LOG_WARNING|syslog.LOG_DAEMON, syslogTag)
	if err != nil {
		log.Fatal(err)
	}

	execCmd := exec.Command("/usr/bin/curl", args[0:]...)
	var buffStdout bytes.Buffer
	var buffStderr bytes.Buffer
	execCmd.Stdout = &buffStdout
	execCmd.Stderr = &buffStderr
	if err := execCmd.Start(); err != nil {
		fmt.Fprintf(sysLog, "execCmd.Start: %v")
		return
	}

	err = execCmd.Wait()

	if len(buffStderr.Bytes()) > 0 {
		fmt.Fprintf(sysLog, "Error: ", strings.TrimSpace(buffStderr.String()))
	}
	if len(buffStdout.Bytes()) > 0 {
		fmt.Fprintf(sysLog, "Info: ", strings.TrimSpace(buffStdout.String()))
	}

	if err != nil {

		if exiterr, ok := err.(*exec.ExitError); ok {
			// The program has exited with an exit code != 0

			// This works on both Unix and Windows. Although package
			// syscall is generally platform dependent, WaitStatus is
			// defined for both Unix and Windows and in both cases has
			// an ExitStatus() method with the same signature.
			if status, ok := exiterr.Sys().(syscall.WaitStatus); ok {
				fmt.Fprintf(sysLog, "Exit Status: %d", status.ExitStatus())
				fmt.Println("Exit Status: %d", status.ExitStatus())
				fmt.Println("Err: ", strings.TrimSpace(buffStderr.String()))
				fmt.Println("Std: ", strings.TrimSpace(buffStdout.String()))
			}
		} else {
			fmt.Fprintf(sysLog, "execCmd.Wait: %v", err)
			return
		}
	}
}
