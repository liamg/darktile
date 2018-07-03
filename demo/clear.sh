#!/bin/bash

echo "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA"
echo "BBBBBBBBBBBBBBBBBBBBBBBBBBBBBBB"
echo "BBBBBBBBBBBBBBBBBBBBBBBBBBBBBBB"
echo "BBBBBBBBBBBBBBBBBBBBBBBBBBBBBBB"
echo -n "BBBBBBBBBBBBBBBBBBBBBBBBBBBBBBB"

echo -ne "\x1b[A" # up
echo -ne "\x1b[D" # left
echo -ne "\x1b[A" # up
echo -ne "\x1b[D" # left
echo -ne "\x1b[A" # up
echo -ne "\x1b[D" # left
sleep 2
echo -ne "\x1b[J0"

