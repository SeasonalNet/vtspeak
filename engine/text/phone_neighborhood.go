package text

import "fmt"

// Paul2013PhoneNeighborhood carries only the directly recovered features for
// one phone and its immediate neighbors. It is not a decision-tree vector.
type Paul2013PhoneNeighborhood struct {
	Previous    TreePhoneFeatures
	Current     TreePhoneFeatures
	Next        TreePhoneFeatures
	HasPrevious bool
	HasNext     bool
}

// BuildPaul2013PhoneNeighborhoods maps a phone sequence to identity/stress
// features and records which adjacent phones exist. The caller supplies one
// already-selected pronunciation span; this helper does not join across word
// or punctuation boundaries or assign tree-specific feature positions.
func BuildPaul2013PhoneNeighborhoods(phones []CMUPhone) ([]Paul2013PhoneNeighborhood, error) {
	features := make([]TreePhoneFeatures, len(phones))
	for index, phone := range phones {
		feature, err := phone.TreeFeatures()
		if err != nil {
			return nil, fmt.Errorf("phone %d: %w", index, err)
		}
		features[index] = feature
	}

	neighborhoods := make([]Paul2013PhoneNeighborhood, len(features))
	for index, current := range features {
		neighborhoods[index].Current = current
		if index > 0 {
			neighborhoods[index].Previous = features[index-1]
			neighborhoods[index].HasPrevious = true
		}
		if index+1 < len(features) {
			neighborhoods[index].Next = features[index+1]
			neighborhoods[index].HasNext = true
		}
	}
	return neighborhoods, nil
}
