.section .text
.set noreorder
.set nomacro
.set noat

.global _start

_start:
    nop
    nop
    j      FUN_a0000390
    nop

.incbin "data_chunk.bin"
.org 0x390

FUN_a0000390:
    lui     $t1,0xa861
    ori     $t1,$t1,0xa00
    ori     $s1,$zero,0x4
    sw      $s1,0x30($t1)
    ori     $s1,$zero,0x4
    sw      $s1,0x50($t1)
    ori     $s1,$zero,0x4
    sw      $s1,0x70($t1)
    ori     $s1,$zero,0x0
    sw      $s1,0x0($t1)
    ori     $s1,$zero,0x0
    sw      $s1,0xc($t1)
    ori     $s1,$zero,0x0
    sw      $s1,0x20($t1)
    ori     $s1,$zero,0x927f
    sw      $s1,0x30($t1)
    lui     $s1,0x1
    sw      $s1,0x40($t1)
    ori     $s1,$zero,0xc27f
    sw      $s1,0x50($t1)
    ori     $s1,$zero,0x0
    sw      $s1,0x60($t1)
    lui     $t1,0xa861
    ori     $t1,$t1,0x1a00
    lw      $s1,0x0($t1)
    lui     $t1,0x200
    and     $t2,$s1,$t1
    bgtz    $t2,LAB_a0000414
    lui     $t1,0xa861
    ori     $t1,$t1,0xa00
    ori     $s1,$zero,0x2197
    sw      $s1,0x70($t1)
    j       LAB_a0000480

LAB_a0000414:
    lui        $t1,0xa861
    ori        $t1,$t1,0x1a00
    lw         $s1,0x0($t1)
    lui        $t1,0x3
    and        $t2,$s1,$t1
    srl        $t1,$t2,0x10
    beq        $t1,$zero,LAB_a0000448
    ori        $t2,$zero,0x1
    sub        $t1,$t1,$t2
    beq        $t1,$zero,LAB_a000045c
    ori        $t2,$zero,0x1
    sub        $t1,$t1,$t2
    beq        $t1,$zero,LAB_a0000470

LAB_a0000448:
    lui        $t1,0xa861
    ori        $t1,$t1,0xa00
    ori        $s1,$zero,0xb197
    sw         $s1,0x70($t1)
    j          LAB_a0000480

LAB_a000045c:
    lui       $t1,0xa861
    ori        $t1,$t1,0xa00
    ori        $s1,$zero,0x7197
    sw         $s1,0x70($t1)
    j          LAB_a0000480

LAB_a0000470:
    lui        $t1,0xa861
    ori        $t1,$t1,0xa00
    ori        $s1,$zero,0xb197
    sw         $s1,0x70($t1)

LAB_a0000480:
    ori     $s1,$zero,0x1f
    sw      $s1,0x80($t1)
    ori     $s1,$zero,0x8000
    sw      $s1,0x90($t1)
    mtc0    $zero,$13
    nop
    lui     $t0,0x40
    mtc0    $t0,$12
    nop
    lui     $gp,0x9406
    ori     $gp,$gp,0x5b8
    lui     $sp,0x9480
    ori     $sp,$sp,0x0
    ori     $t5,$zero,0x800

LAB_a00004b8:
    addi    $t5,$t5,-0x1
    bne     $t5,$zero,LAB_a00004b8
    nop
    jal     FUN_a00032c8
    nop
    jal     maybe_set_kseg0_uncached
    nop
    jal     init_dcache
    nop
    jal     init_icache
    nop
    lui     $v0,0x9406
    ori     $v0,$v0,0x5b8
    lui     $v1,0x9406
    ori     $v1,$v1,0x960

LAB_a00004f4:
    sw      $zero,0x0($v0)
    sltu    $at,$v0,$v1
    bne     $at,$zero,LAB_a00004f4
    addiu   $v0,$v0,0x4
    lui     $t0,0xb000
    ori     $t0,$t0,0x34f0
    lui     $t1,0x9406
    ori     $t1,$t1,0x100
    lui     $t2,0x9406
    ori     $t2,$t2,0x5b8
    lw      $t3,0x0($t0)

LAB_a0000520:
    addiu      $t0,$t0,0x4
    sw         $t3,0x0($t1)
    addiu      $t1,$t1,0x4
    sltu       $at,$t1,$t2
    bne        $at,$zero,LAB_a0000520
    lw         $t3,0x0($t0)
    jal        FUN_a0003290
    nop
    jal        main
    nop
    j          LAB_a0000550
    nop

LAB_a0000550:
    nop
    b          LAB_a0000550
    nop
    nop
