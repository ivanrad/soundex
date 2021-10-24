// Package soundex implements Soundex and Difference functions to evaluate
// similarity of two strings.
package soundex

// A-Z to Soundex digit
const charToSoundex = "01230120022455012623010202"

const (
	stdMode = iota
	altMode
)

// Encoding represents an implementation of soundex phonetic algorithm.
type Encoding struct {
	mode int
}

func newEncoding(mode int) *Encoding {
	e := new(Encoding)
	e.mode = mode
	return e
}

// StdSoundex implements standard American soundex phonetic algorithm.
var StdSoundex = newEncoding(stdMode)

// AltSoundex offers an alternative implementation of soundex phonetic
// algorithm that *should* closely match PostgreSQL's fuzzystrmatch Soundex()
// and Difference() functions in most scenarios.
var AltSoundex = newEncoding(altMode)

func isAlpha(r rune) bool {
	// is it an ASCII alphabet char
	if r >= 'a' && r <= 'z' || r >= 'A' && r <= 'Z' {
		return true
	}
	return false
}

func toUpper(c byte) byte {
	if c >= 'a' && c <= 'z' {
		return c - 32
	}
	return c
}

func soundex(input string, mode int) string {
	sndex := [4]byte{}
	i := 0
	var prev byte
loop:
	for _, r := range input {
		if !isAlpha(r) {
			if mode == altMode {
				// skip over non-alpha in alt soundex mode
				continue
			}
			break
		}
		c := toUpper(byte(r))
		switch {
		case i == 0:
			sndex[i] = c
			i++
			prev = charToSoundex[c-'A']
		case i < 4:
			cur := charToSoundex[c-'A']
			if cur == '0' { // vowel-ish sound
				if mode == stdMode && (c == 'H' || c == 'W') { // skip if H,W
					continue loop
				}
				// for A,E,I,O,U,Y allow adjacent digits to be repeated
				prev = '0'
			} else if cur != prev { // are adjacent digits different
				sndex[i] = cur
				i++
				prev = cur
			}
		case i == 4:
			break loop
		}
	}
	if mode == altMode && i == 0 {
		return ""
	}
	// pad with zeros
	for i < 4 {
		sndex[i] = '0'
		i++
	}
	return string(sndex[:])
}

// Soundex returns the Soundex code for a given input string.
func (e *Encoding) Soundex(input string) string {
	return soundex(input, e.mode)
}

func difference(a, b string, mode int) int {
	if a == b {
		return 4
	}
	sa := []byte(soundex(a, mode))
	sb := []byte(soundex(b, mode))
	if len(sa) == 0 && len(sb) == 0 {
		return 4
	}
	d := 0
	for i := 0; i < len(sa) && i < len(sb); i++ {
		if sa[i] == sb[i] {
			d++
		}
	}
	return d
}

// Difference compares the soundex code for the two strings and returns an
// integer value from 0 to 4 indicating the degree of difference or similarity
// between the two.  A value of 0 indicates weak or no similarity, while 4
// indicates strong similarity.
func (e *Encoding) Difference(a, b string) int {
	return difference(a, b, e.mode)
}
