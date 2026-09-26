set pagination off
set confirm off
set debuginfod enabled off

break *0x100197f0
commands
  silent
  printf "legacy unit-index reader first iteration ebx=%u edi=%#x esi=%#x\n", $ebx, $edi, $esi
  printf "reader state: unit_count=%u record_stride=%u file_fields=%#x/%#x/%#x/%#x\n", *(unsigned int *)($edi+0x4c), *(unsigned short *)($edi+0x52), *(unsigned int *)($esi+4), *(unsigned int *)($esi+8), *(unsigned int *)($esi+0x10), *(unsigned int *)($esi+0x14)
  printf "unit destinations: attr_a=%#x signature=%#x attr_b=%#x group1=%#x/%#x/%#x group2=%#x/%#x/%#x group3=%#x/%#x/%#x\n", *(unsigned int *)($edi+0x48), *(unsigned int *)($edi+0x44), *(unsigned int *)($edi+0x40), *(unsigned int *)($edi+0x1c), *(unsigned int *)($edi+0x28), *(unsigned int *)($edi+0x34), *(unsigned int *)($edi+0x20), *(unsigned int *)($edi+0x2c), *(unsigned int *)($edi+0x38), *(unsigned int *)($edi+0x24), *(unsigned int *)($edi+0x30), *(unsigned int *)($edi+0x3c)
  x/24wx $esi-16
  x/24wx $edi+0x40
  bt
  kill
end

continue
