import ghidra.app.decompiler.DecompInterface;
import ghidra.app.decompiler.DecompileResults;
import ghidra.program.model.listing.Function;
import ghidra.program.model.listing.FunctionIterator;
import ghidra.util.task.ConsoleTaskMonitor;
import java.io.File;
import java.io.PrintWriter;

public class DumpAllFunctions extends ghidra.app.script.GhidraScript {
    @Override
    protected void run() throws Exception {
        String[] args = getScriptArgs();
        if (args.length != 2) {
            throw new IllegalArgumentException(
                "usage: DumpAllFunctions.java <pseudocode-output> <inventory-output>");
        }

        DecompInterface decompiler = new DecompInterface();
        decompiler.openProgram(currentProgram);
        int total = 0;
        int completed = 0;
        int failed = 0;
        try (PrintWriter pseudocode = new PrintWriter(new File(args[0]), "UTF-8");
             PrintWriter inventory = new PrintWriter(new File(args[1]), "UTF-8")) {
            pseudocode.println("/* Ghidra pseudocode for all discovered functions.");
            pseudocode.println(" * Program: " + currentProgram.getName());
            pseudocode.println(" * Decompiled output is approximate, not original source. */");
            inventory.println("status\tname\tentry\tbody_bytes\tthunk\texternal\tdetail");

            FunctionIterator functions = currentProgram.getFunctionManager().getFunctions(true);
            while (functions.hasNext()) {
                monitor.checkCancelled();
                Function function = functions.next();
                total++;
                String entry = function.getEntryPoint().toString();
                long bodyBytes = function.getBody().getNumAddresses();
                String name = function.getName().replace('\t', ' ');
                boolean thunk = function.isThunk();
                boolean external = function.isExternal();
                pseudocode.println("\n/* ===== " + name + " @ " + entry + " ===== */");
                if (external) {
                    inventory.printf("EXTERNAL\t%s\t%s\t%d\t%s\ttrue\t\n",
                        name, entry, bodyBytes, thunk);
                    pseudocode.println("/* External function; no local body. */");
                    continue;
                }

                DecompileResults result = decompiler.decompileFunction(
                    function, 120, new ConsoleTaskMonitor());
                if (result.decompileCompleted()) {
                    completed++;
                    inventory.printf("DECOMPILED\t%s\t%s\t%d\t%s\tfalse\t\n",
                        name, entry, bodyBytes, thunk);
                    pseudocode.println(result.getDecompiledFunction().getC());
                } else {
                    failed++;
                    String detail = result.getErrorMessage() == null
                        ? "unknown decompiler failure"
                        : result.getErrorMessage().replace('\t', ' ').replace('\n', ' ');
                    inventory.printf("FAILED\t%s\t%s\t%d\t%s\tfalse\t%s\n",
                        name, entry, bodyBytes, thunk, detail);
                    pseudocode.println("/* DECOMPILATION FAILED: " + detail + " */");
                }
                if (total % 100 == 0) {
                    println("decompile progress: functions=" + total
                        + " completed=" + completed + " failed=" + failed);
                }
            }
            inventory.println("TOTAL\t" + total + "\tDECOMPILED\t" + completed
                + "\tFAILED\t" + failed);
            println("all-function pseudocode: functions=" + total
                + " completed=" + completed + " failed=" + failed);
        } finally {
            decompiler.dispose();
        }
    }
}
