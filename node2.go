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

func (t tree2) String(output *output1) string {
	var sb strings.Builder
	for _, n := range t {
		n.write(&sb, 0, output)
		sb.WriteString("\n")
	}
	return sb.String()
}

func (n *node2) feature(output *output1) (string, []any) {
	format := "%sFeature:%s %s"
	args := []any{"", "", n.step.title}
	if output.colorful {
		args[0] = bold
		args[1] = noBold
	}
	return format, args
}

func (n *node2) background(output *output1) (string, []any) {
	format := "\n%s%sBackground:%s"
	args := []any{output.indentStep, "", ""}
	if output.colorful {
		args[1] = bold
		args[2] = noBold
	}
	return format, args
}

func (n *node2) scenario(output *output1) (string, []any) {
	format := "\n%s%sScenario:%s %s"
	args := []any{output.indentStep, "", "", n.step.title}
	if output.colorful {
		args[1] = bold
		args[2] = noBold
	}
	return format, args
}

func (n *node2) given(output *output1) (string, []any) {
	format := "%s%sGiven%s %s"
	args := []any{strings.Repeat(output.indentStep, 2), "", "", n.step.title}
	if output.colorful {
		args[1] = cyan
		args[2] = noColor
	}
	return format, args
}

func (n *node2) when(output *output1) (string, []any) {
	format := "%s%sWhen%s %s"
	args := []any{strings.Repeat(output.indentStep, 2), "", "", n.step.title}
	if output.colorful {
		args[1] = green
		args[2] = noColor
	}
	return format, args
}

func (n *node2) then(output *output1) (string, []any) {
	format := "%s%sThen%s %s"
	args := []any{strings.Repeat(output.indentStep, 2), "", "", n.step.title}
	if output.colorful {
		args[1] = yellow
		args[2] = noColor
	}
	return format, args
}

func (n *node2) write(sb *strings.Builder, indent int, output *output1) {
	m := map[featureStepKind]func(output *output1) (string, []any){
		isFeature:    n.feature,
		isBackground: n.background,
		isScenario:   n.scenario,
		isGiven:      n.given,
		isWhen:       n.when,
		isThen:       n.then,
	}

	if f, ok := m[n.step.kind]; ok {
		format, args := f(output)

		if output.printFilenames {
			format += "\t%s:%d"
			args = append(args, strings.TrimPrefix(n.step.file, basePath), n.step.lineNo)
		}

		format += "\n"
		sb.WriteString(fmt.Sprintf(format, args...))
	}

	if n.step.kind == isTable {
		lines := strings.Split(n.step.title, "\n")
		for _, l := range lines {
			sb.WriteString(
				fmt.Sprintf("%s%s\n",
					strings.Repeat(output.indentStep, 3), l,
				),
			)
		}
	}

	for _, c := range n.children {
		c.write(sb, indent+1, output)
	}
}
