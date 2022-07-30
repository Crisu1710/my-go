package myHelper

import (
	"github.com/AlecAivazis/survey/v2"
)

func Selector(label string, opts []string) string {
	var res string
	prompt := &survey.Select{
		Message: label,
		Options: opts,
	}
	survey.AskOne(prompt, &res)
	return res
}
