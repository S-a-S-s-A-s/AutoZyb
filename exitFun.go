package main

import (
	"github.com/eiannone/keyboard"
	"log"
	"os"
)

func escExit(preURL, token string) {
	err := keyboard.Open()
	if err != nil {
		log.Fatal(err)
	}

	//fmt.Println("Press Esc to exit...")

	for {
		_, key, err := keyboard.GetKey()
		if err != nil {
			log.Fatal(err)
		}
		//如果是k键，就退出
		if key == keyboard.KeySpace || key == keyboard.KeyEsc {
			//fmt.Printf("Key pressed: %v\n", key)

			dequeueorder(preURL, token)
			keyboard.Close()
			os.Exit(0)
		}

		//fmt.Printf("Key pressed: %v\n", key)
	}
}
