package text

import (
	"errors"
	"fmt"
	"sort"
	"strings"
)

const Paul2013PronunciationFeatureCount = 15

// Paul2013PronunciationFeatures is a partially or fully produced classifier
// row. Available marks values established by recovered feature producers;
// callers can distinguish unported inputs from valid zero-valued features.
type Paul2013PronunciationFeatures struct {
	Values    [Paul2013PronunciationFeatureCount]int16
	Available uint16
}

// MissingPositions lists classifier positions without a recovered producer
// in this row.
func (features Paul2013PronunciationFeatures) MissingPositions() []int {
	missing := make([]int, 0, Paul2013PronunciationFeatureCount)
	for position := 0; position < Paul2013PronunciationFeatureCount; position++ {
		if features.Available&(1<<position) == 0 {
			missing = append(missing, position)
		}
	}
	return missing
}

// Complete reports whether every classifier feature has a recovered value.
func (features Paul2013PronunciationFeatures) Complete() bool {
	return features.Available == (1<<Paul2013PronunciationFeatureCount)-1
}

// paul2013PronunciationContextWords is the sorted 99-entry string table at
// PTR_s_about_100784c8. FUN_10006ca0 returns the matching index plus two.
var paul2013PronunciationContextWords = [...]string{
	"about", "above", "across", "after", "again", "against", "agin", "ago",
	"along", "although", "among", "amongst", "around", "as", "aside", "at",
	"away", "back", "because", "been", "before", "behind", "being", "below",
	"beneath", "beside", "besides", "between", "beyond", "but", "by", "de",
	"despite", "down", "during", "except", "far", "for", "from", "had",
	"has", "have", "in", "including", "inside", "inter", "into", "like",
	"made", "make", "makes", "many", "much", "near", "nearby", "need",
	"needed", "neither", "next", "of", "off", "on", "once", "only", "onto",
	"opposite", "other", "out", "outside", "over", "past", "per", "post",
	"quite", "rather", "round", "said", "say", "says", "since", "so", "sure",
	"than", "that", "through", "throughout", "till", "to", "toward", "towards",
	"under", "unless", "unlike", "until", "up", "upon", "via", "with", "without",
}

// Paul2013PronunciationWordClass ports FUN_10006ca0: empty input maps to 0,
// an unmatched non-empty word maps to 1, and one of the 99 table entries maps
// to its sorted ordinal plus two. These are opaque model features, not
// linguistic parts of speech.
func Paul2013PronunciationWordClass(surface string) int16 {
	if surface == "" {
		return 0
	}
	index := sort.SearchStrings(paul2013PronunciationContextWords[:], surface)
	if index == len(paul2013PronunciationContextWords) || paul2013PronunciationContextWords[index] != surface {
		return 1
	}
	return int16(index + 2)
}

// Paul2013PronunciationContinuityFlag ports the final feature store in
// FUN_10006ae0. The DLL sign-extends the center surface's first byte and tests
// bit 7 in DAT_1007e188 at that signed offset. In vt_pau.dll, the reachable
// lookup entries with bit 7 set are exactly ASCII 'A' through 'Z'; the
// preceding 128 bytes contain no set bit for high-bit input bytes.
func Paul2013PronunciationContinuityFlag(surface string) int16 {
	if len(surface) == 0 {
		return 0
	}
	if surface[0] >= 'A' && surface[0] <= 'Z' {
		return 1
	}
	return 0
}

