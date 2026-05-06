#!/usr/bin/env bash
dd if=boot.bin of=boot-chunk.bin bs=1 count=$(wc -c bootloader.bin | cut -d' ' -f1)

hexdump -C bootloader.bin > new.txt
hexdump -C original.bin > old.txt

diff new.txt old.txt
