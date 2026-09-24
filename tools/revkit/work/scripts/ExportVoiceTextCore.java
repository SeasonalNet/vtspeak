import ghidra.app.decompiler.DecompInterface;
import ghidra.app.decompiler.DecompileResults;
import ghidra.program.model.listing.Function;
import ghidra.util.task.ConsoleTaskMonitor;
import java.io.File;
import java.io.PrintWriter;

public class ExportVoiceTextCore extends ghidra.app.script.GhidraScript {
    private static final String[] SELECTED_FUNCTIONS = {
        "VT_LOADTTS_EXT_ENG",
        "VT_UNLOADTTS_EXT_ENG",
        "VT_CheckLicense_ENG",
        "VT_TextToFile_ENG",
        "VT_TextToBuffer_ENG",
        "FUN_1001e830",
        "FUN_1001e930",
        "FUN_1001f1f0",
        "FUN_1001f600",
        "FUN_1001c990",
        "FUN_10026ab0",
        "FUN_10026870",
        "FUN_1001e0c0",
        "FUN_10024d20",
        "FUN_10024e20",
        "FUN_10025370",
        "FUN_10025400",
        "FUN_10025500",
        "FUN_100255a0",
        "FUN_10025560",
        "FUN_10025600"
    };

    @Override
    protected void run() throws Exception {
        File output = new File("/work/reports/vt_pau-core-pseudocode.c");
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
            println("Wrote selected engine paths to " + output.getPath());
        } finally {
            decompiler.dispose();
        }
    }
}
