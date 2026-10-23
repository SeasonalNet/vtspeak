# 2006 pitch contour boundary smoothing: assembly clarification

Static REA/Ghidra 12.1.4 analysis of original `vt_eng.dll`, SHA-256
`00fc9375d08bd8cec303845c992481d79c8f616d1bfb6390a5322cd06f69273d`.
Evidence is preserved in `pitch-smoothing-evidence.json.gz`. No execution,
patching, annotations, or reconstruction code changes occurred.

Questions: identify the unknown math helper (answered by FSQRT); identify the
seven-sample aggregation (answered: median); recover the contour arithmetic
obscured by pseudocode stack errors (answered for ordinary finite positive
median inputs); runtime handling of degenerate contours (unresolved).

## Helper identities and decompiler boundary

`0x10059084` loads a double through `0x1005afb5` (`__fload_withFB`) and falls
through to the separately analyzed procedure at `0x1005908d`. Its ordinary
nonnegative path executes `FSQRT` at `0x100590ac`. This is a square-root helper,
not an unknown higher-level pitch function. The split Ghidra procedure extent
and integer return type obscure the x87 ST0 result in caller pseudocode.

`0x1000f9c0` copies a float vector, sorts it descending, and returns element
`count/2`. With its seven-element calls it returns the median. Its temporary
buffer is 20 floats; this investigation concerns the observed count-seven
callers and makes no arbitrary-count API safety claim.

## Recovered boundary operation

`0x1000fa70` processes consecutive pairs of 12-bin signed-16-bit contours.
Let X and Y denote the earlier and later contour. The source pointer is read
from utterance context `+0x54`, at `phrase*0x10 + 0x18`; contour records are
0x28 bytes apart. X bins begin at record `+4`, Y at the following record `+4`.
The signed-16-bit count at context `(phrase+1)*0x10` controls count-minus-one
iterations. The routine returns one without entering this loop for count <=1.

The following formulas are deductions from the assembly's stack offsets, x87
operations, and read constants, not original source:

1. Compute `mX = median(sqrt(X[i]*X[i]/2), i=5..11)` and
   `mY = median(sqrt(Y[i]*Y[i]/2), i=0..6)`.
2. Set `m=(mX+mY)/2`, `a=sqrt(m/mX)`, `b=sqrt(m/mY)`.
3. For `i=0..11`, scale:
   `X[i] *= 1 + i*(a-1)/11`,
   `Y[i] *= b + i*(1-b)/11`.
4. Let A be scaled X[5], B be scaled Y[6]. Replace the fourteen bins
   `X[5..11], Y[0..6]` with a straight line from A to B inclusive:
   `A + (B-A)*j/13`, j=0..13.
5. Convert each of the two twelve-bin vectors back through x87 `__ftol`
   after adding 0.5, storing low 16 bits. No saturation is visible here.

This combines level adjustment with a linear bridge across the boundary.
For positive contours the square-root-of-square samples equal their positive
magnitudes divided by sqrt(2); the report retains the literal operations to
avoid changing floating rounding or inventing a signal interpretation.

The bridge formula describes its mathematical effect. Assembly evaluates it
as `B + ((B-A)/13)*(j-13)` and reloads X[5] as A each iteration. The first
write back to X[5] can therefore affect later samples through floating
rounding. A future byte-exact port must preserve that evaluation order and
single-precision stores, rather than substituting the algebraically equivalent
expression above.

The pair loop advances by 0x28 bytes, so an interior contour can be modified
again when it becomes the earlier member of the next pair. The operation is
ordered and in place, rather than independent smoothing of all boundaries
against untouched originals.

## Why the pseudocode's literal division by zero is misleading

At baseline stack pointer S, X occupies S+0x60..0x8c, Y occupies S+0x30..0x5c.
The transformed seven-bin vectors occupy S+0x90 and S+0xac. Caller-cleaned
arguments temporarily move ESP by 8 or 16 bytes. The median results live at
S+0x10 and S+0x0c; their mean at S+0x18. Instructions `0x1000fbc8..0x1000fbd3`
load that mean and divide by the second median at S+0x0c before the second
sqrt. Pseudocode loses this stack alias and prints an apparent division by
literal zero; the machine instruction uses a memory denominator.

The denominator can still be zero if an input median is zero. This routine
contains no visible guard around those divisions. This pass did not execute
zero, negative, nonfinite, or overflow cases, and does not establish the
expected vendor input range or the runtime error result in those cases.

## Evidence

| Observation | Evidence ID |
| --- | --- |
| Full smoothing pseudocode and assembly | `ev_fda183f97871c582a45c800b08a26622d96804f5220af352ae67cb329979c290` |
| Seven-element median | `ev_7567e274e5bdda996636f8d71e3664b18315d5322c7aa7fade742116439ba4a8` |
| Square-root entry | `ev_2d10d60e38b504bd459acaa862d8c6e669e09d38754adea827933c1e690e7ae1` |
| Entry fallthrough and FSQRT | `ev_8fe32daf30905b8abd67d5676c86c63b160697c9e0b441441d566a82df807ea4` |
| Floating argument load | `ev_eab62ee6e1c7dcde04546631477c4f565fcd8444406ebbe86a2caf34764fc763` |
| Constants 13, 11, 1, 2, 0.5 at 0x10061160 | `ev_5ca92566ad9ec353d3582e9ebe7bf5ce04a12ad7ee563ddd10983234ee9838da` |

The square-root implementation has additional special-value and error paths;
the statement above identifies its arithmetic operation without claiming a
complete reconstruction of that C runtime's exception or errno behavior.
