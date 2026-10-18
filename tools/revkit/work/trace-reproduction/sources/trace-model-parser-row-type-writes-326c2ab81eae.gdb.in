set pagination off
set confirm off
set debuginfod enabled off
set can-use-hw-watchpoints 1
handle SIGSEGV nostop noprint pass

python
import gdb

row_type_watches = []

class ModelRowTypeWatch(gdb.Breakpoint):
    def __init__(self, row_address, row_index):
        self.row_address = row_address
        self.row_index = row_index
        super().__init__(
            "*(unsigned int *)0x%x" % row_address,
            gdb.BP_WATCHPOINT,
            gdb.WP_WRITE,
            internal=False,
        )

    def stop(self):
        value = int(gdb.parse_and_eval("*(unsigned int *)0x%x" % self.row_address))
        pc = int(gdb.parse_and_eval("$pc"))
        gdb.write(
            "MODEL_ROW_TYPE_WRITE index=%d value=%08x pc=%08x\n"
            % (self.row_index, value, pc)
        )
        gdb.execute("bt 6")
        return False

class ModelRowTypeWatchSetup(gdb.Breakpoint):
    def stop(self):
        state = int(gdb.parse_and_eval("*(unsigned int *)($esp + 4)"))
        count = int(gdb.parse_and_eval("*(short *)0x%x" % state))
        gdb.write(
            "MODEL_ROW_TYPE_WATCH state=%08x first_index=%d\n" % (state, count)
        )
        for index in range(max(0, count), min(100, max(0, count) + 4)):
            row_address = state + 0x14 + index * 0x94 + 0x2c
            row_type_watches.append(ModelRowTypeWatch(row_address, index))
        self.enabled = False
        return False

ModelRowTypeWatchSetup("*0x1003d3d0")
end

break *0x1003e240
commands
  silent
  set $state = *(unsigned int *)($esp + 4)
  set $count = *(short *)$state
  printf "MODEL_ROW_TYPE_FINAL count=%d\n", $count
  set $index = 0
  while $index < $count && $index < 100
    set $row = $state + 0x14 + $index * 0x94
    printf "MODEL_ROW_TYPE_VALUE index=%d class=%02x type=%08x source=%d:%d\n", $index, *(unsigned char *)($row + 0x24), *(unsigned int *)($row + 0x2c), *(int *)$row, *(int *)($row + 4)
    set $index = $index + 1
  end
  continue
end

hbreak *0x408187
commands
  silent
  printf "MODEL_ROW_TYPE_TRACE_READY\n"
  continue
end

continue
