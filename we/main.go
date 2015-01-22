package main

import (
	"flag"
	"fmt"
	"github.com/bkeroack/libwordentropy"
	"log"
	"os"
)

var count = flag.Int("count", 1, "number of passphrases to generate")
var length = flag.Int("length", 4, "number of words per passphrase")
var prude = flag.Bool("prude", false, "filter offensive words")
var no_spaces = flag.Bool("no_spaces", false, "no spaces between words")
var add_number = flag.Bool("add_number", false, "add random digit to passphrase (password requirement workaround)")
var add_symbol = flag.Bool("add_symbol", false, "add random symbol to passphrase (password requirement workaround)")
var wordlist_path = flag.String("wordlist_path", "../data/part-of-speech.txt", "path to POS wordlist")
var offensive_path = flag.String("offensive_path", "../data/offensive.txt", "path to offensive wordlist (optional)")
var verbose = flag.Bool("verbose", false, "verbose output")

func msg(m string) {
	if *verbose {
		log.Printf(m)
	}
}

func init() {
	flag.Parse()
	if *count < 1 || *count > 99 {
		log.Fatalf("invalid count: %v\n", *count)
	}
	if *length < 1 || *length > 99 {
		log.Fatalf("invalid length: %v\n", *length)
	}
	if _, err := os.Stat(*wordlist_path); err != nil {
		log.Fatalf("wordlist error: %v\n", err)
	}
	if _, err := os.Stat(*offensive_path); err != nil {
		msg(fmt.Sprintf("warning: offensive path error: %v\n", err))
		*prude = false
		*offensive_path = ""
	}
}

func main() {
	msg("loading word list...\n")
	wo := wordentropy.WordListOptions{
		Wordlist: *wordlist_path,
	}
	if *prude {
		wo.Offensive = *offensive_path
	}
	g, err := wordentropy.LoadGenerator(&wo)
	if err != nil {
		log.Fatalf("error loading wordlist: %v\n", err)
	}

	o := wordentropy.GenerateOptions{
		Count:      uint(*count),
		Length:     uint(*length),
		Prudish:    *prude,
		No_spaces:  *no_spaces,
		Add_digit:  *add_number,
		Add_symbol: *add_symbol,
	}

	msg(fmt.Sprintf("options: %v\n", o))

	p, err := g.GeneratePassphrases(&o)
	if err != nil {
		log.Fatalf("error generating passphrases: %v\n", err)
	}

	msg("passphrases:\n")
	for i := range p {
		fmt.Printf("%v\n", p[i])
	}
}
