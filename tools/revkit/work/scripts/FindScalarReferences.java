import ghidra.program.model.listing.Function;
import ghidra.program.model.listing.Instruction;
import ghidra.program.model.scalar.Scalar;

public class FindScalarReferences extends ghidra.app.script.GhidraScript {
    @Override
    protected void run() throws Exception {
        long[] targets = new long[getScriptArgs().length];
        for (int i = 0; i < targets.length; i++) {
            targets[i] = Long.decode(getScriptArgs()[i]);
        }
        for (Instruction instruction : currentProgram.getListing().getInstructions(true)) {
            for (int operand = 0; operand < instruction.getNumOperands(); operand++) {
                for (Object object : instruction.getOpObjects(operand)) {
                    if (object instanceof Scalar scalar) {
                        for (long target : targets) {
                            if (scalar.getUnsignedValue() == target) {
                                Function function = currentProgram.getFunctionManager()
                                    .getFunctionContaining(instruction.getAddress());
                                println("constant=" + Long.toHexString(target) + " " +
                                    instruction.getAddress() + " " +
                                    (function == null ? "<no function>" : function.getName()) +
                                    " : " + instruction);
                            }
                        }
                    }
                }
            }
        }
    }
}