// Paul2013PronunciationPhoneContextFeature ports the ordered
// case-insensitive word-shape checks in FUN_10007160. Its string helper uses
// the embedded character map, so the native offsets and minimum-length checks
// are kept alongside each comparison instead of reducing them to suffix
// rules.
func Paul2013PronunciationPhoneContextFeature(surface string) int16 {
	if end := strings.IndexByte(surface, 0); end >= 0 {
		surface = surface[:end]
	}
	if surface == "" {
		return 0
	}
	length := len(surface)
	characterMap := Paul2013EmbeddedKeyTables().CharacterMap
	suffixMatches := func(width int, pattern string) bool {
		if length < width {
			return false
		}
		tail := surface[length-width:]
		return len(tail) <= len(pattern) && Paul2013MappedCStringEqual(
			[]byte(tail), []byte(pattern[:len(tail)]), characterMap,
		)
	}
	classTwoSuffixes := [...]struct {
		minimumLength int
		width         int
		suffix        string
	}{
		{9, 6, "liness"},
		{7, 4, "ance"}, {7, 4, "tion"}, {8, 5, "tions"},
		{7, 4, "sion"}, {8, 5, "sions"}, {7, 4, "ness"},
		{9, 6, "nesses"}, {7, 4, "ment"}, {7, 5, "ments"}, {7, 4, "ship"},
		{8, 5, "ships"}, {7, 3, "ity"}, {7, 4, "ties"},
		{7, 4, "ings"},
	}
	for _, candidate := range classTwoSuffixes {
		if length >= candidate.minimumLength && suffixMatches(candidate.width, candidate.suffix) {
			return 2
		}
	}
	if paul2013FUN1000FFA0(surface) && paul2013FUN1000FFD0(surface) {
		return 2
	}
	if (length > 6 && (suffixMatches(4, "tial") || suffixMatches(4, "cial") || suffixMatches(4, "less"))) ||
		(length > 5 && (suffixMatches(3, "ful") || suffixMatches(3, "cal") || suffixMatches(3, "fic") || suffixMatches(3, "ian") || suffixMatches(3, "ous"))) ||
		(length > 6 && (suffixMatches(4, "tive") || suffixMatches(4, "sive"))) ||
		(length > 8 && suffixMatches(6, "tional")) {
		return 3
	}
	if length > 4 && suffixMatches(2, "ly") {
		return 4
	}
	if length > 5 && suffixMatches(3, "ing") {
		return 5
	}
	if length > 5 && suffixMatches(3, "ist") {
		return 6
	}
	if length > 6 && suffixMatches(4, "ists") {
		return 6
	}
	if length > 4 && suffixMatches(2, "ed") {
		return 7
	}
	if length > 4 && suffixMatches(1, "s") {
		return 8
	}
	return 1
}

// paul2013FUN1000FFA0 and paul2013FUN1000FFD0 preserve the two native
// character-table predicates called by FUN_10007160 after its class-two
// suffix checks. The effective table uses signed-byte indexing.
func paul2013FUN1000FFA0(surface string) bool {
	if surface == "" {
		return false
	}
	attributes := Paul2013ExceptionCharacterAttributes()
	for index := 0; index < len(surface); index++ {
		if attributes[surface[index]]&0xc0 == 0 {
			return false
		}
	}
	return true
}

func paul2013FUN1000FFD0(surface string) bool {
	if surface == "" {
		return false
	}
	attributes := Paul2013ExceptionCharacterAttributes()
	for index := 0; index < len(surface); index++ {
		attribute := attributes[surface[index]]
		if attribute&0xc0 != 0 && attribute&0x40 != 0 {
			return false
		}
	}
	return true
}

// These sorted string tables are passed to FUN_100560a0 by FUN_10006ce0.
// Their source pointers, entry counts, and return classes are preserved in the
// names/comments; the strings are stored inline here for independent runtime
// classification.
var paul2013TripletClass2Words = []string{
	"'d", "'ll", "can", "can't", "cannot", "canst", "could", "couldest", "couldn't", "dare",
	"did", "do", "does", "may", "mayest", "mayst", "might", "must", "mustn't", "shall",
	"shalt", "shan't", "should", "shouldest", "shouldn't", "shouldst", "will", "won't", "would", "wouldn't",
} // PTR_10078044, count DAT_100780c4 (30), result 2.

var paul2013TripletClass3Words = []string{
	"'m", "'re", "am", "are", "be", "been", "being", "is", "was", "were",
} // PTR_10078010, count DAT_100780bc (10), result 3.

var paul2013TripletClass4Words = []string{
	"&", "and", "because", "both", "but", "either", "neither", "nor", "or", "versus", "vs", "vs.",
} // PTR_1007844c, count DAT_1007847c (12), result 4.

var paul2013TripletClass12Words = []string{
	"a", "all", "an", "another", "any", "both", "each", "either", "every", "half", "no", "some", "such", "the", "these", "this", "those",
} // PTR_10078480, count DAT_100784c4 (17), result 12.

// Paul2013TripletClass12Contains reports membership in the sorted 17-entry
// class-12 table used by the Paul 2013 counted-row phrase-window rules.
func Paul2013TripletClass12Contains(surface []byte) bool {
	return paul2013TripletWordsContainMapped(paul2013TripletClass12Words, surface)
}

// Paul2013TripletTableContains checks a recovered sorted triplet-class table
// by its opaque numeric class identifier using the mapped-string comparator
// passed by native mode 0x49 call sites.
func Paul2013TripletTableContains(class int16, surface []byte) bool {
	for _, group := range paul2013TripletClassGroups {
		if group.class == class {
			return paul2013TripletWordsContainMapped(group.words, surface)
		}
	}
	return false
}

