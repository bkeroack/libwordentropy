// Pseudo-grammatical English passphrase generation library
package wordentropy

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"
	"sync"
)

const (
	count_max        = 99
	count_default    = 4
	length_max       = 99
	length_default   = 5
	fragment_max     = 99
	fragment_default = 4
)

var grammar_rules = map[string][]string{ // word_type -> "can be followed by..."
	"snoun":        []string{"adverb", "verb", "pronoun", "conjunction"},
	"pnoun":        []string{"adverb", "verb", "pronoun", "conjunction"},
	"verb":         []string{"snoun", "pnoun", "preposition", "adjective", "conjunction", "sarticle", "particle"},
	"adjective":    []string{"snoun", "pnoun"},
	"adverb":       []string{"verb"},
	"preposition":  []string{"snoun", "pnoun", "adverb", "adjective", "verb"},
	"pronoun":      []string{"verb", "adverb", "conjunction"},
	"conjunction":  []string{"snoun", "pnoun", "pronoun", "verb", "sarticle", "particle"},
	"sarticle":     []string{"snoun", "adjective"},
	"particle":     []string{"pnoun", "adjective"},
	"interjection": []string{"snoun", "pnoun", "preposition", "adjective", "conjunction", "sarticle", "particle"},
}

var default_symbols = []string{"!", "@", "#", "$", "%", "^", "&", "*", "(", ")", "-", "+", "_", "="}
var word_types = []string{"snoun", "pnoun", "verb", "adjective", "adverb", "preposition", "pronoun", "conjunction", "sarticle", "particle", "interjection"}

// Options for loading word list. Wordlist is required, Offensive is optional.
// Wordlist must be formatted according to http://wordlist.aspell.net/pos-readme
// Offensive list must be ASCII/UTF8, one word per line
type WordListOptions struct {
	Wordlist  string // path to POS wordlist (required)
	Offensive string // "offensive" wordlist for optional filtering
}

// Load wordlist from disk and return a pointer to a Generator object.
func LoadGenerator(o *WordListOptions) (*Generator, error) {
	g := Generator{}
	err := g.LoadWords(o)
	if err != nil {
		return nil, err
	}
	return &g, nil
}

// Top-level Generator object
type Generator struct {
	word_map   map[string][]string
	offensive  map[string]uint
	options    *GenerateOptions
	sync.Mutex // Used only for loading/parsing word list
}

// Options for passphrase generation. All fields have sane defaults, none are required.
type GenerateOptions struct {
	Count                 uint     // Number of passphrases to generate
	Length                uint     // Length in words of each passphrase
	Magic_fragment_length uint     // Number of words per fragment
	Prudish               bool     // Filter out words in "offensive" wordlist
	No_spaces             bool     // Do not add spaces between words
	Add_digit             bool     // Add a random digit to the end of each passphrase
	Add_symbol            bool     // Add a random symbol to the end of each passphrase
	Symbols               []string // Slice of valid symbols to use with the Add_symbol option
}

func (g *Generator) random_word(word_type string, o *GenerateOptions) string {
	grw := func(words []string) (string, bool) {
		word := random_choice(words)
		_, ok := g.offensive[word]
		return word, ok
	}

	if words, ok := g.word_map[word_type]; ok {
		word, off := grw(words)
		if o.Prudish && off {
			log.Printf("Got offensive word: %v\n", word)
			i := 0
			for i = 0; off && i < 10; i++ {
				word, off = grw(words)
				if off {
					log.Printf("Got offensive word (retry): %v\n", word)
				}
			}
			if i >= 10 {
				log.Printf("Gave up trying to get non-offensive word!")
				word = ""
			}
		}
		return word
	} else {
		log.Printf("WARNING: random_word couldn't find word_type in word_map: %v\n", word_type)
		return "()"
	}
}

// A fragment is an autonomous run of words constructed using grammar rules
func (g *Generator) generate_fragment(o *GenerateOptions) []string {
	fragment_length := o.Magic_fragment_length
	fragment_slice := make([]string, fragment_length)
	prev_type_index := random_range(int64(len(word_types) - 1))       // Random initial word type
	fragment_slice[0] = g.random_word(word_types[prev_type_index], o) // Random initial word
	this_word_type := ""
	for i := uint(1); i < fragment_length; i++ {
		// Get random allowed word type by type of the previous word
		next_word_type_count := int32(len(grammar_rules[word_types[prev_type_index]]) - 1)
		if next_word_type_count > 0 { //rand.Int31n cannot take zero as a param
			this_word_type = grammar_rules[word_types[prev_type_index]][random_range(int64(next_word_type_count))]
		} else {
			this_word_type = grammar_rules[word_types[prev_type_index]][0]
		}
		fragment_slice[i] = g.random_word(this_word_type, o) //Random word of the allowed random type
		for j, v := range word_types {                       // Update previous word type with current word type for next iteration
			if v == this_word_type {
				prev_type_index = int64(j)
			}
		}
	}
	return fragment_slice
}

func (g *Generator) generate_passphrase(o *GenerateOptions) []string {
	iterations := o.Length / o.Magic_fragment_length
	phrase_slice := make([]string, 1)

	phrase_slice = append(phrase_slice, g.generate_fragment(o)...)
	if iterations >= 1 {
		for i := uint(1); i <= iterations; i++ {
			phrase_slice = append(phrase_slice, g.random_word("conjunction", o))
			phrase_slice = append(phrase_slice, g.generate_fragment(o)...)
		}
	}
	return phrase_slice
}

