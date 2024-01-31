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

func (t tree) String(suite *SpecSuite) string {
	var sb strings.Builder
	for _, n := range t {
		n.write(&sb, 0, suite)
		sb.WriteString("\n")
	}
	return sb.String()
}

func (n *node) write(sb *strings.Builder, indent int, suite *SpecSuite) {
	if n.step.printed {
		return
	}

	n.step.printed = true

	icon := "✔ "

	if n.step.block == isIt { //nolint:nestif
		if suite.printFilenames {
			if n.step.failed {
				if n.step.failedAt == 1 {
					icon = "⨯ "
				} else {
					icon = "s "
				}
			}

			sb.WriteString(
				fmt.Sprintf("%s%s%s\t%s:%d\n",
					strings.Repeat("\t", indent),
					icon,
					n.step.title,
					strings.TrimPrefix(n.step.file, basePath), n.step.lineNo,
				),
			)
			return
		}

		if !n.step.executed && n.step.failedAt == 0 {
			icon = "s "
		}

		if n.step.failed {
			if n.step.failedAt == 1 {
				icon = "⨯ "
			} else if n.step.failedAt != 0 {
				icon = "s "
			}
		}

		sb.WriteString(
			fmt.Sprintf("%s%s%s\n",
				strings.Repeat("\t", indent),
				icon,
				n.step.title,
			),
		)

		return
	}

	if suite.printFilenames {
		sb.WriteString(
			fmt.Sprintf("%s%s\t%s:%d\n",
				strings.Repeat("\t", indent),
				n.step.title,
				strings.TrimPrefix(n.step.file, basePath), n.step.lineNo,
			),
		)
	} else {
		sb.WriteString(
			fmt.Sprintf("%s%s\n",
				strings.Repeat("\t", indent),
				n.step.title,
			),
		)
	}

	for _, c := range n.children {
		c.write(sb, indent+1, suite)
	}
}
