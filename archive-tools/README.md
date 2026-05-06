# Cisco CNU Archive Tools

This is a utility for working with Cisco CNU archive files, specifically the Version 3.0 format used the Java based Cisco IP phones (7941, 7961, 7970, 79x5, and others)

## Usage

```
Usage:
  ./archive-tools [command]

Available Commands:
  completion      Generate the autocompletion script for the specified shell
  help            Help about any command
  list            List contents of firmware file
  pack            Pack firmware file
  parse-signature Parse firmware digital signature
  unpack          Unpack firmware file
  unsign          Remove digital signature from firmware
```

### List 

Lists the contents of an archive file

```sh
./archive-tools list <archive filename>
```

### Unpack 

Unpack the contents of an archive into a target directory

```sh
./archive-tools unpack <archive filename> <output directory>
```

### Unpack 

Pack the contents of an target directory into an unsigned archive file

```sh
./archive-tools pack <archive filename> <input directory>
```

### Unsign

Strip the signature from a signed firmware image file

```
./archive-tools unsign <signed firmware filename> <output filename> 
```

### Parse Signature

Parse the signed firmware signature block and print the results to the console

```
./archive-tools parse-signature <signed firmware filename>
```

### Prior art

* https://git.deuxfleurs.fr/distorsion/telephones/src/branch/main/firmware_extraction
* https://github.com/kbdfck/cnu-fpu
