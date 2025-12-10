#!/usr/bin/env bash
go test -v -run Test_CPU 2>&1 \
  | tee /dev/stderr \
  | sed -E 's/ \([0-9.]+s\)//g; /^ok/d' \
  | sed 's/[^[:print:]\t\n]//g' \
  | sed -E 's/^[[:space:]]+$//' \
  > /tmp/cpu_test_output.txt

awk '
/<!-- TEST_OUTPUT_START -->/ { print; print "```txt"; system("cat /tmp/cpu_test_output.txt"); print "```"; skip=1; next }
/<!-- TEST_OUTPUT_END -->/ { skip=0 }
!skip { print }
' README.md > README.tmp && mv README.tmp README.md

rm /tmp/cpu_test_output.txt
