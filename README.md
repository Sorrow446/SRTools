# SRTools
Toolkit for modding Saints Row 2022 written in Go.
![](https://i.imgur.com/ib6Akqt.png)
[Windows, Linux, macOS](https://github.com/Sorrow446/SRTools/releases)

# Setup
[lz4 binary](https://github.com/lz4/lz4/releases/latest) is needed if using Windows.

# Usage
Only extracting is currently supported. vpp_pc and str2_pc pakfiles are supported.

Extract "dlc_01.vpp_pc" to "SRTools_extracted":   
`srtools_x64.exe unpack -i dlc_01.vpp_pc`

Extract "dlc_01.vpp_pc" and "dlc_preorder.vpp_pc" to "G:\sr" with 25 threads:   
`srtools_x64.exe unpack -i dlc_01.vpp_pc dlc_preorder.vpp_pc -o G:\sr -t 25`
