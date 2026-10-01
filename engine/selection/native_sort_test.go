package selection

import (
	"bufio"
	"math"
	"os"
	"strconv"
	"strings"
	"testing"
)

func TestSortPaul2013NativeMatchesCapturedTailPermutations(t *testing.T) {
	for _, vector := range []struct {
		name string
		path string
	}{
		{name: "partition", path: "../../tools/revkit/work/stage20/native-tail-sort-order.tsv"},
		{name: "heap fallback", path: "../../tools/revkit/work/stage20/native-heap-tail-sort-order.tsv"},
	} {
		t.Run(vector.name, func(t *testing.T) {
			inputIDs, scoreBits, outputIDs := loadCapturedNativeTailSortVector(t, vector.path)
			if vector.name == "heap fallback" {
				for index, got := range scoreBits {
					want := math.Float32bits(3)
					switch index {
					case 0:
						want = math.Float32bits(0)
					case 37:
						want = math.Float32bits(1)
					case 74:
						want = math.Float32bits(2)
					}
					if got != want {
						t.Fatalf("forced heap score at position %d = %#08x, want %#08x", index, got, want)
					}
				}
			}
			type candidate struct {
				id    uint32
				score float32
			}
			candidates := make([]candidate, len(inputIDs))
			for index := range inputIDs {
				candidates[index] = candidate{id: inputIDs[index], score: math.Float32frombits(scoreBits[index])}
			}
			sortPaul2013Native(candidates, func(left, right candidate) bool { return left.score < right.score })
			for index, want := range outputIDs {
				if candidates[index].id != want {
					t.Errorf("native tail sort position %d = %d, want captured unit %d", index, candidates[index].id, want)
				}
			}

			engineCandidates := make([]Paul2013ScoredCandidate, len(inputIDs))
			for index, id := range inputIDs {
				engineCandidates[index] = Paul2013ScoredCandidate{
					Continuity:          Paul2013CandidateContinuity{Unit: UnitRef{Bank: "gen", Index: id}},
					SecondaryOrderScore: math.Float32frombits(scoreBits[index]),
				}
			}
			shortlist, err := SortPaul2013CandidateTailAndCap(engineCandidates, 0)
			if err != nil {
				t.Fatal(err)
			}
			if len(shortlist) != 30 {
				t.Fatalf("engine tail sorter returned %d candidates, want 30", len(shortlist))
			}
			for index, want := range outputIDs[:30] {
				if got := shortlist[index].Continuity.Unit.Index; got != want {
					t.Errorf("engine shortlist position %d = %d, want captured unit %d", index, got, want)
				}
			}
		})
	}
}

func loadCapturedNativeTailSortVector(t *testing.T, path string) ([]uint32, []uint32, []uint32) {
	t.Helper()
	file, err := os.Open(path)
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()

	var inputIDs, outputIDs []uint32
	var scoreBits []uint32
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		fields := strings.Fields(scanner.Text())
		if len(fields) == 0 || strings.HasPrefix(fields[0], "#") {
			continue
		}
		values := make([]uint32, 0, len(fields)-1)
		for _, field := range fields[1:] {
			base := 10
			if fields[0] == "score_bits" {
				base = 16
			}
			value, parseErr := strconv.ParseUint(field, base, 32)
			if parseErr != nil {
				t.Fatalf("parse %s value %q: %v", fields[0], field, parseErr)
			}
			values = append(values, uint32(value))
		}
		switch fields[0] {
		case "pre_ids":
			inputIDs = values
		case "score_bits":
			scoreBits = values
		case "post_ids":
			outputIDs = values
		default:
			t.Fatalf("unknown native-sort fixture row %q", fields[0])
		}
	}
	if err := scanner.Err(); err != nil {
		t.Fatal(err)
	}
	if len(inputIDs) != 75 || len(scoreBits) != len(inputIDs) || len(outputIDs) != len(inputIDs) {
		t.Fatalf("native-sort fixture %s has input/scores/output lengths %d/%d/%d, want 75/75/75", path, len(inputIDs), len(scoreBits), len(outputIDs))
	}
	return inputIDs, scoreBits, outputIDs
}
