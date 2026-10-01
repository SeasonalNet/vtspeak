package text

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

// Paul2013ModelPronunciationRow is the caller-visible portion of one native
// 0x554-byte token result row consumed by FUN_100049b0 and FUN_100068b0.
// Offsets and slot order follow FUN_1000d450; this projection assigns no
// linguistic meaning to the type or path bytes.
type Paul2013ModelPronunciationRow struct {
	Index             uint16
	Type              byte
	Surface           []byte
	AlternativeCount  int
	PathControlBytes  []byte
	PhoneAlternatives [][]byte
}

// Paul2013NameContextSelection records one direct FUN_100049b0 selection.
// DirectOutput is true when FUN_10003ff0 writes PhoneString itself;
// otherwise SelectorCode is the path byte requested from FUN_10010640.
type Paul2013NameContextSelection struct {
	Applicable   bool
	DirectOutput bool
	SelectorCode byte
	Alternative  int
	PhoneString  []byte
}

// SelectPaul2013ModelNameContextPronunciation ports direct word/context
// branches in FUN_100049b0. For `anti`, a following class-4 surface requests
// code 0x0e and other positions request 0x13. `conflicts` followed by `with`
// requests '*'; `supply` followed by a class-3 surface requests 0x13;
// `converts` requests 0x16 before `from` or after `the`/`a`, and '*' otherwise;
// `wound` followed by `up`, `down`, or `through` requests '&'; and `minute`
// requests 0x0e before `amount` or 0x13 after `the`/`a` or in `minute by
// minute` contexts. Failed gates for words handled by this helper fall through
// to the unported FUN_10004300 path and fail closed. `number` requests 0x0f
// when preceded by `less` or `more`,
// or followed by `with`; other contexts request 0x13.
// `resume`/`resumes` after at least two rows and followed by sentence
// punctuation request '*'/'%' only when the previous row is outside class 3
// and the row two positions back is class 11. `increase`/`increases` after `the`/`an` or before a class-3
// surface request 0x16/0x13. `transform`/`transforms` before `into` two or
// three rows later request '*'/'%'; `use` after `still` requests '%'; `wind`
// before `up`/`down`/`through` requests '%'; and `de` beside an uppercase-
// initial surface requests 0x0c. `house` beside an uppercase-initial surface
// or before `arrest`, `number`, or `numbers`, and `job hunting` each select
// 0x13; `laden` selects 0x14
// before `'s`, after `bin`, or at the start with an uppercase initial, and
// 0x0e otherwise; `dove into` selects '&'. `interstate(s) to` selects `%`/`*`.
// `dogged` searches its preceding three-row window against the recovered
// ten-phrase table; the match position and previous-row gate choose 0x0e or `(`.
// `does` selects `*` after `much` or before `not`. The early `August`
// shortcut selects 0x1e after `an` or `the`, beside a comma-initial row, at
// the final row, before terminal punctuation in the penultimate row, or next
// to a one-byte surface with the recovered 0xc0 character-attribute bits; it
// also selects 0x1e before a class-4 surface, for its recovered following
// class/pair gate and following-pair predicate, and for its two-row
// all-characters attribute check. Other August contexts select 0x07 after
// these recovered gates, matching the native helper's default selector.
// `invalid` selects 0x0e after class 3 or after class 3 plus `n't`/`not`.
// `lied` selects `&` for its lowercase/uppercase, following-word, and adverb
// gates; `have`/`had`/`been` predecessors select `(` at the remaining gate.
// `lives on` selects `*`, except after `put`/`lay` where it selects 0x16;
// `live on TV` selects 0x0e; `live(s) it up`, `live(s) up to`, and
// `live(s) to *self` select `%`/`*`; the recovered following-literal set and
// FUN_10007160 feature 4 select `%`/`*` for `live`/`lives`. The recovered
// previous-surface exception set selects 0x0e/0x16 for those forms. Other
// contexts fall through. For article `a`, FUN_10003ff0 returns the one-byte
// 0x1e phone for its recovered article, class, punctuation, word-pair, and
// character-table gates, and 0x07 when those gates do not match.
// `separate`/`separates` after `be` or a class-3 triplet predecessor request
// 0x0e; following `from`/`into`/`out`/`up`, or those first two plus `and` two
// rows later, request `*`/`%`. `perfect` followed by
// sentence punctuation and `polish Jewish` request
// 0x0e. `learned` requests '&' at the end of the row list or before
// sentence punctuation. `present` with `to` two rows later and a further
// token row present requests '%'. For `can`, the recovered precedence selects
// 0x13 after `steel`, 0x12 before `be`, 0x13 for class-12/context-table
// neighbors and punctuation, and 0x12 otherwise.
// For `close`, the directly recovered preceding-word rules cover `I`, `you`,
// `we`, and `they` (selector '%'), `feel` (0x0e), `fisherman`/`fishermans`
// (0x13), plus `to a` two rows earlier ('%'). Following `than`, `by`, `on`,
// `to`, `upon`, `call`, or `calls` selects 0x0e; `at hand` and `in with` also
// select 0x0e. The `it up` and `down on` pairs select '%'.
// Other `in` and `up` contexts continue to the generic fallback. Other
// `close` contexts still require the unported fallback.
// `brain reading` selects 0x0e. `contest`/`contests` followed by a class-3
// surface select 0x16/0x13 respectively. `mouth to mouth` selects 0x13. `closer`
// selects 0x0f before `than`, `by`, `on`, `to`, `upon`, `call`,
// `calls`, `look`, or `looks`, after a following surface classified as class
// 2 by FUN_10007160, and in `at hand` or `in with`. Its other recognized
// contexts remain fail-closed. For unhandled words, directly recovered
// `FUN_10004300` gates include, in native order, a non-`nice` word before `of`
// (0x13, 0x16, 0x14, then 0x15), a word before `'s` (0x14, 0x15, 0x13, then
// 0x16), compound prefixes before a same-index hyphen row (0x0d), the
// `less`/`more`/`so`/`very` class-2 then class-8 choice, words following the
// recovered adverb set (optional '(' after triplet class 3, then %, *, or
// &), the FUN_10010690 preceding-word gate (class 4/7, then 0x13, guarded by
// class 5/6 alternatives), and the FUN_10010790 next-word gate with its
// class-10 cross-alternative guard before `%`. It also ports FUN_10010840's
// recovered word-pair predicate, the `have`/`has`/`had` parenthesis choice, and
// the FUN_10010760 verb + long `-ing` apostrophe choice. Other fallback
// branches include the class-7 triplet route through class 4/7 or 0x13, plus
// `how` immediately before the current row or two rows earlier when another
// row follows and the current surface is class 8 (0x0e).
// Other direct name/context branches remain unsupported. A selected code must be present
// in the current row's path stream.
func SelectPaul2013ModelNameContextPronunciation(
	tokenRows []byte,
	rowIndex int,
) (Paul2013NameContextSelection, error) {
	rows, err := ReadPaul2013ModelPronunciationRows(tokenRows)
	if err != nil {
		return Paul2013NameContextSelection{}, fmt.Errorf("read model pronunciation rows: %w", err)
	}
	if rowIndex < 0 || rowIndex >= len(rows) {
		return Paul2013NameContextSelection{}, fmt.Errorf("name/context row %d is outside %d pronunciation rows", rowIndex, len(rows))
	}
	current := rows[rowIndex]
	selectDirectPhone := func(phone byte) Paul2013NameContextSelection {
		return Paul2013NameContextSelection{
			Applicable: true, DirectOutput: true, Alternative: -1, PhoneString: []byte{phone},
		}
	}
	selectFallback := func() (Paul2013NameContextSelection, error) {
		selection, found, err := selectPaul2013ModelFallbackPathChoice(rows, rowIndex)
		if err != nil {
			return Paul2013NameContextSelection{}, err
		}
		if !found {
			return Paul2013NameContextSelection{}, nil
		}
		return selection, nil
	}
	selector := byte(0)
	switch {
	case Paul2013ContextMappedCStringEqual(current.Surface, []byte("a")):
		// FUN_100049b0 sends the article `a` through FUN_10003ff0 before
		// general scoring. Its direct gates write one byte to the output.
		previousArticle := rowIndex > 0 && paul2013NameContextSurfaceIsOneOf(
			rows[rowIndex-1].Surface, "an", "the",
		)
		nextClassFour := rowIndex+1 < len(rows) &&
			Paul2013TripletTableContains(4, rows[rowIndex+1].Surface)
		nextSurface := []byte(nil)
		if rowIndex+1 < len(rows) {
			nextSurface = cString(rows[rowIndex+1].Surface)
		}
		attributes := Paul2013ExceptionCharacterAttributes()
		nextClassTwoOrThree := rowIndex+1 < len(rows) &&
			(Paul2013TripletTableContains(3, nextSurface) || Paul2013ContextWordPairSpecial(nextSurface, nil))
		nextOf := rowIndex+1 < len(rows) && Paul2013ContextMappedCStringEqual(nextSurface, []byte("of"))
		nextCanOrMust := rowIndex+1 < len(rows) &&
			paul2013NameContextSurfaceIsOneOf(nextSurface, "can", "must")
		nextInitialHasHighAttribute := len(nextSurface) > 0 && attributes[nextSurface[0]]&0x80 != 0
		nextClassWordGate := (nextClassTwoOrThree || nextOf) &&
			(!nextCanOrMust || nextInitialHasHighAttribute)
		nextPairSpecial := rowIndex+2 < len(rows) &&
			Paul2013ContextWordPairSpecial(nextSurface, cString(rows[rowIndex+2].Surface))
		previousStartsComma := false
		if rowIndex > 0 {
			previous := cString(rows[rowIndex-1].Surface)
			previousStartsComma = len(previous) > 0 && previous[0] == ','
		}
		commaPair := (rowIndex == 0 || previousStartsComma) && len(nextSurface) > 0 && nextSurface[0] == ','
		neighborHasSingleByteClass := func(surface []byte) bool {
			value := cString(surface)
			return len(value) == 1 && attributes[value[0]]&0xc0 != 0
		}
		previousClassCharacter := rowIndex > 0 && neighborHasSingleByteClass(rows[rowIndex-1].Surface)
		nextClassCharacter := rowIndex+1 < len(rows) && neighborHasSingleByteClass(rows[rowIndex+1].Surface)
		previousAttributeSpecial := false
		if rowIndex > 0 && rowIndex+1 < len(rows) {
			previousSurface := cString(rows[rowIndex-1].Surface)
			nextSurface := cString(rows[rowIndex+1].Surface)
			previousAttributeSpecial = len(previousSurface) > 0 && attributes[previousSurface[0]]&0x80 != 0 &&
				paul2013NameContextAllCharactersHaveAttribute(previousSurface, attributes, 0x40) &&
				!paul2013NameContextAllCharactersHaveAttribute(nextSurface, attributes, 0x40)
		}
		terminalRow := rowIndex == len(rows)-1
		terminalPunctuation := rowIndex+2 == len(rows) && rowIndex+1 < len(rows) &&
			len(nextSurface) > 0 && bytes.IndexByte([]byte(".?!,"), nextSurface[0]) >= 0
		if previousArticle || nextClassFour || nextClassWordGate || commaPair || previousAttributeSpecial ||
			nextPairSpecial || previousClassCharacter || nextClassCharacter || terminalRow || terminalPunctuation {
			return selectDirectPhone(0x1e), nil
		}
		return selectDirectPhone(0x07), nil
	case Paul2013ContextMappedCStringEqual(current.Surface, []byte("anti")):
		selector = 0x13
		if rowIndex+1 < len(rows) && Paul2013TripletTableContains(4, rows[rowIndex+1].Surface) {
			selector = 0x0e
		}
	case Paul2013ContextMappedCStringEqual(current.Surface, []byte("conflicts")):
		if rowIndex+1 >= len(rows) || !Paul2013ContextMappedCStringEqual(rows[rowIndex+1].Surface, []byte("with")) {
			return selectFallback()
		}
		selector = '*'
	case paul2013NameContextSurfaceIsOneOf(current.Surface, "contest", "contests"):
		if rowIndex+1 >= len(rows) || !Paul2013TripletTableContains(3, rows[rowIndex+1].Surface) {
			return selectFallback()
		}
		if Paul2013ContextMappedCStringEqual(current.Surface, []byte("contests")) {
			selector = 0x13
		} else {
			selector = 0x16
		}
	case Paul2013ContextMappedCStringEqual(current.Surface, []byte("supply")):
		if rowIndex+1 >= len(rows) || !Paul2013TripletTableContains(3, rows[rowIndex+1].Surface) {
			return selectFallback()
		}
		selector = 0x13
	case Paul2013ContextMappedCStringEqual(current.Surface, []byte("converts")):
		nextFrom := rowIndex+1 < len(rows) && Paul2013ContextMappedCStringEqual(rows[rowIndex+1].Surface, []byte("from"))
		previousArticle := rowIndex > 0 && paul2013NameContextSurfaceIsOneOf(rows[rowIndex-1].Surface, "the", "a")
		if !nextFrom && !previousArticle {
			selector = '*'
		} else {
			selector = 0x16
		}
	case Paul2013ContextMappedCStringEqual(current.Surface, []byte("does")):
		previousMuch := rowIndex > 0 && Paul2013ContextMappedCStringEqual(rows[rowIndex-1].Surface, []byte("much"))
		nextNot := rowIndex+1 < len(rows) && Paul2013ContextMappedCStringEqual(rows[rowIndex+1].Surface, []byte("not"))
		if !previousMuch && !nextNot {
			return selectFallback()
		}
		selector = '*'
	case Paul2013ContextMappedCStringEqual(current.Surface, []byte("august")):
		previousArticle := rowIndex > 0 && paul2013NameContextSurfaceIsOneOf(
			rows[rowIndex-1].Surface, "an", "the",
		)
		previousComma := rowIndex > 0 && bytes.HasPrefix(cString(rows[rowIndex-1].Surface), []byte{','})
		nextSurface := []byte(nil)
		if rowIndex+1 < len(rows) {
			nextSurface = cString(rows[rowIndex+1].Surface)
		}
		attributes := Paul2013ExceptionCharacterAttributes()
		nextComma := len(nextSurface) > 0 && nextSurface[0] == ','
		nextClassFour := rowIndex+1 < len(rows) && Paul2013TripletTableContains(4, nextSurface)
		nextClassOrPairGate := false
		if rowIndex+1 < len(rows) {
			nextClassOrPairGate = Paul2013TripletTableContains(3, nextSurface) ||
				Paul2013TripletTableContains(2, nextSurface) ||
				Paul2013ContextMappedCStringEqual(nextSurface, []byte("of"))
			if nextClassOrPairGate {
				nextCanOrMust := paul2013NameContextSurfaceIsOneOf(nextSurface, "can", "must")
				nextClassOrPairGate = !nextCanOrMust || attributes[nextSurface[0]]&0x80 != 0
			}
		}
		nextPairSpecial := rowIndex+2 < len(rows) &&
			Paul2013ContextWordPairSpecial(nextSurface, cString(rows[rowIndex+2].Surface))
		previousAttributeSpecial := false
		if rowIndex > 0 && rowIndex+1 < len(rows) {
			previousSurface := cString(rows[rowIndex-1].Surface)
			previousAttributeSpecial = len(previousSurface) > 0 && attributes[previousSurface[0]]&0x80 != 0 &&
				paul2013NameContextAllCharactersHaveAttribute(previousSurface, attributes, 0x40) &&
				!paul2013NameContextAllCharactersHaveAttribute(nextSurface, attributes, 0x40)
		}
		terminalRow := rowIndex == len(rows)-1
		terminalPunctuation := rowIndex+2 == len(rows) && len(nextSurface) > 0 &&
			bytes.IndexByte([]byte(".?!,"), nextSurface[0]) >= 0
		neighborHasSingleByteClass := func(surface []byte) bool {
			value := cString(surface)
			return len(value) == 1 && attributes[value[0]]&0xc0 != 0
		}
		previousClassCharacter := rowIndex > 0 && neighborHasSingleByteClass(rows[rowIndex-1].Surface)
		nextClassCharacter := rowIndex+1 < len(rows) && neighborHasSingleByteClass(rows[rowIndex+1].Surface)
		if previousArticle || previousComma || nextComma || nextClassFour || nextClassOrPairGate ||
			nextPairSpecial || previousAttributeSpecial || terminalRow || terminalPunctuation ||
			previousClassCharacter || nextClassCharacter {
			selector = 0x1e
		} else {
			selector = 0x07
		}
	case Paul2013ContextMappedCStringEqual(current.Surface, []byte("dogged")):
		matchPosition := paul2013ModelContextPhrasePosition(
			rows, rowIndex, 'L', 3, paul2013DoggedContextPhrases[:],
		)
		if matchPosition < 0 {
			selector = 0x0e
		} else if matchPosition > 0 && rowIndex > 0 &&
			(Paul2013TripletClass12Contains(rows[rowIndex-1].Surface) ||
				Paul2013ContextMappedCStringEqual(rows[rowIndex-1].Surface, []byte("that"))) {
			selector = 0x0e
		} else {
			selector = '('
		}
	case Paul2013ContextMappedCStringEqual(current.Surface, []byte("lied")):
		currentStartsLowercaseL := len(cString(current.Surface)) > 0 && cString(current.Surface)[0] == 'l'
		previousUpperInitial := rowIndex > 0 && Paul2013PronunciationContinuityFlag(
			string(cString(rows[rowIndex-1].Surface)),
		) != 0
		nextLiteralContext := rowIndex+1 < len(rows) && paul2013NameContextSurfaceIsOneOf(
			rows[rowIndex+1].Surface,
			"about", "along", "around", "at", "back", "before", "behind", "by",
			"close", "down", "in", "into", "off", "on", "out", "over", "to", "up", "with",
		)
		previousAdverb := rowIndex > 0 && paul2013NameContextIsFUN10010730Word(rows[rowIndex-1].Surface)
		previousAuxiliary := rowIndex > 0 && paul2013NameContextSurfaceIsOneOf(
			rows[rowIndex-1].Surface, "have", "had", "been",
		)
		switch {
		case currentStartsLowercaseL && previousUpperInitial, nextLiteralContext, previousAdverb:
			selector = '&'
		case previousAuxiliary:
			selector = '('
		default:
			return selectFallback()
		}
	case Paul2013ContextMappedCStringEqual(current.Surface, []byte("wound")):
		if rowIndex+1 >= len(rows) || !paul2013NameContextSurfaceIsOneOf(rows[rowIndex+1].Surface, "up", "down", "through") {
			return selectFallback()
		}
		selector = '&'
	case Paul2013ContextMappedCStringEqual(current.Surface, []byte("minute")):
		nextAmount := rowIndex+1 < len(rows) && Paul2013ContextMappedCStringEqual(rows[rowIndex+1].Surface, []byte("amount"))
		previousArticle := rowIndex > 0 && paul2013NameContextSurfaceIsOneOf(rows[rowIndex-1].Surface, "the", "a")
		forwardMinuteByMinute := rowIndex+2 < len(rows) &&
			Paul2013ContextMappedCStringEqual(rows[rowIndex+1].Surface, []byte("by")) &&
			Paul2013ContextMappedCStringEqual(rows[rowIndex+2].Surface, []byte("minute"))
		backwardMinuteByMinute := rowIndex >= 2 &&
			Paul2013ContextMappedCStringEqual(rows[rowIndex-1].Surface, []byte("by")) &&
			Paul2013ContextMappedCStringEqual(rows[rowIndex-2].Surface, []byte("minute"))
		switch {
		case nextAmount:
			selector = 0x0e
		case previousArticle, forwardMinuteByMinute, backwardMinuteByMinute:
			selector = 0x13
		default:
			return selectFallback()
		}
	case Paul2013ContextMappedCStringEqual(current.Surface, []byte("number")):
		previousLessOrMore := rowIndex > 0 && paul2013NameContextSurfaceIsOneOf(rows[rowIndex-1].Surface, "less", "more")
		nextWith := rowIndex+1 < len(rows) && Paul2013ContextMappedCStringEqual(rows[rowIndex+1].Surface, []byte("with"))
		switch {
		case previousLessOrMore || nextWith:
			selector = 0x0f
		default:
			selector = 0x13
		}
	case Paul2013ContextMappedCStringEqual(current.Surface, []byte("perfect")):
		if rowIndex+1 >= len(rows) || !paul2013NameContextIsSentencePunctuation(rows[rowIndex+1].Surface) {
			return selectFallback()
		}
		selector = 0x0e
	case Paul2013ContextMappedCStringEqual(current.Surface, []byte("polish")):
		if rowIndex+1 >= len(rows) || !Paul2013ContextMappedCStringEqual(rows[rowIndex+1].Surface, []byte("jewish")) {
			return selectFallback()
		}
		selector = 0x0e
	case Paul2013ContextMappedCStringEqual(current.Surface, []byte("learned")):
		if rowIndex+1 < len(rows) && !paul2013NameContextIsSentencePunctuation(rows[rowIndex+1].Surface) {
			return selectFallback()
		}
		selector = '&'
	case Paul2013ContextMappedCStringEqual(current.Surface, []byte("present")):
		if rowIndex+3 >= len(rows) || !Paul2013ContextMappedCStringEqual(rows[rowIndex+2].Surface, []byte("to")) {
			return selectFallback()
		}
		selector = '%'
	case Paul2013ContextMappedCStringEqual(current.Surface, []byte("can")):
		previousSteel := rowIndex > 0 && Paul2013ContextMappedCStringEqual(rows[rowIndex-1].Surface, []byte("steel"))
		nextExists := rowIndex+1 < len(rows)
		nextBe := nextExists && Paul2013ContextMappedCStringEqual(rows[rowIndex+1].Surface, []byte("be"))
		previousContext := rowIndex > 0 && paul2013NameContextIsClass12OrContextWord(rows[rowIndex-1].Surface)
		nextContext := nextExists && paul2013NameContextIsContextWord(rows[rowIndex+1].Surface)
		nextPunctuation := nextExists && paul2013NameContextIsSentencePunctuation(rows[rowIndex+1].Surface)
		switch {
		case previousSteel:
			selector = 0x13
		case nextBe:
			selector = 0x12
		case previousContext, nextContext, nextPunctuation, !nextExists:
			selector = 0x13
		default:
			selector = 0x12
		}
	case Paul2013ContextMappedCStringEqual(current.Surface, []byte("close")):
		selector, err = selectPaul2013ModelCloseContext(rows, rowIndex)
		if err != nil {
			return selectFallback()
		}
	case Paul2013ContextMappedCStringEqual(current.Surface, []byte("closer")):
		selector, err = selectPaul2013ModelCloserContext(rows, rowIndex)
		if err != nil {
			return selectFallback()
		}
	case Paul2013ContextMappedCStringEqual(current.Surface, []byte("mouth")):
		if rowIndex+2 >= len(rows) ||
			!Paul2013ContextMappedCStringEqual(rows[rowIndex+1].Surface, []byte("to")) ||
			!Paul2013ContextMappedCStringEqual(rows[rowIndex+2].Surface, []byte("mouth")) {
			return selectFallback()
		}
		selector = 0x13
	case Paul2013ContextMappedCStringEqual(current.Surface, []byte("reading")):
		if rowIndex == 0 || !Paul2013ContextMappedCStringEqual(rows[rowIndex-1].Surface, []byte("brain")) {
			return selectFallback()
		}
		selector = 0x0e
	case paul2013NameContextSurfaceIsOneOf(current.Surface, "record", "records"):
		previousSpecial := rowIndex > 0 && paul2013TripletWordsContainMapped(
			paul2013TripletApostropheSLeftClass3Words, rows[rowIndex-1].Surface,
		)
		previousNumber := rowIndex > 0 && Paul2013IsNumberWord(rows[rowIndex-1].Surface)
		if previousSpecial {
			selector = '%'
		} else {
			if !previousNumber && (rowIndex+1 >= len(rows) ||
				(!Paul2013ContextMappedCStringEqual(rows[rowIndex+1].Surface, []byte("it")) &&
					!Paul2013ContextMappedCStringEqual(rows[rowIndex+1].Surface, []byte("to")))) {
				return selectFallback()
			}
			if Paul2013ContextMappedCStringEqual(current.Surface, []byte("records")) {
				selector = 0x13
			} else {
				selector = 0x16
			}
		}
	case paul2013NameContextSurfaceIsOneOf(current.Surface, "refuse", "refuses"):
		if rowIndex+1 >= len(rows) || !Paul2013ContextMappedCStringEqual(rows[rowIndex+1].Surface, []byte("to")) {
			return selectFallback()
		}
		if Paul2013ContextMappedCStringEqual(current.Surface, []byte("refuses")) {
			selector = '%'
		} else {
			selector = '*'
		}
	case Paul2013ContextMappedCStringEqual(current.Surface, []byte("import")):
		previousArticleOrExport := rowIndex > 0 && paul2013NameContextSurfaceIsOneOf(
			rows[rowIndex-1].Surface, "export", "the", "an",
		)
		nextBank := rowIndex+1 < len(rows) &&
			Paul2013ContextMappedCStringEqual(rows[rowIndex+1].Surface, []byte("bank"))
		if !previousArticleOrExport && !nextBank {
			return selectFallback()
		}
		selector = 0x13
	case paul2013NameContextSurfaceIsOneOf(current.Surface, "live", "lives"):
		nextOn := rowIndex+1 < len(rows) && Paul2013ContextMappedCStringEqual(rows[rowIndex+1].Surface, []byte("on"))
		if Paul2013ContextMappedCStringEqual(current.Surface, []byte("lives")) && nextOn {
			previousPutOrLay := rowIndex > 0 && paul2013NameContextSurfaceIsOneOf(
				rows[rowIndex-1].Surface, "put", "lay",
			)
			if previousPutOrLay {
				selector = 0x16
			} else {
				selector = '*'
			}
		} else if nextOn && rowIndex+2 < len(rows) &&
			Paul2013ContextMappedCStringEqual(rows[rowIndex+2].Surface, []byte("TV")) {
			if Paul2013ContextMappedCStringEqual(current.Surface, []byte("live")) {
				selector = 0x0e
			} else {
				selector = 0x16
			}
		} else if rowIndex+1 < len(rows) && paul2013NameContextSurfaceIsOneOf(
			rows[rowIndex+1].Surface,
			"at", "by", "down", "for", "in", "near", "off", "on", "out", "through",
			"to", "together", "up", "upon", "with", "call", "calls",
		) {
			if Paul2013ContextMappedCStringEqual(current.Surface, []byte("lives")) {
				selector = '*'
			} else {
				selector = '%'
			}
		} else if rowIndex+2 < len(rows) &&
			((Paul2013ContextMappedCStringEqual(rows[rowIndex+1].Surface, []byte("it")) &&
				Paul2013ContextMappedCStringEqual(rows[rowIndex+2].Surface, []byte("up"))) ||
				(Paul2013ContextMappedCStringEqual(rows[rowIndex+1].Surface, []byte("up")) &&
					Paul2013ContextMappedCStringEqual(rows[rowIndex+2].Surface, []byte("to"))) ||
				(Paul2013ContextMappedCStringEqual(rows[rowIndex+1].Surface, []byte("to")) &&
					paul2013NameContextEndsWithSelf(rows[rowIndex+2].Surface))) {
			if Paul2013ContextMappedCStringEqual(current.Surface, []byte("live")) {
				selector = '%'
			} else {
				selector = '*'
			}
		} else if rowIndex+1 < len(rows) &&
			Paul2013PronunciationPhoneContextFeature(string(cString(rows[rowIndex+1].Surface))) == 4 {
			if Paul2013ContextMappedCStringEqual(current.Surface, []byte("lives")) {
				selector = '*'
			} else {
				selector = '%'
			}
		} else if paul2013NameContextLivePreviousException(rows, rowIndex) {
			if Paul2013ContextMappedCStringEqual(current.Surface, []byte("lives")) {
				selector = 0x16
			} else {
				selector = 0x0e
			}
		} else {
			return selectFallback()
		}
	case paul2013NameContextSurfaceIsOneOf(current.Surface, "increase", "increases"):
		previousArticle := rowIndex > 0 && paul2013NameContextSurfaceIsOneOf(
			rows[rowIndex-1].Surface, "the", "an",
		)
		nextClass3 := rowIndex+1 < len(rows) && Paul2013TripletTableContains(3, rows[rowIndex+1].Surface)
		if !previousArticle && !nextClass3 {
			return selectFallback()
		}
		if Paul2013ContextMappedCStringEqual(current.Surface, []byte("increases")) {
			selector = 0x13
		} else {
			selector = 0x16
		}
	case Paul2013ContextMappedCStringEqual(current.Surface, []byte("house")):
		previousUpperInitial := rowIndex > 0 && Paul2013PronunciationContinuityFlag(
			string(cString(rows[rowIndex-1].Surface)),
		) != 0
		nextUpperInitial := rowIndex+1 < len(rows) && Paul2013PronunciationContinuityFlag(
			string(cString(rows[rowIndex+1].Surface)),
		) != 0
		if previousUpperInitial || nextUpperInitial {
			selector = 0x13
			break
		}
		if rowIndex+1 >= len(rows) || !paul2013NameContextSurfaceIsOneOf(
			rows[rowIndex+1].Surface, "arrest", "number", "numbers",
		) {
			return selectFallback()
		}
		selector = 0x13
	case Paul2013ContextMappedCStringEqual(current.Surface, []byte("dove")):
		if rowIndex+1 >= len(rows) || !Paul2013ContextMappedCStringEqual(rows[rowIndex+1].Surface, []byte("into")) {
			return selectFallback()
		}
		selector = '&'
	case Paul2013ContextMappedCStringEqual(current.Surface, []byte("elaborate")):
		if rowIndex+1 >= len(rows) || !Paul2013ContextMappedCStringEqual(rows[rowIndex+1].Surface, []byte("into")) {
			return selectFallback()
		}
		selector = '&'
	case paul2013NameContextSurfaceIsOneOf(current.Surface, "transform", "transforms"):
		nextInto := rowIndex+2 < len(rows) && Paul2013ContextMappedCStringEqual(rows[rowIndex+2].Surface, []byte("into"))
		laterInto := rowIndex+3 < len(rows) && Paul2013ContextMappedCStringEqual(rows[rowIndex+3].Surface, []byte("into"))
		if !nextInto && !laterInto {
			return selectFallback()
		}
		if Paul2013ContextMappedCStringEqual(current.Surface, []byte("transforms")) {
			selector = '%'
		} else {
			selector = '*'
		}
	case Paul2013ContextMappedCStringEqual(current.Surface, []byte("use")):
		if rowIndex == 0 || !Paul2013ContextMappedCStringEqual(rows[rowIndex-1].Surface, []byte("still")) {
			return selectFallback()
		}
		selector = '%'
	case Paul2013ContextMappedCStringEqual(current.Surface, []byte("wind")):
		if rowIndex+1 >= len(rows) || !paul2013NameContextSurfaceIsOneOf(
			rows[rowIndex+1].Surface, "up", "down", "through",
		) {
			return selectFallback()
		}
		selector = '%'
	case Paul2013ContextMappedCStringEqual(current.Surface, []byte("de")):
		previousUpperInitial := rowIndex > 0 && Paul2013PronunciationContinuityFlag(
			string(cString(rows[rowIndex-1].Surface)),
		) != 0
		nextUpperInitial := rowIndex+1 < len(rows) && Paul2013PronunciationContinuityFlag(
			string(cString(rows[rowIndex+1].Surface)),
		) != 0
		if !previousUpperInitial && !nextUpperInitial {
			return selectFallback()
		}
		selector = 0x0c
	case Paul2013ContextMappedCStringEqual(current.Surface, []byte("job")):
		if rowIndex+1 >= len(rows) || !Paul2013ContextMappedCStringEqual(rows[rowIndex+1].Surface, []byte("hunting")) {
			return selectFallback()
		}
		selector = 0x13
	case Paul2013ContextMappedCStringEqual(current.Surface, []byte("laden")):
		nextPossessive := rowIndex+1 < len(rows) && Paul2013ContextMappedCStringEqual(rows[rowIndex+1].Surface, []byte("'s"))
		previousBin := rowIndex > 0 && Paul2013ContextMappedCStringEqual(rows[rowIndex-1].Surface, []byte("bin"))
		initialUppercase := rowIndex == 0 && Paul2013PronunciationContinuityFlag(
			string(cString(current.Surface)),
		) != 0
		if nextPossessive || previousBin || initialUppercase {
			selector = 0x14
		} else {
			selector = 0x0e
		}
	case paul2013NameContextSurfaceIsOneOf(current.Surface, "resume", "resumes"):
		if rowIndex < 2 || rowIndex+1 >= len(rows) ||
			!paul2013NameContextIsSentencePunctuation(rows[rowIndex+1].Surface) {
			return selectFallback()
		}
		previousClass3 := Paul2013TripletTableContains(3, rows[rowIndex-1].Surface)
		previousPreviousClass11 := Paul2013TripletTableContains(11, rows[rowIndex-2].Surface)
		if previousClass3 || !previousPreviousClass11 {
			return selectFallback()
		}
		if Paul2013ContextMappedCStringEqual(current.Surface, []byte("resumes")) {
			selector = '%'
		} else {
			selector = '*'
		}
	case paul2013NameContextSurfaceIsOneOf(current.Surface, "separate", "separates"):
		previousBe := rowIndex > 0 && Paul2013ContextMappedCStringEqual(rows[rowIndex-1].Surface, []byte("be"))
		previousClass3 := rowIndex > 0 && Paul2013TripletTableContains(3, rows[rowIndex-1].Surface)
		if previousBe || previousClass3 {
			selector = 0x0e
			break
		}
		nextForward := rowIndex+1 < len(rows) && paul2013NameContextSurfaceIsOneOf(
			rows[rowIndex+1].Surface, "from", "into", "out", "up",
		)
		twoRowsForward := rowIndex+2 < len(rows) && paul2013NameContextSurfaceIsOneOf(
			rows[rowIndex+2].Surface, "from", "into", "and",
		)
		if !nextForward && !twoRowsForward {
			return selectFallback()
		}
		if Paul2013ContextMappedCStringEqual(current.Surface, []byte("separates")) {
			selector = '%'
		} else {
			selector = '*'
		}
	case Paul2013ContextMappedCStringEqual(current.Surface, []byte("subject")):
		if rowIndex+1 >= len(rows) || !Paul2013ContextMappedCStringEqual(rows[rowIndex+1].Surface, []byte("to")) {
			return selectFallback()
		}
		previousClass3 := rowIndex > 0 && Paul2013TripletTableContains(3, rows[rowIndex-1].Surface)
		previousPreviousClass3 := rowIndex > 1 && Paul2013TripletTableContains(3, rows[rowIndex-2].Surface)
		previousJobOrNot := rowIndex > 0 && paul2013NameContextSurfaceIsOneOf(rows[rowIndex-1].Surface, "job", "not")
		if previousClass3 || (previousPreviousClass3 && previousJobOrNot) {
			selector = 0x0e
		} else {
			selector = '%'
		}
	case paul2013NameContextSurfaceIsOneOf(current.Surface, "interstate", "interstates"):
		if rowIndex+1 < len(rows) && Paul2013ContextMappedCStringEqual(rows[rowIndex+1].Surface, []byte("to")) {
			if len(cString(current.Surface)) == len("interstate") {
				selector = '%'
			} else {
				selector = '*'
			}
			break
		}
		if !Paul2013ContextMappedCStringEqual(current.Surface, []byte("interstate")) {
			return selectFallback()
		}
		if rowIndex+1 >= len(rows) ||
			Paul2013PronunciationContinuityFlag(string(cString(current.Surface))) == 0 ||
			Paul2013PronunciationContinuityFlag(string(cString(rows[rowIndex+1].Surface))) == 0 {
			return selectFallback()
		}
		selector = 0x0e
	case Paul2013ContextMappedCStringEqual(current.Surface, []byte("invalid")):
		previousClass3 := rowIndex > 0 && Paul2013TripletTableContains(3, rows[rowIndex-1].Surface)
		previousPreviousClass3 := rowIndex > 1 && Paul2013TripletTableContains(3, rows[rowIndex-2].Surface)
		previousNegator := rowIndex > 0 && paul2013NameContextSurfaceIsOneOf(
			rows[rowIndex-1].Surface, "n't", "not",
		)
		if !previousClass3 && !(previousPreviousClass3 && previousNegator) {
			return selectFallback()
		}
		selector = 0x0e
	default:
		return selectFallback()
	}
	if current.AlternativeCount < 2 {
		return selectFallback()
	}
	alternative, phone, found, err := FindPaul2013PronunciationForPathCode(
		current.PathControlBytes, current.PhoneAlternatives, selector,
	)
	if err != nil {
		return Paul2013NameContextSelection{}, fmt.Errorf("select name/context pronunciation by path code %#02x: %w", selector, err)
	}
	if !found {
		return selectFallback()
	}
	return Paul2013NameContextSelection{
		Applicable: true, SelectorCode: selector,
		Alternative: alternative, PhoneString: phone,
	}, nil
}

