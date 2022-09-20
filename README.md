# wabin: WebAssembly Binary Format in Go

wabin includes WebAssembly an WebAssembly data model and binary encoder. Most
won't use this library. It mainly supports advanced manipulation of WebAssembly
binaries prior to instantiation with wazero. Notably, this has no dependencies,
so is cleaner to use in Go projects.
