// +build all

package recovery

import (
	"testing"
)

func BenchmarkCreate(b *testing.B) {
	secret := "shoulder artefact abstract position deny example shoulder myth orchard wolf jewel coconut tourist wrong cram"

	for i := 0; i < b.N; i++ {
		_, err := Create(2, 3, secret)
		if err != nil {
			b.Fatal(err)
		}
	}
}
