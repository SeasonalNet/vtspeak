# Reimplementation language decision

**Status:** Go is the selected language for the reimplementation for now.
This records the starting implementation choice; the engine has not yet been
implemented.

## Decision

Write the native reimplementation engine in Go. The goal remains an engine
that can understand the supported mid-2000s and early-2010s VoiceText data and
produce speech. This language decision does not establish which individual
voices or package variants are supported, nor does it claim synthesis parity.
Those boundaries belong in the implementation plan and must be established
against original-engine evidence.

Use Go's standard library for the initial core, including
[`encoding/binary`](https://pkg.go.dev/encoding/binary) for explicit binary
data handling. No third-party runtime
dependency has been identified as necessary for the currently understood
work: reading binary resources, managing fixed-width values and sample buffers,
performing synthesis operations, and writing PCM output. Revisit that only if
reverse engineering identifies a concrete format or algorithm need that the
standard library does not address.

## Fidelity and validation

The language does not provide byte-for-byte parity automatically. The Go
implementation must reproduce the observed operation order, integer widths,
rounding, lookup behavior, edge cases, and output layout. Validate each
implemented stage against the corresponding original-engine behavior and
retain byte comparisons for outputs where parity is established. Keep
observations, inferences, and remaining compatibility gaps explicit.

Go's garbage collector and allocation behavior are implementation concerns to
measure in the actual engine. They do not prevent matching output bytes. The
[Go language specification](https://go.dev/ref/spec) defines integer overflow
behavior, but that does not establish what any particular engine operation
should do. Use fixed-width types and explicit arithmetic where the recovered
behavior calls for them; do not rely on a language default as a substitute for
evidence.

## Optional native ABI boundary

Keep open the option to implement a bounded component in Rust or C++ if later
evidence shows a specific need, such as a performance bottleneck or an
operation that is materially easier to express and validate there. Go can call
such a component through [`cgo`](https://pkg.go.dev/cmd/cgo) using a
C-compatible ABI. Rust can export an
[`extern "C"`](https://doc.rust-lang.org/reference/items/external-blocks.html#abi)
interface; C++ should expose a small `extern "C"` wrapper around its C++
implementation. Do not depend on Rust's native ABI or expose C++
compiler-specific types and classes across this boundary. See Microsoft's
[C++ binary compatibility notes](https://learn.microsoft.com/en-us/cpp/porting/binary-compat-2015-2017)
for an example of the toolchain compatibility constraints involved.

If introduced, keep the boundary coarse and explicit: fixed-width scalar
types, byte buffers with lengths and capacities, documented ownership, and
status/error results. Do not allow panics, C++ exceptions, or retained Go
pointers to cross it. The native component would add compiler, linker, and
platform build requirements, so add one only to solve a demonstrated need.

## Reconsideration

Revisit the language choice if implementation evidence shows that Go cannot
meet a concrete fidelity, resource, or platform requirement without
disproportionate complexity. Make that comparison against the recovered engine
behavior and representative fixtures, not general language preferences.
