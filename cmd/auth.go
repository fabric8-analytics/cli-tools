package cmd

import (
	"errors"
	"fmt"
	"os"
	"regexp"
	"strings"

	auth "github.com/fabric8-analytics/cli-tools/auth"
	"github.com/fatih/color"
	"github.com/manifoldco/promptui"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	snykURL = "https://app.snyk.io/redhat/snyk-token"
)

var showUUID bool
var snykToken string

type promtVars struct {
	Name        string
	Description string
}

// authCmd represents the auth command
var authCmd = &cobra.Command{
	Use:   "auth",
	Short: "Links uuid with Snyk token.",
	Long: fmt.Sprintf(`Command maps Snyk Token with UUID and Outputs 'crda-key' for further Authentication.
	
	To get "Snyk Token" Please click here: %s`, snykURL),
	Run: main,
}

func init() {
	rootCmd.AddCommand(authCmd)
	authCmd.Flags().StringVarP(&snykToken, "snyk-token", "t", "", "Authenticate with Snyk Token to unlock Verbose stack anaylses.")
}

func main(cmd *cobra.Command, args []string) {
	log.Debug().Msgf("Executing Auth command.")
	if snykToken == "" {
		snykToken = promptForToken()
	}
	requestParams := auth.RequestServerType{
		UserID:          viper.GetString("crda-key"),
		SynkToken:       snykToken,
		ThreeScaleToken: viper.GetString("auth-token"),
		Host:            viper.GetString("host"),
	}
	userID := auth.RequestServer(requestParams)

	fmt.Fprint(os.Stdout, "Successfully Registered. \n\n")
	green := color.New(color.FgHiGreen, color.Bold).SprintFunc()
	fmt.Fprintln(os.Stdout,
		fmt.Sprintf(green("\t crda-key: ")+"%s\n\n", color.GreenString(userID)),
		fmt.Sprintf("This token is confidential, Please keep it safe!"),
	)
	viper.Set("crda-key", userID)
	viper.WriteConfig()
	log.Debug().Msgf("Successfully Executed Auth command.")
}

func promptForToken() string {
	validate := func(input string) error {
		input = strings.TrimSpace(input)
		if input == "" {
			return nil
		}
		r := regexp.MustCompile("^[a-fA-F0-9]{8}-[a-fA-F0-9]{4}-4[a-fA-F0-9]{3}-[8|9|aA|bB][a-fA-F0-9]{3}-[a-fA-F0-9]{12}$")
		isValid := r.MatchString(input)
		if !isValid {
			return errors.New("Invalid Snyk Token")
		}
		return nil
	}
	fmt.Fprintln(os.Stdout,
		fmt.Sprintf("To get Snyk Token, Please click %s\n", snykURL),
	)
	promtValues := &promtVars{
		Name:        "Snyk Token",
		Description: "[Press Enter to continue]",
	}
	templates := &promptui.PromptTemplates{
		Prompt:  "{{ .Name }} {{ .Description | faint}}:  ({{ .HeatUnit | red | italic }}) ",
		Valid:   "{{ .Name | green|bold }} {{.Description | faint }}: ",
		Invalid: "{{ .Name | red }} {{ .Description | faint}}: ",
		Success: "{{ .Name | bold }} {{ .Description | faint}}",
	}
	prompt := promptui.Prompt{
		Label:       promtValues,
		Validate:    validate,
		Templates:   templates,
		Default:     "",
		HideEntered: true,
	}
	snykToken, err := prompt.Run()
	if err != nil {
		log.Fatal().Msgf("Unable to read Snyk Token. Please try again. %v\n", err)
	}
	return snykToken
}
