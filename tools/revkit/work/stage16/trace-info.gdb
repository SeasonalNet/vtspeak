set pagination off
set confirm off
set debuginfod enabled off
handle SIGSEGV nostop noprint pass

break *0x1001da50
commands
  silent
  disable 1
  set $info = (char *)malloc(256)
  set $req = 0
  while $req <= 26
    set $word = 0
    while $word < 64
      set *(unsigned int *)($info + $word * 4) = 0
      set $word = $word + 1
    end
    set $info_result = ((int (*)(int, char *, void *, int))0x1002a690)($req, (char *)0, (void *)$info, 256)
    if ($req == 0) || ($req == 3) || ($req == 23)
      printf "INFO request=%d result=%d value=%s\n", $req, $info_result, $info
    else
      printf "INFO request=%d result=%d value=%d\n", $req, $info_result, *(int *)$info
    end
    set $req = $req + 1
  end
  set $info_result = ((int (*)(int, char *, void *, int))0x1002a690)(27, (char *)0, (void *)$info, 256)
  printf "INFO invalid-request result=%d\n", $info_result
  set $info_result = ((int (*)(int, char *, void *, int))0x1002a690)(10, (char *)0, (void *)0, 4)
  printf "INFO null-value result=%d\n", $info_result
  set $info_result = ((int (*)(int, char *, void *, int))0x1002a690)(10, (char *)0, (void *)$info, 1)
  printf "INFO short-int-size result=%d\n", $info_result
  set $info_result = ((int (*)(int, char *, void *, int))0x1002a690)(0, (char *)0, (void *)$info, 1)
  printf "INFO short-string-size result=%d\n", $info_result
  set $info_result = ((int (*)(int, char *, void *, int))0x1002a690)(0, (char *)0, (void *)$info, 11)
  printf "INFO build-date-size-11 result=%d\n", $info_result
  set $info_result = ((int (*)(int, char *, void *, int))0x1002a690)(0, (char *)0, (void *)$info, 12)
  printf "INFO build-date-size-12 result=%d value=%s\n", $info_result, $info
  set $info_result = ((int (*)(int, char *, void *, int))0x1002a690)(3, (char *)0, (void *)$info, 3)
  printf "INFO db-directory-size-3 result=%d\n", $info_result
  set $info_result = ((int (*)(int, char *, void *, int))0x1002a690)(3, (char *)0, (void *)$info, 4)
  printf "INFO db-directory-size-4 result=%d value=%s\n", $info_result, $info
  set $info_result = ((int (*)(int, char *, void *, int))0x1002a690)(10, (char *)0, (void *)$info, 0)
  printf "INFO int-size-0 result=%d value=%d\n", $info_result, *(int *)$info
  continue
end

continue
