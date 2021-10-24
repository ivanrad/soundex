package soundex

import "testing"

type testpair struct {
	input, soundex string
}

var soundexPairs []testpair = []testpair{
	{"a", "A000"},
	{"A", "A000"},
	{"AA", "A000"},
	{"aa", "A000"},
	{"AAA", "A000"},
	{"aaa", "A000"},
	{"AAAA", "A000"},
	{"aaaa", "A000"},
	// Wikipedia
	{"Ashcroft", "A261"},
	{"Ashcraft", "A261"},
	{"Honeyman", "H555"},
	{"Pfister", "P236"},
	{"Robert", "R163"},
	{"Rupert", "R163"},
	{"Rubin", "R150"},
	{"Tymczak", "T522"},

	{"Aardvark", "A631"},
	{"Ahead", "A300"},
	{"Craft", "C613"},
	{"Checker", "C260"},
	{"Flavor", "F416"},
	{"Flavour", "F416"},
	{"Herbert", "H616"},
	{"John", "J500"},
	{"Philip", "P410"},
	{"Phillip", "P410"},
	{"Smith", "S530"},
	{"smith", "S530"},
	{"sMITH", "S530"},
	{"SMITH", "S530"},
	{"Smythe", "S530"},

	{"Z", "Z000"},
	{"ZZ", "Z000"},
	{"ZZZ", "Z000"},
	{"ZZZZ", "Z000"},
	{"R", "R000"},
	{"RH", "R000"},
	{"RHR", "R000"},
	{"RHRHR", "R000"},
	{"RHRHRHARHHHR", "R600"},
	{"RHRHRHARHHHAR", "R660"},
	{"RHRHRHARHHWWR", "R600"},
	{"RHRHRHARHHWWRAAR", "R660"},
	{"RW", "R000"},
	{"RWR", "R000"},
	{"RWRWR", "R000"},
	{"RAWR", "R600"},
	{"RAWARAWR", "R660"},
	{"H erbert", "H000"},
	{"heR bert", "H600"},
	{"HerB ert", "H610"},
	{"", "0000"},
	{" ", "0000"},
	{"  ", "0000"},
}

var soundexAltPairs = []testpair{
	{"a", "A000"},
	{"A", "A000"},
	{"AA", "A000"},
	{"aa", "A000"},
	{"AAA", "A000"},
	{"aaa", "A000"},
	{"AAAA", "A000"},
	{"aaaa", "A000"},

	{"Ashcroft", "A226"},
	{"Ashcraft", "A226"},
	{"Honeyman", "H555"},
	{"Pfister", "P236"},
	{"Robert", "R163"},
	{"Rupert", "R163"},
	{"Rubin", "R150"},
	{"Tymczak", "T522"},

	{"Aardvark", "A631"},
	{"Ahead", "A300"},
	{"Craft", "C613"},
	{"Checker", "C260"},
	{"Flavor", "F416"},
	{"Flavour", "F416"},
	{"Herbert", "H616"},
	{"John", "J500"},
	{"Philip", "P410"},
	{"Phillip", "P410"},
	{"Smith", "S530"},
	{"smith", "S530"},
	{"sMITH", "S530"},
	{"SMITH", "S530"},
	{"Smythe", "S530"},

	{"Z", "Z000"},
	{"ZZ", "Z000"},
	{"ZZZ", "Z000"},
	{"ZZZZ", "Z000"},
	{"R", "R000"},
	{"RH", "R000"},
	{"RHR", "R600"},
	{"RHRHR", "R660"},
	{"RHRHRHARHHHR", "R666"},
	{"RHRHRHARHHHAR", "R666"},
	{"RHRHRHARHHWWR", "R666"},
	{"RHRHRHARHHWWRAAR", "R666"},
	{"RW", "R000"},
	{"RWR", "R600"},
	{"RWRWR", "R660"},
	{"RAWR", "R600"},
	{"RAWARAWR", "R660"},
	{"H erbert", "H616"},
	{"heR bert", "H616"},
	{"HerB ert", "H616"},
	{"", ""},
	{" ", ""},
	{"  ", ""},
}

func toSoundexDigitHelper(c byte) byte {
	switch c {
	case 'B', 'F', 'P', 'V':
		return '1'
	case 'C', 'G', 'J', 'K', 'Q', 'S', 'X', 'Z':
		return '2'
	case 'D', 'T':
		return '3'
	case 'L':
		return '4'
	case 'M', 'N':
		return '5'
	case 'R':
		return '6'
	default: //'A','E','I','O','U','Y','H','W'
		return '0'
	}
}
func TestLookupTable(t *testing.T) {
	n := len(charToSoundex)
	if n != 26 {
		t.Fatalf("len(charToSoundex) = %v; want 26", n)
	}
	for i := 'A'; i <= 'Z'; i++ {
		want := toSoundexDigitHelper(byte(i))
		got := charToSoundex[i-'A']
		if got != want {
			t.Errorf("charToSoundex[%d] = %v; want %v", i, got, want)
		}
	}
}

func TestStdSoundexEncoding(t *testing.T) {
	for _, p := range soundexPairs {
		got := StdSoundex.Soundex(p.input)
		if got != p.soundex {
			t.Errorf("soundex(%q) = %v; wanted %v", p.input, got, p.soundex)
		}
	}
}

func TestAltSoundexEncoding(t *testing.T) {
	for _, p := range soundexAltPairs {
		got := AltSoundex.Soundex(p.input)
		if got != p.soundex {
			t.Errorf("SoundexAlt(%q) = %v; want %v", p.input, got, p.soundex)
		}
	}
}

func TestDifference(t *testing.T) {
	testCases := []struct {
		a, b          string
		diff, altdiff int
	}{
		// actually, pg's fuzzystrmatch diff("", "") == 3
		{"", "", 4, 4},
		{" ", "  ", 4, 4},
		{"cow", "cowper", 2, 2},
		{"a", "aaaa", 4, 4},
		{"Herbert", "H erbert", 1, 4},
		{"Philip", "Phillip", 4, 4},
		{"Rubin", "Robert", 2, 2},
		{"Smith", "Smythe", 4, 4},
		{"RW", "RWR", 4, 3},
		{"RWR", "RWAR", 3, 4},
	}

	for _, tc := range testCases {
		got := StdSoundex.Difference(tc.a, tc.b)
		if got != tc.diff {
			t.Errorf("Difference(%q, %q) = %v; want %v", tc.a, tc.b, got, tc.diff)
		}
		got = AltSoundex.Difference(tc.a, tc.b)
		if got != tc.altdiff {
			t.Errorf("DifferenceAlt(%q, %q) = %v; want %v", tc.a, tc.b, got, tc.diff)
		}
	}
}
