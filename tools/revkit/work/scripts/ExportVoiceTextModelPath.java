import ghidra.app.decompiler.DecompInterface;
import ghidra.app.decompiler.DecompileResults;
import ghidra.program.model.listing.Function;
import ghidra.util.task.ConsoleTaskMonitor;
import java.io.File;
import java.io.PrintWriter;

public class ExportVoiceTextModelPath extends ghidra.app.script.GhidraScript {
    private static final String[] SELECTED_FUNCTIONS = {
        "FUN_1001c510", "FUN_10028620", "FUN_10012380", "FUN_10012a60",
        "FUN_10012580", "FUN_100126e0", "FUN_10012790", "FUN_100125d0",
        "FUN_100128c0", "FUN_10001050", "FUN_10001900", "FUN_1001ad20",
        "FUN_1001ae70", "FUN_10019f40", "FUN_10019e80", "FUN_10019700",
        "FUN_10019600", "FUN_100199f0", "FUN_10019b10", "FUN_10024d20",
        "FUN_10024e20", "FUN_10025440", "FUN_10063ed6"
    };

    @Override
    protected void run() throws Exception {
        File output = new File("/work/reports/vt_pau-model-loader.c");
        DecompInterface decompiler = new DecompInterface();
        decompiler.openProgram(currentProgram);
        try (PrintWriter writer = new PrintWriter(output, "UTF-8")) {
            writer.println("/* Ghidra decompiler output: approximate C pseudocode, not original source. */");
            writer.println("/* Program: " + currentProgram.getName() + " */");
            for (String name : SELECTED_FUNCTIONS) {
                Function function = null;
                for (Function candidate : currentProgram.getFunctionManager().getFunctions(true)) {
                    if (candidate.getName().equals(name)) {
                        function = candidate;
                        break;
                    }
                }
                writer.println();
                writer.println("/* ===== " + name + " ===== */");
                if (function == null) {
                    writer.println("/* FUNCTION NOT FOUND */");
                    continue;
                }
                writer.println("/* Entry: " + function.getEntryPoint() + " */");
                DecompileResults result = decompiler.decompileFunction(
                    function, 90, new ConsoleTaskMonitor());
                if (result.decompileCompleted()) {
                    writer.println(result.getDecompiledFunction().getC());
                } else {
                    writer.println("/* DECOMPILATION FAILED: " + result.getErrorMessage() + " */");
                }
            }
            println("Wrote model-loading and file-format paths to " + output.getPath());
        } finally {
            decompiler.dispose();
        }
    }
}
