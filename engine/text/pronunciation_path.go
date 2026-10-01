package text

import "fmt"

const (
	paul2013PathGroupDelimiter = byte('d')
	paul2013PathTerminator     = byte(0xff)
	paul2013PathGroupCapacity  = 20
)

// paul2013PathCodeClasses is the signed-short map at DAT_100783ec. The DLL
// indexes it by a signed path byte. Values of -1 are preserved because the
// native membership check compares them like any other class value.
var paul2013PathCodeClasses = [...]int16{
	-1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1,
	0, 1, 2, 3, -1, -1, -1, 4, 5, 6, 7,
	-1, -1, -1, -1, -1, 8, -1, -1, -1, -1, -1, -1, -1, -1,
	9, 10, 11, 10, 9, 12, -1, -1, -1,
}

// SelectPaul2013PronunciationByPathMarker ports FUN_10003f10. It searches the
// path rows for selector, treating '%' as a match for ')', then prefers the
// first row containing byte 0x17, and otherwise uses row zero. It returns the
// phone string aligned with that path row. The selector's meaning and its
// producer are not established. An empty path list returns false; mismatched
// path and phone row counts are rejected.
func SelectPaul2013PronunciationByPathMarker(
	paths [][]byte,
	phoneRows [][]byte,
	selector byte,
) ([]byte, bool, error) {
	if len(paths) == 0 {
		if len(phoneRows) != 0 {
			return nil, false, fmt.Errorf("received %d phone rows for no path rows", len(phoneRows))
		}
		return nil, false, nil
	}
	if len(paths) != len(phoneRows) {
		return nil, false, fmt.Errorf("received %d path rows and %d phone rows", len(paths), len(phoneRows))
	}
	selected := -1
	for rowIndex, path := range paths {
		for _, value := range paul2013PathRowBytes(path) {
			if value == selector || selector == ')' && value == '%' {
				return append([]byte(nil), paul2013PhoneRowBytes(phoneRows[rowIndex])...), true, nil
			}
			if value == 0x17 && selected == -1 {
				selected = rowIndex
			}
		}
	}
	if selected == -1 {
		selected = 0
	}
	return append([]byte(nil), paul2013PhoneRowBytes(phoneRows[selected])...), true, nil
}

func paul2013PathRowBytes(path []byte) []byte {
	for index, value := range path {
		if value == paul2013PathTerminator {
			return path[:index]
		}
	}
	return path
}

func paul2013PhoneRowBytes(phones []byte) []byte {
	for index, value := range phones {
		if value == 0 {
			return phones[:index]
		}
	}
	return phones
}

// ParsePaul2013PronunciationPath ports FUN_100105d0 and FUN_10010550. A path
// is divided into groups at 'd'; within each group '(' and ')' become '&' and
// '%' and repeated instances of either marker are collapsed. The 24-byte
// native group record reserves four bytes before a 20-byte, 0xff-terminated
// path string, so at most 19 path bytes fit in one group.
func ParsePaul2013PronunciationPath(path []byte) ([][]byte, error) {
	groups := make([][]byte, 0, 1)
	group := make([]byte, 0, paul2013PathGroupCapacity-1)
	flush := func() {
		groups = append(groups, append([]byte(nil), group...))
		group = group[:0]
	}

	for index, value := range path {
		if value == paul2013PathTerminator {
			if index != len(path)-1 {
				return nil, fmt.Errorf("pronunciation path has bytes after terminator at offset %d", index)
			}
			flush()
			return groups, nil
		}
		if value == paul2013PathGroupDelimiter {
			flush()
			continue
		}
		switch value {
		case '(':
			if containsPronunciationPathByte(group, '&') {
				continue
			}
			value = '&'
		case ')':
			if containsPronunciationPathByte(group, '%') {
				continue
			}
			value = '%'
		}
		if len(group) >= paul2013PathGroupCapacity-1 {
			return nil, fmt.Errorf("pronunciation path group at offset %d exceeds %d bytes", index, paul2013PathGroupCapacity-1)
		}
		group = append(group, value)
	}

	flush()
	return groups, nil
}

