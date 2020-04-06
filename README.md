## Caveats

- The only implemented ApostropheMode is [DOUBLE_REQUIRED](https://unicode-org.github.io/icu-docs/apidoc/released/icu4j/com/ibm/icu/text/MessagePattern.ApostropheMode.html#DOUBLE_REQUIRED)
- Supported numeric types are `[u]int[8|16|32|64]`. Additionally, `string` is supported as long as it is in `integral[.fraction]` format.
- Plural offset must be non-negative integer.
- The supported arguments are `{arg}`, `{arg, select}`, `{arg, plural}` and `{arg, selectordinal}`.
