import ghidra.program.model.address.Address;
import ghidra.program.model.listing.Function;
import ghidra.app.decompiler.DecompInterface;
import ghidra.app.decompiler.DecompileResults;
import ghidra.util.task.ConsoleTaskMonitor;

public class InspectDecodeTable extends ghidra.app.script.GhidraScript {
    @Override
    protected void run() throws Exception {
        Address table = toAddr(0x10001e70L);
        DecompInterface decompiler = new DecompInterface();
        decompiler.openProgram(currentProgram);
        for (int index = 0; index < 9; index++) {
            int target = currentProgram.getMemory().getInt(table.add(index * 4));
            Address address = toAddr(Integer.toUnsignedLong(target));
            Function function = currentProgram.getFunctionManager().getFunctionContaining(address);
            if (function == null) {
                disassemble(address);
                function = createFunction(address, "decode_mode_" + index);
            }
            println(index + " -> 0x" + Long.toHexString(Integer.toUnsignedLong(target))
                + (function == null ? " (no function)" : " (" + function.getName() + ")"));
            if (function != null) {
                DecompileResults result = decompiler.decompileFunction(
                    function, 60, new ConsoleTaskMonitor());
                if (result.decompileCompleted()) {
                    println(result.getDecompiledFunction().getC());
                } else {
                    println("Decompile error: " + result.getErrorMessage());
                }
            }
        }
        decompiler.dispose();
    }
}
