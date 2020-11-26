/*
Copyright © 2020 Deepak Sharma <deepshar@redhat.com>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	guid "github.com/google/uuid"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// configCmd represents the config command
var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Links uuid with Snyk token.",
	Long: `config registers Snyk Token 
	with UUID on crda server and outputs UUID.`,
	Args: validateArgs,
	Run:  main,
}

func init() {
	rootCmd.AddCommand(configCmd)

}
func validateArgs(cmd *cobra.Command, args []string) error {
	if len(args) < 1 {
		return errors.New("snyk token is missing")
	}
	return nil
}

func requestServer(synkToken string) string {
	ThreeScale := viper.GetString("3scaleToken")
	url := viper.GetString("server") + "/user?user_key=" + string(ThreeScale)
	uuid := guid.New().String()
	payload, err := json.Marshal(map[string]string{
		"snyk_api_token": synkToken,
		"user_id":        uuid,
	})

	client := &http.Client{}
	req, err := http.NewRequest(http.MethodPut, url, bytes.NewBuffer(payload))

	if err != nil {
		fmt.Println(err)
		panic("Unable to build request")
	}
	req.Header.Add("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		panic("Unable to request server")
	}
	defer res.Body.Close()
	if int(res.StatusCode) != 200 {
		fmt.Println(err)
		panic("Server Errored. Please try again.")
	}
	return uuid
}

func main(cmd *cobra.Command, args []string) {
	uuid := requestServer(args[0])
	fmt.Println("Successfully Registered.\n")
	fmt.Println("Please update CI env with: \n")
	fmt.Println("user_key: ", uuid, "\n\n")
	fmt.Println("This token is confidential and exculsive to your Snyk Id.")
}
