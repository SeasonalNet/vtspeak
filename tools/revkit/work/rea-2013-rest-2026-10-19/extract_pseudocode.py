#!/usr/bin/env python3
"""Extract one pseudocode listing per procedure from a REA evidence bundle."""

import json
import sys
from pathlib import Path


bundle = Path(sys.argv[1])
output = Path(sys.argv[2])
evidence = json.loads(bundle.read_text())
procedures = {}
for record in evidence["records"]:
    if record.get("operation") != "batch_decompile":
        continue
    for item in record.get("normalized_result", {}).get("items", []):
        address = item.get("procedure", {}).get("address", item.get("address"))
        procedures[address] = item

with output.open("w") as stream:
    stream.write("/* REA/Ghidra pseudocode for vt_pau.dll; not original source. */\n")
    stream.write("/* Source SHA-256: 200945bbb2853cc56b3eaa03770eedf26066b20bfc471ef5b0213387035b37ce */\n\n")
    for address, item in sorted(procedures.items(), key=lambda pair: int(pair[0], 16)):
        procedure = item.get("procedure", {})
        stream.write(f"/* Entry: {address}; {procedure.get('name', 'unresolved')} */\n")
        if item.get("status") == "ok":
            stream.write(item.get("pseudocode", ""))
        else:
            stream.write(f"/* REA status: {item.get('status')}; {item.get('error', 'no pseudocode')} */\n")
        stream.write("\n\n")

print(f"wrote {len(procedures)} unique procedure records to {output}")
