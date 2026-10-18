set pagination off
set confirm off
set debuginfod enabled off
handle SIGSEGV nostop noprint pass

break *0x1001da50
commands
  silent
  set $text = *(unsigned int *)($esp + 8)
  set $speaker = *(int *)($esp + 16)
  set $pitch = *(int *)($esp + 20)
  set $speed = *(int *)($esp + 24)
  set $volume = *(int *)($esp + 28)
  set $pause = *(int *)($esp + 32)
  set $dict = *(int *)($esp + 36)
  set $texttype = *(int *)($esp + 40)
  disable 1
  set $buffer = (char *)malloc(65536)
  set $length = (int *)malloc(4)
  set *$length = 65536
  set $baseline_result = ((int (*)(int, char *, char *, int *, int, int, int, int, int, int, int, int, int))0x1001dc40)(0, (char *)$text, $buffer, $length, 0, 0, $speaker, $pitch, $speed, $volume, $pause, $dict, $texttype)
  printf "BUFFER_CAPACITY baseline_return=%d output_len=%d\n", $baseline_result, *$length
  dump binary memory /work/stage16/capacity-baseline.bin $buffer $buffer + *$length
  set $capacity_index = 0
  while $capacity_index < 5
    if $capacity_index == 0
      set $capacity = 0
    end
    if $capacity_index == 1
      set $capacity = 1
    end
    if $capacity_index == 2
      set $capacity = 23605
    end
    if $capacity_index == 3
      set $capacity = 23606
    end
    if $capacity_index == 4
      set $capacity = 23607
    end
    set $baseline_byte = *(unsigned char *)($buffer + $capacity)
    set $sentinel = $baseline_byte ^ 255
    set *(unsigned char *)($buffer + $capacity) = $sentinel
    set *$length = $capacity
    set $result = ((int (*)(int, char *, char *, int *, int, int, int, int, int, int, int, int, int))0x1001dc40)(0, (char *)$text, $buffer, $length, 0, 0, $speaker, $pitch, $speed, $volume, $pause, $dict, $texttype)
    set $observed_byte = *(unsigned char *)($buffer + $capacity)
    printf "BUFFER_CAPACITY supplied=%d return=%d output_len=%d boundary_before=%u boundary_after=%u\n", $capacity, $result, *$length, $sentinel, $observed_byte
    if $capacity_index == 0
      dump binary memory /work/stage16/capacity-0.bin $buffer $buffer + *$length
    end
    if $capacity_index == 1
      dump binary memory /work/stage16/capacity-1.bin $buffer $buffer + *$length
    end
    if $capacity_index == 2
      dump binary memory /work/stage16/capacity-23605.bin $buffer $buffer + *$length
    end
    if $capacity_index == 3
      dump binary memory /work/stage16/capacity-23606.bin $buffer $buffer + *$length
    end
    if $capacity_index == 4
      dump binary memory /work/stage16/capacity-23607.bin $buffer $buffer + *$length
    end
    set $capacity_index = $capacity_index + 1
  end
  continue
end

continue
