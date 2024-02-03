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

func (t tree) String(output *output1) string {
	var sb strings.Builder
	for _, n := range t {
		n.write(&sb, 0, output)
		sb.WriteString("\n")
	}
	return sb.String()
}

func (n *node) write(sb *strings.Builder, indent int, output *output1) { //nolint:cyclop
	var (
		format = "%s%s%s%s%s%s%s"
		args   = []any{
			strings.Repeat(output.indentStep, indent),
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
			args[4] = bold
			args[6] = noColor
		}
		if n.step.block == isIt {
			args[1] = green
			args[3] = noColor

			args[4] = gray
			args[6] = noColor
		}
	}

	if output.durations && n.step.block == isIt {
		format += " (%dms)"
		args = append(args, n.step.timeSpent.Milliseconds())
	}

	if n.step.block == isIt {
		args[2] = "✔ "

		if !n.step.executed && n.step.failedAt == 0 {
			if output.colorful {
				args[1] = yellow
			}
			args[2] = "s "
		}
	}

	if n.step.block == isIt && n.step.failed {
		if n.step.failedAt == 1 {
			if output.colorful {
				args[1] = red
			}
			args[2] = "⨯ "
		} else {
			args[2] = "s "
		}
	}

	if output.printFilenames {
		format += "\t%s:%d"
		args = append(args, strings.TrimPrefix(n.step.file, basePath), n.step.lineNo)
	}

	sb.WriteString(fmt.Sprintf(format, args...))
	sb.WriteString("\n")

	for _, c := range n.children {
		c.write(sb, indent+1, output)
	}
}
