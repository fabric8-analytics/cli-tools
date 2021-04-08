package utils

import (
	"fmt"

	"github.com/manifoldco/promptui"
	"github.com/rs/zerolog/log"
)

func GetTelemetryConsent() bool {
	fmt.Println("CRDA CLI is constantly improving and we would like to know more about usage")
	prompt := promptui.Prompt{
		Label:       "Would you like to contribute anonymous usage statistics [y/n]",
		HideEntered: true,
	}
	userInput, err := prompt.Run()

	if err != nil {
		log.Fatal().Msgf(fmt.Sprintf("Prompt failed %v\n", err))
	}

	userResponse := Find(userInput)

	return userResponse
}

func Find(val string) bool {
	yes := []string{"y", "Y", "1"}
	no := []string{"n", "N", "0"}

	for _, item := range yes {
		if item == val {
			return true
		}
	}
	for _, item := range no {
		if item == val {
			return false
		}
	}
	return GetTelemetryConsent()
}
