package pefile

import (
	"debug/pe"
	"encoding/binary"
	"io"
	"strconv"
	"unicode/utf16"
)

// File represents a PE file. It wraps a pe.File to provide access to more
// headers and elements.
type File struct {
	*pe.File
}

// Open opens the named PE file
func Open(name string) (*File, error) {
	p, err := pe.Open(name)
	return &File{p}, err
}

// New initializes a File from a ReaderAt
func New(r io.ReaderAt) (*File, error) {
	p, err := pe.NewFile(r)
	return &File{p}, err
}

// Resource holds the full name and data of a data entry in a resource directory structure.
// The name represents all 3 parts of the tree, separated by /, <type>/<name>/<language> with
// For example: "3/1/1033" for a resources with ID names, or "10/SOMERES/1033" for a named
// resource in language 1033.
type Resource struct {
	Name string
	Data []byte
}

// GetResources returns a list of resources
func (f *File) GetResources() ([]Resource, error) {
	rsrc := f.Section(".rsrc")

	b, err := rsrc.Data()
	if err != nil {
		return nil, err
	}

	// Recursively parse the directory tree starting with the root directory (p=0).
	r := parseDir(b, 0, "", rsrc.SectionHeader.VirtualAddress)
	return r, nil
}

// Recursively parses a IMAGE_RESOURCE_DIRECTORY in slice b starting at position p
// building on path prefix. virtual is needed to calculate the position of the data
// in the resource
func parseDir(b []byte, p int, prefix string, virtual uint32) []Resource {
	var resources []Resource

	// Skip Characteristics, Timestamp, Major, Minor in the directory

	numberOfNamedEntries := int(binary.LittleEndian.Uint16(b[p+12 : p+14]))
	numberOfIdEntries := int(binary.LittleEndian.Uint16(b[p+14 : p+16]))
	n := numberOfNamedEntries + numberOfIdEntries

	// Iterate over all entries in the current directory record
	for i := 0; i < n; i++ {
		o := 8*i + p + 16
		name := int(binary.LittleEndian.Uint32(b[o : o+4]))
		offsetToData := int(binary.LittleEndian.Uint32(b[o+4 : o+8]))
		path := prefix
		if name&0x80000000 > 0 { // Named entry if the high bit is set in the name
			dirString := name & (0x80000000 - 1)
			length := int(binary.LittleEndian.Uint16(b[dirString : dirString+2]))
			s := UTF16ToString(b[dirString+2 : dirString+2+length*2])
			path += s
		} else { // ID entry
			path += strconv.Itoa(name)
		}

		if offsetToData&0x80000000 > 0 { // Ptr to other directory if high bit is set
			subdir := offsetToData & (0x80000000 - 1)

			// Recursively get the resources from the sub dirs
			l := parseDir(b, subdir, path+"/", virtual)
			resources = append(resources, l...)
			continue
		}

		// Leaf, ptr to the data entry. Read IMAGE_RESOURCE_DATA_ENTRY
		offset := int(binary.LittleEndian.Uint32(b[offsetToData : offsetToData+4]))
		length := int(binary.LittleEndian.Uint32(b[offsetToData+4 : offsetToData+8]))

		// The offset in IMAGE_RESOURCE_DATA_ENTRY is relative to the virual address.
		// Calculate the address in the file
		offset -= int(virtual)
		data := b[offset : offset+length]

		// Add Resource to the list
		resources = append(resources, Resource{Name: path, Data: data})
	}
	return resources
}

// UTF16ToString converts a UTF16-LE byte slice into a string
func UTF16ToString(b []byte) string {
	var r []uint16
	for {
		if len(b) < 2 {
			break
		}
		v := binary.LittleEndian.Uint16(b[0:2])
		r = append(r, v)
		b = b[2:]
	}
	return string(utf16.Decode(r))
}
