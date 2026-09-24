import ghidra.app.decompiler.DecompInterface;
import ghidra.app.decompiler.DecompileResults;
import ghidra.program.model.listing.Function;
import ghidra.util.task.ConsoleTaskMonitor;
import java.io.FileWriter;
import java.io.PrintWriter;

public class DumpDecompileToFile extends ghidra.app.script.GhidraScript {
    @Override
    protected void run() throws Exception {
        String[] args = getScriptArgs();
        DecompInterface decompiler = new DecompInterface();
        decompiler.openProgram(currentProgram);
        try (PrintWriter out = new PrintWriter(new FileWriter(args[0]))) {
            for (int i = 1; i < args.length; i++) {
                Function function = currentProgram.getFunctionManager()
                    .getFunctionContaining(toAddr(Long.decode(args[i])));
                out.println("===== " + args[i] + " =====");
                if (function == null) {
                    out.println("Function not found");
                    continue;
                }
                out.println("Function: " + function.getName() + " @ " + function.getEntryPoint());
                DecompileResults result = decompiler.decompileFunction(
                    function, 90, new ConsoleTaskMonitor());
                if (result.decompileCompleted()) {
                    out.println(result.getDecompiledFunction().getC());
                } else {
                    out.println("Decompile error: " + result.getErrorMessage());
                }
            }
        } finally {
            decompiler.dispose();
        }
    }
}
