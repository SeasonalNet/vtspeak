set pagination off
set confirm off
set debuginfod enabled off
handle SIGSEGV nostop noprint pass

break *0x1001da50
commands
  silent
  disable 1
  set $text_digits = ((char *(*)(unsigned int))0x1001d9c0)(128)
  set {char}($text_digits + 0) = 52
  set {char}($text_digits + 1) = 50
  set {char}($text_digits + 2) = 46
  set {char}($text_digits + 3) = 0
  set $path_digits = ((char *(*)(unsigned int))0x1001d9c0)(128)
  set {char}($path_digits + 0) = 90
  set {char}($path_digits + 1) = 58
  set {char}($path_digits + 2) = 47
  set {char}($path_digits + 3) = 119
  set {char}($path_digits + 4) = 111
  set {char}($path_digits + 5) = 114
  set {char}($path_digits + 6) = 107
  set {char}($path_digits + 7) = 47
  set {char}($path_digits + 8) = 115
  set {char}($path_digits + 9) = 116
  set {char}($path_digits + 10) = 97
  set {char}($path_digits + 11) = 103
  set {char}($path_digits + 12) = 101
  set {char}($path_digits + 13) = 49
  set {char}($path_digits + 14) = 54
  set {char}($path_digits + 15) = 47
  set {char}($path_digits + 16) = 109
  set {char}($path_digits + 17) = 105
  set {char}($path_digits + 18) = 45
  set {char}($path_digits + 19) = 116
  set {char}($path_digits + 20) = 101
  set {char}($path_digits + 21) = 120
  set {char}($path_digits + 22) = 116
  set {char}($path_digits + 23) = 45
  set {char}($path_digits + 24) = 100
  set {char}($path_digits + 25) = 105
  set {char}($path_digits + 26) = 103
  set {char}($path_digits + 27) = 105
  set {char}($path_digits + 28) = 116
  set {char}($path_digits + 29) = 115
  set {char}($path_digits + 30) = 45
  set {char}($path_digits + 31) = 50
  set {char}($path_digits + 32) = 48
  set {char}($path_digits + 33) = 50
  set {char}($path_digits + 34) = 54
  set {char}($path_digits + 35) = 48
  set {char}($path_digits + 36) = 57
  set {char}($path_digits + 37) = 50
  set {char}($path_digits + 38) = 54
  set {char}($path_digits + 39) = 0
  set $raw_digits = ((unsigned int (*)(char *, char *, int, int, int, int, int, int, int))0x1001de90)($text_digits, $path_digits, 1, -1, -1, -1, -1, -1, -1)
  printf "MAKEINFO_TEXT case=digits raw_eax=%#x text=%s path=%s\n", $raw_digits, $text_digits, $path_digits
  set $text_alphabet = ((char *(*)(unsigned int))0x1001d9c0)(128)
  set {char}($text_alphabet + 0) = 65
  set {char}($text_alphabet + 1) = 46
  set {char}($text_alphabet + 2) = 66
  set {char}($text_alphabet + 3) = 46
  set {char}($text_alphabet + 4) = 67
  set {char}($text_alphabet + 5) = 46
  set {char}($text_alphabet + 6) = 0
  set $path_alphabet = ((char *(*)(unsigned int))0x1001d9c0)(128)
  set {char}($path_alphabet + 0) = 90
  set {char}($path_alphabet + 1) = 58
  set {char}($path_alphabet + 2) = 47
  set {char}($path_alphabet + 3) = 119
  set {char}($path_alphabet + 4) = 111
  set {char}($path_alphabet + 5) = 114
  set {char}($path_alphabet + 6) = 107
  set {char}($path_alphabet + 7) = 47
  set {char}($path_alphabet + 8) = 115
  set {char}($path_alphabet + 9) = 116
  set {char}($path_alphabet + 10) = 97
  set {char}($path_alphabet + 11) = 103
  set {char}($path_alphabet + 12) = 101
  set {char}($path_alphabet + 13) = 49
  set {char}($path_alphabet + 14) = 54
  set {char}($path_alphabet + 15) = 47
  set {char}($path_alphabet + 16) = 109
  set {char}($path_alphabet + 17) = 105
  set {char}($path_alphabet + 18) = 45
  set {char}($path_alphabet + 19) = 116
  set {char}($path_alphabet + 20) = 101
  set {char}($path_alphabet + 21) = 120
  set {char}($path_alphabet + 22) = 116
  set {char}($path_alphabet + 23) = 45
  set {char}($path_alphabet + 24) = 97
  set {char}($path_alphabet + 25) = 108
  set {char}($path_alphabet + 26) = 112
  set {char}($path_alphabet + 27) = 104
  set {char}($path_alphabet + 28) = 97
  set {char}($path_alphabet + 29) = 98
  set {char}($path_alphabet + 30) = 101
  set {char}($path_alphabet + 31) = 116
  set {char}($path_alphabet + 32) = 45
  set {char}($path_alphabet + 33) = 50
  set {char}($path_alphabet + 34) = 48
  set {char}($path_alphabet + 35) = 50
  set {char}($path_alphabet + 36) = 54
  set {char}($path_alphabet + 37) = 48
  set {char}($path_alphabet + 38) = 57
  set {char}($path_alphabet + 39) = 50
  set {char}($path_alphabet + 40) = 54
  set {char}($path_alphabet + 41) = 0
  set $raw_alphabet = ((unsigned int (*)(char *, char *, int, int, int, int, int, int, int))0x1001de90)($text_alphabet, $path_alphabet, 1, -1, -1, -1, -1, -1, -1)
  printf "MAKEINFO_TEXT case=alphabet raw_eax=%#x text=%s path=%s\n", $raw_alphabet, $text_alphabet, $path_alphabet
  set $text_mixed = ((char *(*)(unsigned int))0x1001d9c0)(128)
  set {char}($text_mixed + 0) = 72
  set {char}($text_mixed + 1) = 101
  set {char}($text_mixed + 2) = 108
  set {char}($text_mixed + 3) = 108
  set {char}($text_mixed + 4) = 111
  set {char}($text_mixed + 5) = 32
  set {char}($text_mixed + 6) = 52
  set {char}($text_mixed + 7) = 50
  set {char}($text_mixed + 8) = 46
  set {char}($text_mixed + 9) = 0
  set $path_mixed = ((char *(*)(unsigned int))0x1001d9c0)(128)
  set {char}($path_mixed + 0) = 90
  set {char}($path_mixed + 1) = 58
  set {char}($path_mixed + 2) = 47
  set {char}($path_mixed + 3) = 119
  set {char}($path_mixed + 4) = 111
  set {char}($path_mixed + 5) = 114
  set {char}($path_mixed + 6) = 107
  set {char}($path_mixed + 7) = 47
  set {char}($path_mixed + 8) = 115
  set {char}($path_mixed + 9) = 116
  set {char}($path_mixed + 10) = 97
  set {char}($path_mixed + 11) = 103
  set {char}($path_mixed + 12) = 101
  set {char}($path_mixed + 13) = 49
  set {char}($path_mixed + 14) = 54
  set {char}($path_mixed + 15) = 47
  set {char}($path_mixed + 16) = 109
  set {char}($path_mixed + 17) = 105
  set {char}($path_mixed + 18) = 45
  set {char}($path_mixed + 19) = 116
  set {char}($path_mixed + 20) = 101
  set {char}($path_mixed + 21) = 120
  set {char}($path_mixed + 22) = 116
  set {char}($path_mixed + 23) = 45
  set {char}($path_mixed + 24) = 109
  set {char}($path_mixed + 25) = 105
  set {char}($path_mixed + 26) = 120
  set {char}($path_mixed + 27) = 101
  set {char}($path_mixed + 28) = 100
  set {char}($path_mixed + 29) = 45
  set {char}($path_mixed + 30) = 50
  set {char}($path_mixed + 31) = 48
  set {char}($path_mixed + 32) = 50
  set {char}($path_mixed + 33) = 54
  set {char}($path_mixed + 34) = 48
  set {char}($path_mixed + 35) = 57
  set {char}($path_mixed + 36) = 50
  set {char}($path_mixed + 37) = 54
  set {char}($path_mixed + 38) = 0
  set $raw_mixed = ((unsigned int (*)(char *, char *, int, int, int, int, int, int, int))0x1001de90)($text_mixed, $path_mixed, 1, -1, -1, -1, -1, -1, -1)
  printf "MAKEINFO_TEXT case=mixed raw_eax=%#x text=%s path=%s\n", $raw_mixed, $text_mixed, $path_mixed
  continue
end

continue
