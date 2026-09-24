import ghidra.app.decompiler.DecompInterface;
import ghidra.app.decompiler.DecompileResults;
import ghidra.program.model.listing.Function;
import ghidra.util.task.ConsoleTaskMonitor;

public class DecompileNamedFunctions extends ghidra.app.script.GhidraScript {
    @Override
    protected void run() throws Exception {
        DecompInterface decompiler = new DecompInterface();
        decompiler.openProgram(currentProgram);
        try {
            for (String name : getScriptArgs()) {
                println("===== " + name + " =====");
                Function function;
                if (name.startsWith("0x")) {
                    function = currentProgram.getFunctionManager().getFunctionContaining(
                        toAddr(Long.decode(name)));
                }
                else {
                    function = null;
                    for (Function candidate : currentProgram.getFunctionManager().getFunctions(true)) {
                        if (candidate.getName().equals(name)) {
                            function = candidate;
                            break;
                        }
                    }
                }
                if (function == null) {
                    println("Function not found");
                    continue;
                }
                println("Function: " + function.getName() + " @ " + function.getEntryPoint());
                DecompileResults result = decompiler.decompileFunction(
                    function, 90, new ConsoleTaskMonitor());
                if (result.decompileCompleted()) {
                    println(result.getDecompiledFunction().getC());
                } else {
                    println("Decompile error: " + result.getErrorMessage());
                }
            }
        } finally {
            decompiler.dispose();
        }
    }
}
