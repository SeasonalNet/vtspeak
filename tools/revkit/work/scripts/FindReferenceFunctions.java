import ghidra.program.model.address.Address;
import ghidra.program.model.listing.Function;
import ghidra.program.model.symbol.Reference;
import ghidra.program.model.symbol.ReferenceIterator;

public class FindReferenceFunctions extends ghidra.app.script.GhidraScript {
    @Override
    protected void run() throws Exception {
        for (String value : getScriptArgs()) {
            Address target = toAddr(Long.decode(value));
            println("===== references to " + target + " =====");
            ReferenceIterator references = currentProgram.getReferenceManager().getReferencesTo(target);
            while (references.hasNext()) {
                Reference reference = references.next();
                Function function = currentProgram.getFunctionManager()
                    .getFunctionContaining(reference.getFromAddress());
                println(reference.getFromAddress() + " " +
                    (function == null ? "<no function>" : function.getName() + " @ " + function.getEntryPoint()));
            }
        }
    }
}
