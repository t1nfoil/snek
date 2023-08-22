package main

import (
	"errors"
	"fmt"
	"math/rand"

	"golang.org/x/term"
)

type SnekKey int

const (
	SNEK_UP    SnekKey = 1
	SNEK_DOWN  SnekKey = 2
	SNEK_LEFT  SnekKey = 3
	SNEK_RIGHT SnekKey = 4
)

type Point struct {
	x, y      int
	character rune
}

type Screen struct {
	snekMeals     int
	buffer        []Point
	sizeX, sizeY  int
	snek          []Point
	snekDirection SnekKey
	screenFill    rune
}

func (s *Screen) getTermSize() error {
	var err error
	s.sizeX, s.sizeY, err = term.GetSize(0)
	if err != nil {
		return err
	}
	return nil
}

func (s *Screen) init() error {
	err := s.getTermSize()
	if err != nil {
		return err
	}

	s.snekMeals = 0
	s.screenFill = rune(' ')
	s.snek = make([]Point, 3)
	s.snek[0].x = s.sizeX/2 + 5
	s.snek[0].y = s.sizeY / 2
	s.snek[0].character = rune('@')
	s.snek[1].x = s.snek[0].x + 1
	s.snek[1].y = s.snek[0].y
	s.snek[1].character = rune('@')
	s.snek[2].x = s.snek[1].x + 1
	s.snek[2].y = s.snek[1].y
	s.snek[2].character = rune('@')

	s.buffer = make([]Point, s.sizeX*s.sizeY)

	for x := 0; x < s.sizeX; x++ {
		for y := 0; y < s.sizeY; y++ {
			r := rand.Intn(230)
			if r == 0 {
				s.setBufferPoint(x, y, rune('*'))
			} else if r == 1 {
				s.setBufferPoint(x, y, rune('o'))
			} else {
				s.setBufferPoint(x, y, s.screenFill)
			}
		}
	}

	return nil
}

func (s *Screen) setBufferPoint(x, y int, character rune) {
	index := x*s.sizeY + y
	s.buffer[index].x = x
	s.buffer[index].y = y
	s.buffer[index].character = character
}

func (s *Screen) getBufferPoint(x, y int) rune {
	index := x*s.sizeY + y
	return s.buffer[index].character
}

func (s *Screen) drawScreen() {
	fmt.Printf("\033[?25l")
	fmt.Printf("\033[0;0H")
	for x := 0; x < s.sizeX; x++ {
		for y := 0; y < s.sizeY; y++ {
			fmt.Printf("\033[%d;%dH%c", y, x, s.getBufferPoint(x, y))
		}
	}
	for _, snekPoint := range s.snek {
		fmt.Printf("\033[%d;%dH%c", snekPoint.y, snekPoint.x, snekPoint.character)
	}

}