func paul2013TripletWordsContainMapped(words []string, surface []byte) bool {
	rows := make([][]byte, len(words))
	for index, word := range words {
		rows[index] = []byte(word)
	}
	return FindPaul2013SortedCString(rows, cString(surface), paul2013ContextCharacterWeights()) >= 0
}

var paul2013TripletClass6Words = []string{"n't", "never", "not"} // PTR_10078038, count DAT_100780c0 (3), result 6.

var paul2013TripletClass7Words = []string{
	"her", "hers", "his", "its", "mine", "my", "our", "ours", "their", "theirs", "your", "yours",
} // PTR_10078658, count DAT_10078688 (12), result 7.

var paul2013TripletClass8Words = []string{
	"'em", "'emselves", "'im", "he", "herself", "him", "himself", "it", "itself", "me", "myself", "one",
	"oneself", "ourselves", "she", "some", "them", "themselves", "they", "us", "we", "you", "yourself", "yourselves",
} // PTR_1007868c, count DAT_100786ec (24), result 8.

var paul2013TripletClass11Words = []string{
	"how", "that", "what", "whatever", "when", "whence", "whenever", "where", "whereby", "whereever", "which", "whichever", "who", "whoever", "whom", "whose", "whoseever", "why",
} // PTR_1007871c, count DAT_10078764 (18), result 11.

var paul2013TripletClass9Words = []string{
	"billion", "dozen", "eight", "eighteen", "eighty", "eleven", "fifteen", "fifties", "fifty", "first", "five", "fives", "forties", "forty", "four", "fourteen", "half", "hundred", "hundreds", "million", "millions", "nine", "nineteen", "ninety", "octillion", "one", "quadrillion", "quintillion", "second", "septillion", "seven", "seventeen", "seventy", "sextillion", "six", "sixteen", "sixties", "sixty", "ten", "third", "thirteen", "thirties", "thirty", "thousand", "thousands", "three", "trillion", "twelve", "twenties", "twenty", "two", "zero", "zillion",
} // PTR_10078768, count DAT_1007883c (53), result 9.

var paul2013TripletApostropheSLeftClass3Words = []string{
	"anybody", "everybody", "he", "i", "it", "nobody", "she", "they", "we", "you",
} // PTR_100786f0, count DAT_10078718 (10); FUN_10006ce0 returns 3 for middle "'s".

// Paul2013B800SuffixWordTableMatches looks up the two sorted pointer tables
// consulted by FUN_1000b800's apostrophe-d special case. Matching uses the
// DLL's mapped C-string ordering; the returned booleans preserve table
// membership without assigning linguistic meaning to either table.
func Paul2013B800SuffixWordTableMatches(surface []byte) (apostropheSLeftClass3, class11 bool) {
	weights := paul2013ContextCharacterWeights()
	leftRows := make([][]byte, len(paul2013TripletApostropheSLeftClass3Words))
	for index, word := range paul2013TripletApostropheSLeftClass3Words {
		leftRows[index] = []byte(word)
	}
	class11Rows := make([][]byte, len(paul2013TripletClass11Words))
	for index, word := range paul2013TripletClass11Words {
		class11Rows[index] = []byte(word)
	}
	return FindPaul2013SortedCString(leftRows, surface, weights) >= 0,
		FindPaul2013SortedCString(class11Rows, surface, weights) >= 0
}

// Paul2013ContextMappedCStringEqual compares C strings with the local DLL's
// mapped-character weights used by its context handlers.
func Paul2013ContextMappedCStringEqual(left, right []byte) bool {
	return ComparePaul2013MappedCString(left, right, paul2013ContextCharacterWeights()) == 0
}

// Paul2013ContextWordPairSpecial ports FUN_10008440's recovered two-string
// predicate. Its conditions are exposed as a boolean because the native
// helper's return value is only used as a gate by its callers.
func Paul2013ContextWordPairSpecial(left, right []byte) bool {
	left = cString(left)
	right = cString(right)
	weights := paul2013ContextCharacterWeights()
	if len(right) == 0 {
		rows := make([][]byte, len(paul2013TripletClass2Words))
		for index, word := range paul2013TripletClass2Words {
			rows[index] = []byte(word)
		}
		return FindPaul2013SortedCString(rows, left, weights) >= 0
	}
	if (Paul2013ContextMappedCStringEqual(left, []byte("wo")) ||
		Paul2013ContextMappedCStringEqual(left, []byte("sha")) ||
		Paul2013ContextMappedCStringEqual(left, []byte("ca"))) &&
		Paul2013ContextMappedCStringEqual(right, []byte("n't")) {
		return true
	}
	for _, word := range [...]string{"ought", "have", "has", "had", "need"} {
		if Paul2013ContextMappedCStringEqual(left, []byte(word)) {
			return Paul2013ContextMappedCStringEqual(right, []byte("'s"))
		}
	}
	return false
}

