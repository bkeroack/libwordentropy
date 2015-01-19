package wordentropy

import (
	"crypto/rand"
	"log"
	"math/big"
)

func random_range(max int64) int64 {
	max_big := *big.NewInt(max)
	n, err := rand.Int(rand.Reader, &max_big)
	if err != nil {
		log.Fatalf("ERROR: cannot get random integer!\n")
	}
	return n.Int64()
}

func random_choice(l []string) string {
	return l[random_range(int64(len(l)-1))]
}

func random_digit() string {
	return random_choice([]string{"0", "1", "2", "3", "4", "5", "6", "7", "8", "9"})
}
