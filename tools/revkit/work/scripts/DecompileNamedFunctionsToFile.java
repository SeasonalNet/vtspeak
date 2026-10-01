import ghidra.app.decompiler.DecompInterface;
import ghidra.app.decompiler.DecompileResults;
import ghidra.program.model.listing.Function;
import ghidra.util.task.ConsoleTaskMonitor;
import java.nio.charset.StandardCharsets;
import java.nio.file.Files;
import java.nio.file.Path;

public class DecompileNamedFunctionsToFile extends ghidra.app.script.GhidraScript {
    @Override
    protected void run() throws Exception {
        String[] args = getScriptArgs();
        if (args.length < 2) {
            throw new IllegalArgumentException("Expected output path and at least one function");
        }

        StringBuilder report = new StringBuilder();
        DecompInterface decompiler = new DecompInterface();
        decompiler.openProgram(currentProgram);
        try {
            for (int i = 1; i < args.length; i++) {
                String name = args[i];
                report.append("===== ").append(name).append(" =====\n");
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
                    report.append("Function not found\n\n");
                    continue;
                }

                report.append("Function: ").append(function.getName())
                    .append(" @ ").append(function.getEntryPoint()).append('\n');
                DecompileResults result = decompiler.decompileFunction(
                    function, 90, new ConsoleTaskMonitor());
                if (result.decompileCompleted()) {
                    report.append(result.getDecompiledFunction().getC()).append('\n');
                }
                else {
                    report.append("Decompile error: ").append(result.getErrorMessage()).append('\n');
                }
                report.append('\n');
            }
        }
        finally {
            decompiler.dispose();
        }

        Path output = Path.of(args[0]);
        Files.writeString(output, report.toString(), StandardCharsets.UTF_8);
        println("Wrote decompilation report to " + output);
    }
}