func paul2013NameContextAllCharactersHaveAttribute(
	surface []byte,
	attributes [256]byte,
	mask byte,
) bool {
	surface = cString(surface)
	if len(surface) == 0 {
		return false
	}
	for _, character := range surface {
		if attributes[character]&mask == 0 {
			return false
		}
	}
	return true
}

func paul2013NameContextEndsWithSelf(surface []byte) bool {
	surface = cString(surface)
	if len(surface) <= len("self") {
		return false
	}
	return ComparePaul2013MappedCString(
		surface[len(surface)-len("self"):], []byte("self"), paul2013ContextCharacterWeights(),
	) == 0
}

func paul2013NameContextLivePreviousException(
	rows []Paul2013ModelPronunciationRow,
	rowIndex int,
) bool {
	if rowIndex <= 0 || rowIndex >= len(rows) {
		return false
	}
	previous := rows[rowIndex-1].Surface
	if paul2013NameContextSurfaceIsOneOf(previous, "Window", "Windows", "xbox", "the", "a") {
		return true
	}
	return rowIndex > 1 &&
		Paul2013ContextMappedCStringEqual(rows[rowIndex-2].Surface, []byte("Larry")) &&
		Paul2013ContextMappedCStringEqual(previous, []byte("King"))
}

var paul2013DoggedContextPhrases = [...]string{
	"'ve", "had", "had not", "hadn't", "has", "has not", "hasn't", "have", "have not", "haven't",
} // PTR_DAT_100783a4, ten entries; source pointer order retained.

