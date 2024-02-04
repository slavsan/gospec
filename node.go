package gospec

import (
	"fmt"
	"strings"
)

type node struct {
	step     *step
	children []*node
}

type tree []*node

func (t tree) String(output *output1, suite *SpecSuite) string {
	var sb strings.Builder
	for _, n := range t {
		n.write(&sb, 0, output, suite)
		sb.WriteString("\n")
	}
	return sb.String()
}

func (n *node) write(sb *strings.Builder, indent int, output *output1, suite *SpecSuite) { //nolint:cyclop
	var (
		format = "%s%s%s%s%s%s%s%s"
		args   = []any{
			strings.Repeat(output.indentStep, indent),
			"", // only
			"", // green,
			"", // icon
			"", // noColor,
			"", // gray,
			n.step.title,
			"", // noColor,
		}
	)

	if output.colorful {
		if n.step.block == isDescribe {
			args[5] = bold
			args[7] = noColor
		}
		if n.step.block == isIt {
			args[2] = green
			args[4] = noColor

			args[5] = gray
			args[7] = noColor
		}
	}

	if output.durations && n.step.block == isIt {
		format += " (%dms)"
		args = append(args, n.step.timeSpent.Milliseconds())
	}

	if n.step.block == isIt {
		args[3] = "✔ "
	}

	if n.step.block == isIt && ((n.step.t != nil && n.step.t.Skipped()) || suite.skipped(n.step.index)) {
		if output.colorful {
			args[2] = cyan
			args[5] = cyan
		}
		args[3] = "[skip] "
	}

	if n.step.block == isIt && ((n.step.t != nil && n.step.t.Failed()) || suite.failed(n.step.index)) {
		if output.colorful {
			args[2] = red
			args[5] = red
		}
		args[3] = "⨯ "
	}

	if n.step.block == isIt && n.step.only {
		if output.colorful {
			args[1] = fmt.Sprintf("%s[only]%s ", yellow, noColor)
		} else {
			args[1] = "[only] "
		}
	}

	if output.printFilenames {
		format += "\t%s:%d"
		args = append(args, strings.TrimPrefix(n.step.file, basePath), n.step.lineNo)
	}

	sb.WriteString(fmt.Sprintf(format, args...))
	sb.WriteString("\n")

	for _, c := range n.children {
		c.write(sb, indent+1, output, suite)
	}
}
