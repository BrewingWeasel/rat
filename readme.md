Rat is a cat clone that supports smart concatenation and basic syntax highlighting. It uses toml files to configure rules for merging files together.

Example use case:
# file1
```
#!/bin/bash
echo "The first cool file"
```
# file2

```
#!/bin/bash
echo "The second cool file"
```
# Running rat file1 file2
```
#!/bin/bash
echo "The first cool file"
echo "The second cool file"
```
Rat merged the two shebangs into just one file.

If you just want the syntax highlighting or want other features like git differences, use [bat](https://github.com/sharkdp/bat), the much more complete and featured inspiration for this project. At it's current state, I recommend bat over rat in pretty much all use cases.