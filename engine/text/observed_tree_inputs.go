package text

import "errors"

// Paul2013DurationTreeMetadata holds the five non-phone values read by the
// nine-short duration input builder. The text frontend derives these from
// FUN_10013c00 rows and FUN_100135d0 boundary state; this low-level type keeps
// the recovered values explicit for trace assembly.
type Paul2013DurationTreeMetadata struct {
	RecordByte1   int8
	RecordByte2   int8
	PositionState uint8
	ContextCount  uint8
	AuxiliaryByte int8
}

// Paul2013DurationTreeNeighbor supplies either a decoded phone or an explicit
// boundary identity ordinal for a low-level duration-vector assembly. The
// helper does not model the state markers used by FUN_100135d0 to choose the
// legacy boundary value.
type Paul2013DurationTreeNeighbor struct {
	Phone                   *CMUPhone
	BoundaryIdentityOrdinal int16
}

func (neighbor Paul2013DurationTreeNeighbor) identityOrdinal(name string) (int16, error) {
	if neighbor.Phone != nil {
		if neighbor.BoundaryIdentityOrdinal != 0 {
			return 0, errors.New(name + " duration neighbor must contain a phone or boundary ordinal, not both")
		}
		features, err := neighbor.Phone.TreeFeatures()
		if err != nil {
			return 0, err
		}
		return features.IdentityOrdinal, nil
	}
	if neighbor.BoundaryIdentityOrdinal <= 0 {
		return 0, errors.New(name + " duration neighbor requires a phone or positive boundary ordinal")
	}
	return neighbor.BoundaryIdentityOrdinal, nil
}

// BuildObservedPaul2013DurationTreeInput builds all nine shorts when supplied
// the current phone, its neighbors or boundary sentinels, and the five
// non-phone values. It is a low-level helper; the text frontend derives these
// inputs from the phone sequence and row producer.
func BuildObservedPaul2013DurationTreeInput(
	previous Paul2013DurationTreeNeighbor,
	current CMUPhone,
	next Paul2013DurationTreeNeighbor,
	metadata Paul2013DurationTreeMetadata,
) ([9]int16, error) {
	var input [9]int16
	if metadata.PositionState < 1 || metadata.PositionState > 3 {
		return input, errors.New("duration position state must be 1, 2, or 3")
	}
	previousOrdinal, err := previous.identityOrdinal("previous")
	if err != nil {
		return input, err
	}
	currentFeatures, err := current.TreeFeatures()
	if err != nil {
		return input, err
	}
	nextOrdinal, err := next.identityOrdinal("next")
	if err != nil {
		return input, err
	}
	input[0] = currentFeatures.IdentityOrdinal
	input[1] = previousOrdinal
	input[2] = nextOrdinal
	input[3] = currentFeatures.StressValue
	input[4] = int16(metadata.RecordByte1)
	input[5] = int16(metadata.RecordByte2)
	input[6] = int16(metadata.PositionState)
	input[7] = int16(metadata.ContextCount)
	input[8] = int16(metadata.AuxiliaryByte)
	return input, nil
}

// BuildObservedPaul2013InteriorDurationTreeInput is the concrete-phone
// convenience form of BuildObservedPaul2013DurationTreeInput.
func BuildObservedPaul2013InteriorDurationTreeInput(
	previous, current, next CMUPhone,
	metadata Paul2013DurationTreeMetadata,
) ([9]int16, error) {
	return BuildObservedPaul2013DurationTreeInput(
		Paul2013DurationTreeNeighbor{Phone: &previous},
		current,
		Paul2013DurationTreeNeighbor{Phone: &next},
		metadata,
	)
}

