libwordentropy - Random pseudo-grammatical passphrase generator
===============================================================

https://www.wordentropy.org

[API Documentation](http://godoc.org/github.com/bkeroack/libwordentropy)

Installation:

```bash
$ go get github.com/bkeroack/libwordentropy
```

Usage:

```go
import (
	"log"
	"wordentropy"
)

g, err := wordentropy.LoadGenerator(&wordentropy.WordListOptions{
	Wordlist: "data/part-of-speech.txt",
})
if err != nil {
	log.Fatalf("Error loading wordlist: %v\n", err)
}

p, err := g.GeneratePassphrases(nil)  //default options

for i := range p{
	log.Printf("Passphrase: %v\n", p[i])	
}
```

Speed:

The majority of execution overhead is in loading and parsing the wordlist from disk (done by ``LoadGenerator()``)--in the range of several hundred milliseconds. After loading the wordlist, passphrase generation is performed in memory and is very fast.

Using go test -bench on my Macbook with default passphrase settings, each call to ``GeneratePassphrases()`` completes in submillisecond time (in many cases less than 1/10 millisecond).

Command Line Tool:

```bash
$ cd libwordentropy/we
$ go build
$ ./we --help
Usage of ./we:
  -add_number=false: add random digit to passphrase (password requirement workaround)
  -add_symbol=false: add random symbol to passphrase (password requirement workaround)
  -count=1: number of passphrases to generate
  -length=4: number of words per passphrase
  -no_spaces=false: no spaces between words
  -offensive_path="../data/offensive.txt": path to offensive wordlist (optional)
  -prude=false: filter offensive words
  -verbose=false: verbose output
  -wordlist_path="../data/part-of-speech.txt": path to POS wordlist
```