func (s *Screen) processSnek() error {
	snekHeadIndex := len(s.snek) - 1
	snekHeadX := s.snek[snekHeadIndex].x
	snekHeadY := s.snek[snekHeadIndex].y

	if s.snekDirection == SNEK_UP {
		if snekHeadY-1 <= 0 || s.getBufferPoint(snekHeadX, snekHeadY-1) == rune('*') {
			return errors.New("yur snek is ded")
		}

		for i := 0; i < snekHeadIndex; i++ {
			if s.snek[i].x == snekHeadX && s.snek[i].y == snekHeadY-1 {
				return errors.New("yur snek is ded")
			}
		}

		if s.getBufferPoint(snekHeadX, snekHeadY-1) == rune('o') {
			s.snek = append(s.snek, Point{})
			s.snek[snekHeadIndex+1].x = snekHeadX
			s.snek[snekHeadIndex+1].y = snekHeadY - 1
			s.snek[snekHeadIndex+1].character = rune('@')
			s.setBufferPoint(snekHeadX, snekHeadY-1, s.screenFill)
			s.snekMeals++
			return nil
		}
		s.setBufferPoint(s.snek[0].x, s.snek[0].y, s.screenFill)

		s.snek = s.snek[1:]

		snekHead := Point{
			x:         snekHeadX,
			y:         snekHeadY - 1,
			character: rune('@'),
		}
		s.snek = append(s.snek, snekHead)
	}

	if s.snekDirection == SNEK_DOWN {
		if snekHeadY+1 >= s.sizeY || s.getBufferPoint(snekHeadX, snekHeadY+1) == rune('*') {
			return errors.New("yur snek is ded")
		}

		for i := 0; i < snekHeadIndex; i++ {
			if s.snek[i].x == snekHeadX && s.snek[i].y == snekHeadY+1 {
				return errors.New("yur snek is ded")
			}
		}

		if s.getBufferPoint(snekHeadX, snekHeadY+1) == rune('o') {
			s.snek = append(s.snek, Point{})
			s.snek[snekHeadIndex+1].x = snekHeadX
			s.snek[snekHeadIndex+1].y = snekHeadY + 1
			s.snek[snekHeadIndex+1].character = rune('@')
			s.setBufferPoint(snekHeadX, snekHeadY+1, s.screenFill)
			s.snekMeals++
			return nil
		}
		s.setBufferPoint(s.snek[0].x, s.snek[0].y, s.screenFill)
		s.snek = s.snek[1:]

		snekHead := Point{
			x:         snekHeadX,
			y:         snekHeadY + 1,
			character: rune('@'),
		}
		s.snek = append(s.snek, snekHead)
	}

	if s.snekDirection == SNEK_LEFT {
		if snekHeadX-1 <= 0 || s.getBufferPoint(snekHeadX-1, snekHeadY) == rune('*') {
			return errors.New("yur snek is ded")
		}

		for i := 0; i < snekHeadIndex; i++ {
			if s.snek[i].x == snekHeadX-1 && s.snek[i].y == snekHeadY {
				return errors.New("yur snek is ded")
			}
		}

		if s.getBufferPoint(snekHeadX-1, snekHeadY) == rune('o') {
			s.snek = append(s.snek, Point{})
			s.snek[snekHeadIndex+1].x = snekHeadX - 1
			s.snek[snekHeadIndex+1].y = snekHeadY
			s.snek[snekHeadIndex+1].character = rune('@')
			s.setBufferPoint(snekHeadX-1, snekHeadY, s.screenFill)
			s.snekMeals++
			return nil
		}
		s.setBufferPoint(s.snek[0].x, s.snek[0].y, s.screenFill)
		s.snek = s.snek[1:]

		snekHead := Point{
			x:         snekHeadX - 1,
			y:         snekHeadY,
			character: rune('@'),
		}
		s.snek = append(s.snek, snekHead)
	}

	if s.snekDirection == SNEK_RIGHT {
		if snekHeadX+1 >= s.sizeX || s.getBufferPoint(snekHeadX+1, snekHeadY) == rune('*') {
			return errors.New("yur snek is ded")
		}

		for i := 0; i < snekHeadIndex; i++ {
			if s.snek[i].x == snekHeadX+1 && s.snek[i].y == snekHeadY {
				return errors.New("yur snek is ded")
			}
		}

		if s.getBufferPoint(snekHeadX+1, snekHeadY) == rune('o') {
			s.snek = append(s.snek, Point{})
			s.snek[snekHeadIndex+1].x = snekHeadX + 1
			s.snek[snekHeadIndex+1].y = snekHeadY
			s.snek[snekHeadIndex+1].character = rune('@')
			s.setBufferPoint(snekHeadX+1, snekHeadY, s.screenFill)
			s.snekMeals++
			return nil
		}
		s.setBufferPoint(s.snek[0].x, s.snek[0].y, s.screenFill)

		s.snek = s.snek[1:]

		snekHead := Point{
			x:         snekHeadX + 1,
			y:         snekHeadY,
			character: rune('@'),
		}
		s.snek = append(s.snek, snekHead)
	}
	return nil
}
