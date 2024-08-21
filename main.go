package main

import (
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"time"

	"github.com/eiannone/keyboard"
)

type field struct {
	blocks [20][10]bool

	figures [19]string

	is_use   bool
	is_start bool

	cur_fig [3]int
	lines   int
}

func (f *field) moveLeft(blocks *[20][10]bool) {
	for j := 0; j < 4; j++ {
		for i := 0; i < 4; i++ {
			if f.figures[f.cur_fig[2]][j*4+i] == '@' && (f.cur_fig[1]+i == 0 || blocks[f.cur_fig[0]+j][f.cur_fig[1]+i-1]) {
				f.cur_fig[1]++
				break
			}
		}
	}
	f.cur_fig[1]--
}

func (f *field) moveRight(blocks *[20][10]bool) {
	for j := 3; j > -1; j-- {
		for i := 3; i > -1; i-- {
			if f.figures[f.cur_fig[2]][j*4+i] == '@' && (f.cur_fig[1]+i == 9 || blocks[f.cur_fig[0]+j][f.cur_fig[1]+i+1]) {
				f.cur_fig[1]--
				break
			}
		}
	}
	f.cur_fig[1]++
}

func (f *field) moveDown(blocks *[20][10]bool) {
	for j := 3; j > -1; j-- {
		for i := 3; i > -1; i-- {
			if f.figures[f.cur_fig[2]][j*4+i] == '@' && (f.cur_fig[0]+j == 19 || blocks[f.cur_fig[0]+j+1][f.cur_fig[1]+i]) {
				f.cur_fig[0]--
				break
			}
		}
	}

	f.cur_fig[0]++
}

func (f *field) checkForRotate() {
	for j := 0; j < 4; j++ {
		for i := 0; i < 4; i++ {
			if f.figures[f.cur_fig[2]][j*4+i] == '@' && f.cur_fig[1]+i < 0 {
				println(f.cur_fig[1] + i)
				f.cur_fig[1]++
				break
			}
		}
	}
	for j := 3; j > -1; j-- {
		for i := 3; i > -1; i-- {
			if f.figures[f.cur_fig[2]][j*4+i] == '@' && f.cur_fig[1]+i > 9 {
				println(f.cur_fig[1] + i)
				f.cur_fig[1]--
				break
			}
		}
	}
}

func (f *field) rotate() {
	f.cur_fig[2] = ([]int{0, 0, 2, 1, 4, 3, 6, 5, 8, 9, 10, 7, 12, 13, 14, 11, 16, 17, 18, 15})[f.cur_fig[2]+1]

	f.checkForRotate()
}

func (f *field) create() {
	f.figures[0] = ".....@@..@@....." // O

	f.figures[1] = "....@@@@........" // I
	f.figures[2] = "..@...@...@...@."

	f.figures[3] = "......@@.@@....." // S
	f.figures[4] = "..@...@@...@...."

	f.figures[5] = ".....@@...@@...." // Z
	f.figures[6] = "...@..@@..@....."

	f.figures[7] = ".....@@@.@......" // L
	f.figures[8] = "..@...@...@@...."
	f.figures[9] = "...@.@@@........"
	f.figures[10] = ".@@...@...@....."

	f.figures[11] = ".....@@@...@...." // J
	f.figures[12] = "..@@..@...@....."
	f.figures[13] = ".@...@@@........"
	f.figures[14] = "..@...@..@@....."

	f.figures[15] = ".....@@@..@....." // T
	f.figures[16] = "..@...@@..@....."
	f.figures[17] = "..@..@@@........"
	f.figures[18] = "..@..@@...@....."
}

func (f *field) checkLines() {
	var n []int
	for j, value := range f.blocks {
		for i := 0; i < 10; i++ {
			if !value[i] {
				break
			}
			if i == 9 {
				n = append(n, j)
			}
		}
	}
	fmt.Println(n)

	for _, value := range n {
		f.lines++
		for i := value; i > 0; i-- {
			f.blocks[i] = f.blocks[i-1]
		}
	}
}

func (f *field) move(turn string, disp *display) {
	blocks := f.blocks
	if !f.is_use {
		if f.blocks[0][5] {
			f.is_start = false
		}
		f.cur_fig = [3]int{0, 3, rand.Intn(19)}
		f.is_use = true
	} else {
		f.moveDown(&blocks)
	}
	if turn == "a" {
		f.moveLeft(&blocks)
	} else if turn == "d" {
		f.moveRight(&blocks)
	} else if turn == "w" {
		f.rotate()
	} else if turn == "s" {
		f.moveDown(&blocks)
	}
	for j := 3; j > -1; j-- {
		for i := 3; i > -1; i-- {
			if f.figures[f.cur_fig[2]][j*4+i] == '@' {
				if f.cur_fig[0]+j < 19 && !f.blocks[f.cur_fig[0]+j+1][f.cur_fig[1]+i] {
					blocks[f.cur_fig[0]+j][f.cur_fig[1]+i] = true
				} else {
					for j := 3; j > -1; j-- {
						for i := 3; i > -1; i-- {
							if f.figures[f.cur_fig[2]][j*4+i] == '@' && f.cur_fig[0]+j < 20 {
								f.blocks[f.cur_fig[0]+j][f.cur_fig[1]+i] = true
								blocks[f.cur_fig[0]+j][f.cur_fig[1]+i] = true
							}
						}
					}
					f.is_use = false
				}
			}
		}
	}

	f.checkLines()
	disp.draw(&blocks, f.lines)
}

type display struct {
	width, height int
}

func (disp *display) draw(f *[20][10]bool, lines int) {
	var screen string

	for i, j := 0, 0; j < disp.height; i++ {
		if i == 0 {
			screen += "<!"
			i++
		} else if i == 22 {
			screen += "!>"
			i++
		} else if i < 22 {
			if f[j][(i-2)/2] {
				screen += "[]"
			} else {
				screen += " ."
			}
			i++
		} else if j == 0 && i == 25 {
			screen += fmt.Sprintf("Линий уничтоженно: %d", lines)
		} else if j == 2 && i == 25 {
			screen += "a, s, d - управление, w - поворот, q - выход"
		}
		// конец строки
		if i == disp.width {
			screen += "\n"
			i = -1
			j++
		}
	}
	screen += "<!====================!>"
	fmt.Println(screen)
}

func reader(fld *field, disp *display) {
	for {
		char, _, err := keyboard.GetSingleKey()
		if err != nil {
			panic(err)
		}
		if char == 'q' {
			os.Exit(3)
		}
		fmt.Printf("You pressed: %q\r\n", char)
		cmd := exec.Command("clear")
		cmd.Stdout = os.Stdout
		cmd.Run()
		fld.move(string(char), disp)
	}
}

func main() {
	disp := display{50, 20}
	fld := field{}
	fld.create()
	rand.Seed(time.Now().Unix())

	fld.is_start = true
	go reader(&fld, &disp)

	for fld.is_start {

		time.Sleep(time.Second)
		cmd := exec.Command("clear")
		cmd.Stdout = os.Stdout
		cmd.Run()
		fld.move("56", &disp)
	}

	fmt.Println("Конец!")
}