// paul2013ModelContextPhrasePosition ports the bounded phrase scan used by
// FUN_1000cf00 for the left/right row windows. It returns the matched
// single-token index or the final token index of a multi-token phrase, or -1
// when no phrase matches. The dogged caller uses left mode, a three-row window,
// and mapped mode 0x49 comparisons.
func paul2013ModelContextPhrasePosition(
	rows []Paul2013ModelPronunciationRow,
	rowIndex int,
	direction byte,
	windowLimit int,
	phrases []string,
) int {
	if rowIndex < 0 || rowIndex >= len(rows) || windowLimit <= 0 || windowLimit > 20 {
		return -1
	}
	start, end := rowIndex, rowIndex
	switch direction {
	case 'L':
		start = rowIndex - windowLimit
		if start < 0 {
			start = 0
		}
	case 'R':
		start = rowIndex + 1
		end = start + windowLimit
		if end > len(rows) {
			end = len(rows)
		}
	default:
		return -1
	}
	window := rows[start:end]
	for position := range window {
		for _, phrase := range phrases {
			parts := bytes.Split([]byte(phrase), []byte{' '})
			if position+len(parts) > len(window) {
				continue
			}
			matched := true
			for partIndex, part := range parts {
				if !Paul2013ContextMappedCStringEqual(window[position+partIndex].Surface, part) {
					matched = false
					break
				}
			}
			if matched {
				return position + len(parts) - 1
			}
		}
	}
	return -1
}

