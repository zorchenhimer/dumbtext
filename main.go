package main

import (
	"fmt"
	"os"
	"math/rand"
	"strings"
	"time"
)

func main() {
	if len(os.Args) == 1 {
		return
	}

	// keep track of last two cases.  We don't want long sequences of one case.
	last1, last2 := false, false

	rng := rand.New(rand.NewSource(time.Now().Unix()))
	words := []string{}
	for i, a := range os.Args [1:] {
		a = strings.ToLower(a)
		word := []string{}

		for _, c := range a {
			cs := string(c)
			this := false

			// Don't bother with non-alpha characters
			if !((c >= 'A' && c <= 'Z') || (c >= 'a' && c <= 'z')) {
				word = append(word, cs)
				continue
			}

			// if last two were lower, upper next (and vice versa), execept if
			// first char.
			if i > 0 && last1 == last2 {
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
			word = append(word, cs)
		}
		words = append(words, strings.Join(word, ""))
	}
	fmt.Println(strings.Join(words, " "))
}
