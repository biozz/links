# Links

An opinionated tool to efficiently manage large amount of bookmarks with aliases.

## Usage

```
go run . serve
```

- Go to admin interface and add a device
- Use the device id and token to login
- Go to help page and follow the instructions on how to setup your custom search engine

On the New form, tags and default arguments are separated by spaces. Default
arguments are stored as a JSON array (empty when omitted) and used when an alias
is entered without arguments. For example, a `gl` alias with URL
`https://gitlab.something/%s/%s/-/merge_requests` and defaults `extranet extrota`
opens `https://gitlab.something/extranet/extrota/-/merge_requests` when you enter
`gl`. Explicit arguments replace the defaults.