// Load and parse word list into memory.
func (g *Generator) LoadWords(o *WordListOptions) error {
	var err error

	g.Lock()
	defer g.Unlock()

	if o.Wordlist != "" {
		g.word_map, err = load_wordmap(o.Wordlist)
		if err != nil {
			return err
		}
	} else {
		return errors.New("Wordlist path is required")
	}

	if o.Offensive != "" {
		g.offensive, err = load_offensive_words(o.Offensive)
		if err != nil {
			return err
		}
	}

	return nil
}

func (g *Generator) check_options(o *GenerateOptions) error {
	if o == nil {
		o = &GenerateOptions{}
	}
	if len(g.word_map) == 0 {
		return fmt.Errorf("Empty wordlist, call LoadWords() first")
	}
	if o.Count > count_max {
		return fmt.Errorf("Count exceeds max: %v", count_max)
	}
	if o.Count == 0 {
		o.Count = count_default
	}
	if o.Length > length_max {
		return fmt.Errorf("Length exceeds max: %v", length_max)
	}
	if o.Length == 0 {
		o.Length = length_default
	}
	if o.Magic_fragment_length > fragment_max {
		return fmt.Errorf("Fragment length exceeds max: %v", fragment_max)
	}
	if o.Magic_fragment_length == 0 {
		o.Magic_fragment_length = fragment_default
	}
	if len(o.Symbols) == 0 {
		o.Symbols = default_symbols
	}
	return nil
}

// Generate and return passphrases according to options provided.
func (g *Generator) GeneratePassphrases(options *GenerateOptions) ([]string, error) {
	// Generate count passphrase slices
	// Merge each passphrase slice into a single string
	// Split string by spaces (individual random "words" can actually be multiword phrases)
	// Truncate slice to length words
	// Merge truncated slice back into string
	// Return slice of strings (final random passphrases)

	err := g.check_options(options)
	if err != nil {
		return nil, err
	}
	passphrases := make([]string, options.Count)

	var sep string
	if options.No_spaces {
		sep = ""
	} else {
		sep = " "
	}
	for i := uint(0); i < options.Count; i++ {
		ps := g.generate_passphrase(options)
		pj := strings.Join(ps, " ")
		ps = strings.Split(pj, " ")
		ps = ps[:options.Length+1]
		pp := strings.TrimSpace(strings.Join(ps, sep))
		if options.Add_digit {
			pp += random_digit()
		}
		if options.Add_symbol {
			pp += random_choice(options.Symbols)
		}
		passphrases[i] = pp
	}
	return passphrases, nil
}

func load_offensive_words(p string) (map[string]uint, error) {
	offensive := make(map[string]uint)

	f, err := os.Open(p)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		l := scanner.Text()
		offensive[strings.TrimSpace(l)] = 1
	}
	return offensive, nil
}

//Load word list into a mapping of word type to words of that type
func load_wordmap(p string) (map[string][]string, error) {

	word_map := map[string][]string{
		"snoun":        []string{},
		"pnoun":        []string{},
		"verb":         []string{},
		"adjective":    []string{},
		"adverb":       []string{},
		"preposition":  []string{},
		"pronoun":      []string{},
		"conjunction":  []string{},
		"sarticle":     []string{},
		"particle":     []string{},
		"interjection": []string{},
	}

	file, err := os.Open(p)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		word_type := ""
		plural := false
		line := scanner.Text()
		line_array := strings.Split(line, "\t")
		if len(line_array) != 2 {
			log.Printf("Bad string array length: %v, string: %v", len(line_array), line)
			continue
		}
		word := line_array[0]
		pos_tag := line_array[1]
		if strings.Contains(pos_tag, "N") || strings.Contains(pos_tag, "D") || strings.Contains(pos_tag, "I") {
			if strings.Contains(pos_tag, "P") {
				plural = true
			}
		}
		if strings.Contains(pos_tag, "D") || strings.Contains(pos_tag, "I") {
			if plural {
				word_type = "particle"
			} else {
				word_type = "sarticle"
			}
		} else if strings.Contains(pos_tag, "N") || strings.Contains(pos_tag, "h") || strings.Contains(pos_tag, "o") {
			if plural {
				word_type = "pnoun"
			} else {
				word_type = "snoun"
			}
		} else if strings.Contains(pos_tag, "V") || strings.Contains(pos_tag, "t") || strings.Contains(pos_tag, "i") {
			word_type = "verb"
		} else if strings.Contains(pos_tag, "A") {
			word_type = "adjective"
		} else if strings.Contains(pos_tag, "v") {
			word_type = "adverb"
		} else if strings.Contains(pos_tag, "C") {
			word_type = "conjunction"
		} else if strings.Contains(pos_tag, "p") || strings.Contains(pos_tag, "P") {
			word_type = "preposition"
		} else if strings.Contains(pos_tag, "r") {
			word_type = "pronoun"
		} else if strings.Contains(pos_tag, "!") {
			word_type = "interjection"
		} else {
			log.Printf("Unknown word type! word: %v; pos: %v\n", word, pos_tag)
			continue
		}
		if len(word) > 0 {
			word_map[word_type] = append(word_map[word_type], word)
		} else {
			log.Printf("WARNING: got zero length word: line: %v (interpreted type: %v)", line, word_type)
		}

	}

	return word_map, nil
}

// Get parsed wordlist as map of word type to words of that type
func (g *Generator) GetWordMap() map[string][]string {
	return g.word_map
}
