package soundex_test

import (
	"fmt"

	"github.com/ivanrad/soundex"
)

func ExampleStdSoundex_Soundex() {
	s := soundex.StdSoundex.Soundex("Ashcroft")
	fmt.Println(s)
	// Output:
	// A261
}

func ExampleAltSoundex_Soundex() {
	s := soundex.AltSoundex.Soundex("Ashcroft")
	fmt.Println(s)
	// Output:
	// A226
}

func ExampleStdSoundex_Difference() {
	d := soundex.StdSoundex.Difference("Cow", "Cowper")
	fmt.Println(d)
	// Output:
	// 2
}
