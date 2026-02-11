# mail-sdk-go

Go SDK for [Nulz Mail](https://nulz.lol). Stdlib only, minimal.

## Install

```bash
go get github.com/nulz-rip/mail-sdk-go@latest
```

## Set API key

**Windows CMD**
```cmd
set NULZ_API_KEY=your-key
```

**PowerShell**
```powershell
$env:NULZ_API_KEY="your-key"
```

**Linux / macOS**
```bash
export NULZ_API_KEY=your-key
```

## Quick use

```go
package main

import (
    "context"
    "github.com/nulz-rip/mail-sdk-go/nulzmail"
)

func main() {
    client := nulzmail.New() // uses NULZ_API_KEY
    ctx := context.Background()

    inbox, _ := client.CreateInbox(ctx)
    code, msg, _ := client.WaitForCode(ctx, inbox.ID, nulzmail.WaitOpts{})
    // use code, msg
}
```

## Base URL override

```go
client := nulzmail.New("your-key")
client.SetBaseURL("https://custom.example.com")
```

## Dev

```bash
go test ./...
```

## Release

Tag a version:

```bash
git tag v0.1.0
git push --tags
```

Users can pin: `go get github.com/nulz-rip/mail-sdk-go@v0.1.0`

## Contributing

Fork, create a branch, open a PR.

## License

MIT License

Copyright (c) 2025 Nulz

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
