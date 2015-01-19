libwordentropy - Random pseudo-grammatical passphrase generator
===============================================================

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
	Wordlist_path: "data/part-of-speech.txt",
})
if err != nil {
	log.Fatalf("Error loading wordlist: %v\n", err)
}

p, err := g.GeneratePassphrases(nil)

for i := range p{
	log.Printf("Passphrase: %v\n", p[i])	
}
```