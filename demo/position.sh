#!/bin/bash

echo "ABCDEFGHIJKLM"
echo "NOPQRSTUVWXYZ"
echo "ABCDEFGHIJKLM"
echo "NOPQRSTUVWXYZ"
echo "ABCDEFGHIJKLM"
echo "NOPQRSTUVWXYZ"
echo "ABCDEFGHIJKLM"
echo "NOPQRSTUVWXYZ"
sleep 1
echo -ne "\x1b[A" # up
sleep 1
echo -ne "\x1b[D" # left
sleep 1
echo -ne "\x1b[C" # right
sleep 1
echo -ne "\x1b[C" # right
sleep 1
echo -ne "\x1b[C" # right
sleep 1
echo -ne "\x1b[A" # up
sleep 1
echo -ne "\x1b[A" # up
sleep 1
echo -ne "\x1b[B" # down
sleep 1
echo -n "123"
sleep 1
echo -ne "\x1b[E" # line down (col 0)
sleep 1
echo -n "ZZZ"
sleep 1
echo -ne "\x1b[F" # line up (col 0)
sleep 1
echo -n "XXX"
sleep 1
echo -ne "\x1b[E" # line down (col 0)
sleep 1
echo