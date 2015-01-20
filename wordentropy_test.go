package wordentropy

import (
	"testing"
)

func TestPassphrases(t *testing.T) {

	var ops GenerateOptions

	g, err := LoadGenerator(&WordListOptions{
		Wordlist:  "data/part-of-speech.txt",
		Offensive: "data/offensive.txt",
	})
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

func BenchmarkPassphraseGeneration(b *testing.B) {

	g, err := LoadGenerator(&WordListOptions{
		Wordlist:  "data/part-of-speech.txt",
		Offensive: "data/offensive.txt",
	})
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

func BenchmarkWordlistLoading(b *testing.B) {
	wo := WordListOptions{
		Wordlist:  "data/part-of-speech.txt",
		Offensive: "data/offensive.txt",
	}
	for i := 0; i < b.N; i++ {
		_, err := LoadGenerator(&wo)
		if err != nil {
			b.Fatalf("Error loading wordlist: %v\n", err)
		}
	}
}
