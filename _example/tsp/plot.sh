#!/bin/bash
set -eo pipefail

PNGNAME=$(date -r stats.csv +stats.%Y%m%d.%H%M%S.png)

gnuplot > "$PNGNAME" <<EOF
set datafile separator ","
set key autotitle columnhead
set terminal png
set xlabel 'time(ms)'
plot 'stats.csv' using 5:2 with lines, '' using 5:3 with lines, '' using 5:4 with lines
EOF

echo written "$PNGNAME"
display "$PNGNAME"
