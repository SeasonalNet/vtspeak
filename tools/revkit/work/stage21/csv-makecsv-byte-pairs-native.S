.text
.globl probe_start
.globl probe_end
.type probe_start, @function
probe_start:
    pushl %ebp
    movl %esp, %ebp
    pushl %ebx
    pushl %esi
    pushl %edi
    subl $20, %esp

    call .Lpc
.Lpc:
    popl %ecx
    movl ptr_fields-.Lpc(%ecx), %esi
    movl %esi, -20(%ebp)
    movl ptr_field-.Lpc(%ecx), %edi
    movl %edi, -32(%ebp)
    movl ptr_raw-.Lpc(%ecx), %ebx
    movl ptr_results-.Lpc(%ecx), %eax
    movl %eax, -16(%ebp)
    movl $0, -24(%ebp)
    movl %eax, %edx
    movl $0, 917506(%edx)

.Louter:
    movl $0, -28(%ebp)
.Linner:
    movl -32(%ebp), %edi
    movl -16(%ebp), %edx
    movl -24(%ebp), %eax
    movb %al, 917504(%edx)
    movl -28(%ebp), %eax
    movb %al, 917505(%edx)
    movl -24(%ebp), %eax
    movb %al, 0(%edi)
    movl -28(%ebp), %eax
    movb %al, 1(%edi)
    movb $0, 2(%edi)
    movb $0x5a, 0(%ebx)
    movb $0xa5, 9(%ebx)

    pushl $8
    leal 1(%ebx), %eax
    pushl %eax
    pushl $1
    movl -20(%ebp), %esi
    pushl %esi
    movl $0x10016c50, %eax
    call *%eax
    addl $16, %esp
    movswl %ax, %eax

    movl -24(%ebp), %ecx
    shll $8, %ecx
    addl -28(%ebp), %ecx
    imull $14, %ecx, %ecx
    movl -16(%ebp), %edi
    addl %ecx, %edi
    movl %eax, 0(%edi)

    leal 1(%ebx), %esi
    xorl %ecx, %ecx
.Lcopy:
    movb (%esi,%ecx,1), %al
    movb %al, 4(%edi,%ecx,1)
    incl %ecx
    cmpl $8, %ecx
    jne .Lcopy
    movb 0(%ebx), %al
    movb %al, 12(%edi)
    movb 9(%ebx), %al
    movb %al, 13(%edi)

    movl -24(%ebp), %eax
    shll $8, %eax
    addl -28(%ebp), %eax
    incl %eax
    movl -16(%ebp), %edx
    movl %eax, 917506(%edx)

    incl -28(%ebp)
    cmpl $256, -28(%ebp)
    jl .Linner
    incl -24(%ebp)
    cmpl $256, -24(%ebp)
    jl .Louter

    xorl %eax, %eax
    addl $20, %esp
    popl %edi
    popl %esi
    popl %ebx
    leave
    ret
probe_end:
.size probe_start, probe_end-probe_start

.globl ptr_fields
ptr_fields:
    .long 0
.globl ptr_field
ptr_field:
    .long 0
.globl ptr_raw
ptr_raw:
    .long 0
.globl ptr_results
ptr_results:
    .long 0
