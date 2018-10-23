#!/bin/bash

echo "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA"
echo "BBBBBBBBBBBBBBBBBBBBBBBBBBBBBBB"
echo "CCCCCCCCCCCCCCCCCCCCCCCCCCCCCCC"
echo "DDDDDDDDDDDDDDDDDDDDDDDDDDDDDDD"
echo "EEEEEEEEEEEEEEEEEEEEEEEEEEEEEEE"

echo -ne "\x1b[A" # up
echo -ne "\x1b[A" # up
echo -ne "\x1b[C" # right
echo -ne "\x1b[C" # rightt
sleep 2
echo -ne "\x1b[2L" # insert 2 lines
sleep 2
echo -ne "\x1b[B" # down
echo -ne "\x1b[B" # down
sleep 2