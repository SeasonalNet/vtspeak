set pagination off
set confirm off
set debuginfod enabled off
handle SIGSEGV nostop noprint pass
set logging file /work/corpus-parity/stage20/hello-hi-range-scope-baseline/live-range-scope.gdb.log
set logging overwrite on
set logging enabled on
source /work/stage20/trace-2013-hello-hi-prefix-range-scope-common.gdb