// ApplyObservedPaul2013OnsetIdentityFeature writes the controlled consonant's
// identity ordinal to position 0 of the selected onset-scalar input. It does
// not select the consonant's tree family or fill the other input fields.
func ApplyObservedPaul2013OnsetIdentityFeature(input []int16, phone CMUPhone) error {
	if len(input) < 1 {
		return errors.New("onset-scalar tree input must contain at least one value")
	}
	if phone.Vowel {
		return errors.New("observed onset identity field is established for consonants")
	}
	features, err := phone.TreeFeatures()
	if err != nil {
		return err
	}
	input[0] = features.IdentityOrdinal
	return nil
}

// ApplyObservedPaul2013InteriorDurationPhoneFeatures fills the directly
// mapped phone-class fields for an interior duration-tree input: current,
// previous, and next identity at positions 0–2, then current stress at
// position 3. Boundary sentinels and positions 4–8 remain caller-owned.
func ApplyObservedPaul2013InteriorDurationPhoneFeatures(
	input []int16,
	previous, current, next CMUPhone,
) error {
	if len(input) < 4 {
		return errors.New("duration tree input must contain at least four values")
	}
	previousFeatures, err := previous.TreeFeatures()
	if err != nil {
		return err
	}
	currentFeatures, err := current.TreeFeatures()
	if err != nil {
		return err
	}
	nextFeatures, err := next.TreeFeatures()
	if err != nil {
		return err
	}
	input[0] = currentFeatures.IdentityOrdinal
	input[1] = previousFeatures.IdentityOrdinal
	input[2] = nextFeatures.IdentityOrdinal
	input[3] = currentFeatures.StressValue
	return nil
}

// ApplyObservedPaul2013PairedScalarFeatures writes the directly observed
// second/first phone identity ordinals into positions 0 and 1. For AH as the
// second phone, it also writes the observed stress value at position 3.
// Other positions are left untouched because their producers are unresolved.
func ApplyObservedPaul2013PairedScalarFeatures(input []int16, first, second CMUPhone) error {
	if len(input) < 2 {
		return errors.New("paired-scalar tree input must contain at least two values")
	}
	firstFeatures, err := first.TreeFeatures()
	if err != nil {
		return err
	}
	secondFeatures, err := second.TreeFeatures()
	if err != nil {
		return err
	}
	if second.Label == "AH" {
		if len(input) < 4 {
			return errors.New("AH paired-scalar tree input must contain at least four values")
		}
	}
	input[0] = secondFeatures.IdentityOrdinal
	input[1] = firstFeatures.IdentityOrdinal
	if second.Label == "AH" {
		input[3] = secondFeatures.StressValue
	}
	return nil
}

// ApplyObservedPaul2013AHPairFeatures is the AH-specific compatibility helper.
// It requires the four-field template used by the captured stress probes.
func ApplyObservedPaul2013AHPairFeatures(input []int16, first, second CMUPhone) error {
	if len(input) < 4 {
		return errors.New("paired-scalar tree input must contain at least four values")
	}
	if second.Label != "AH" {
		return errors.New("observed paired-scalar stress field is established only for AH")
	}
	return ApplyObservedPaul2013PairedScalarFeatures(input, first, second)
}

// ApplyObservedPaul2013AHStressFields fills the fields established for the
// third scalar tree and its paired 12-value vector tree. Both inputs receive
// AH stress at position 2; the vector also receives the preceding scalar
// result at position 11. It does not evaluate either tree.
func ApplyObservedPaul2013AHStressFields(
	scalarInput, vectorInput []int16,
	phone CMUPhone,
	precedingScalarResult int16,
) error {
	if len(scalarInput) < 3 {
		return errors.New("observed AH scalar tree input must contain at least three values")
	}
	if len(vectorInput) < 12 {
		return errors.New("observed AH vector tree input must contain at least twelve values")
	}
	if phone.Label != "AH" {
		return errors.New("observed scalar/vector stress fields are established only for AH")
	}
	features, err := phone.TreeFeatures()
	if err != nil {
		return err
	}
	scalarInput[2] = features.StressValue
	vectorInput[2] = features.StressValue
	vectorInput[11] = precedingScalarResult
	return nil
}
