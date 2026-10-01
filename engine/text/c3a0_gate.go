package text

import (
	"bytes"
)

// Paul2013C3A0Gate ports FUN_10008cc0's source eligibility predicate for
// byte input. The character-class bits remain opaque and are used only as
// observed by the native branches. High-bit bytes read zero from the verified
// signed-index prefix and therefore fail the native eligibility predicate.
func Paul2013C3A0Gate(source []byte) (bool, error) {
	if nul := bytes.IndexByte(source, 0); nul >= 0 {
		source = source[:nul]
	}
	if len(source) == 0 {
		return false, nil
	}

	characterMap := Paul2013EmbeddedKeyTables().CharacterMap
	if paul2013MappedPrefixEqual(source, []byte("Mac"), 3, characterMap) ||
		paul2013MappedPrefixEqual(source, []byte("Mc"), 2, characterMap) {
		return false, nil
	}

	attributes := Paul2013ExceptionCharacterAttributes()
	caseState := int16(0)
	hasMixedClassTransition := false
	for _, value := range source {
		attribute := attributes[value]
		if attribute&0x80 != 0 {
			if caseState == -1 {
				hasMixedClassTransition = true
				break
			}
			caseState = 1
		} else if attribute&0x40 != 0 {
			if caseState == 1 {
				hasMixedClassTransition = true
				break
			}
			caseState = -1
		}
	}
	if !hasMixedClassTransition {
		return false, nil
	}

	for index, value := range source {
		if attributes[value]&0xc0 != 0 {
			continue
		}
		if value != '\'' {
			return false, nil
		}
		if index+2 < len(source) && source[index+1] == 's' && attributes[source[index+2]]&0x80 != 0 {
			continue
		}
		if index > 0 && source[index-1] == 's' && index+1 < len(source) && attributes[source[index+1]]&0x80 != 0 {
			continue
		}
		return false, nil
	}
	return false, nil
}

func paul2013MappedPrefixEqual(source, prefix []byte, length int, characterMap [256]byte) bool {
	if len(source) < length || len(prefix) < length {
		return false
	}
	for index := 0; index < length; index++ {
		if characterMap[source[index]] != characterMap[prefix[index]] {
			return false
		}
	}
	return true
}
