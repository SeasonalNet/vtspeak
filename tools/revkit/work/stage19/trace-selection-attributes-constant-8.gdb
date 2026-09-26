set pagination off
set confirm off
set debuginfod enabled off
handle SIGSEGV nostop noprint pass
set logging file /probe/index-attributes-gdb.log
set logging overwrite on
set logging enabled on
set $attribute_calls = 0

break *0x100182e0
commands
  silent
  set $attribute_calls = $attribute_calls + 1
  if $attribute_calls <= 40
    set $unit = *(unsigned int *)($esp + 4)
    set $model = *(unsigned int *)($esp + 24)
    set $target = *(unsigned int *)($esp + 16)
    set $column_a = *(unsigned int *)($model + 0x68)
    set $column_b = *(unsigned int *)($model + 0x60)
    set $category_table = *(unsigned int *)($model + 0x28)
    set $pair_table = *(unsigned int *)($model + 0x18)
    set $a = *(unsigned char *)($column_a + $unit)
    set $b = *(unsigned char *)($column_b + $unit)
    set $target_a = *(unsigned char *)$target
    set $category = *(unsigned char *)($category_table + $a)
    set $row = *(unsigned int *)($pair_table + $b * 4)
    printf "INDEX_ATTR unit=%u attr_a=%u category_a=%u attr_b=%u target_b=%u pair_cost=%g\n", $unit, $a, $category, $b, $target_a, *(float *)($row + $target_a * 4)
  end
  continue
end

continue
