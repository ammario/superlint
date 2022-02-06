# superlint

`superlint` is a linting system configured by user-defined Go code. Instead of a bespoke (and poorly documented) matching
language, `superlint` lets the user define arbitrary matching functions.

`superlint` rules are **language-agonistic** and run **codebase-wide**. They're both fast (no default AST) and capable of
encorcing arbitrarily complex rules. A ruleset can create an AST if it so desires.

For example, `superlint` can catch:
* That each Go binary has an accompanying Make entry
* That each http Handler has an accompanying test

## Basic Usage

```go
package rules

func main() {
	
}
```

```
$ go build -buildmode=plugin rules/ -o rules.so
$ superlint rules.so
```



## Architecture

`superlint` loads your ruleset as a Go plugin.