import ghidra.app.decompiler.DecompInterface;
import ghidra.app.decompiler.DecompileResults;
import ghidra.program.model.listing.Function;
import ghidra.program.model.listing.FunctionIterator;
import ghidra.util.task.ConsoleTaskMonitor;

import java.io.FileWriter;
import java.io.PrintWriter;
import java.util.regex.Matcher;
import java.util.regex.Pattern;

public class ScanStructOffsets extends ghidra.app.script.GhidraScript {
    private static final Pattern FIELD_ACCESS = Pattern.compile(
        "[^\\n]*(?:param_\\d+|iVar\\d+|uVar\\d+) \\+ 0x(?:c|10|3c|48|54|60|68|8c|90)[^\\n]*"
    );

    @Override
    protected void run() throws Exception {
        String outputPath = getScriptArgs()[0];
        DecompInterface decompiler = new DecompInterface();
        decompiler.openProgram(currentProgram);
        int total = 0;
        int matched = 0;
        try (PrintWriter out = new PrintWriter(new FileWriter(outputPath))) {
            FunctionIterator functions = currentProgram.getFunctionManager().getFunctions(true);
            while (functions.hasNext() && !monitor.isCancelled()) {
                Function function = functions.next();
                total++;
                DecompileResults result = decompiler.decompileFunction(
                    function, 20, new ConsoleTaskMonitor());
                if (!result.decompileCompleted()) {
                    continue;
                }
                String code = result.getDecompiledFunction().getC();
                Matcher matcher = FIELD_ACCESS.matcher(code);
                StringBuilder hits = new StringBuilder();
                while (matcher.find()) {
                    hits.append(matcher.group()).append('\n');
                }
                if (hits.length() != 0) {
                    matched++;
                    out.println("===== " + function.getName() + " @ " + function.getEntryPoint()
                        + " =====");
                    out.print(hits);
                }
            }
            out.println("FUNCTIONS_SCANNED=" + total);
            out.println("FUNCTIONS_MATCHED=" + matched);
        } finally {
            decompiler.dispose();
        }
    }
}
