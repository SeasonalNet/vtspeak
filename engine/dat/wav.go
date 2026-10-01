package dat

import (
	"encoding/binary"
	"errors"
	"io"
)

// WAV wraps mono 16 kHz signed 16-bit PCM in a 44-byte RIFF/WAVE header.
func WAV(pcm []byte) ([]byte, error) {
	maxInt := uint64(^uint(0) >> 1)
	if len(pcm)%2 != 0 || uint64(len(pcm)) > uint64(^uint32(0))-36 ||
		uint64(len(pcm)) > maxInt-44 {
		return nil, errors.New("PCM length is invalid for a RIFF/WAVE file")
	}
	wav := make([]byte, 44, 44+len(pcm))
	writeWAVHeader(wav, uint32(len(pcm)))
	return append(wav, pcm...), nil
}

// WAVFromSamples wraps signed 16-bit mono samples in the observed 16 kHz
// WAVE profile without first allocating an intermediate byte slice.
func WAVFromSamples(samples []int16) ([]byte, error) {
	dataBytes := uint64(len(samples)) * 2
	maxInt := uint64(^uint(0) >> 1)
	if dataBytes > uint64(^uint32(0))-36 || dataBytes > maxInt-44 {
		return nil, errors.New("PCM length is invalid for a RIFF/WAVE file")
	}
	wav := make([]byte, 44+int(dataBytes))
	writeWAVHeader(wav[:44], uint32(dataBytes))
	for index, sample := range samples {
		binary.LittleEndian.PutUint16(wav[44+index*2:], uint16(sample))
	}
	return wav, nil
}

// WriteWAV writes signed 16-bit mono PCM in the observed 16 kHz WAVE profile
// without assembling a second copy of the complete file in memory. The data
// length is known up front, so the finalized RIFF header is written first.
func WriteWAV(destination io.Writer, pcm []byte) error {
	if destination == nil {
		return errors.New("WAVE destination is nil")
	}
	if len(pcm)%2 != 0 || uint64(len(pcm)) > uint64(^uint32(0))-36 {
		return errors.New("PCM length is invalid for a RIFF/WAVE file")
	}
	var header [44]byte
	writeWAVHeader(header[:], uint32(len(pcm)))
	if err := writeAll(destination, header[:]); err != nil {
		return err
	}
	return writeAll(destination, pcm)
}

// WriteWAVFromSamples streams signed samples through a fixed-size conversion
// buffer, avoiding both the intermediate PCM byte slice and the complete WAVE
// allocation used by WAVFromSamples.
func WriteWAVFromSamples(destination io.Writer, samples []int16) error {
	if destination == nil {
		return errors.New("WAVE destination is nil")
	}
	dataBytes := uint64(len(samples)) * 2
	if dataBytes > uint64(^uint32(0))-36 {
		return errors.New("PCM length is invalid for a RIFF/WAVE file")
	}
	var header [44]byte
	writeWAVHeader(header[:], uint32(dataBytes))
	if err := writeAll(destination, header[:]); err != nil {
		return err
	}
	var buffer [8192]byte
	for offset := 0; offset < len(samples); {
		count := len(buffer) / 2
		if remaining := len(samples) - offset; remaining < count {
			count = remaining
		}
		chunk := buffer[:count*2]
		for index, sample := range samples[offset : offset+count] {
			binary.LittleEndian.PutUint16(chunk[index*2:], uint16(sample))
		}
		if err := writeAll(destination, chunk); err != nil {
			return err
		}
		offset += count
	}
	return nil
}

func writeAll(destination io.Writer, value []byte) error {
	for len(value) != 0 {
		written, err := destination.Write(value)
		if written < 0 || written > len(value) {
			return errors.New("WAVE destination returned an invalid write count")
		}
		value = value[written:]
		if err != nil {
			return err
		}
		if written == 0 {
			return io.ErrShortWrite
		}
	}
	return nil
}

func writeWAVHeader(header []byte, dataBytes uint32) {
	copy(header[0:4], "RIFF")
	binary.LittleEndian.PutUint32(header[4:8], 36+dataBytes)
	copy(header[8:12], "WAVE")
	copy(header[12:16], "fmt ")
	binary.LittleEndian.PutUint32(header[16:20], 16)
	binary.LittleEndian.PutUint16(header[20:22], 1)
	binary.LittleEndian.PutUint16(header[22:24], 1)
	binary.LittleEndian.PutUint32(header[24:28], 16000)
	binary.LittleEndian.PutUint32(header[28:32], 32000)
	binary.LittleEndian.PutUint16(header[32:34], 2)
	binary.LittleEndian.PutUint16(header[34:36], 16)
	copy(header[36:40], "data")
	binary.LittleEndian.PutUint32(header[40:44], dataBytes)
}
