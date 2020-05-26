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

var linksConfig = pkg.LinkConfig{
	Host:          "http://localhost:1313/",
	ExternalLinks: false,
	RespectTree:   true,
	FilterIn:      []string{},
	FilterOut:     []string{},
}

// linksCmd represents the links command
var linksCmd = &cobra.Command{
	Use:   "links",
	Short: "craw a website and fetch all links",
	Run: func(cmd *cobra.Command, args []string) {
		links := pkg.FetchLinks(linksConfig)
		// linkSlice := strings.Join(links, " ")
		if len(linksConfig.FilterIn) > 0 {
			filtered := make([]string, 0, len(links))
			for _, current := range links {
				if isStringIn(current, linksConfig.FilterIn...) {
					filtered = append(filtered, current)
				}
			}
			links = filtered
		}
		if len(linksConfig.FilterOut) > 0 {
			filtered := make([]string, 0, len(links))
			for _, current := range links {
				if !isStringIn(current, linksConfig.FilterOut...) {
					filtered = append(filtered, current)
				}
			}
			links = filtered
		}
		fmt.Println(strings.Join(links, " "))
	},
}

func isStringIn(item string, filters ...string) bool {
	for _, filter := range filters {
		if strings.Contains(item, filter) {
			return true
		}
	}
	return false
}

func init() {
	rootCmd.AddCommand(linksCmd)
	linksCmd.Flags().StringVarP(&linksConfig.Host, "address", "a", linksConfig.Host, "host to start crawling")
	linksCmd.Flags().BoolVarP(&linksConfig.ExternalLinks, "enable-external", "e", linksConfig.ExternalLinks, "crawl external links")
	linksCmd.Flags().BoolVarP(&linksConfig.RespectTree, "respect-tree", "t", linksConfig.RespectTree, "respect url hierarchy")
	linksCmd.Flags().StringSliceVarP(&linksConfig.FilterIn, "filter-in", "i", linksConfig.FilterIn, "url paths to include. If empty will include all")
	linksCmd.Flags().StringSliceVarP(&linksConfig.FilterOut, "filter-out", "o", linksConfig.FilterOut, "url paths to exclude. If empty will not exclude any")
}
