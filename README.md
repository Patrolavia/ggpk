# GGPK Helpers for Path of Exile

## Synopsis

```sh
# list files in Content.ggpk
list Content.ggpk

# extract all files from Content.ggpk to folder destination
extract -d destination -r Content.ggpk /

# defrag Content.ggpk, this will create a new ggpk file nam
defrag Content.ggpk

# Verify checksum of all files in Content.ggpk
check Content.ggpk
```

## Defragment

Defrag tool does not do it's work on the position. It creates another file named `result.ggpk`.

It also puts all directory record together, so we have bigger chance to read a child node without doing additional hardware I/O. Also, if GGG caches records in memory, this can benefits program initial speed a little.

## License

Any version of MIT, GPL or LGPL.
