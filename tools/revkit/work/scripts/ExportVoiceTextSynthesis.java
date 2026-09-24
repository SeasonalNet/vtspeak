import ghidra.app.decompiler.DecompInterface;
import ghidra.app.decompiler.DecompileResults;
import ghidra.program.model.listing.Function;
import ghidra.util.task.ConsoleTaskMonitor;
import java.io.File;
import java.io.PrintWriter;

public class ExportVoiceTextSynthesis extends ghidra.app.script.GhidraScript {
    private static final String[] SELECTED_FUNCTIONS = {
        "VT_TextToFile_ENG", "FUN_1001e830", "FUN_1001e930", "FUN_1001f1f0",
        "FUN_1001f600", "FUN_1001c990", "FUN_10017510", "FUN_10013380",
        "FUN_100138c0", "FUN_10022dc0", "FUN_10026750", "FUN_10016c90",
        "FUN_100130e0", "FUN_10026630", "FUN_100267d0", "FUN_1002c9b0",
        "FUN_1002bd90", "FUN_1002d230", "FUN_10024680", "FUN_10024060",
        "FUN_100242a0", "FUN_10018770", "FUN_10023af0", "FUN_10023c70",
        "FUN_1002c8b0", "FUN_1002bbd0", "FUN_10001b30", "FUN_1002d010",
        "FUN_1002cdf0", "FUN_1001b200", "FUN_1002d1e0"
    };

    @Override
    protected void run() throws Exception {
        File output = new File("/work/reports/vt_pau-synthesis-pseudocode.c");
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
            println("Wrote model-backed synthesis paths to " + output.getPath());
        } finally {
            decompiler.dispose();
        }
    }
}
