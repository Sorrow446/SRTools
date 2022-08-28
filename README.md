# SRTools
Toolkit for modding Saints Row 2022 written in Go.
![](https://i.imgur.com/ib6Akqt.png)
[Windows, Linux, and macOS binaries](https://github.com/Sorrow446/SRTools/releases)

# Setup
[lz4 binary](https://github.com/lz4/lz4/releases/latest) is needed if using Windows.

# Usage
Only extracting is currently supported. vpp_pc and str2_pc pakfiles are supported.

Extract "dlc_01.vpp_pc" to "SRTools_extracted":   
`srtools_x64.exe unpack -i dlc_01.vpp_pc`

Extract "dlc_01.vpp_pc" and "dlc_preorder.vpp_pc" to "G:\sr" with 25 threads:   
`srtools_x64.exe unpack -i dlc_01.vpp_pc dlc_preorder.vpp_pc -o G:\sr -t 25`

```
Usage: srtools_x64.exe --inpaths INPATHS [--outpath OUTPATH] [--threads THREADS] COMMAND

Positional arguments:
  COMMAND

Options:
  --inpaths INPATHS, -i INPATHS
                         Input paths of packfiles.
  --outpath OUTPATH, -o OUTPATH
                         Output path. Path will be made if it doesn't already exist.
  --threads THREADS, -t THREADS
                         Max threads (1-100). [default: 10]
  --help, -h             display this help and exit
```