// Paul2013IsNumberWord ports FUN_10008550's mode-0x49 lookup in the 53-entry
// number-word table used by the native model-context pass.
func Paul2013IsNumberWord(surface []byte) bool {
	rows := make([][]byte, len(paul2013TripletClass9Words))
	for index, word := range paul2013TripletClass9Words {
		rows[index] = []byte(word)
	}
	return FindPaul2013SortedCString(rows, surface, paul2013ContextCharacterWeights()) >= 0
}

var paul2013TripletClassGroups = [...]struct {
	class int16
	words []string
}{
	{2, paul2013TripletClass2Words},
	{3, paul2013TripletClass3Words},
	{4, paul2013TripletClass4Words},
	{12, paul2013TripletClass12Words},
	{6, paul2013TripletClass6Words},
	{7, paul2013TripletClass7Words},
	{8, paul2013TripletClass8Words},
	{11, paul2013TripletClass11Words},
	{9, paul2013TripletClass9Words},
}

// FUN_10006ce0 passes this NUL-terminated byte string to FUN_10064440 for
// each middle-word byte after its class-9 character-attribute scan.
const paul2013TripletClass10Punctuation = "`~!@#$%^&*()-_+=|\\[]{};:'\"<>?,./"

func paul2013ContainsWord(words []string, surface string) bool {
	index := sort.SearchStrings(words, surface)
	return index < len(words) && words[index] == surface
}

func paul2013TripletFallbackClass(surface string) int16 {
	end := strings.IndexByte(surface, 0)
	if end < 0 {
		end = len(surface)
	}
	for i := 0; i < end; i++ {
		if surface[i] >= '0' && surface[i] <= '9' {
			return 9
		}
	}
	for i := 0; i < end; i++ {
		if strings.IndexByte(paul2013TripletClass10Punctuation, surface[i]) >= 0 {
			return 10
		}
	}
	return 1
}

// Paul2013KnownPronunciationTripletFeature returns feature values established
// by recovered fast paths in FUN_10006ce0. A null or empty middle word returns
// zero; the exact case-sensitive middle word "you" returns eight before
// either neighboring word is examined. It also ports recovered literal
// conditions from FUN_10006ce0, including several three-word cases and the
// left/middle "let" / "'s" case. It then searches the recovered sorted middle
// word tables in native order and applies the classifier's digit, punctuation,
// and default-character fallbacks.
func Paul2013KnownPronunciationTripletFeature(left, middle, right string) (int16, bool) {
	if middle == "" {
		return 0, true
	}
	if middle == "you" {
		return 8, true
	}
	if (middle == "wo" || middle == "sha" || middle == "ca") && right == "job" {
		return 2, true
	}
	if (middle == "ought" || middle == "have" || middle == "has" || middle == "had" || middle == "need") && right == "to" {
		return 2, true
	}
	if (left == "ought" || left == "have" || left == "has" || left == "had" || left == "need") && middle == "to" {
		return 2, true
	}
	if middle == "'s" && left == "let" {
		return 8, true
	}
	groupStart := 0
	if middle == "'s" {
		if paul2013ContainsWord(paul2013TripletClass8Words, left) ||
			paul2013ContainsWord(paul2013TripletApostropheSLeftClass3Words, left) ||
			paul2013ContainsWord(paul2013TripletClass11Words, left) {
			return 3, true
		}
		// The native apostrophe-s path bypasses classes 2, 3, 4, 12, and 6.
		groupStart = 5
	}
	for _, group := range paul2013TripletClassGroups[groupStart:] {
		if paul2013ContainsWord(group.words, middle) {
			return group.class, true
		}
	}
	if Paul2013PronunciationWordClass(middle) > 1 {
		return 5, true
	}
	return paul2013TripletFallbackClass(middle), true
}

