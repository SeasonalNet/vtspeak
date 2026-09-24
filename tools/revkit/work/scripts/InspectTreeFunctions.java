import ghidra.program.model.address.Address;
import ghidra.program.model.listing.Function;
import ghidra.program.model.symbol.Reference;

public class InspectTreeFunctions extends ghidra.app.script.GhidraScript {
    @Override
    protected void run() throws Exception {
        Address start = toAddr(0x10001500L);
        Address end = toAddr(0x10002000L);
        for (Function function : currentProgram.getFunctionManager().getFunctions(start, true)) {
            if (function.getEntryPoint().compareTo(end) >= 0) {
                break;
            }
            println("FUNCTION " + function.getName() + " @ " + function.getEntryPoint());
            for (Reference reference : getReferencesTo(function.getEntryPoint())) {
                Function caller = currentProgram.getFunctionManager().getFunctionContaining(
                    reference.getFromAddress());
                if (caller != null) {
                    println("  XREF from " + caller.getName() + " @ " + caller.getEntryPoint()
                        + " at " + reference.getFromAddress());
                }
            }
        }
    }
}
