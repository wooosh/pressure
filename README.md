# Pressure
Pressure is a utility aimed at simplifying compression and decompression.
## Syntax
Pressure is designed to have a consistent simple syntax.
```
# Basic Syntax
pressure <src> <dst>

# Decompression Examples
pressure file.tar outputdir
pressure code.zip code

# Compression Examples
pressure code code.zip
pressure src assets tests test.zip # Not currently implemented
```
## Supported Formats
Currently only tar is supported. More formats such as lzma tar, zip and rar are planned to be implemented.

## Planned Feaures
- [ ] tar.gz compression
- [ ] RAR support
- [ ] Compressing multiple files into one archive
- [ ] Zip support
- [ ] Better error handling
- [ ] Tests
- [ ] List files option
