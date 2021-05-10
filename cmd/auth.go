package cmd

import (
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/fabric8-analytics/cli-tools/pkg/telemetry"

	"github.com/fabric8-analytics/cli-tools/auth"
	"github.com/fatih/color"
	"github.com/manifoldco/promptui"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	snykURL = "https://app.snyk.io/redhat/snyk-token"
)

var snykToken string

type promptVars struct {
	Name        string
	Description string
}

// authCmd represents the auth command
var authCmd = &cobra.Command{
	Use:   "auth",
	Short: "Links uuid with Snyk token.",
	Long: fmt.Sprintf(`Command maps Snyk Token with UUID and Outputs 'crda_key' for further Authentication.
	
	To get "Snyk Token" Please click here: %s`, snykURL),
	RunE: runAuth,
}

func init() {
	rootCmd.AddCommand(authCmd)
	authCmd.Flags().StringVarP(&snykToken, "snyk-token", "t", "", "Authenticate with Snyk Token to unlock Verbose stack analyses.")
}

// runAuth is controller for Auth command.
func runAuth(cmd *cobra.Command, _ []string) error {
	log.Debug().Msgf("Executing Auth command.")
	askTelemetryConsent()
	var err error
	if snykToken == "" {
		snykToken, err = promptForToken()
		if err != nil {
			telemetry.SetExitCode(cmd.Context(), 1)
			return err
		}
	}
	client, err := telemetry.GetContextProperty(ctx, "client")
	if err != nil {
		return err
	}
	requestParams := auth.RequestServerType{
		UserID:          viper.GetString("crda_key"),
		SynkToken:       snykToken,
		ThreeScaleToken: viper.GetString("auth_token"),
		Host:            viper.GetString("host"),
		Client:          client,
	}
	userID, err := auth.RequestServer(cmd.Context(), requestParams)
	if err != nil {
		telemetry.SetExitCode(cmd.Context(), 1)
		return err
	}
	fmt.Print("Successfully Registered. \n\n")
	green := color.New(color.FgHiGreen, color.Bold).SprintFunc()
	fmt.Println(fmt.Sprintf(green("crda_key: ")+"%s\n", color.GreenString(userID)))
	fmt.Println("This key is confidential, Please keep it safe!")

	viper.Set("crda_key", userID)
	err = viper.WriteConfig()
	if err != nil {
		telemetry.SetExitCode(cmd.Context(), 1)
		return err
	}
	telemetry.SetExitCode(cmd.Context(), 0)
	log.Debug().Msgf("Successfully Executed Auth command.")
	return nil
}

func promptForToken() (string, error) {
	validate := func(input string) error {
		input = strings.TrimSpace(input)
		if input == "" {
			return nil
		}
		r := regexp.MustCompile("^[a-fA-F0-9]{8}-[a-fA-F0-9]{4}-4[a-fA-F0-9]{3}-[8|9aAbB][a-fA-F0-9]{3}-[a-fA-F0-9]{12}$")
		isValid := r.MatchString(input)
		if !isValid {
			return errors.New("invalid snyk token")
		}
		return nil
	}
	fmt.Println(fmt.Sprintf("To get Snyk Token, Please click %s\n", snykURL))
	promptValues := &promptVars{
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
		Label:       promptValues,
		Validate:    validate,
		Templates:   templates,
		Default:     "",
		HideEntered: true,
	}
	snykToken, err := prompt.Run()
	if err != nil {
		log.Error().Msgf("Unable to read Snyk Token. Please try again. %v", err)
		return "", err
	}
	return snykToken, nil
}
