package main

import (
	"context"
	"encoding/binary"
	"errors"
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"

	"vtspeak/engine/dat"
	"vtspeak/engine/selection"
	"vtspeak/engine/synthesis"
	"vtspeak/engine/voice"
)

func main() {
	dataRoot := flag.String("data-root", "", "local 2013 Paul M16 data root")
	unitList := flag.String("units", "", "comma-separated bank:index sequence, for example gen:12,gen:13")
	outputPath := flag.String("output", "", "new WAVE output path (must not already exist)")
	renderMode := flag.String("render-mode", "concat", "renderer: concat, periods, or segments")
	pitch := flag.Int("pitch", -1, "pitch control; -1 selects the engine default")
	speed := flag.Int("speed", -1, "speed control; -1 selects the engine default")
	volume := flag.Int("volume", -1, "volume control; -1 selects the engine default")
	applyGain := flag.Bool("apply-observed-gain", false, "apply the observed volume gain")
	flag.Parse()
	options := renderOptions{
		mode: *renderMode, pitch: *pitch, speed: *speed,
		volume: *volume, applyGain: *applyGain,
	}
	if err := runWithOptions(*dataRoot, *unitList, *outputPath, options); err != nil {
		fmt.Fprintln(os.Stderr, "vtconcat:", err)
		os.Exit(1)
	}
}

func run(dataRoot, unitList, outputPath string) error {
	return runWithOptions(dataRoot, unitList, outputPath, renderOptions{
		mode: "concat", pitch: -1, speed: -1, volume: -1,
	})
}

type renderOptions struct {
	mode      string
	pitch     int
	speed     int
	volume    int
	applyGain bool
}

func runWithOptions(dataRoot, unitList, outputPath string, options renderOptions) error {
	if dataRoot == "" || outputPath == "" {
		return errors.New("-data-root and -output are required")
	}
	units, err := parseUnitRefs(unitList)
	if err != nil {
		return err
	}
	model, err := voice.OpenPaul2013(dataRoot)
	if err != nil {
		return fmt.Errorf("open Paul model: %w", err)
	}
	defer model.Close()
	controls, err := makeControls(options)
	if err != nil {
		return err
	}
	renderer, err := makeRenderer(options.mode, model, options.applyGain)
	if err != nil {
		return err
	}
	pcm, err := renderer.Render(context.Background(), units, controls)
	if err != nil {
		return fmt.Errorf("render selected units: %w", err)
	}
	pcmData := make([]byte, 0, len(pcm.Samples)*2)
	for _, sample := range pcm.Samples {
		pcmData = binary.LittleEndian.AppendUint16(pcmData, uint16(sample))
	}
	wave, err := dat.WAV(pcmData)
	if err != nil {
		return fmt.Errorf("write WAVE container: %w", err)
	}
	file, err := os.OpenFile(outputPath, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0o644)
	if err != nil {
		return fmt.Errorf("create output without overwriting: %w", err)
	}
	if _, err := file.Write(wave); err != nil {
		_ = file.Close()
		_ = os.Remove(outputPath)
		return fmt.Errorf("write WAVE output: %w", err)
	}
	if err := file.Close(); err != nil {
		_ = os.Remove(outputPath)
		return fmt.Errorf("close WAVE output: %w", err)
	}
	fmt.Printf("wrote %s (%d samples at %d Hz)\n", outputPath, len(pcm.Samples), pcm.SampleRate)
	return nil
}

func makeControls(options renderOptions) (synthesis.Controls, error) {
	pitch, err := controlValue("pitch", options.pitch)
	if err != nil {
		return synthesis.Controls{}, err
	}
	speed, err := controlValue("speed", options.speed)
	if err != nil {
		return synthesis.Controls{}, err
	}
	volume, err := controlValue("volume", options.volume)
	if err != nil {
		return synthesis.Controls{}, err
	}
	return synthesis.Controls{Pitch: pitch, Speed: speed, Volume: volume}, nil
}

func controlValue(name string, value int) (int32, error) {
	value64 := int64(value)
	if value64 < -1<<31 || value64 > 1<<31-1 {
		return 0, fmt.Errorf("-%s value %d is outside the signed 32-bit range", name, value)
	}
	return int32(value), nil
}

func makeRenderer(mode string, units synthesis.UnitReader, applyGain bool) (synthesis.Renderer, error) {
	switch mode {
	case "concat":
		return synthesis.ConcatenatingRenderer{Units: units, ApplyObservedGain: applyGain}, nil
	case "periods":
		return synthesis.PeriodResamplingRenderer{Units: units, ApplyObservedGain: applyGain}, nil
	case "segments":
		return synthesis.UPMSegmentPlanRenderer{Units: units, ApplyObservedGain: applyGain}, nil
	default:
		return nil, fmt.Errorf("unknown -render-mode %q (want concat, periods, or segments)", mode)
	}
}

func parseUnitRefs(value string) ([]selection.UnitRef, error) {
	if strings.TrimSpace(value) == "" {
		return nil, errors.New("-units must contain at least one bank:index reference")
	}
	parts := strings.Split(value, ",")
	units := make([]selection.UnitRef, 0, len(parts))
	for position, part := range parts {
		bank, indexText, ok := strings.Cut(part, ":")
		if !ok || bank == "" || indexText == "" || strings.Contains(indexText, ":") {
			return nil, fmt.Errorf("unit reference %d %q must use bank:index form", position, part)
		}
		index, err := strconv.ParseUint(indexText, 10, 32)
		if err != nil {
			return nil, fmt.Errorf("unit reference %d has invalid index %q: %w", position, indexText, err)
		}
		units = append(units, selection.UnitRef{Bank: bank, Index: uint32(index)})
	}
	return units, nil
}
