package duration

import (
	"errors"

	"vtspeak/engine/text"
)

// ResolvePaul2013PronunciationExceptionFromModelRows evaluates the full-model
// FUN_10007520 gate, extracts the caller-selected exception surface sequence,
// performs the loaded exception lookup and phone decoding, then applies a
// match through FUN_1000ca50's model-row writes. Candidate rows must be the
// contiguous sequence beginning at gateRowIndex; upstream row mutations still
// determine which sequence native code reaches.
func (engine *Engine) ResolvePaul2013PronunciationExceptionFromModelRows(
	parserRows []byte,
	model []byte,
	gateRowIndex int,
	candidateRowIndexes []int,
) (
	text.Paul2013ModelContextExceptionGate,
	[]byte,
	text.ExceptionMatch,
	[][]text.CMUPhone,
	bool,
	error,
) {
	if engine == nil || engine.exceptions == nil {
		return text.Paul2013ModelContextExceptionGate{}, nil, text.ExceptionMatch{}, nil, false,
			errors.New("Paul 2013 pronunciation exceptions are nil or unloaded")
	}
	gate, err := text.EvaluatePaul2013ModelContextExceptionGate(model, parserRows, gateRowIndex)
	if err != nil {
		return text.Paul2013ModelContextExceptionGate{}, nil, text.ExceptionMatch{}, nil, false, err
	}
	if !gate.Dispatch {
		return gate, append([]byte(nil), model...), text.ExceptionMatch{}, nil, false, nil
	}
	if len(candidateRowIndexes) == 0 || candidateRowIndexes[0] != gateRowIndex {
		return gate, nil, text.ExceptionMatch{}, nil, false,
			errors.New("exception candidate rows must begin at the gate row")
	}
	for offset, rowIndex := range candidateRowIndexes {
		if rowIndex != gateRowIndex+offset {
			return gate, nil, text.ExceptionMatch{}, nil, false,
				errors.New("exception candidate rows must be contiguous and ordered")
		}
	}
	inputs, err := text.ExtractPaul2013ModelExceptionRowInputs(model, candidateRowIndexes)
	if err != nil {
		return gate, nil, text.ExceptionMatch{}, nil, false, err
	}
	surfaces := make([]string, len(inputs))
	retryHPrefix := make([]bool, len(inputs))
	for index, input := range inputs {
		surfaces[index] = input.Surface
		retryHPrefix[index] = input.RetryHPrefix
	}
	match, phones, found, err := engine.ResolvePaul2013PronunciationExceptionWithDispatchGate(
		gate.Input, surfaces, retryHPrefix,
	)
	if err != nil {
		return gate, nil, text.ExceptionMatch{}, nil, false, err
	}
	if !found {
		return gate, append([]byte(nil), model...), match, phones, false, nil
	}
	updatedModel, _, err := text.ApplyPaul2013ExceptionMatchToModelRows(model, gateRowIndex, match)
	if err != nil {
		return gate, nil, text.ExceptionMatch{}, nil, false, err
	}
	return gate, updatedModel, match, phones, true, nil
}

// ResolvePaul2013PronunciationExceptionFromModelContextRows derives the
// contiguous candidate prefix from the model arena using FUN_10008dc0's
// normalized one-to-four component bound, then composes the dispatch gate,
// X-prefix retry, phone decoding, and FUN_1000ca50 writes. This is the
// caller-facing path when native candidate rows are not already available.
func (engine *Engine) ResolvePaul2013PronunciationExceptionFromModelContextRows(
	parserRows []byte,
	model []byte,
	gateRowIndex int,
) (
	text.Paul2013ModelContextExceptionGate,
	[]byte,
	text.ExceptionMatch,
	[][]text.CMUPhone,
	bool,
	error,
) {
	if engine == nil || engine.exceptions == nil {
		return text.Paul2013ModelContextExceptionGate{}, nil, text.ExceptionMatch{}, nil, false,
			errors.New("Paul 2013 pronunciation exceptions are nil or unloaded")
	}
	gate, err := text.EvaluatePaul2013ModelContextExceptionGate(model, parserRows, gateRowIndex)
	if err != nil {
		return text.Paul2013ModelContextExceptionGate{}, nil, text.ExceptionMatch{}, nil, false, err
	}
	if !gate.Dispatch {
		return gate, append([]byte(nil), model...), text.ExceptionMatch{}, nil, false, nil
	}
	candidates, err := text.BuildPaul2013NativeExceptionCandidateRows(model, gateRowIndex)
	if err != nil {
		return gate, nil, text.ExceptionMatch{}, nil, false, err
	}
	if len(candidates) == 0 {
		return gate, append([]byte(nil), model...), text.ExceptionMatch{}, nil, false, nil
	}
	indexes := make([]int, len(candidates))
	for index, candidate := range candidates {
		indexes[index] = candidate.RowIndex
	}
	return engine.ResolvePaul2013PronunciationExceptionFromModelRows(
		parserRows, model, gateRowIndex, indexes,
	)
}
