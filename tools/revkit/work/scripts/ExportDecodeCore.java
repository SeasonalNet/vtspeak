import ghidra.app.decompiler.DecompInterface;
import ghidra.app.decompiler.DecompileResults;
import ghidra.program.model.listing.Function;
import ghidra.util.task.ConsoleTaskMonitor;
import java.io.File;
import java.io.PrintWriter;

public class ExportDecodeCore extends ghidra.app.script.GhidraScript {
    private static final long[] ENTRIES = {
        0x10001b30L, 0x10001ea0L, 0x10002170L, 0x100020a0L,
        0x100020d0L, 0x10002100L, 0x10001fb0L, 0x10001c59L,
        0x10001c75L, 0x10001cb1L, 0x10001ce7L, 0x10001d20L,
        0x1002c120L
    };

    @Override
    protected void run() throws Exception {
        File output = new File("/work/reports/vt_pau-dat-decoder.c");
        DecompInterface decompiler = new DecompInterface();
        decompiler.openProgram(currentProgram);
        try (PrintWriter writer = new PrintWriter(output, "UTF-8")) {
            writer.println("/* Ghidra decompiler output: approximate C pseudocode, not original source. */");
            for (long entry : ENTRIES) {
                Function function = currentProgram.getFunctionManager()
                    .getFunctionContaining(toAddr(entry));
                if (function == null) {
                    disassemble(toAddr(entry));
                    function = createFunction(toAddr(entry), "decode_" + Long.toHexString(entry));
                }
                writer.println("\n/* ===== 0x" + Long.toHexString(entry) + " ===== */");
                if (function == null) {
                    writer.println("/* FUNCTION NOT FOUND */");
                    continue;
                }
                writer.println("/* Entry: " + function.getEntryPoint() + " " + function.getName() + " */");
                DecompileResults result = decompiler.decompileFunction(
                    function, 90, new ConsoleTaskMonitor());
                if (result.decompileCompleted()) {
                    writer.println(result.getDecompiledFunction().getC());
                } else {
                    writer.println("/* DECOMPILATION FAILED: " + result.getErrorMessage() + " */");
                }
            }
            println("Wrote DAT decoder pseudocode to " + output.getPath());
        } finally {
            decompiler.dispose();
        }
    }
}