// RankPaul2013PronunciationPathGroups ports the path-membership test in
// FUN_100104e0 and the stable maximum-score choice in FUN_100068b0. Each
// classifier output is compared with every code in each parsed path group via
// the DLL's embedded code-to-class table. A group gains one point per matching
// output; ties keep the earliest group. The caller must still evaluate
// candidate feature rows in native order and associate each output with its
// path group.
func RankPaul2013PronunciationPathGroups(
	groups [][]byte,
	classifierOutputs []int16,
) (selected int, scores []int, err error) {
	if len(groups) == 0 {
		return 0, nil, fmt.Errorf("rank pronunciation paths: no path groups")
	}
	for groupIndex, group := range groups {
		for _, pathCode := range group {
			if pathCode != paul2013PathTerminator && int(pathCode) >= len(paul2013PathCodeClasses) {
				return 0, nil, fmt.Errorf("pronunciation path group %d has code 0x%02x outside the recovered table", groupIndex, pathCode)
			}
		}
	}
	scores = make([]int, len(groups))
	for _, output := range classifierOutputs {
		class := int16(int8(uint8(output)))
		for groupIndex, group := range groups {
			if groupContainsClass(group, class) {
				scores[groupIndex]++
			}
		}
	}
	for groupIndex := 1; groupIndex < len(scores); groupIndex++ {
		if scores[groupIndex] > scores[selected] {
			selected = groupIndex
		}
	}
	return selected, scores, nil
}

// Paul2013PronunciationPathClass returns the first byte's class mapping used
// by FUN_100068b0 for classifier feature 13. present is false for an empty
// path group because the native loop skips groups whose first byte is 0xff.
func Paul2013PronunciationPathClass(group []byte) (class int16, present bool, err error) {
	if len(group) == 0 {
		return 0, false, nil
	}
	pathCode := group[0]
	if int(pathCode) >= len(paul2013PathCodeClasses) {
		return 0, false, fmt.Errorf("pronunciation path code 0x%02x is outside the recovered table", pathCode)
	}
	return paul2013PathCodeClasses[pathCode], true, nil
}

// FindPaul2013PronunciationForPathCode ports the alternative scan in
// FUN_10010640. It maps the requested code and every code in each alternative
// through DAT_100783ec, then returns the first alternative whose path row
// contains the same class. The class comparison includes mapped -1 entries,
// matching the native table comparison. Alternative order is preserved.
//
// The caller supplies the recovered path-control stream and aligned phone
// strings. This helper performs the lookup mechanics used by the native
// name/context shortcut; it does not detect when that shortcut applies.
func FindPaul2013PronunciationForPathCode(
	pathControlBytes []byte,
	phoneStrings [][]byte,
	requestedCode byte,
) (alternativeIndex int, phoneString []byte, found bool, err error) {
	if requestedCode >= 0x80 || int(requestedCode) >= len(paul2013PathCodeClasses) {
		return 0, nil, false, fmt.Errorf("requested pronunciation path code 0x%02x is outside the recovered signed-byte table", requestedCode)
	}
	if len(pathControlBytes) == 0 {
		return 0, nil, false, fmt.Errorf("pronunciation path-control stream is empty")
	}
	if pathControlBytes[len(pathControlBytes)-1] != paul2013PathTerminator {
		return 0, nil, false, fmt.Errorf("pronunciation path-control stream has no final 0xff terminator")
	}
	groups, err := ParsePaul2013PronunciationPath(pathControlBytes)
	if err != nil {
		return 0, nil, false, fmt.Errorf("parse pronunciation path-control stream: %w", err)
	}
	if len(groups) != len(phoneStrings) {
		return 0, nil, false, fmt.Errorf("pronunciation path stream has %d alternatives for %d phone strings", len(groups), len(phoneStrings))
	}
	requestedClass := paul2013PathCodeClasses[requestedCode]
	for groupIndex, group := range groups {
		for _, pathCode := range group {
			if pathCode == paul2013PathTerminator {
				break
			}
			if int(pathCode) >= len(paul2013PathCodeClasses) {
				return 0, nil, false, fmt.Errorf("pronunciation alternative %d has path code 0x%02x outside the recovered table", groupIndex, pathCode)
			}
			if paul2013PathCodeClasses[pathCode] == requestedClass {
				return groupIndex, append([]byte(nil), paul2013PhoneRowBytes(phoneStrings[groupIndex])...), true, nil
			}
		}
	}
	return 0, nil, false, nil
}

func groupContainsClass(group []byte, class int16) bool {
	for _, pathCode := range group {
		if pathCode == paul2013PathTerminator {
			break
		}
		if paul2013PathCodeClasses[pathCode] == class {
			return true
		}
	}
	return false
}

func containsPronunciationPathByte(path []byte, value byte) bool {
	for _, candidate := range path {
		if candidate == value {
			return true
		}
	}
	return false
}
