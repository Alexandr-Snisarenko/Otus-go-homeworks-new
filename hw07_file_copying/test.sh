#!/usr/bin/env bash
set -xeuo pipefail

go build -o go-cp

./go-cp -from testdata/input.txt -to out.txt -rewrite
cmp out.txt testdata/out_offset0_limit0.txt

./go-cp -from testdata/input.txt -to out.txt -limit 10 -rewrite
cmp out.txt testdata/out_offset0_limit10.txt

./go-cp -from testdata/input.txt -to out.txt -limit 1000 -rewrite
cmp out.txt testdata/out_offset0_limit1000.txt

./go-cp -from testdata/input.txt -to out.txt -limit 10000 -rewrite
cmp out.txt testdata/out_offset0_limit10000.txt

./go-cp -from testdata/input.txt -to out.txt -offset 100 -limit 1000 -rewrite
cmp out.txt testdata/out_offset100_limit1000.txt

./go-cp -from testdata/input.txt -to out.txt -offset 6000 -limit 1000 -rewrite
cmp out.txt testdata/out_offset6000_limit1000.txt

rm -f go-cp out.txt
echo "PASS"
