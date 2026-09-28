import ghidra.app.decompiler.DecompInterface;
import ghidra.app.decompiler.DecompileResults;
import ghidra.program.model.address.Address;
import ghidra.program.model.listing.Data;
import ghidra.program.model.listing.Function;
import ghidra.program.model.listing.Listing;
import ghidra.program.model.symbol.Reference;
import ghidra.program.model.symbol.ReferenceIterator;
import ghidra.util.task.ConsoleTaskMonitor;
import java.io.File;
import java.io.PrintWriter;
import java.util.LinkedHashMap;
import java.util.Map;

public class TraceKateDictionaryRefs extends ghidra.app.script.GhidraScript {
    @Override
    protected void run() throws Exception {
        Map<String, Function> targets = new LinkedHashMap<>();
        try (PrintWriter out = new PrintWriter(new File("/work/corpus-parity/kate-dictionary-refs.txt"), "UTF-8")) {
            byte[] phoneTable = new byte[256 * 5];
            currentProgram.getMemory().getBytes(toAddr("1009e9a0"), phoneTable);
            for (int id = 0; id < 256; id++) {
                StringBuilder slot = new StringBuilder();
                for (int j = 0; j < 5; j++) slot.append(String.format("%02x", phoneTable[id * 5 + j] & 0xff));
                out.println("PHONE_ID_SLOT id=" + String.format("%02x", id) + " bytes=" + slot);
            }
            Listing listing = currentProgram.getListing();
            String[] needles = {"hashcont_emb", "hashidx_emb", "engttsdict_emb", "hashparams_emb", "exceptdict"};
            for (Data data : listing.getDefinedData(true)) {
                Object value = data.getValue();
                if (value == null) continue;
                String rendered = value.toString();
                for (String needle : needles) {
                    if (!rendered.contains(needle)) continue;
                    out.println("STRING address=" + data.getAddress() + " value=" + rendered);
                    ReferenceIterator refs = currentProgram.getReferenceManager().getReferencesTo(data.getAddress());
                    while (refs.hasNext()) {
                        Reference ref = refs.next();
                        Function f = currentProgram.getFunctionManager().getFunctionContaining(ref.getFromAddress());
                        out.println("  REF from=" + ref.getFromAddress() + " type=" + ref.getReferenceType() + " function=" + (f == null ? "<none>" : f.getName() + "@" + f.getEntryPoint()));
                        if (f != null) targets.put(f.getName() + "@" + f.getEntryPoint(), f);
                    }
                }
            }
            String[] globals = {"1008cbec", "1008cc10", "1008cc20", "1008cc24", "1008cbf8", "1008cc04", "1008cc1c", "1008cc18", "1008cbf0", "1008cc28", "1008cc08", "1008cc14", "1008cbf4", "1008cc0c", "1008cc00", "1008cbfc", "1009e9a0"};
            for (String global : globals) {
                Address address = toAddr(global);
                ReferenceIterator refs = currentProgram.getReferenceManager().getReferencesTo(address);
                while (refs.hasNext()) {
                    Reference ref = refs.next();
                    Function f = currentProgram.getFunctionManager().getFunctionContaining(ref.getFromAddress());
                    if (f != null) {
                        out.println("GLOBAL_REF address=" + address + " from=" + ref.getFromAddress() + " function=" + f.getName() + "@" + f.getEntryPoint());
                        targets.put(f.getName() + "@" + f.getEntryPoint(), f);
                    }
                }
            }
            String[] targetsOfInterest = {"1000e440", "10003eb0", "1000cd50", "100040a0", "1000bb90"};
            for (String target : targetsOfInterest) {
                Function targetFunction = currentProgram.getFunctionManager().getFunctionAt(toAddr(target));
                if (targetFunction != null) targets.put(targetFunction.getName() + "@" + targetFunction.getEntryPoint(), targetFunction);
                ReferenceIterator refs = currentProgram.getReferenceManager().getReferencesTo(toAddr(target));
                while (refs.hasNext()) {
                    Reference ref = refs.next();
                    Function f = currentProgram.getFunctionManager().getFunctionContaining(ref.getFromAddress());
                    out.println("FUNCTION_REF target=" + target + " from=" + ref.getFromAddress() + " type=" + ref.getReferenceType() + " function=" + (f == null ? "<none>" : f.getName() + "@" + f.getEntryPoint()));
                    if (f != null) targets.put(f.getName() + "@" + f.getEntryPoint(), f);
                }
            }
            DecompInterface decompiler = new DecompInterface();
            decompiler.openProgram(currentProgram);
            try {
                for (Function f : targets.values()) {
                    out.println("\n/* FUNCTION " + f.getName() + " @ " + f.getEntryPoint() + " */");
                    DecompileResults result = decompiler.decompileFunction(f, 120, new ConsoleTaskMonitor());
                    out.println(result.decompileCompleted() ? result.getDecompiledFunction().getC() : "DECOMPILE ERROR: " + result.getErrorMessage());
                }
            } finally {
                decompiler.dispose();
            }
        }
        println("Wrote dictionary string references and caller pseudocode.");
    }
}
