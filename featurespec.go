package gospec

import (
	"strings"
	"testing"
)

type featureStepKind int

const (
	isUndefined featureStepKind = iota
	isFeature
	isBackground
	isScenario
	isGiven
	isWhen
	isThen
	isTable
)

type featureStep struct {
	indent int
	kind   featureStepKind
	title  string
	cb     func()
	// ..
}

type World struct {
	// ..
	suite *FeatureSuite
}

type FeatureSuite struct {
	t               testingInterface
	world           *World
	stack           []*featureStep
	backgroundStack []*featureStep
	suites          [][]*featureStep
	indent          int
	inBackground    bool
	atSuiteIndex    int
}

func NewFeatureSuite(t testingInterface) *FeatureSuite {
	fs := &FeatureSuite{
		t:     t,
		world: &World{},
	}
	fs.world.suite = fs
	return fs
}

func (fs *FeatureSuite) API() (
	func(string, func()),
	func(string, func()),
	func(string, func()),
	func(string, func()),
	func(string, func()),
	func(string, func()),
	func() *World,
	func(columns []string, items interface{}),
) {
	return fs.Feature, fs.Background, fs.Scenario, fs.Given,
		fs.When, fs.Then, fs.World, fs.Table
}

func (fs *FeatureSuite) prevKind() featureStepKind {
	fs.t.Helper()
	if len(fs.stack) == 0 {
		return isUndefined
	}
	return fs.stack[len(fs.stack)-1].kind
}

func (fs *FeatureSuite) Feature(title string, cb func()) {
	fs.t.Helper()
	if fs.prevKind() != isUndefined {
		fs.t.Errorf("invalid position for `Feature` function, it must be at top level")
		return
	}

	s := &featureStep{
		kind:  isFeature,
		title: title,
	}
	fs.pushStack(s)

	cb()

	fs.popBackgroundFromStackIfExists()
	fs.popStack(s)
	fs.backgroundStack = []*featureStep{}

	if len(fs.stack) > 0 {
		fs.t.Errorf("expected stack to be empty but it has %d steps", len(fs.stack))
		return
	}

	fs.start()
}

func (fs *FeatureSuite) pushStack(s *featureStep) {
	fs.t.Helper()
	fs.stack = append(fs.stack, s)
}

func (fs *FeatureSuite) pushToBackgroundStack(s *featureStep) {
	fs.t.Helper()
	fs.backgroundStack = append(fs.backgroundStack, s)
}

func (fs *FeatureSuite) popBackgroundFromStackIfExists() {
	fs.t.Helper()
	if len(fs.stack) == 0 {
		return
	}

	lastStep := fs.stack[len(fs.stack)-1]
	if lastStep.kind == isBackground {
		fs.stack = fs.stack[:len(fs.stack)-1]
	}
}

func (fs *FeatureSuite) popStack(s *featureStep) {
	fs.t.Helper()
	if len(fs.stack) == 0 {
		fs.t.Errorf("unexpected empty stack")
		return
	}

	lastStep := fs.stack[len(fs.stack)-1]
	if lastStep != s {
		fs.t.Errorf("unexpected step")
		return
	}

	fs.stack = fs.stack[:len(fs.stack)-1]
}

func (fs *FeatureSuite) popStackUntilStep(s *featureStep) {
	fs.t.Helper()
	if len(fs.stack) == 0 {
		fs.t.Errorf("unexpected empty stack")
		return
	}

	index := fs.findIndexOfStep(s)
	if index < 0 {
		return
	}

	if index+1 > len(fs.stack) {
		fs.t.Errorf("out of bound index")
		return
	}

	fs.stack = fs.stack[:index+1]
}

func (fs *FeatureSuite) findIndexOfStep(s *featureStep) int {
	fs.t.Helper()
	if len(fs.stack) == 0 {
		return -1
	}

	for i := len(fs.stack) - 1; i >= 0; i-- {
		if fs.stack[i] == s {
			return i
		}
	}

	return -1
}

func (fs *FeatureSuite) Background(title string, cb func()) {
	fs.t.Helper()
	if fs.prevKind() != isFeature {
		fs.t.Errorf("invalid position for `Background` function, it must be inside a `Feature` call")
		return
	}

	s := &featureStep{
		kind:  isBackground,
		title: title,
	}

	fs.inBackground = true
	fs.pushToBackgroundStack(s)

	cb()

	fs.inBackground = false
}