func selectPaul2013ModelFallbackPathChoice(
	rows []Paul2013ModelPronunciationRow,
	rowIndex int,
) (Paul2013NameContextSelection, bool, error) {
	current := rows[rowIndex]
	selectCodes := func(pathCodes []byte) (Paul2013NameContextSelection, bool, error) {
		for _, pathCode := range pathCodes {
			alternative, phone, found, err := FindPaul2013PronunciationForPathCode(
				current.PathControlBytes, current.PhoneAlternatives, pathCode,
			)
			if err != nil {
				return Paul2013NameContextSelection{}, false, fmt.Errorf("select fallback pronunciation by path code %#02x: %w", pathCode, err)
			}
			if found {
				return Paul2013NameContextSelection{
					Applicable: true, SelectorCode: pathCode,
					Alternative: alternative, PhoneString: phone,
				}, true, nil
			}
		}
		return Paul2013NameContextSelection{}, false, nil
	}
	selectCode := func(pathCode byte) (Paul2013NameContextSelection, bool, error) {
		return selectCodes([]byte{pathCode})
	}
	if rowIndex+1 < len(rows) {
		next := rows[rowIndex+1].Surface
		if !Paul2013ContextMappedCStringEqual(current.Surface, []byte("nice")) &&
			Paul2013ContextMappedCStringEqual(next, []byte("of")) {
			selection, found, err := selectCodes([]byte{0x13, 0x16, 0x14, 0x15})
			if err != nil || found {
				return selection, found, err
			}
		}
		if Paul2013ContextMappedCStringEqual(next, []byte("'s")) {
			selection, found, err := selectCodes([]byte{0x14, 0x15, 0x13, 0x16})
			if err != nil || found {
				return selection, found, err
			}
		}
	}
	if rowIndex+1 < len(rows) &&
		(rowIndex == 0 || rows[rowIndex-1].Index != current.Index) &&
		paul2013NameContextSurfaceIsOneOf(current.Surface, "de", "inter", "re") &&
		Paul2013ContextMappedCStringEqual(rows[rowIndex+1].Surface, []byte("-")) &&
		current.Index == rows[rowIndex+1].Index {
		selection, found, err := selectCode(0x0d)
		if err != nil {
			return selection, false, err
		}
		if !found {
			return Paul2013NameContextSelection{}, false, fmt.Errorf("compound-prefix fallback row %d has no 0x0d path", rowIndex)
		}
		if found {
			return selection, found, err
		}
	}
	if rowIndex > 0 && paul2013NameContextIsFUN10010700Word(rows[rowIndex-1].Surface) {
		selection, found, err := selectPaul2013ModelFallbackPathClass(current, 2)
		if err != nil {
			return selection, false, err
		}
		if found {
			return selection, found, err
		}
		selection, found, err = selectPaul2013ModelFallbackPathClass(current, 8)
		if err != nil {
			return selection, false, err
		}
		if !found {
			return Paul2013NameContextSelection{}, false, fmt.Errorf("FUN_10010700 fallback row %d has neither class-2 nor class-8 path", rowIndex)
		}
		return selection, true, nil
	}
	if rowIndex > 0 && paul2013NameContextIsFUN10010730Word(rows[rowIndex-1].Surface) {
		if rowIndex > 2 {
			feature, known := Paul2013KnownPronunciationTripletFeature(
				string(rows[rowIndex-3].Surface), string(rows[rowIndex-2].Surface), "",
			)
			if known && feature == 3 {
				selection, found, err := selectCode('(')
				if err != nil || found {
					return selection, found, err
				}
			}
		}
		selection, found, err := selectCodes([]byte{'%', '*', '&'})
		if err != nil {
			return selection, false, err
		}
		if !found {
			return Paul2013NameContextSelection{}, false, fmt.Errorf("FUN_10010730 fallback row %d has no percent, asterisk, or ampersand path", rowIndex)
		}
		if found {
			return selection, found, err
		}
	}
	if rowIndex > 0 && paul2013NameContextIsFUN10010690Word(rows[rowIndex-1].Surface) {
		otherClass56, err := paul2013NameContextHasClassOutsideSelected(current, []int16{4, 7}, []int16{5, 6})
		if err != nil {
			return Paul2013NameContextSelection{}, false, err
		}
		if otherClass56 {
			return Paul2013NameContextSelection{}, false, fmt.Errorf("FUN_10010690 fallback row %d has an ambiguous class-5/6 alternative", rowIndex)
		}
		selection, found, err := selectPaul2013ModelFallbackPathClasses(current, []int16{4, 7})
		if err != nil {
			return selection, false, err
		}
		if found {
			return selection, true, nil
		}
		selection, found, err = selectCode(0x13)
		if err != nil {
			return selection, false, err
		}
		if !found {
			return Paul2013NameContextSelection{}, false, fmt.Errorf("FUN_10010690 fallback row %d has no class-4/7 or 0x13 path", rowIndex)
		}
		return selection, true, nil
	}
	if rowIndex > 1 {
		feature, known := Paul2013KnownPronunciationTripletFeature(
			string(rows[rowIndex-2].Surface), string(rows[rowIndex-1].Surface), string(current.Surface),
		)
		if known && feature == 7 {
			otherClass56, err := paul2013NameContextHasClassOutsideSelected(current, []int16{4, 7}, []int16{5, 6})
			if err != nil {
				return Paul2013NameContextSelection{}, false, err
			}
			if otherClass56 {
				return Paul2013NameContextSelection{}, false, fmt.Errorf("class-7 triplet fallback row %d has an ambiguous class-5/6 alternative", rowIndex)
			}
			selection, found, err := selectPaul2013ModelFallbackPathClasses(current, []int16{4, 7})
			if err != nil {
				return selection, false, err
			}
			if found {
				return selection, true, nil
			}
			selection, found, err = selectCode(0x13)
			if err != nil {
				return selection, false, err
			}
			if !found {
				return Paul2013NameContextSelection{}, false, fmt.Errorf("class-7 triplet fallback row %d has no class-4/7 or 0x13 path", rowIndex)
			}
			return selection, true, nil
		}
	}
	if rowIndex+1 < len(rows) && paul2013NameContextIsFUN10010790Word(rows[rowIndex+1].Surface) {
		otherClass10, err := paul2013NameContextHasClass10OutsideClass9Or12(current)
		if err != nil {
			return Paul2013NameContextSelection{}, false, err
		}
		if otherClass10 {
			return Paul2013NameContextSelection{}, false, fmt.Errorf("FUN_10010790 fallback row %d has an ambiguous class-10 alternative", rowIndex)
		}
		selection, found, err := selectCode('%')
		if err != nil {
			return selection, false, err
		}
		if !found {
			return Paul2013NameContextSelection{}, false, fmt.Errorf("FUN_10010790 fallback row %d has no percent path", rowIndex)
		}
		return selection, true, nil
	}
	if rowIndex > 0 {
		previous := rows[rowIndex-1].Surface
		if paul2013NameContextSurfaceIsOneOf(previous, "have", "has", "had") {
			selection, found, err := selectCode('(')
			if err != nil {
				return selection, false, err
			}
			if found {
				return selection, true, nil
			}
		}
		if paul2013NameContextFUN10010840(previous, nil) {
			selection, found, err := selectCode('%')
			if err != nil {
				return selection, false, err
			}
			if !found {
				return Paul2013NameContextSelection{}, false, fmt.Errorf("FUN_10010840 fallback row %d has no percent path", rowIndex)
			}
			return selection, true, nil
		}
		previousPairSpecial := rowIndex > 1 &&
			paul2013NameContextFUN10010840(rows[rowIndex-2].Surface, previous)
		if previousPairSpecial {
			selection, found, err := selectCode('%')
			if err != nil {
				return selection, false, err
			}
			if !found {
				return Paul2013NameContextSelection{}, false, fmt.Errorf("FUN_10010840 pair fallback row %d has no percent path", rowIndex)
			}
			return selection, true, nil
		}
		if (rowIndex < 2 || !previousPairSpecial) &&
			paul2013NameContextIsFUN10010760Word(previous) {
			surface := cString(current.Surface)
			if len(surface) > 5 && Paul2013ContextMappedCStringEqual(surface[len(surface)-3:], []byte("ing")) {
				selection, found, err := selectCode('\'')
				if err != nil {
					return selection, false, err
				}
				if !found {
					return Paul2013NameContextSelection{}, false, fmt.Errorf("FUN_10010760 fallback row %d has no apostrophe path", rowIndex)
				}
				return selection, true, nil
			}
		}
		if rowIndex < 2 || !previousPairSpecial {
			previousHow := Paul2013ContextMappedCStringEqual(previous, []byte("how"))
			previousPreviousHow := rowIndex > 1 &&
				Paul2013ContextMappedCStringEqual(rows[rowIndex-2].Surface, []byte("how"))
			if previousHow || previousPreviousHow {
				if rowIndex+1 >= len(rows) || !Paul2013TripletTableContains(8, current.Surface) {
					return Paul2013NameContextSelection{}, false, fmt.Errorf("how-context fallback row %d requires a following row and current class-8 surface", rowIndex)
				}
				selection, found, err := selectCode(0x0e)
				if err != nil {
					return selection, false, err
				}
				if !found {
					return Paul2013NameContextSelection{}, false, fmt.Errorf("how-context fallback row %d has no 0x0e path", rowIndex)
				}
				return selection, true, nil
			}
		}
	}
	return Paul2013NameContextSelection{}, false, nil
}

