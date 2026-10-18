set pagination off
set confirm off
set debuginfod enabled off
handle SIGSEGV nostop noprint pass

break *0x10023350
commands
  silent
  set $context = *(int *)($esp + 4)
  if $context == 1
    set $state = *(unsigned int *)($esp + 8)
    set $model = *(unsigned int *)($esp + 12)
    set $record = $state + 0xae894 + $context * 0xfc
    set $count_ptr = *(unsigned int *)($model + 0x23 * 4)
    set $units_ptr = *(unsigned int *)($model + 0x25 * 4)
    set $class_count = *(unsigned short *)$record
    printf "EXPANSION_ENTRY context=%d classes=%u count_table=%#x unit_lists=%#x\n", $context, $class_count, $count_ptr, $units_ptr
    set $i = 0
    while $i < $class_count && $i < 30
      set $class_id = *(unsigned int *)($record + 4 + $i * 4)
      set $unit_count = *(unsigned short *)($count_ptr + $class_id * 2)
      set $unit_list = *(unsigned int *)($units_ptr + $class_id * 4)
      printf "CLASS_EXPANSION context=%d slot=%d class=%u units=%u ids=", $context, $i, $class_id, $unit_count
      set $j = 0
      while $j < $unit_count && $j < 100
        printf "%u ", *(unsigned int *)($unit_list + $j * 4)
        set $j = $j + 1
      end
      printf "\n"
      set $i = $i + 1
    end
  end
  continue
end

continue
