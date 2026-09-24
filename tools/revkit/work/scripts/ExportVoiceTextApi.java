import ghidra.app.decompiler.DecompInterface;
import ghidra.app.decompiler.DecompileResults;
import ghidra.program.model.listing.Function;
import ghidra.util.task.ConsoleTaskMonitor;
import java.io.File;
import java.io.PrintWriter;
import java.util.ArrayList;
import java.util.List;

public class ExportVoiceTextApi extends ghidra.app.script.GhidraScript {
    @Override
    protected void run() throws Exception {
        File output = new File("/work/reports/vt_pau-exported-api.c");
        List<Function> apiFunctions = new ArrayList<>();
        for (Function function : currentProgram.getFunctionManager().getFunctions(true)) {
            String name = function.getName();
            if (!function.isExternal() && (name.startsWith("VT_") || name.startsWith("VTDTTS_"))) {
                apiFunctions.add(function);
            }
        }

        DecompInterface decompiler = new DecompInterface();
        decompiler.openProgram(currentProgram);
        try (PrintWriter writer = new PrintWriter(output, "UTF-8")) {
            writer.println("/* Ghidra decompiler output; reconstructed pseudocode, not original source. */");
            writer.println("/* Program: " + currentProgram.getName() + " */");
            for (Function function : apiFunctions) {
                writer.println();
                writer.println("/* ===== " + function.getName() + " @ " + function.getEntryPoint() + " ===== */");
                DecompileResults result = decompiler.decompileFunction(
                    function, 90, new ConsoleTaskMonitor());
                if (result.decompileCompleted()) {
                    writer.println(result.getDecompiledFunction().getC());
                } else {
                    writer.println("/* DECOMPILATION FAILED: " + result.getErrorMessage() + " */");
                }
            }
            println("Wrote " + apiFunctions.size() + " API functions to " + output.getPath());
        } finally {
            decompiler.dispose();
        }
    }
}