func paul2013NameContextFUN10010840(left, right []byte) bool {
	if len(cString(left)) == 0 {
		return false
	}
	if Paul2013ContextMappedCStringEqual(left, []byte("let")) || Paul2013ContextWordPairSpecial(left, right) {
		return true
	}
	if Paul2013ContextWordPairSpecial(left, nil) && Paul2013ContextMappedCStringEqual(right, []byte("not")) {
		return true
	}
	if len(cString(right)) != 0 && Paul2013ContextWordPairSpecial(right, nil) {
		return true
	}
	return Paul2013ContextMappedCStringEqual(right, []byte("do"))
}

func paul2013NameContextIsFUN10010760Word(surface []byte) bool {
	return paul2013NameContextSurfaceIsOneOf(
		surface, "begin", "continue", "enjoy", "finish", "give up", "hate",
		"keep on", "love", "mind", "start", "stop",
	)
}

func paul2013NameContextHasClass10OutsideClass9Or12(row Paul2013ModelPronunciationRow) (bool, error) {
	groups, err := ParsePaul2013PronunciationPath(row.PathControlBytes)
	if err != nil {
		return false, fmt.Errorf("parse FUN_10010900 path classes: %w", err)
	}
	if len(groups) != len(row.PhoneAlternatives) {
		return false, fmt.Errorf("FUN_10010900 path stream has %d alternatives for %d phone strings", len(groups), len(row.PhoneAlternatives))
	}
	class9Or12Alternative := -1
	for alternative, group := range groups {
		for _, pathCode := range group {
			if pathCode == paul2013PathTerminator {
				break
			}
			if int(pathCode) >= len(paul2013PathCodeClasses) {
				return false, fmt.Errorf("FUN_10010900 alternative %d has path code 0x%02x outside the recovered table", alternative, pathCode)
			}
			class := paul2013PathCodeClasses[pathCode]
			if class == 9 || class == 12 {
				class9Or12Alternative = alternative
				break
			}
		}
		if class9Or12Alternative >= 0 {
			break
		}
	}
	for alternative, group := range groups {
		for _, pathCode := range group {
			if pathCode == paul2013PathTerminator {
				break
			}
			if int(pathCode) >= len(paul2013PathCodeClasses) {
				return false, fmt.Errorf("FUN_10010900 alternative %d has path code 0x%02x outside the recovered table", alternative, pathCode)
			}
			if paul2013PathCodeClasses[pathCode] == 10 && alternative != class9Or12Alternative {
				return true, nil
			}
		}
	}
	return false, nil
}

