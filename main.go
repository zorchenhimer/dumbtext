package main

import (
	"fmt"
	"os"
	"io"
	"math/rand"
	"strings"
	"time"
	"bytes"

	"github.com/alexflint/go-arg"
)

type Options struct {
	Input []string `arg:"positional" placeholder:"text" help:"Text to transform"`
	InputFile string `arg:"-i,--input" placeholder:"input.txt" help:"File containing text to transform"`
}

func (o *Options) Description() string {
	return "Shuffle the case of text"
}

func main() {
	opts := &Options{}
	arg.MustParse(opts)

	err := run(opts)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}


func run(opts *Options) error {
	if opts.Input != nil {
		fmt.Println(convert(true, strings.Join(opts.Input, " ")))
		return nil
	}

	if opts.InputFile != "" {
		//raw, err := os.ReadFile(opts.InputFile)
		file, err := os.Open(opts.InputFile)
		if err != nil {
			return fmt.Errorf("Unable to read file: %w", err)
		}
		return streamConvert(file, os.Stdout)
	}

	return streamConvert(os.Stdin, os.Stdout)
}

func streamConvert(input io.Reader, output io.Writer) error {
	buf := &bytes.Buffer{}
	buf.Grow(2048)
	var err error
	var n int64
	first := true

	for err != io.EOF {
		n, err = io.Copy(buf, input)
		if err != nil && err != io.EOF {
			return fmt.Errorf("Error reading input: %w", err)
		}
		if n == 0 {
			return nil
		}

		fmt.Fprint(output, convert(first, string(buf.Bytes()[:buf.Len()])))
		buf.Reset()
		first = false
	}

	return nil
}

func convert(first bool, input string) string {
	// keep track of last two cases.  We don't want long sequences of one case.
	last1, last2 := false, false

	rng := rand.New(rand.NewSource(time.Now().Unix()))
	chars := []string{}

	for _, c := range input {
		cs := string(c)
		this := false

		// Don't bother with non-alpha characters
		if !((c >= 'A' && c <= 'Z') || (c >= 'a' && c <= 'z')) {
			//word = append(word, cs)
			chars = append(chars, cs)
			continue
		}

		// if last two were lower, upper next (and vice versa), execept if
		// first char.
		if !first && last1 == last2 {
			// in here so we get both upper and lowercase without dropping
			// to the else if below
			if !last1 {
				cs = strings.ToUpper(cs)
				this = true
			}
		} else if rng.Int() % 2 == 0 {
			cs = strings.ToUpper(cs)
			this = true
		}

		last2, last1 = last1, this
		//word = append(word, cs)
		chars = append(chars, cs)
		first = false
	}
	return strings.Join(chars, "")
}
