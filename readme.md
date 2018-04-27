#### Brook-Go

Read and write octet- and bitstreams.

For .NET-platforms, use [Brook-Dotnet](https://github.com/Piot/brook-dotnet).

#### Usage

##### Octet streams

```go
octetStream := outstream.New()
octetStream.WriteUint16(51942)
```

##### Bit streams

```go
octetStream := outstream.New()
bitStream := outbitstream.New(octetStream)
writeError := bitStream.WriteBits(0xcafe9, 20)
```