func paul2013NameContextHasClassOutsideSelected(
	row Paul2013ModelPronunciationRow,
	selectedClasses []int16,
	guardClasses []int16,
) (bool, error) {
	groups, err := ParsePaul2013PronunciationPath(row.PathControlBytes)
	if err != nil {
		return false, fmt.Errorf("parse fallback cross-class path rows: %w", err)
	}
	if len(groups) != len(row.PhoneAlternatives) {
		return false, fmt.Errorf("fallback path stream has %d alternatives for %d phone strings", len(groups), len(row.PhoneAlternatives))
	}
	selectedAlternative := -1
	for alternative, group := range groups {
		for _, pathCode := range group {
			if pathCode == paul2013PathTerminator {
				break
			}
			if int(pathCode) >= len(paul2013PathCodeClasses) {
				return false, fmt.Errorf("fallback alternative %d has path code 0x%02x outside the recovered table", alternative, pathCode)
			}
			if paul2013ContainsClass(selectedClasses, paul2013PathCodeClasses[pathCode]) {
				selectedAlternative = alternative
				break
			}
		}
		if selectedAlternative >= 0 {
			break
		}
	}
	for alternative, group := range groups {
		for _, pathCode := range group {
			if pathCode == paul2013PathTerminator {
				break
			}
			if int(pathCode) >= len(paul2013PathCodeClasses) {
				return false, fmt.Errorf("fallback alternative %d has path code 0x%02x outside the recovered table", alternative, pathCode)
			}
			if alternative != selectedAlternative && paul2013ContainsClass(guardClasses, paul2013PathCodeClasses[pathCode]) {
				return true, nil
			}
		}
	}
	return false, nil
}

