// snek by andy

package main

import (
	"fmt"
	"time"

	"github.com/eiannone/keyboard"
)

func main() {
	var screen Screen
	done := make(chan bool, 1)
	screen.init()
	screen.drawScreen()

	keysEvents, err := keyboard.GetKeys(10)
	if err != nil {
		panic(err)
	}
	defer func() {
		_ = keyboard.Close()
	}()

	go func() {
		tickerTime := 350
		ticker := time.NewTicker(time.Duration(tickerTime) * time.Millisecond)
		for range ticker.C {
			err := screen.processSnek()
			if err != nil {
				fmt.Println(err.Error())
				done <- true
				return
			}
			if screen.snekMeals%5 == 0 && screen.snekMeals != 0 {
				if tickerTime-75 > 0 {
					tickerTime = tickerTime - 75
				}
				screen.snekMeals = 0
				ticker.Reset(time.Duration(tickerTime) * time.Millisecond)
			}
			screen.drawScreen()
		}
	}()

	for {
		select {

		case event := <-keysEvents:
			if event.Err != nil {
				panic(event.Err)
			}
			if event.Key == keyboard.KeyEsc {
				done <- true
				return
			}
			if event.Key == keyboard.KeyArrowUp {
				screen.snekDirection = SNEK_UP
			}
			if event.Key == keyboard.KeyArrowDown {
				screen.snekDirection = SNEK_DOWN
			}
			if event.Key == keyboard.KeyArrowLeft {
				screen.snekDirection = SNEK_LEFT
			}
			if event.Key == keyboard.KeyArrowRight {
				screen.snekDirection = SNEK_RIGHT
			}
		case <-done:
			return
		}
	}

}
