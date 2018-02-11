// Copyright © 2018 Chris Warner
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
	"log"
	"math/rand"
	"net/http"
	"time"

	"github.com/cwarner818/giota"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	addresses *[]string
	limit     *int
)

// confirmsCmd represents the confirms command
var confirmsCmd = &cobra.Command{
	Use:   "confirms",
	Short: "Returns the percentage of confirmed transactions sent to the specified address(es)",
	Long: `confirms takes a random sampling of all transactions sent to the specified
address(es) and then checks to see if they are confirmed according to 
the latest milestone.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Set the default http client to use
		httpClient := &http.Client{
			Timeout: viper.GetDuration("timeout"),
		}

		api := giota.NewAPI(viper.GetString("node"), httpClient)

		//fmt.Println("Using node:", viper.GetString("node"))
		nodeInfo, err := api.GetNodeInfo()
		if err != nil {
			log.Fatal(err)
		}

		milestone := nodeInfo.LatestMilestone

		addrTrytes, err := toAddress(*addresses)
		if err != nil {
			log.Fatal(err)
		}
		txnsResponse, err := api.FindTransactions(&giota.FindTransactionsRequest{
			Addresses: addrTrytes,
		})

		if err != nil {
			log.Fatal("Error getting transaction list: ", err)
		}
		fmt.Println(len(txnsResponse.Hashes), "results returned")

		if *limit > len(txnsResponse.Hashes) {
			*limit = len(txnsResponse.Hashes)
		}
		hashes := make([]giota.Trytes, *limit)

		for i := 0; i < *limit; i++ {
			r := rand.Intn(len(txnsResponse.Hashes))
			hashes[i] = txnsResponse.Hashes[r]
			txnsResponse.Hashes = append(txnsResponse.Hashes[:r],
				txnsResponse.Hashes[r+1:]...)
		}

		//pp.Print(txnsResponse.Hashes[:5])
		inclusionResponse, err := api.GetInclusionStates(txnsResponse.Hashes[:*limit],
			[]giota.Trytes{milestone})

		if err != nil {
			log.Fatal("Error getting confirmation status: ", err)
		}

		var counter int
		for _, v := range inclusionResponse.States {
			if v {
				counter++
			}
		}
		confirmed := float64(counter) / float64(len(inclusionResponse.States))
		fmt.Printf("Transaction confirmation: %0.2f%%\n", confirmed*100)
	},
}

func toTrytes(input []string) ([]giota.Trytes, error) {
	var output []giota.Trytes
	for _, t := range input {
		trytes, err := giota.ToTrytes(t)
		if err != nil {
			return nil, err
		}
		output = append(output, trytes)
	}
	return output, nil
}
func toAddress(input []string) ([]giota.Address, error) {
	var output []giota.Address
	for _, t := range input {
		trytes, err := giota.ToAddress(t)
		if err != nil {
			return nil, err
		}
		output = append(output, trytes)
	}
	return output, nil
}

func init() {
	RootCmd.AddCommand(confirmsCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// confirmsCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	limit = confirmsCmd.Flags().IntP("limit", "l", 500, "Limit checking to this number of transactions")

	addresses = confirmsCmd.Flags().StringSliceP("address", "a", nil, "Address to get confirmation information for")

	rand.Seed(time.Now().UnixNano())
}
