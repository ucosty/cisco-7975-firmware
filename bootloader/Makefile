# Toolchain definition
CROSS_COMPILE = mips-unknown-elf-
AS      = $(CROSS_COMPILE)as
LD      = $(CROSS_COMPILE)ld
OBJCOPY = $(CROSS_COMPILE)objcopy

# Files
INPUT_S  = bootloader.S
OUTPUT_ELF = bootloader.elf
OUTPUT_BIN = bootloader.bin

ASFLAGS = -march=mips3 -mabi=32 -EB
LDFLAGS =-T linker.ld

all: $(OUTPUT_BIN)

$(OUTPUT_BIN): $(INPUT_S)
	$(AS) $(ASFLAGS) $(INPUT_S) -o $(OUTPUT_ELF)
	$(OBJCOPY) -j .text -O binary $(OUTPUT_ELF) $(OUTPUT_BIN)

clean:
	rm -f $(OUTPUT_ELF) $(OUTPUT_BIN)

check:
	./check.sh

.PHONY: all clean
