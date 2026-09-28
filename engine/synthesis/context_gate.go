package synthesis

import (
	"errors"
	"fmt"
)

// ContextGateInput supplies the observed row and index values consumed by
// FUN_1002d230. Row bytes and mode values stay numeric; this does not assign
// linguistic meanings to the packed categories.
type ContextGateInput struct {
	Rows                 [][7]byte
	Modes                []byte
	CurrentIndex         int
	LeftIndex            int
	RightIndex           int
	TimelineSpeedControl int32
}

// ContextMultipliers reports the independently measured left/right context
// multipliers. Zero suppresses a side; eligible sides use one or two.
type ContextMultipliers struct {
	Left  uint8
	Right uint8
}

// Paul2013ContextMultipliers ports the context-side gate predicates observed
// in FUN_1002d230. It does not gather neighbor units or blend their samples.
func Paul2013ContextMultipliers(input ContextGateInput) (ContextMultipliers, error) {
	if len(input.Rows) == 0 || len(input.Modes) != len(input.Rows) {
		return ContextMultipliers{}, errors.New("context gate rows and modes must have the same nonzero length")
	}
	if input.CurrentIndex < 0 || input.CurrentIndex >= len(input.Rows) {
		return ContextMultipliers{}, errors.New("current context index is outside the row range")
	}
	current := input.Rows[input.CurrentIndex]
	mode := input.Modes[input.CurrentIndex]
	result := ContextMultipliers{Left: 2, Right: 2}

	leftSuppressed := input.CurrentIndex == 0 || input.TimelineSpeedControl == 0 ||
		current[5]&0x80 != 0 && current[5]&0x38 > 0x17
	if !leftSuppressed {
		if input.LeftIndex < 0 || input.LeftIndex >= len(input.Rows) {
			return ContextMultipliers{}, fmt.Errorf("left context index %d is outside the row range", input.LeftIndex)
		}
		left := input.Rows[input.LeftIndex]
		switch {
		case mode == 2 && input.LeftIndex == input.CurrentIndex:
			leftSuppressed = true
		case input.LeftIndex == len(input.Rows)-1:
			leftSuppressed = true
		case mode != 2 && input.LeftIndex+1 == input.CurrentIndex && left[6]&0x80 != 0:
			leftSuppressed = true
		case mode != 2 && left[6]&0x80 == 0:
			leftSuppressed = true
		case mode != 2 && input.Rows[input.CurrentIndex-1][6]&0x80 == 0:
			result.Left = 1
		}
	}
	if leftSuppressed {
		result.Left = 0
	}

	rightSuppressed := input.CurrentIndex == len(input.Rows)-1 ||
		current[5]&0x40 != 0 && current[5]&0x07 > 2
	if !rightSuppressed {
		if input.RightIndex < 0 || input.RightIndex >= len(input.Rows) {
			return ContextMultipliers{}, fmt.Errorf("right context index %d is outside the row range", input.RightIndex)
		}
		switch {
		case mode == 1 && input.RightIndex == input.CurrentIndex:
			rightSuppressed = true
		case input.RightIndex == 0:
			rightSuppressed = true
		case mode != 1 && input.CurrentIndex+1 == input.RightIndex && current[6]&0x80 != 0:
			rightSuppressed = true
		case mode != 1 && input.Rows[input.RightIndex-1][6]&0x80 == 0:
			rightSuppressed = true
		case mode != 1 && current[6]&0x80 == 0:
			result.Right = 1
		}
	}
	if rightSuppressed {
		result.Right = 0
	}
	return result, nil
}

// HasEligibleSide reports whether the trace gate would construct a context
// row for at least one side.
func (multipliers ContextMultipliers) HasEligibleSide() bool {
	return multipliers.Left != 0 || multipliers.Right != 0
}
