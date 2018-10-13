# pefile

Small library and tool to list and extract resources from Portable Executable (PE) files.

## Installing the tool

```shell
go get -u github.com/folbricht/pefile/cmd/pe
```

## Running the tool

List available resources in the file:

```shell
# pe list-resources /path/to/file.exe
3/1/1033
3/2/1033
5/105/1033
5/106/1033
10/SOMERESOURCE/1033
14/103/1033
```

Extract a single resource to STDOUT:

```shell
# pe extract-resource /path/to/file.exe 10/SOMERESOURCE/1033
```

Extract a single resource to a file:

```shell
# pe extract-resource /path/to/file.exe 10/SOMERESOURCE/1033 /tmp/someresource.bin
```

Extract all resources to a directory, preserving the directory layout from the PE file.

```shell
# pe extract-resources /path/to/file.exe /tmp/resources
```

## Using the library

```go
f, err := pefile.Open("/path/to/file.exe")
if err != nil {
  return err
}
defer f.Close()
resources, err := f.GetResources()
if err != nil {
  return err
}
for _, r := range resources {
  fmt.Println("Name:", r.Name, ", Size:", len(r.Data))
}
```

## References

- Godoc - [https://godoc.org/github.com/folbricht/pefile](https://godoc.org/github.com/folbricht/pefile)
- Peering Inside the PE: A Tour of the Win32 Portable Executable File Format - [https://msdn.microsoft.com/en-us/library/ms809762.aspx](https://msdn.microsoft.com/en-us/library/ms809762.aspx)
- PE Format - [https://docs.microsoft.com/en-us/windows/desktop/debug/pe-format](https://docs.microsoft.com/en-us/windows/desktop/debug/pe-format)
- Go PE library - [https://golang.org/pkg/debug/pe/](https://golang.org/pkg/debug/pe/)
