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
	"bufio"
	"os"
	"io"
	"fmt"
	"regexp"
	"strings"
	"log"

	"github.com/spf13/cobra"
	// "github.com/davecgh/go-spew/spew"
	"github.com/mattn/go-shellwords"
)

// genCmd represents the gen command
var genCmd = &cobra.Command{
	Use:   "gen",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: HelloWorld,
}
var (
	wrapper string
	envFile string
)

func init() {
	RootCmd.AddCommand(genCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// genCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// genCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	genCmd.Flags().StringVarP(&wrapper, "wrapper", "w", fmt.Sprintf("%s curl --", os.Args[0]), "Overwrite default curl wrapper")
	genCmd.Flags().StringVarP(&envFile, "env-file", "e", "", "Extra env file to import")
}

func HelloWorld(cmd *cobra.Command, args []string) {
	reader := bufio.NewReader(os.Stdin)
	line := 1
	envMatch, _ := regexp.Compile(`^[A-Z][A-Z0-9_]+=`)
	cronMatch, _ := regexp.Compile(`^[0-9/\*]+\ [0-9/\*]+\ [0-9/\*]+\ [0-9/\*]+\ [0-9/\*]+`)
	envs := make([]string, 0)
	crons :=make([][]string, 0)
	for {
		input, err := reader.ReadString('\n')
		if err != nil && err == io.EOF {
			break
		}

		if envMatch.MatchString(input) == true {
			log.Printf("ENV  %2d: %s", line, input)
			envs = append(envs, input)
		} else if cronMatch.MatchString(input) == true {
			log.Printf("CRON %2d: %s", line, input)
			args, _ := shellwords.Parse(input)
			crons = append(crons, args)
		}
		line++
	}

	for _, item := range crons {
		// spew.Dump(item)
		for _, time := range item[0:5] {
			fmt.Printf("%s ", time)
		}
		if envFile != "" {
			fmt.Printf(". %s; ", envFile)
		}
		for _, env := range envs {
			results := strings.SplitN(strings.TrimSpace(env), "=", 2)
			fmt.Printf("%s=\"%s\"; ", results[0], results[1])
		}
		fmt.Printf("%s ", wrapper)
		for _, arg := range item[5:] {
			fmt.Printf("\"%s\" ", strings.Replace(arg, "\"", "\\\"", -1))
		}
		fmt.Println("")
	}

}
