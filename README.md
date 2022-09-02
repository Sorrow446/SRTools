# SRTools
Toolkit for modding Saints Row 2022 written in Go.
![](https://i.imgur.com/ib6Akqt.png)
[Windows, Linux, and macOS binaries](https://github.com/Sorrow446/SRTools/releases)

# Setup
[lz4 binary](https://github.com/lz4/lz4/releases/latest) is needed if using Windows.

# Usage
vpp_pc and str2_pc packfiles are supported.


## Extract
Extract "dlc_01.vpp_pc" to "SRTools_extracted":   
`srtools_x64.exe unpack -i dlc_01.vpp_pc`

Extract "dlc_01.vpp_pc" and "dlc_preorder.vpp_pc" to "G:\sr" with 20 threads:   
`srtools_x64.exe unpack -i dlc_01.vpp_pc dlc_preorder.vpp_pc -o G:\sr -t 20`

## Pack
**Experimental. May cause the game to black screen on some boots.**    
Src folder must have the same structure created by the unpacker.

Pack "SRTools_extracted" to "SRTools_packed.vpp_pc":   
`srtools_x64.exe unpack -i dlc_01.vpp_pc`


```
Usage: sr_tools_x64.exe --inpaths INPATHS [--outpath OUTPATH] [--threads THREADS] [--nocompression] COMMAND

Positional arguments:
  COMMAND

Options:
  --inpaths INPATHS, -i INPATHS
                         Input path(s).
  --outpath OUTPATH, -o OUTPATH
                         Output path. Path will be made if it doesn't already exist.
  --threads THREADS, -t THREADS
                         Max threads (1-50). Be careful; memory intensive. [default: 10]
  --nocompression, -n    Don't compress any files when packing. Might be a bit more stable.
  --help, -h             display this help and exit
```
