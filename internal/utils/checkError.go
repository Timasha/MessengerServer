package utils

import "log"

func CheckErr(text string, err error) {
	if err != nil {
		log.Printf(text, err)
	}
}
