/*
Copyright © 2020 Svetoslav Blyahov (svetlio.blyahoff@gmail.com)

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/

// The `cmd` package is in charge of handling the 'front-end' part of the app.
// It takes commands from the CLI, calls the `app` with them and waits for the
// `app` to provide the data that needs  to be printed back to the user.
package cmd

import (
	"github.com/betasve/mstd/app"
	"github.com/spf13/cobra"
	"log"
)

// rootCmd represents the base command when called without any subcommands.
// We choose to have this command only print some help instructions.
var rootCmd = &cobra.Command{
	Use:   "mstd", // TODO: Replace references of `mstd` with the value of args[0]
	Short: "A CLI for Microsoft To Do",
	Long: `This command line interface is for interacting with your account
in Microsoft's To Do application https://todo.microsoft.com/tasks/.

To begin your journey you have to:
1. Configure your personal API key (https://github.com/kiblee/tod0/blob/master/GET_KEY.md)
2. Set personal API key in your config file
3. Start using it (mstd --help)`,
}

// The method that's attached to execute the `root` command. As mentioned above
// it just executes `rootCmd` and logs any possible errors.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}

// Going in the init() method (each time the root command is invoked, and sub-
// commands are provided) we initialize our `app` and get env vars for config
// file path that the user possibly did set.
func init() {
	cobra.OnInitialize(app.InitAppConfig)

	rootCmd.PersistentFlags().StringVar(&app.CfgFilePath, "config", "", "config file (default is $HOME/.mstd.yaml)")
}