func (fs *FeatureSuite) Scenario(title string, cb func()) {
	fs.t.Helper()
	if fs.prevKind() != isFeature && fs.prevKind() != isBackground {
		fs.t.Errorf("invalid position for `Scenario` function, it must be inside a `Feature` call")
		return
	}
	s := &featureStep{
		kind:  isScenario,
		title: title,
	}
	fs.pushStack(s)

	cb()

	if len(fs.stack) > 0 {
		fs.copyStack()
		fs.popStackUntilStep(s)
	}

	fs.popStack(s)
}

func (fs *FeatureSuite) Given(title string, cb func()) {
	fs.t.Helper()

	s := &featureStep{
		kind:  isGiven,
		title: title,
		cb:    cb,
	}
	if fs.inBackground {
		fs.pushToBackgroundStack(s)
	} else {
		fs.pushStack(s)
	}
}

func (fs *FeatureSuite) When(title string, cb func()) {
	fs.t.Helper()

	s := &featureStep{
		kind:  isWhen,
		title: title,
		cb:    cb,
	}

	if fs.inBackground {
		fs.pushToBackgroundStack(s)
	} else {
		fs.pushStack(s)
	}
}

func (fs *FeatureSuite) Then(title string, cb func()) {
	fs.t.Helper()

	s := &featureStep{
		kind:  isThen,
		title: title,
		cb:    cb,
	}
	if fs.inBackground {
		fs.pushToBackgroundStack(s)
	} else {
		fs.pushStack(s)
	}
}

func (fs *FeatureSuite) copyStack() {
	fs.t.Helper()
	if len(fs.stack) <= 0 {
		return
	}

	var suite []*featureStep
	for _, s := range fs.stack[:1] {
		suite = append(suite, s)
	}
	for _, s := range fs.backgroundStack {
		suite = append(suite, s)
	}
	for _, s := range fs.stack[1:] {
		suite = append(suite, s)
	}
	fs.suites = append(fs.suites, suite)
}

func (fs *FeatureSuite) Table(columns []string, items interface{}) {
	fs.t.Helper()
	//fs.steps = append(fs.steps, &featureStep{
	//	kind: isTable,
	//	// title: title,
	//	// cb:    cb,
	//	cb: func() {
	//		// ..
	//		// fmt.Println()
	//		// fmt.Printf("      XXXX\n")
	//		items2 := reflect.ValueOf(items)
	//		// fmt.Printf("ITEMS: %#v\n", items2)
	//		// fmt.Println()
	//
	//		if items2.Kind() != reflect.Slice {
	//			panic("EXPECTED SLICE...\n")
	//			return
	//		}
	//
	//		columnWidths := make(map[string]int, items2.Len())
	//		_ = columnWidths
	//
	//		for _, x := range columns {
	//			columnWidths[x] = len(x)
	//		}
	//
	//		rows := []map[string]string{}
	//
	//		for i := 0; i < items2.Len(); i++ {
	//			item := items2.Index(i)
	//			if item.Kind() == reflect.Struct {
	//				row := map[string]string{}
	//				v := reflect.Indirect(item)
	//				for j := 0; j < v.NumField(); j++ {
	//					name := v.Type().Field(j).Name
	//					value := v.Field(j).Interface()
	//					max, ok := columnWidths[name]
	//					if !ok {
	//						continue
	//					}
	//					switch z := value.(type) {
	//					case string:
	//						if len(z) > max {
	//							columnWidths[name] = len(z)
	//						}
	//						row[name] = z
	//					case float64, float32:
	//						ff := fmt.Sprintf("%.2f", z)
	//						if len(ff) > max {
	//							columnWidths[name] = len(ff)
	//						}
	//						row[name] = ff
	//						// ff, err := strconv.ParseFloat(z)
	//						// _ = ff
	//						// _ = err
	//					case int, int8, int16, int32, int64:
	//						ff := fmt.Sprintf("%d", z)
	//						if len(ff) > max {
	//							columnWidths[name] = len(ff)
	//						}
	//						row[name] = ff
	//						// ff, err := strconv.ParseFloat(z)
	//						// _ = ff
	//						// _ = err
	//					}
	//					// fmt.Println(name, value)
	//				}
	//				rows = append(rows, row)
	//			}
	//		}
	//		var sb strings.Builder
	//		sb.WriteString("      |")
	//		for _, c := range columns {
	//			_ = c
	//			sb.WriteString(fmt.Sprintf(" %-"+strconv.Itoa(columnWidths[c])+"s ", c))
	//			//
	//			sb.WriteString("|")
	//		}
	//		sb.WriteString("\n")
	//
	//		for _, r := range rows {
	//			_ = r
	//			sb.WriteString("      |")
	//			for _, c := range columns {
	//				sb.WriteString(fmt.Sprintf(" %-"+strconv.Itoa(columnWidths[c])+"s ", r[c]))
	//				_ = c
	//
	//				sb.WriteString("|")
	//			}
	//			sb.WriteString("\n")
	//		}
	//
	//		fmt.Printf("%s", sb.String())
	//		// fmt.Printf("ROWS: %#v\n", rows)
	//		// fmt.Printf("COLUMN WIDTHS: %#v\n", columnWidths)
	//		// fmt.Println()
	//	},
	//})
}

