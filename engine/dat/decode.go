// Package dat decodes the framed DAT unit payloads observed in the 2013 M16 Paul voice.
package dat

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io"
)

const (
	defaultFrameSamples = 256
	maxFrameSamples     = 1 << 16
	maxUnaryQuotient    = 1 << 20
)

type bitReader struct {
	data []byte
	pos  int
}

func (r *bitReader) bit() (uint32, error) {
	if r.pos >= len(r.data)*8 {
		return 0, io.ErrUnexpectedEOF
	}
	b := (r.data[r.pos/8] >> (7 - r.pos%8)) & 1
	r.pos++
	return uint32(b), nil
}

func (r *bitReader) rice(width uint32) (uint32, error) {
	if width > 31 {
		return 0, fmt.Errorf("unsupported Rice remainder width %d", width)
	}
	var quotient uint32
	for {
		b, err := r.bit()
		if err != nil {
			return 0, err
		}
		if b == 1 {
			break
		}
		quotient++
		if quotient > maxUnaryQuotient {
			return 0, errors.New("DAT unary quotient exceeds safety limit")
		}
	}
	var remainder uint32
	for range width {
		b, err := r.bit()
		if err != nil {
			return 0, err
		}
		remainder = remainder<<1 | b
	}
	return quotient<<width | remainder, nil
}

func prior(frame []int32, history [3]int32, distance int) int32 {
	if len(frame) >= distance {
		return frame[len(frame)-distance]
	}
	return history[3-distance+len(frame)]
}

// Decode returns little-endian signed 16-bit PCM for one complete DAT payload.
// The supported behavior is bounded by the observed 2013 M16 Paul decoder path.
func Decode(payload []byte) ([]byte, error) {
	r := bitReader{data: payload}
	frameSamples := uint32(defaultFrameSamples)
	var outputShift uint32
	var means [4]int32
	var history [3]int32
	limit := len(payload) * 128
	if limit < defaultFrameSamples {
		limit = defaultFrameSamples
	}
	var pcm []byte
	for {
		mode, err := r.rice(2)
		if err != nil {
			return nil, fmt.Errorf("DAT control: %w", err)
		}
		switch mode {
		case 4:
			if len(pcm) == 0 {
				return nil, errors.New("DAT payload terminated without samples")
			}
			return pcm, nil
		case 5:
			width, err := r.rice(2)
			if err != nil {
				return nil, fmt.Errorf("DAT frame-size width: %w", err)
			}
			frameSamples, err = r.rice(width)
			if err != nil {
				return nil, fmt.Errorf("DAT frame size: %w", err)
			}
			if frameSamples == 0 || frameSamples > maxFrameSamples {
				return nil, fmt.Errorf("invalid DAT frame size %d", frameSamples)
			}
			continue
		case 6:
			outputShift, err = r.rice(2)
			if err != nil {
				return nil, fmt.Errorf("DAT output shift: %w", err)
			}
			if outputShift > 15 {
				return nil, fmt.Errorf("invalid DAT output shift %d", outputShift)
			}
			continue
		case 0, 1, 2, 3, 8:
		default:
			return nil, fmt.Errorf("unsupported DAT control value %d", mode)
		}
		var residualWidth uint32
		if mode != 8 {
			residualWidth, err = r.rice(3)
			if err != nil {
				return nil, fmt.Errorf("DAT residual width: %w", err)
			}
			if residualWidth > 30 {
				return nil, fmt.Errorf("invalid DAT residual width %d", residualWidth)
			}
		}
		meanSum := int32(int64(means[0]) + int64(means[1]) + int64(means[2]) + int64(means[3]) + 2)
		if meanSum < 0 {
			meanSum += 3
		}
		meanOffset := meanSum >> 2
		if outputShift != 0 {
			meanOffset = (meanOffset >> (outputShift - 1)) >> 1
		}
		frame := make([]int32, 0, frameSamples)
		for range frameSamples {
			var value int32
			if mode != 8 {
				coded, readErr := r.rice(residualWidth + 1)
				if readErr != nil {
					return nil, fmt.Errorf("DAT residual: %w", readErr)
				}
				residual := int32(coded >> 1)
				if coded&1 != 0 {
					residual = -residual - 1
				}
				prediction := meanOffset
				switch mode {
				case 1:
					prediction = prior(frame, history, 1)
				case 2:
					prediction = 2*prior(frame, history, 1) - prior(frame, history, 2)
				case 3:
					prediction = 3*(prior(frame, history, 1)-prior(frame, history, 2)) + prior(frame, history, 3)
				}
				value = prediction + residual
			}
			frame = append(frame, value)
		}
		var total int64 = int64(frameSamples / 2)
		for _, value := range frame {
			total += int64(value)
		}
		negative := total < 0
		if negative {
			total = -total
		}
		mean := int32(total / int64(frameSamples))
		if negative {
			mean = -mean
		}
		means = [4]int32{means[1], means[2], means[3], mean << outputShift}
		for _, value := range frame {
			history = [3]int32{history[1], history[2], value}
			scaled := value << outputShift
			if scaled > 32767 {
				scaled = 32767
			}
			pcm = binary.LittleEndian.AppendUint16(pcm, uint16(scaled))
		}
		if len(pcm)/2 > limit {
			return nil, errors.New("decoded DAT output exceeded safety limit")
		}
	}
}
