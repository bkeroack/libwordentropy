package wordentropy

import (
	"testing"
)

func TestPassphrases(t *testing.T) {

	var ops GenerateOptions

	g := Generator{}
	word_ops := WordMapOptions{
		Wordlist_path:  "data/part-of-speech.txt",
		Offensive_path: "data/offensive.txt",
	}

	err := g.LoadWords(&word_ops)
	if err != nil {
		t.Fatalf("Could not load wordlist: %v", err)
	}

	for i := 0; i < 20; i++ {
		ops.Length = uint(random_range(int64(20)))
		ops.Count = uint(random_range(int64(20)))
		_, err := g.GeneratePassphrases(&ops)
		if err != nil {
			t.Fatalf("Error generating passphrases (i: %v): %v", i, err)
		}
	}
}

func BenchmarkPassphrases(b *testing.B) {
	g := Generator{}
	word_ops := WordMapOptions{
		Wordlist_path:  "data/part-of-speech.txt",
		Offensive_path: "data/offensive.txt",
	}

	err := g.LoadWords(&word_ops)
	if err != nil {
		b.Fatalf("Could not load wordlist: %v", err)
	}
	for i := 0; i < b.N; i++ {
		_, err := g.GeneratePassphrases(&GenerateOptions{})
		if err != nil {
			b.Fatalf("Error generating passphrases (i: %v): %v", i, err)
		}
	}
}
