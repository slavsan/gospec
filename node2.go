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

func (t tree2) String(suite *FeatureSuite, output *output2) string {
	var sb strings.Builder
	for _, n := range t {
		n.write(&sb, 0, output, suite)
		sb.WriteString("\n")
	}
	return sb.String()
}

func (n *node2) feature(output *output2, _ *FeatureSuite) (string, []any) {
	format := "%sFeature:%s %s"
	args := []any{"", "", n.step.title}
	if output.colorful {
		args[0] = bold
		args[1] = noBold
	}
	return format, args
}

func (n *node2) background(output *output2, suite *FeatureSuite) (string, []any) {
	format := "\n%s%sBackground:%s"
	args := []any{suite.indentStep, "", ""}
	if output.colorful {
		args[1] = bold
		args[2] = noBold
	}
	return format, args
}

func (n *node2) scenario(output *output2, suite *FeatureSuite) (string, []any) {
	format := "\n%s%sScenario:%s %s"
	args := []any{suite.indentStep, "", "", n.step.title}
	if output.colorful {
		args[1] = bold
		args[2] = noBold
	}
	return format, args
}

func (n *node2) given(output *output2, suite *FeatureSuite) (string, []any) {
	format := "%s%sGiven%s %s"
	args := []any{strings.Repeat(suite.indentStep, 2), "", "", n.step.title}
	if output.colorful {
		args[1] = cyan
		args[2] = noColor
	}
	return format, args
}

func (n *node2) when(output *output2, suite *FeatureSuite) (string, []any) {
	format := "%s%sWhen%s %s"
	args := []any{strings.Repeat(suite.indentStep, 2), "", "", n.step.title}
	if output.colorful {
		args[1] = green
		args[2] = noColor
	}
	return format, args
}

func (n *node2) then(output *output2, suite *FeatureSuite) (string, []any) {
	format := "%s%sThen%s %s"
	args := []any{strings.Repeat(suite.indentStep, 2), "", "", n.step.title}
	if output.colorful {
		args[1] = yellow
		args[2] = noColor
	}
	return format, args
}

func (n *node2) write(sb *strings.Builder, indent int, output *output2, suite *FeatureSuite) {
	m := map[featureStepKind]func(output *output2, suite *FeatureSuite) (string, []any){
		isFeature:    n.feature,
		isBackground: n.background,
		isScenario:   n.scenario,
		isGiven:      n.given,
		isWhen:       n.when,
		isThen:       n.then,
	}

	if f, ok := m[n.step.kind]; ok {
		format, args := f(output, suite)

		if output.printFilenames {
			format += "\t%s:%d"
			args = append(args, strings.TrimPrefix(n.step.file, basePath), n.step.lineNo)
		}

		format += "\n"
		sb.WriteString(fmt.Sprintf(format, args...))
	}

	if n.step.kind == isTable {
		sb.WriteString(
			fmt.Sprintf("%s\n", n.step.title),
		)
	}

	for _, c := range n.children {
		c.write(sb, indent+1, output, suite)
	}
}
