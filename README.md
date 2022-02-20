# progress-bar

[![Go CI](https://github.com/elulcao/progress-bar/actions/workflows/go.yaml/badge.svg)](https://github.com/elulcao/progress-bar/actions/workflows/go.yaml)
[![CodeQl CI](https://github.com/elulcao/progress-bar/actions/workflows/codeql-analysis.yaml/badge.svg)](https://github.com/elulcao/progress-bar/actions/workflows/codeql-analysis.yaml)

---

<p
    align="center">
    <img
        src="./assets/demo-01.gif"
        alt="Demo 01"
        width="600"
        height="400"
    />
</p>

---

Another simple progress bar, but with a package installer flavor. The progress bar displays the
progress like the well known package manager. This version adapts well to the windows size, small
or big window size the progress bar will adapt to the change.

Characters used by default are "#" and "." but can be changed by updating variables:

```go
var (
 doneStr    = "#"
 ongoingStr = "."
)
```

## Usage

The `progress-bar.go` should be modified to your needs, currently implementation is limited to
printing numbers from 1 to 100. In a real world scenario, the `count` variable should be used to
be updated by another routine in your code.

```go
 for count = 1; count <= total; count++ {
  renderPBar()
  time.Sleep(time.Second)
  fmt.Println(count) // Update here to your needs
 }
```
