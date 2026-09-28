// vtdecode extracts one 2013 Paul unit as a standalone WAV for analysis.
package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"

	"vtspeak/engine/dat"
)

func run(args []string) error {
	flags := flag.NewFlagSet("vtdecode", flag.ContinueOnError)
	indexPath := flags.String("index", "", "2013 Paul unit index")
	datPath := flags.String("dat", "", "matching DAT bank")
	unitText := flags.String("unit", "", "zero-based unit number")
	outputPath := flags.String("output", "", "output WAV path")
	if err := flags.Parse(args); err != nil {
		return err
	}
	if *indexPath == "" || *datPath == "" || *unitText == "" || *outputPath == "" || flags.NArg() != 0 {
		return fmt.Errorf("usage: vtdecode -index FILE -dat FILE -unit NUMBER -output FILE")
	}
	unit, err := strconv.ParseUint(*unitText, 10, 32)
	if err != nil {
		return fmt.Errorf("invalid unit number: %w", err)
	}
	index, err := os.ReadFile(*indexPath)
	if err != nil {
		return err
	}
	bank, err := os.Open(*datPath)
	if err != nil {
		return err
	}
	defer bank.Close()
	info, err := bank.Stat()
	if err != nil {
		return err
	}
	payload, err := dat.UnitPayload(index, bank, info.Size(), uint32(unit))
	if err != nil {
		return err
	}
	pcm, err := dat.Decode(payload)
	if err != nil {
		return err
	}
	wav, err := dat.WAV(pcm)
	if err != nil {
		return err
	}
	return os.WriteFile(*outputPath, wav, 0o644)
}

func main() {
	if err := run(os.Args[1:]); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
