// Package tex emits the paper's generated LaTeX fragments. Every file starts
// with a provenance header (generator, config hash, input hashes) and must
// never be edited by hand.
package tex

import (
	"fmt"
	"io"
	"strings"
)

var escaper = strings.NewReplacer(
	`\`, `\textbackslash{}`,
	`&`, `\&`,
	`%`, `\%`,
	`$`, `\$`,
	`#`, `\#`,
	`_`, `\_`,
	`{`, `\{`,
	`}`, `\}`,
	`~`, `\textasciitilde{}`,
	`^`, `\textasciicircum{}`,
)

// Escape makes s safe for LaTeX text mode.
func Escape(s string) string { return escaper.Replace(s) }

// Header writes the do-not-edit provenance comment block.
func Header(w io.Writer, generator, configHash string, inputs map[string]string) {
	fmt.Fprintf(w, "%% GENERATED FILE — do not edit by hand (CLAUDE.md rule: *.gen.tex are never hand-edited).\n")
	fmt.Fprintf(w, "%% generator: %s\n", generator)
	fmt.Fprintf(w, "%% config sha256: %s\n", configHash)
	keys := make([]string, 0, len(inputs))
	for k := range inputs {
		keys = append(keys, k)
	}
	sortStrings(keys)
	for _, k := range keys {
		fmt.Fprintf(w, "%% input %s sha256: %s\n", k, inputs[k])
	}
}

// Table writes a booktabs tabular wrapped in a table environment. Cells are
// escaped by the caller only if they contain raw content; use Escape.
func Table(w io.Writer, caption, label, colspec string, header []string, rows [][]string) {
	fmt.Fprintf(w, "\\begin{table}[t]\n\\centering\n\\small\n")
	fmt.Fprintf(w, "\\begin{tabular}{%s}\n\\toprule\n", colspec)
	fmt.Fprintf(w, "%s \\\\\n\\midrule\n", strings.Join(header, " & "))
	for _, r := range rows {
		fmt.Fprintf(w, "%s \\\\\n", strings.Join(r, " & "))
	}
	fmt.Fprintf(w, "\\bottomrule\n\\end{tabular}\n")
	fmt.Fprintf(w, "\\caption{%s}\n\\label{%s}\n\\end{table}\n", caption, label)
}

func sortStrings(xs []string) {
	for i := 1; i < len(xs); i++ {
		for j := i; j > 0 && xs[j] < xs[j-1]; j-- {
			xs[j], xs[j-1] = xs[j-1], xs[j]
		}
	}
}
