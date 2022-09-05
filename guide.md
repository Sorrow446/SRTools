## Extract
Extract a vpp_pc or str2_pc packfile with folder structure.

`unpack -i dlc_01.vpp_pc -o G:\sr`    
The -i arg supports multiple input paths (duplicates will be filtered).

## Pack
**Experimental. May cause the game to black screen on some boots.**    
Pack files into a vpp_pc or str2_pc packfile.
  
`pack -i SRTools_extracted -o packed.vpp_pc`    
Input folder must have the same structure created by the unpacker.

## Convert

### Scribe
Edit scribe localisation file.

1. Export scribe to JSON.    
`convert -i activity.en.scribe_pad -o activity.en.json`
2. Open activity.en.json in a text editor to change the strings (see the text key and type string in each entry).    
3. Convert JSON to scribe.    
`convert -i activity.en.json-o activity.en.scribe_pad`
