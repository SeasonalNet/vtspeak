set pagination off
set confirm off
set debuginfod enabled off
handle SIGSEGV nostop noprint pass

set $duplicate_audit_done = 0

# Scan the sorted five-byte class table once after model load. Report contiguous
# equal-key runs without changing the model or the query being synthesized.
break *0x10023de1
commands
  silent
  if !$duplicate_audit_done
    set $duplicate_audit_done = 1
    set $duplicate_table = *(unsigned int *)($ebp + 0x14) + 0x88
    set $duplicate_keys = *(unsigned int *)($duplicate_table + 8)
    set $duplicate_count = *(unsigned int *)$duplicate_table
    python
import gdb
table = int(gdb.parse_and_eval("$duplicate_keys"))
count = int(gdb.parse_and_eval("$duplicate_count"))
keys = bytes(gdb.selected_inferior().read_memory(table, count * 5))
groups = []
index = 0
while index < count:
    key = keys[index * 5:index * 5 + 5]
    end = index + 1
    while end < count and keys[end * 5:end * 5 + 5] == key:
        end += 1
    if end - index > 1:
        groups.append((index, end - index, key))
    index = end
gdb.write("CLASS_DUPLICATE_AUDIT key_count=%d\\n" % count)
for first, run, key in groups[:20]:
    gdb.write("CLASS_DUPLICATE_RANGE first_index=%d count=%d key=%s\\n" %
              (first, run, list(key)))
gdb.write("CLASS_DUPLICATE_SUMMARY groups=%d records_in_groups=%d max_group=%d printed=%d\\n" %
          (len(groups), sum(run for _, run, _ in groups),
           max((run for _, run, _ in groups), default=1), min(20, len(groups))))
    end
  end
  continue
end

continue
