ENTRY(_start)

SECTIONS
{
    . = 0xa0000000;

    .text : {
        *(.text.init)
        *(.text)
    }

    .rodata : { *(.rodata) }
    .data   : { *(.data) }
    .bss    : { *(.bss)  *(COMMON) }
}