func (fs *FeatureSuite) World() *World {
	fs.t.Helper()
	return fs.world
}

func buildSuiteTitleForFeature(suite []*featureStep) string {
	var sb strings.Builder
	for i, s := range suite {
		if s.kind == isFeature || s.kind == isBackground || s.kind == isScenario {
			if i != 0 {
				sb.WriteString("/")
			}
			sb.WriteString(strings.TrimSpace(s.title))
		}
	}
	return sb.String()
}

func (fs *FeatureSuite) start() {
	for _, suite := range fs.suites[fs.atSuiteIndex:] {
		fs.atSuiteIndex++
		fs.t.Run(buildSuiteTitleForFeature(suite), func(t *testing.T) {
			for _, s := range suite {
				if s.cb != nil {
					s.cb()
				}
			}
		})
	}
}

//func debugFeature(fs *FeatureSuite) {
//	fmt.Printf("FEATURE: %#v\n", fs)
//	for _, s := range fs.steps {
//		fmt.Printf("STEP: %#v\n", s)
//	}
//}

//func printFeature(fs *FeatureSuite) {
//	// var (
//	// 	inBackground bool
//	// 	inScenario   bool
//	// )
//	var (
//		lastStep featureStepKind
//	)
//
//	for _, s := range fs.steps {
//		if s.kind == isFeature {
//			fmt.Printf("Feature: %s\n", s.title)
//			lastStep = s.kind
//			continue
//		}
//
//		if s.kind == isBackground {
//			fmt.Printf("\n  Background: %s\n", s.title)
//			// inBackground = true
//			lastStep = s.kind
//			continue
//		}
//
//		if s.kind == isScenario {
//			fmt.Printf("\n  Scenario: %s\n", s.title)
//			// inScenario = true
//			// inBackground = false
//			lastStep = s.kind
//			continue
//		}
//
//		if s.kind == isGiven {
//			if lastStep == s.kind {
//				fmt.Printf("    And %s\n", s.title)
//				continue
//			}
//			fmt.Printf("    Given %s\n", s.title)
//			lastStep = s.kind
//			continue
//			// ..
//		}
//
//		if s.kind == isWhen {
//			if lastStep == s.kind {
//				fmt.Printf("    And %s\n", s.title)
//				continue
//			}
//			fmt.Printf("    When %s\n", s.title)
//			lastStep = s.kind
//			continue
//			// ..
//		}
//
//		if s.kind == isThen {
//			if lastStep == s.kind {
//				fmt.Printf("    And %s\n", s.title)
//				continue
//			}
//			fmt.Printf("    Then %s\n", s.title)
//			lastStep = s.kind
//			continue
//			// ..
//		}
//
//		if s.kind == isTable {
//			// fmt.Printf("      TABLE HERE...\n")
//			s.cb()
//			// if lastStep == s.kind {
//			// 	fmt.Printf("    And %s\n", s.title)
//			// 	continue
//			// }
//			// fmt.Printf("    Then %s\n", s.title)
//			// lastStep = s.kind
//			continue
//			// ..
//		}
//
//		// fmt.Printf("STEP: %#v\n", s)
//	}
//	fmt.Println()
//}
