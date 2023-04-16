# µDiff

<p>
<a href="https://github.com/aymanbagabas/go-udiff/releases"><img src="https://img.shields.io/github/release/aymanbagabas/go-udiff.svg" alt="Latest Release"></a>
<a href="https://pkg.go.dev/github.com/aymanbagabas/go-udiff?tab=doc"><img src="https://godoc.org/github.com/golang/gddo?status.svg" alt="Go Docs"></a>
<a href="https://github.com/aymanbagabas/go-udiff/actions"><img src="https://github.com/aymanbagabas/go-udiff/workflows/build/badge.svg" alt="Build Status"></a>
<a href="https://goreportcard.com/report/github.com/aymanbagabas/go-udiff"><img alt="Go Report Card" src="https://goreportcard.com/badge/github.com/aymanbagabas/go-udiff"></a>
</p>

Micro diff (µDiff) is a Go library that implements the
[Myers'](http://www.xmailserver.org/diff2.pdf) diffing algorithm. It aims to
provide a minimum implementation with no extra dependencies. The
library supports computing diffs, generating [unified format](https://www.gnu.org/software/diffutils/manual/html_node/Unified-Format.html)
diffs, and applying diffs to strings and bytes.

This is merely a copy of the [Golang tools internal diff package](https://github.com/golang/tools/tree/master/internal/diff).
All credit goes to the [Go authors](https://go.dev/AUTHORS).

## Usage

You can import the package using the following command:

```bash
go get github.com/aymanbagabas/go-udiff
```

## Example

Generate a unified diff for strings `a` and `b`.

```go
package main

import (
  "fmt"

  "github.com/aymanbagabas/go-udiff"
)

func main() {
  a := "Hello, world!"
  b := "Hello, Go!"

  edits := udiff.Strings(a, b)
  d, err := udiff.ToUnifiedDiff("a.txt", "b.txt", a, edits)
  if err != nil {
      panic(err)
  }

  fmt.Println(d.String())
}
```

## Alternatives

- [sergi/go-diff](https://github.com/sergi/go-diff) No longer reliable. See [#123](https://github.com/sergi/go-diff/issues/123) and [#141](https://github.com/sergi/go-diff/pull/141).
- [hexops/gotextdiff](https://github.com/hexops/gotextdiff) Takes the same approach but looks like the project is abandoned.
- [sourcegraph/go-diff](https://github.com/sourcegraph/go-diff) It doesn't compute diffs. Great package for parsing and printing unified diffs.

## Contributing

Please send any contributions [upstream](https://github.com/golang/tools). Pull
requests made against [the upsteam diff package](https://github.com/golang/tools/tree/master/internal/diff)
are welcome.

## License

[BSD 3-Clause](./LICENSE-BSD) and [MIT](./LICENSE-MIT).
