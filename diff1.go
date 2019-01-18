/*******************************************************************************
/
/      filename:  diff1.go
/
/   description:  Implements the UNIX env utility.
/
/        author:  Schwartz, Jacob
/      login id:  SP_19_CPS444_09
/
/         class:  CPS 444
/    instructor:  Perugini
/    assignment:  Homework #1
/
/      assigned:  January 14, 2018
/           due:  January 23, 2018
/
/******************************************************************************/

package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"
)

var a_flag, l_flag, m_flag, t_flag *bool
var files [2]*os.File
var readers [2]*bufio.Reader
var lines [2]string
var done [2]error

func remaining(e error, file int, line int) {
	for e == nil {
		line += 1
		fmt.Println(line)
		_, e = readers[file].Peek(1)
	}
}

func getStdIn() (result string) {
	input := bufio.NewReader(os.Stdin)

	for {
		line, err := input.ReadString('\n')

		if err == io.EOF {
			return
		} else {
			result += line
		}
	}
}

func main() {
	a_flag = flag.Bool("a", false, "All")
	l_flag = flag.Bool("l", false, "Leading")
	m_flag = flag.Bool("m", false, "Middle")
	t_flag = flag.Bool("t", false, "Trailing")
	flag.Parse()

	if *a_flag && (*l_flag || *m_flag || *t_flag) {
		_, _ = fmt.Fprintln(os.Stderr, "Option -a cannot be used in combination with other options.")
		os.Exit(9)
	} else if len(flag.Args()) > 2 {
		_, _ = fmt.Fprintf(os.Stderr, "diff: extra operand(s) '%s'\n", flag.Args()[2])
		os.Exit(2)
	} else if len(flag.Args()) < 2 || (flag.Args()[0] == "-" && flag.Args()[1] == "-") {
		os.Exit(0)
	}

	for index, element := range flag.Args() {
		if element == "-" {
			//fmt.Println(element)
			_ = ioutil.WriteFile("tmp.txt", []byte(getStdIn()), 0644)
			file, _ := os.Open("tmp.txt")
			files[index] = file
			readers[index] = bufio.NewReader(file)
		} else if _, err := os.Stat(element); os.IsNotExist(err) {
			_, _ = fmt.Fprintf(os.Stderr, "diff: file %s does not exist\n", element)
			os.Exit(2)
		} else {
			//fmt.Println(element)
			file, _ := os.Open(element)
			files[index] = file
			readers[index] = bufio.NewReader(file)
		}
	}

	index := 0
	cutset := " \t\r\n"

	for {
		index += 1
		lines[0], _ = readers[0].ReadString('\n')
		lines[1], _ = readers[1].ReadString('\n')

		empty := len(strings.Join(strings.Fields(lines[0]), "")) == 0 || len(strings.Join(strings.Fields(lines[1]), "")) == 0

		if *a_flag && !empty {
			lines[0] = strings.Join(strings.Fields(lines[0]), "")
			lines[1] = strings.Join(strings.Fields(lines[1]), "")
		}

		if *l_flag && !empty {
			lines[0] = strings.TrimLeft(lines[0], cutset)
			lines[1] = strings.TrimLeft(lines[1], cutset)
		}

		if *m_flag && !empty {
			line0 := lines[0]
			line1 := lines[1]
			line0_prime := strings.Join(strings.Fields(line0), "")
			line1_prime := strings.Join(strings.Fields(line1), "")
			line0_index := len(strings.TrimRight(lines[0], cutset))
			line1_index := len(strings.TrimRight(lines[1], cutset))
			lines[0] = lines[0][0:strings.Index(lines[0], line0_prime[0:1])] + line0_prime + lines[0][line0_index:]
			lines[1] = lines[1][0:strings.Index(lines[1], line1_prime[0:1])] + line1_prime + lines[1][line1_index:]
		}

		if *t_flag && !empty {
			lines[0] = strings.TrimRight(lines[0], cutset)
			lines[1] = strings.TrimRight(lines[1], cutset)
		}

		if lines[0] != lines[1] {
			fmt.Println(index)
		}

		_, done[0] = readers[0].Peek(1)
		_, done[1] = readers[1].Peek(1)

		if done[0] != nil && done[1] != nil {
			break
		}
	}

	if done[0] != nil {
		remaining(done[0], 0, index)
	} else {
		remaining(done[1], 1, index)
	}

	_ = files[0].Close()
	_ = files[0].Close()

}
