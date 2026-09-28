import ghidra.app.decompiler.DecompInterface;
import ghidra.app.decompiler.DecompileResults;
import ghidra.program.model.address.Address;
import ghidra.program.model.listing.Function;
import ghidra.util.task.ConsoleTaskMonitor;
import java.io.File;
import java.io.PrintWriter;

public class ExportOldKateEntry extends ghidra.app.script.GhidraScript {
    @Override
    protected void run() throws Exception {
        String[] names = {"VT_GetTTSInfo_ENG", "VT_LOADTTS_ENG", "VT_TextToFile_ENG", "VT_UNLOADTTS_ENG", "FUN_10017320", "FUN_10016dd0", "FUN_10016f60", "FUN_10017140", "FUN_10017820", "FUN_10017420", "FUN_10020750", "FUN_100202f0", "FUN_1001f8c0", "FUN_1001f8b0", "FUN_1000e930", "FUN_1000e820", "FUN_1000e880", "FUN_1000e8a0", "FUN_10003550", "FUN_10003710", "FUN_1003cab0", "FUN_1003c3a0", "FUN_10025060", "FUN_100250c0", "FUN_10025090", "FUN_10014f70", "FUN_10022430", "FUN_10022560", "FUN_1001c180", "FUN_1001c440", "FUN_1001c740", "FUN_100546d0", "FUN_10054700", "FUN_100547d0", "FUN_10055050", "FUN_1004edc0", "FUN_1004a0b0", "FUN_10050290", "FUN_100344a0", "FUN_100552e0", "FUN_100512a0", "FUN_10054420", "FUN_10055e60", "FUN_10058300", "FUN_10057930", "FUN_10054220", "FUN_10020510", "FUN_10020460", "FUN_10015370", "FUN_100157b0", "FUN_10015b50", "FUN_10016c00", "FUN_1001f950", "FUN_1000eff0", "FUN_10001000", "FUN_10013820", "FUN_10013970", "FUN_100129f0", "FUN_100137a0", "FUN_10013760", "FUN_1000eb20", "FUN_1000ec70", "FUN_1000ed30", "FUN_1000ee40", "FUN_1000f3e0", "FUN_1000f5a0", "FUN_1000fd30", "FUN_1001b800", "FUN_10001570", "FUN_100019b0", "FUN_10022f60", "FUN_1000b340", "FUN_10021f10"};
        String[] missingAddresses = {"1001c440", "1001c740", "100546d0"};
        for (String value : missingAddresses) {
            Address address = toAddr(value);
            if (currentProgram.getFunctionManager().getFunctionAt(address) == null) {
                createFunction(address, "FUN_" + value);
            }
        }
        DecompInterface decompiler = new DecompInterface();
        decompiler.openProgram(currentProgram);
        try (PrintWriter out = new PrintWriter(new File("/work/corpus-parity/old-kate-api-pseudocode.c"), "UTF-8")) {
            out.println("/* Ghidra pseudocode for the extracted 2006 MSI DLL; names/control flow may be imperfect. */");
            out.println("/* Program: " + currentProgram.getName() + " */");
            for (String name : names) {
                Function f = null;
                for (Function candidate : currentProgram.getFunctionManager().getFunctions(true)) {
                    if (candidate.getName().equals(name)) { f = candidate; break; }
                }
                out.println("\n/* ===== " + name + " ===== */");
                if (f == null) { out.println("/* FUNCTION NOT FOUND */"); continue; }
                out.println("/* Entry: " + f.getEntryPoint() + " */");
                DecompileResults result = decompiler.decompileFunction(f, 120, new ConsoleTaskMonitor());
                if (result.decompileCompleted()) out.println(result.getDecompiledFunction().getC());
                else out.println("/* DECOMPILATION FAILED: " + result.getErrorMessage() + " */");
            }
            println("Wrote old Kate API pseudocode to /work/corpus-parity/old-kate-api-pseudocode.c");
        } finally { decompiler.dispose(); }
    }
}
