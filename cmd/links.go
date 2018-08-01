// Copyright © 2018 Daniel <danielfbm@gmail.com>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"fmt"
	"strings"

	"github.com/danielfbm/linkcrawler/pkg"

	"github.com/spf13/cobra"
)

var linksConfig = pkg.LinkConfig{}

// linksCmd represents the links command
var linksCmd = &cobra.Command{
	Use:   "links",
	Short: "craw a website and fetch all links",
	Run: func(cmd *cobra.Command, args []string) {
		links := pkg.FetchLinks(linksConfig)
		fmt.Println(strings.Join(links, " "))
	},
}

func init() {
	rootCmd.AddCommand(linksCmd)
	linksCmd.Flags().StringVarP(&linksConfig.Host, "address", "a", "http://localhost:1313/", "host to start crawling")
	linksCmd.Flags().BoolVarP(&linksConfig.ExternalLinks, "enable-external", "e", false, "crawl external links")
	linksCmd.Flags().BoolVarP(&linksConfig.RespectTree, "respect-tree", "t", true, "respect url hierarchy")
}
