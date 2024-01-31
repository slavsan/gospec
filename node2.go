package gospec

import (
	"fmt"
	"strings"
)

type node2 struct {
	step     *featureStep
	children []*node2
}

type tree2 []*node2

func (t tree2) String(suite *FeatureSuite) string {
	var sb strings.Builder
	for _, n := range t {
		n.write(&sb, 0, suite)
		sb.WriteString("\n")
	}
	return sb.String()
}

func (n *node2) write(sb *strings.Builder, indent int, suite *FeatureSuite) { //nolint:cyclop
	if n.step.printed {
		return
	}

	n.step.printed = true

	if n.step.kind == isFeature {
		if suite.printFilenames {
			sb.WriteString(
				fmt.Sprintf("Feature: %s\t%s:%d\n",
					n.step.title,
					strings.TrimPrefix(n.step.file, basePath), n.step.lineNo,
				),
			)
		} else {
			sb.WriteString(
				fmt.Sprintf("Feature: %s\n", n.step.title),
			)
		}
	}

	if n.step.kind == isBackground {
		if suite.printFilenames {
			sb.WriteString(
				fmt.Sprintf("\n%sBackground:\t%s:%d\n",
					suite.indentStep,
					strings.TrimPrefix(n.step.file, basePath), n.step.lineNo,
				),
			)
		} else {
			sb.WriteString(fmt.Sprintf("\n%sBackground:\n", suite.indentStep))
		}
	}

	if n.step.kind == isScenario {
		if suite.printFilenames {
			sb.WriteString(
				fmt.Sprintf("\n%sScenario: %s\t%s:%d\n",
					suite.indentStep,
					n.step.title,
					strings.TrimPrefix(n.step.file, basePath), n.step.lineNo,
				),
			)
		} else {
			sb.WriteString(
				fmt.Sprintf("\n%sScenario: %s\n", suite.indentStep, n.step.title),
			)
		}
	}

	if n.step.kind == isGiven {
		if suite.printFilenames {
			sb.WriteString(
				fmt.Sprintf("%sGiven %s\t%s:%d\n",
					strings.Repeat(suite.indentStep, 2),
					n.step.title,
					strings.TrimPrefix(n.step.file, basePath), n.step.lineNo,
				),
			)
		} else {
			sb.WriteString(
				fmt.Sprintf("%sGiven %s\n",
					strings.Repeat(suite.indentStep, 2),
					n.step.title,
				),
			)
		}
	}

	if n.step.kind == isWhen {
		if suite.printFilenames {
			sb.WriteString(
				fmt.Sprintf("%sWhen %s\t%s:%d\n",
					strings.Repeat(suite.indentStep, 2),
					n.step.title,
					strings.TrimPrefix(n.step.file, basePath), n.step.lineNo,
				),
			)
		} else {
			sb.WriteString(
				fmt.Sprintf("%sWhen %s\n",
					strings.Repeat(suite.indentStep, 2),
					n.step.title,
				),
			)
		}
	}

	if n.step.kind == isThen {
		if suite.printFilenames {
			sb.WriteString(
				fmt.Sprintf("%sThen %s\t%s:%d\n",
					strings.Repeat(suite.indentStep, 2),
					n.step.title,
					strings.TrimPrefix(n.step.file, basePath), n.step.lineNo,
				),
			)
		} else {
			sb.WriteString(
				fmt.Sprintf("%sThen %s\n",
					strings.Repeat(suite.indentStep, 2),
					n.step.title,
				),
			)
		}
	}

	if n.step.kind == isTable {
		sb.WriteString(
			fmt.Sprintf("%s\n", n.step.title),
		)
	}

	for _, c := range n.children {
		c.write(sb, indent+1, suite)
	}
}
