package text

import "vtspeak/engine/tree3"

// BuildPaul2013ModelPhoneRows composes the recovered dictionary token-row
// writer and FUN_1000ea20 context-row projection for already-resolved tokens.
// The caller supplies each token's parser-row association and the complete
// 0x94-byte parser rows. PronunciationTree is used for ordinary ambiguous
// rows; unresolved FUN_100049b0 name/context cases remain outside this helper.
// Source parsing, FUN_10007520 row processing, and TPP updates are not run.
func BuildPaul2013ModelPhoneRows(
	model []byte,
	tokens []LexicalToken,
	parserRows []byte,
	tokenParserIndexes []uint16,
	sourceLength int,
	pronunciationTree *tree3.Tree,
) ([]byte, error) {
	tokenRows, err := BuildPaul2013ModelTokenRows(model, tokens, parserRows, tokenParserIndexes)
	if err != nil {
		return nil, err
	}
	selector := SelectPaul2013ModelPronunciationAlternative(pronunciationTree)
	return NormalizePaul2013ModelPhoneRows(tokenRows, parserRows, sourceLength, selector)
}
