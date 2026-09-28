package selection

import (
	"context"
	"errors"
	"fmt"

	"vtspeak/engine/dat"
	"vtspeak/engine/voice"
)

// CandidateRecord pairs a candidate reference with its indexed unit metadata.
// The DAT waveform and UPM timing payload are not read.
type CandidateRecord struct {
	Unit   UnitRef
	Record dat.UnitRecord
}

// ReadCandidateRecords loads metadata for candidates in source order without
// reading or decoding waveform data. The result preserves duplicate
// references if they occur in the input.
func ReadCandidateRecords(ctx context.Context, model *voice.Paul2013, candidates []UnitRef) ([]CandidateRecord, error) {
	if model == nil {
		return nil, errors.New("candidate record lookup has no Paul model")
	}
	result := make([]CandidateRecord, 0, len(candidates))
	for index, candidate := range candidates {
		if err := ctx.Err(); err != nil {
			return nil, err
		}
		bank, ok := model.Banks[candidate.Bank]
		if !ok || bank == nil {
			return nil, fmt.Errorf("candidate %d references unavailable bank %q", index, candidate.Bank)
		}
		unitRecord, err := bank.ReadRecord(candidate.Index)
		if err != nil {
			return nil, fmt.Errorf("read candidate %d (%s:%d): %w", index, candidate.Bank, candidate.Index, err)
		}
		result = append(result, CandidateRecord{Unit: candidate, Record: unitRecord})
	}
	return result, nil
}
