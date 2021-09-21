package main

import "log"

func checkError(err error) {
	if err != nil {
		log.Println("Error", err.Error())
	}
}
