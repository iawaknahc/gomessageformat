## What is this?

This is a Go implementation of [ICU MessageFormat](https://unicode-org.github.io/icu-docs/apidoc/released/icu4j/com/ibm/icu/text/MessageFormat.html).

## How to build

This library depends on icu4c 67.1

On macOS, the simplest way is to install it with brew

```sh
brew install icu4c
```

Since macOS comes with its own icu4c, in order for the Go toolchain to find our installation of icu4c,
we have to set the following environment variable.

```sh
export PKG_CONFIG_PATH="/usr/local/opt/icu4c/lib/pkgconfig"
```

## Example

```golang
package main

import (
	"fmt"

	"github.com/iawaknahc/gomessageformat"
	"golang.org/x/text/language"
)

func main() {
	// Try change numFiles to 0 or 2.
	numFiles := 1
	out, err := messageformat.FormatPositional(
		language.English,
		`{0, plural,
			=0 {There are no files on disk.}
			=1 {There is only 1 file on disk.}
			other {There are # files on disk.}
		}`,
		numFiles,
	)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%s\n", out)
}
```

## Caveats

- The only implemented ApostropheMode is [DOUBLE_REQUIRED](https://unicode-org.github.io/icu-docs/apidoc/released/icu4j/com/ibm/icu/text/MessagePattern.ApostropheMode.html#DOUBLE_REQUIRED)
- Supported numeric types are `[u]int[8|16|32|64]`. Additionally, `string` is supported as long as it is in `integral[.fraction]` format.
- Plural offset must be non-negative integer.
- The supported arguments are `{arg}`, `{arg, select}`, `{arg, plural}` and `{arg, selectordinal}`.
