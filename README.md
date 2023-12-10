# Go RESP

Implements the RESP protocol for version 2 &amp; 3 in order to build Redis clients.

## Roadmap

This is currently in-progress. Nothing has been tested yet, but this is what is
implemented so far.

### RESP 2

- [X] Simple String
- [X] Simple Error
- [X] Integers
- [X] Bulk String
- [X] Array

### RESP 3

- [X] Nulls
- [X] Booleans
- [ ] Doubles
- [ ] Big Number
- [ ] Bulk Error
- [ ] Verbatim String
- [ ] Map
- [ ] Set
- [ ] Push

### General Usage

- [ ] Unmarshal
- [ ] Marshal
- [ ] io.Reader
- [ ] io.Writer
 
## Installation

```
go get -u github.com/maple-tech/go-resp
```

## Usage

Not sure yet

## License

MIT License

Copyright (c) 2023 Maple Technologies

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.

## Authors

- Chris Pikul