// BuildPaul2013PronunciationWordFeatures builds feature positions 0 through 3
// of FUN_10006ae0 from the two preceding and two following tokens. The center
// word is used by the separate triplet and byte-class features. Sequence
// boundaries and empty token surfaces map to zero, matching FUN_10006ca0's
// null/empty branch.
func BuildPaul2013PronunciationWordFeatures(tokens []LexicalToken, center int) ([4]int16, error) {
	var features [4]int16
	if len(tokens) == 0 {
		return features, errors.New("pronunciation word window has no tokens")
	}
	if center < 0 || center >= len(tokens) {
		return features, fmt.Errorf("pronunciation word window center %d is outside %d tokens", center, len(tokens))
	}
	for featureIndex, offset := range [...]int{-2, -1, 1, 2} {
		tokenIndex := center + offset
		if tokenIndex < 0 || tokenIndex >= len(tokens) {
			continue
		}
		features[featureIndex] = Paul2013PronunciationWordClass(tokens[tokenIndex].Surface)
	}
	return features, nil
}

// BuildPaul2013KnownPronunciationFeatures fills the classifier positions
// currently derivable directly from lexical token surfaces. Missing neighbors
// at utterance edges are known empty inputs. Triplet positions 4–7 are fully
// produced by the recovered tables, literal conditions, and character-based
// fallbacks in FUN_10006ce0. Positions 8–12 apply FUN_10007160 to the two
// preceding words, center word, and two following words. Missing words map to
// zero through that helper's null/empty early return.
func BuildPaul2013KnownPronunciationFeatures(tokens []LexicalToken, center int) (Paul2013PronunciationFeatures, error) {
	var features Paul2013PronunciationFeatures
	wordFeatures, err := BuildPaul2013PronunciationWordFeatures(tokens, center)
	if err != nil {
		return features, err
	}
	copy(features.Values[:4], wordFeatures[:])
	features.Available |= 0x000f
	for tripletIndex, middleOffset := range [...]int{-2, -1, 1, 2} {
		middleIndex := center + middleOffset
		left, middle, right := "", "", ""
		if middleIndex >= 0 && middleIndex < len(tokens) {
			middle = tokens[middleIndex].Surface
		}
		if middleIndex-1 >= 0 && middleIndex-1 < len(tokens) {
			left = tokens[middleIndex-1].Surface
		}
		if middleIndex+1 >= 0 && middleIndex+1 < len(tokens) {
			right = tokens[middleIndex+1].Surface
		}
		if value, known := Paul2013KnownPronunciationTripletFeature(left, middle, right); known {
			position := 4 + tripletIndex
			features.Values[position] = value
			features.Available |= 1 << position
		}
	}
	for phoneIndex, offset := range [...]int{-2, -1, 0, 1, 2} {
		wordIndex := center + offset
		if wordIndex >= 0 && wordIndex < len(tokens) {
			features.Values[8+phoneIndex] = Paul2013PronunciationPhoneContextFeature(tokens[wordIndex].Surface)
		}
		features.Available |= 1 << (8 + phoneIndex)
	}
	features.Values[14] = Paul2013PronunciationContinuityFlag(tokens[center].Surface)
	features.Available |= 1 << 14
	return features, nil
}

// BuildPaul2013PronunciationPathFeatures is a convenience wrapper that adds
// the first code's class at position 13. Native candidate generation should
// use BuildPaul2013PronunciationPathCodeFeatures for every code in the group.
func BuildPaul2013PronunciationPathFeatures(
	tokens []LexicalToken,
	center int,
	pathGroup []byte,
) (Paul2013PronunciationFeatures, bool, error) {
	if len(pathGroup) == 0 {
		features, err := BuildPaul2013KnownPronunciationFeatures(tokens, center)
		return features, false, err
	}
	return BuildPaul2013PronunciationPathCodeFeatures(tokens, center, pathGroup[0])
}

// BuildPaul2013PronunciationPathCodeFeatures adds feature position 13 for
// one code inside a path group. FUN_100068b0 evaluates every code before it
// scores the resulting classifier outputs against all path groups.
func BuildPaul2013PronunciationPathCodeFeatures(
	tokens []LexicalToken,
	center int,
	pathCode byte,
) (Paul2013PronunciationFeatures, bool, error) {
	features, err := BuildPaul2013KnownPronunciationFeatures(tokens, center)
	if err != nil {
		return features, false, err
	}
	pathClass, present, err := Paul2013PronunciationPathClass([]byte{pathCode})
	if err != nil || !present {
		return features, present, err
	}
	features.Values[13] = pathClass
	features.Available |= 1 << 13
	return features, true, nil
}
