module github.com/aymanbagabas/go-udiff/cmd

go 1.23.0

replace github.com/aymanbagabas/go-udiff => ./..

require (
	github.com/aymanbagabas/go-udiff v0.3.1
	github.com/mattn/go-isatty v0.0.20
	github.com/spf13/pflag v1.0.5
)

require golang.org/x/sys v0.6.0 // indirect