func selectPaul2013ModelFallbackPathClasses(
	row Paul2013ModelPronunciationRow,
	wantedClasses []int16,
) (Paul2013NameContextSelection, bool, error) {
	groups, err := ParsePaul2013PronunciationPath(row.PathControlBytes)
	if err != nil {
		return Paul2013NameContextSelection{}, false, fmt.Errorf("parse fallback path classes: %w", err)
	}
	if len(groups) != len(row.PhoneAlternatives) {
		return Paul2013NameContextSelection{}, false, fmt.Errorf("fallback path stream has %d alternatives for %d phone strings", len(groups), len(row.PhoneAlternatives))
	}
	for alternative, group := range groups {
		for _, pathCode := range group {
			if pathCode == paul2013PathTerminator {
				break
			}
			if int(pathCode) >= len(paul2013PathCodeClasses) {
				return Paul2013NameContextSelection{}, false, fmt.Errorf("fallback alternative %d has path code 0x%02x outside the recovered table", alternative, pathCode)
			}
			if paul2013ContainsClass(wantedClasses, paul2013PathCodeClasses[pathCode]) {
				return Paul2013NameContextSelection{
					Applicable: true, SelectorCode: pathCode,
					Alternative: alternative,
					PhoneString: append([]byte(nil), paul2013PhoneRowBytes(row.PhoneAlternatives[alternative])...),
				}, true, nil
			}
		}
	}
	return Paul2013NameContextSelection{}, false, nil
}

func paul2013ContainsClass(classes []int16, candidate int16) bool {
	for _, class := range classes {
		if class == candidate {
			return true
		}
	}
	return false
}

func paul2013NameContextIsFUN10010690Word(surface []byte) bool {
	if Paul2013ContextMappedCStringEqual(surface, []byte("of")) ||
		Paul2013PronunciationPhoneContextFeature(string(surface)) == 3 ||
		Paul2013TripletClass12Contains(surface) {
		return true
	}
	for _, word := range [...]string{"each", "her", "his", "in", "its", "my", "our", "their", "your"} {
		if Paul2013ContextMappedCStringEqual(surface, []byte(word)) {
			return true
		}
	}
	return false
}

func paul2013NameContextIsFUN10010790Word(surface []byte) bool {
	if Paul2013ContextMappedCStringEqual(surface, []byte("it")) ||
		Paul2013TripletClass12Contains(surface) ||
		Paul2013TripletTableContains(7, surface) {
		return true
	}
	for _, word := range [...]string{"him", "me", "them", "us"} {
		if Paul2013ContextMappedCStringEqual(surface, []byte(word)) {
			return true
		}
	}
	return false
}

func selectPaul2013ModelFallbackPathClass(
	row Paul2013ModelPronunciationRow,
	wantedClass int16,
) (Paul2013NameContextSelection, bool, error) {
	groups, err := ParsePaul2013PronunciationPath(row.PathControlBytes)
	if err != nil {
		return Paul2013NameContextSelection{}, false, fmt.Errorf("parse fallback path classes: %w", err)
	}
	if len(groups) != len(row.PhoneAlternatives) {
		return Paul2013NameContextSelection{}, false, fmt.Errorf(
			"fallback path stream has %d alternatives for %d phone strings",
			len(groups), len(row.PhoneAlternatives),
		)
	}
	for alternative, group := range groups {
		for _, pathCode := range group {
			if pathCode == paul2013PathTerminator {
				break
			}
			if int(pathCode) >= len(paul2013PathCodeClasses) {
				return Paul2013NameContextSelection{}, false, fmt.Errorf(
					"fallback alternative %d has path code 0x%02x outside the recovered table",
					alternative, pathCode,
				)
			}
			if paul2013PathCodeClasses[pathCode] == wantedClass {
				return Paul2013NameContextSelection{
					Applicable: true, SelectorCode: pathCode,
					Alternative: alternative,
					PhoneString: append([]byte(nil), paul2013PhoneRowBytes(row.PhoneAlternatives[alternative])...),
				}, true, nil
			}
		}
	}
	return Paul2013NameContextSelection{}, false, nil
}

func paul2013NameContextIsFUN10010700Word(surface []byte) bool {
	for _, word := range [...]string{"less", "more", "so", "very"} {
		if Paul2013ContextMappedCStringEqual(surface, []byte(word)) {
			return true
		}
	}
	return false
}

