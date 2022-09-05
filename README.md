# SRTools
Toolkit for modding Saints Row 2022 written in Go.
![](https://i.imgur.com/ib6Akqt.png)
![](https://i.imgur.com/Ui2PV0h.png)
[Windows, Linux, and macOS binaries](https://github.com/Sorrow446/SRTools/releases)

# Setup
[lz4 binary](https://github.com/lz4/lz4/releases/latest) is needed if using Windows.

# Usage
[Click here for guide.](https://github.com/Sorrow446/SRTools/blob/main/guide.md)

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

# Supported Formats
|File Type|Read|Write|
| --- | --- | --- |
|vpp_pc/str2_pc|y|repack only, unstable|
|Scribe (.scribe_pad)|y|y|
