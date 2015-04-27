[![GoDoc](https://godoc.org/github.com/brentp/vcfgo?status.svg)](https://godoc.org/github.com/brentp/vcfgo)
[![Build Status](https://travis-ci.org/brentp/vcfgo.svg)](https://travis-ci.org/brentp/vcfgo)

vcfgo is a golang library to read, write and manipulate files in the variant call format.

# vcfgo
--
    import "github.com/brentp/vcfgo"

Package vcfgo implements a Reader and Writer for variant call format. It eases
reading, filtering modifying VCF's even if they are not to spec. Example:

## Usage

```go
f, _ := os.Open("examples/test.auto_dom.no_parents.vcf")
rdr, err := vcfgo.NewReader(f, false)
if err != nil {
	panic(err)
}
for {
	variant := rdr.Read()
	if variant == nil {
		break
	}
	fmt.Printf("%s\t%d\t%s\t%s\n", variant.Chromosome, variant.Pos, variant.Ref, variant.Alt)
	fmt.Printf("%s", variant.Info["DP"].(int) > 10)
	sample := variant.Samples[0]
	// we can get the PL field as a list (-1 is default in case of missing value)
	fmt.Println("%s", variant.GetGenotypeField(sample, "PL", -1))
	_ = sample.DP
}
fmt.Fprintln(os.Stderr, rdr.Error())
```

## Status

`vcfgo` is well-tested, but still in development. It tries to tolerate, but report
errors; after every `rdr.Read()` call, the caller can check `rdr.Error()`
and get feedback on the errors without stopping execution unless it is explicitly
requested to do so.

Info and sample fields are pre-parsed and stored as `map[string]interface{}` so
callers will have to cast to the appropriate type upon retrieval.

#### type Header

```go
type Header struct {
	SampleNames   []string
	Infos         map[string]*Info
	SampleFormats map[string]*SampleFormat
	Filters       map[string]string
	Extras        map[string]string
	FileFormat    string
	// contid id maps to a map of length, URL, etc.
	Contigs map[string]map[string]string
}
```

Header holds all the type and format information for the variants.

#### func  NewHeader

```go
func NewHeader() *Header
```
NewHeader returns a Header with the requisite allocations.

#### type Info

```go
type Info struct {
	Id          string
	Description string
	Number      string // A G R . ''
	Type        string // STRING INTEGER FLOAT FLAG CHARACTER UNKONWN
}
```

Info holds the Info and Format fields

#### func (*Info) String

```go
func (i *Info) String() string
```
String returns a string representation.

#### type InfoMap

```go
type InfoMap map[string]interface{}
```

InfoMap holds the parsed Info field which can contain floats, ints and lists
thereof.

#### func (InfoMap) String

```go
func (m InfoMap) String() string
```
String returns a string that matches the original info field.

#### type Reader

```go
type Reader struct {
	Header *Header

	LineNumber int64
}
```

Reader holds information about the current line number (for errors) and The VCF
header that indicates the structure of records.

#### func  NewReader

```go
func NewReader(r io.Reader, lazySamples bool) (*Reader, error)
```
NewReader returns a Reader.

#### func (*Reader) Clear

```go
func (vr *Reader) Clear()
```
Clear empties the cache of errors.

#### func (*Reader) Error

```go
func (vr *Reader) Error() error
```
Error() aggregates the multiple errors that can occur into a single object.

#### func (*Reader) Read

```go
func (vr *Reader) Read() *Variant
```
Read returns a pointer to a Variant. Upon reading the caller is assumed to check
Reader.Err()

#### type SampleFormat

```go
type SampleFormat Info
```

SampleFormat holds the type info for Format fields.

#### func (*SampleFormat) String

```go
func (i *SampleFormat) String() string
```
String returns a string representation.

#### type SampleGenotype

```go
type SampleGenotype struct {
	Phased bool
	GT     []int
	DP     int
	GL     []float32
	GQ     int
	MQ     int
	// TODO: add methods for Ref, Alt depth.
	Fields map[string]string
}
```

SampleGenotype holds the information about a sample. Several fields are
pre-parsed, but all fields are kept in Fields as well.

#### func  NewSampleGenotype

```go
func NewSampleGenotype() *SampleGenotype
```
NewSampleGenotype allocates the internals and returns a SampleGenotype

#### func (*SampleGenotype) String

```go
func (sg *SampleGenotype) String(fields []string) string
```
String returns the string representation of the sample field.

#### type VCFError

```go
type VCFError struct {
	Msgs  []string
	Lines []int64
}
```

VCFError satisfies the error interface and allows multiple errors. This is
useful because, for example, on a single line, every sample may have a field
that doesn't match the description in the header. We want to keep parsing but
also let the caller know about the error.

#### func  NewVCFError

```go
func NewVCFError() *VCFError
```
NewVCFError allocates the needed ingredients.

#### func (*VCFError) Add

```go
func (e *VCFError) Add(err error, line int64)
```
Add adds an error and the line number within the vcf where the error took place.

#### func (*VCFError) Clear

```go
func (e *VCFError) Clear()
```
Clear empties the Messages

#### func (*VCFError) Error

```go
func (e *VCFError) Error() string
```
Error returns a string with all errors delimited by newlines.

#### func (*VCFError) IsEmpty

```go
func (e *VCFError) IsEmpty() bool
```
IsEmpty returns true if there no errors stored.

#### type Variant

```go
type Variant struct {
	Chromosome      string
	Pos        		uint64
	Id         		string
	Ref        		string
	Alt        		[]string
	Quality    		float32
	Filter     		string
	Info       		InfoMap
	Format     		[]string
	Samples    		[]*SampleGenotype
	Header     		*Header
	LineNumber 		int64
}
```

Variant holds the information about a single site. It is analagous to a row in a
VCF file.

#### func (*Variant) GetGenotypeField

```go
func (v *Variant) GetGenotypeField(g *SampleGenotype, field string, missing interface{}) (interface{}, error)
```
GetGenotypeField uses the information from the header to parse the correct time
from a genotype field. It returns an interface that can be asserted to the
expected type.

#### func (*Variant) String

```go
func (v *Variant) String() string
```
String gives a string representation of a variant

#### type Writer

```go
type Writer struct {
	io.Writer
	Header *Header
}
```

Writer allows writing VCF files.

#### func  NewWriter

```go
func NewWriter(w io.Writer, h *Header) (*Writer, error)
```
NewWriter returns a writer after writing the header.

#### func (*Writer) WriteVariant

```go
func (w *Writer) WriteVariant(v *Variant)
```
WriteVariant writes a single variant
