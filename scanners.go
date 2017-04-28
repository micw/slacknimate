package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
)

// StdinScanner will scan over STDIN, emitting a result on the returned channel
// every time it is able to successfully read a line.
//
// When it encounters an EOF, it will close the results channel.
func StdinScanner() chan string {
	ch := make(chan string)
	go func() {
		reader := bufio.NewScanner(os.Stdin)
		for reader.Scan() {
			ch <- reader.Text()
		}
		if err := reader.Err(); err != nil {
			fmt.Fprintln(os.Stderr, "reading standard input:", err)
			os.Exit(1)
		}
		close(ch)
	}()
	return ch
}

// LoopingStdinScanner will consume entire STDIN until EOF, and then
// continuously output on the results channel and never close.
//
// As a result, it is only suitable for input that will end, and will continue
// consuming memory while never sending anything if STDIN is a process that
// generates continuous output.
func LoopingStdinScanner() chan string {
	ch := make(chan string)
	go func() {
		var frames []string
		reader := bufio.NewScanner(os.Stdin)
		for reader.Scan() {
			frames = append(frames, reader.Text())
		}
		if err := reader.Err(); err != nil {
			fmt.Fprintln(os.Stderr, "reading standard input:", err)
			os.Exit(1)
		}

		for {
			for _, frame := range frames {
				ch <- frame
			}
		}
	}()
	return ch
}

type movie_frame struct {
    duration int
    frame string
}

func MovieScanner(loop_movie bool) chan string {
	ch := make(chan string)
	go func() {
		var frames []movie_frame
		reader := bufio.NewScanner(os.Stdin)
		var pos int = 0
		var dur int
		var err error
		var text string = ""
		for reader.Scan() {
			var line string = reader.Text()
			if (pos==0) {
				dur, err = strconv.Atoi(line)
				if err != nil {
					fmt.Fprintln(os.Stderr, "reading movie duration:", err)
					os.Exit(1)
				}
			} else {
				if (len(text)>0) {
					text+="\n"
				}
				text+=line
			}
			pos++;
			if (pos>13) {
				pos=0
				frames = append (frames, movie_frame{duration: dur,frame: text})
				text = ""
			}
		}
		if err := reader.Err(); err != nil {
			fmt.Fprintln(os.Stderr, "reading standard input:", err)
			os.Exit(1)
		}

		for {
			for _, frame := range frames {
				ch <- "```"+frame.frame+"```"
			}
			if ! loop_movie {
				break
			}
		}
	}()
	return ch
}