func paul2013NameContextIsFUN10010730Word(surface []byte) bool {
	for _, word := range [...]string{
		"almost", "also", "always", "ever", "frequently", "generally", "hardly",
		"just", "never", "occasionally", "often", "only", "quite", "rarely",
		"really", "scarcely", "seldom", "sometimes", "still", "usually",
	} {
		if Paul2013ContextMappedCStringEqual(surface, []byte(word)) {
			return true
		}
	}
	return false
}

func paul2013NameContextSurfaceIsOneOf(surface []byte, options ...string) bool {
	for _, option := range options {
		if Paul2013ContextMappedCStringEqual(surface, []byte(option)) {
			return true
		}
	}
	return false
}

func paul2013NameContextIsSentencePunctuation(surface []byte) bool {
	return paul2013NameContextSurfaceIsOneOf(surface, ".", ",", "!", "?")
}

func paul2013NameContextIsClass12OrContextWord(surface []byte) bool {
	if Paul2013TripletClass12Contains(surface) {
		return true
	}
	return paul2013NameContextIsContextWord(surface)
}

func paul2013NameContextIsContextWord(surface []byte) bool {
	for _, word := range paul2013PronunciationContextWords {
		if Paul2013ContextMappedCStringEqual(surface, []byte(word)) {
			return true
		}
	}
	return false
}

func selectPaul2013ModelCloseContext(rows []Paul2013ModelPronunciationRow, rowIndex int) (byte, error) {
	if rowIndex > 0 {
		previous := rows[rowIndex-1].Surface
		switch {
		case paul2013NameContextSurfaceIsOneOf(previous, "I", "you", "we", "they"):
			return '%', nil
		case Paul2013ContextMappedCStringEqual(previous, []byte("feel")):
			return 0x0e, nil
		case paul2013NameContextSurfaceIsOneOf(previous, "fisherman", "fishermans"):
			return 0x13, nil
		}
	}
	if rowIndex >= 2 &&
		Paul2013ContextMappedCStringEqual(rows[rowIndex-2].Surface, []byte("to")) &&
		Paul2013ContextMappedCStringEqual(rows[rowIndex-1].Surface, []byte("a")) {
		return '%', nil
	}
	if rowIndex+1 >= len(rows) {
		return 0, fmt.Errorf("close name/context row %d is outside the directly ported cases; other FUN_100049b0 context cases remain unsupported", rowIndex)
	}
	next := rows[rowIndex+1].Surface
	if paul2013NameContextSurfaceIsOneOf(next, "than", "by", "on", "to", "upon", "call", "calls") {
		return 0x0e, nil
	}
	if Paul2013ContextMappedCStringEqual(next, []byte("at")) && rowIndex+2 < len(rows) &&
		Paul2013ContextMappedCStringEqual(rows[rowIndex+2].Surface, []byte("hand")) {
		return 0x0e, nil
	}
	if Paul2013ContextMappedCStringEqual(next, []byte("in")) {
		if rowIndex+2 < len(rows) && Paul2013ContextMappedCStringEqual(rows[rowIndex+2].Surface, []byte("with")) {
			return 0x0e, nil
		}
	}
	if Paul2013ContextMappedCStringEqual(next, []byte("it")) && rowIndex+2 < len(rows) &&
		Paul2013ContextMappedCStringEqual(rows[rowIndex+2].Surface, []byte("up")) {
		return '%', nil
	}
	if Paul2013ContextMappedCStringEqual(next, []byte("down")) && rowIndex+2 < len(rows) &&
		Paul2013ContextMappedCStringEqual(rows[rowIndex+2].Surface, []byte("on")) {
		return '%', nil
	}
	return 0, fmt.Errorf("close name/context row %d is outside the directly ported cases; other FUN_100049b0 context cases remain unsupported", rowIndex)
}

func selectPaul2013ModelCloserContext(rows []Paul2013ModelPronunciationRow, rowIndex int) (byte, error) {
	if rowIndex+1 < len(rows) {
		next := rows[rowIndex+1].Surface
		if paul2013NameContextSurfaceIsOneOf(next, "than", "by", "on", "to", "upon", "call", "calls", "look", "looks") {
			return 0x0f, nil
		}
		if Paul2013PronunciationPhoneContextFeature(string(next)) == 2 {
			return 0x0f, nil
		}
		if rowIndex+2 < len(rows) {
			nextNext := rows[rowIndex+2].Surface
			if (Paul2013ContextMappedCStringEqual(next, []byte("at")) &&
				Paul2013ContextMappedCStringEqual(nextNext, []byte("hand"))) ||
				(Paul2013ContextMappedCStringEqual(next, []byte("in")) &&
					Paul2013ContextMappedCStringEqual(nextNext, []byte("with"))) {
				return 0x0f, nil
			}
		}
	}
	return 0, fmt.Errorf("closer name/context row %d is outside the directly ported cases; other FUN_100049b0 context cases remain unsupported", rowIndex)
}

// ReadPaul2013ModelPronunciationRows decodes the pronunciation fields from a
// contiguous token-row arena. The source-length prefix is not part of rows;
// callers pass the bytes beginning at model offset +2, as FUN_100068b0 does.
// It validates all surface, path, and phone-string terminators before
// returning a detached snapshot.
func ReadPaul2013ModelPronunciationRows(tokenRows []byte) ([]Paul2013ModelPronunciationRow, error) {
	if len(tokenRows)%Paul2013TokenResultRowSize != 0 {
		return nil, fmt.Errorf("model token rows have %d bytes, not a multiple of 0x%x", len(tokenRows), Paul2013TokenResultRowSize)
	}
	rowCount := len(tokenRows) / Paul2013TokenResultRowSize
	if rowCount > paul2013ModelTokenRowLimit {
		return nil, fmt.Errorf("model token rows contain %d rows, native limit is %d", rowCount, paul2013ModelTokenRowLimit)
	}
	rows := make([]Paul2013ModelPronunciationRow, rowCount)
	for rowIndex := range rows {
		start := rowIndex * Paul2013TokenResultRowSize
		row := tokenRows[start : start+Paul2013TokenResultRowSize]
		surfaceArea := row[Paul2013TokenResultRowSurface:Paul2013TokenResultPathControls]
		surfaceEnd := bytes.IndexByte(surfaceArea, 0)
		if surfaceEnd < 0 {
			return nil, fmt.Errorf("model token row %d has no NUL-terminated surface", rowIndex)
		}
		alternativeCount := int(int16(binary.LittleEndian.Uint16(row[Paul2013TokenResultRowCount:])))
		if alternativeCount < 0 || alternativeCount > paul2013TokenResultAlternativeCapacity {
			return nil, fmt.Errorf("model token row %d has unsupported alternative count %d", rowIndex, alternativeCount)
		}
		decoded := Paul2013ModelPronunciationRow{
			Index:            binary.LittleEndian.Uint16(row[Paul2013TokenResultRowIndex:]),
			Type:             row[Paul2013TokenResultRowType],
			Surface:          append([]byte(nil), surfaceArea[:surfaceEnd]...),
			AlternativeCount: alternativeCount,
		}
		if alternativeCount >= 2 {
			pathArea := row[Paul2013TokenResultPathControls:Paul2013TokenResultPhoneStrings]
			pathEnd := bytes.IndexByte(pathArea, paul2013PathTerminator)
			if pathEnd < 0 {
				return nil, fmt.Errorf("model token row %d has no 0xff path terminator", rowIndex)
			}
			decoded.PathControlBytes = append([]byte(nil), pathArea[:pathEnd+1]...)
		}
		decoded.PhoneAlternatives = make([][]byte, alternativeCount)
		for alternativeIndex := range decoded.PhoneAlternatives {
			phoneStart := Paul2013TokenResultPhoneStrings + alternativeIndex*Paul2013TokenResultPhoneStride
			phoneArea := row[phoneStart : phoneStart+Paul2013TokenResultPhoneStride]
			phoneEnd := bytes.IndexByte(phoneArea, 0)
			if phoneEnd < 0 {
				return nil, fmt.Errorf("model token row %d phone alternative %d has no NUL terminator", rowIndex, alternativeIndex)
			}
			decoded.PhoneAlternatives[alternativeIndex] = append([]byte(nil), phoneArea[:phoneEnd]...)
		}
		rows[rowIndex] = decoded
	}
	return rows, nil
}
