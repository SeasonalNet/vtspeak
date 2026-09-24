import ghidra.program.model.address.Address;

import java.nio.ByteBuffer;
import java.nio.ByteOrder;

public class DumpFloatConstants extends ghidra.app.script.GhidraScript {
    @Override
    protected void run() throws Exception {
        for (String value : getScriptArgs()) {
            Address address = toAddr(Long.decode(value));
            byte[] bytes = getBytes(address, 4);
            float result = ByteBuffer.wrap(bytes).order(ByteOrder.LITTLE_ENDIAN).getFloat();
            println(address + " float=" + result);
        }
    }
}
