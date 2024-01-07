package gospec

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"testing"
)

type featureStepKind int

const (
	isFeature featureStepKind = iota
	isBackground
	isScenario
	isGiven
	isWhen
	isThen
	isTable
)

type featureStep struct {
	kind  featureStepKind
	title string
	cb    func()
	// ..
}

type World struct {
	// ..
	suite *FeatureSuite
}

type FeatureSuite struct {
	t     testingInterface
	world *World
	steps []*featureStep
}

func NewFeatureSuite(t *testing.T) *FeatureSuite {
	f := &FeatureSuite{
		t:     t,
		world: &World{},
	}
	f.world.suite = f
	return f
}

func (s *FeatureSuite) Feature(title string, cb func()) {
	s.steps = append(s.steps, &featureStep{
		kind:  isFeature,
		title: title,
		cb:    cb,
	})

	cb()
}

func (s *FeatureSuite) Background(title string, cb func()) {
	s.steps = append(s.steps, &featureStep{
		kind:  isBackground,
		title: title,
		cb:    cb,
	})

	cb()
}

func (s *FeatureSuite) Scenario(title string, cb func()) {
	s.steps = append(s.steps, &featureStep{
		kind:  isScenario,
		title: title,
		cb:    cb,
	})

	cb()
}

func (s *FeatureSuite) Given(title string, cb func()) {
	s.steps = append(s.steps, &featureStep{
		kind:  isGiven,
		title: title,
		cb:    cb,
	})

	cb()
}

func (s *FeatureSuite) When(title string, cb func()) {
	s.steps = append(s.steps, &featureStep{
		kind:  isWhen,
		title: title,
		cb:    cb,
	})
}

func (s *FeatureSuite) Then(title string, cb func()) {
	s.steps = append(s.steps, &featureStep{
		kind:  isThen,
		title: title,
		cb:    cb,
	})
}

func (s *FeatureSuite) Table(columns []string, items interface{}) {
	s.steps = append(s.steps, &featureStep{
		kind: isTable,
		// title: title,
		// cb:    cb,
		cb: func() {
			// ..
			// fmt.Println()
			// fmt.Printf("      XXXX\n")
			items2 := reflect.ValueOf(items)
			// fmt.Printf("ITEMS: %#v\n", items2)
			// fmt.Println()

			if items2.Kind() != reflect.Slice {
				panic("EXPECTED SLICE...\n")
				return
			}

			columnWidths := make(map[string]int, items2.Len())
			_ = columnWidths

			for _, x := range columns {
				columnWidths[x] = len(x)
			}

			rows := []map[string]string{}

			for i := 0; i < items2.Len(); i++ {
				item := items2.Index(i)
				if item.Kind() == reflect.Struct {
					row := map[string]string{}
					v := reflect.Indirect(item)
					for j := 0; j < v.NumField(); j++ {
						name := v.Type().Field(j).Name
						value := v.Field(j).Interface()
						max, ok := columnWidths[name]
						if !ok {
							continue
						}
						switch z := value.(type) {
						case string:
							if len(z) > max {
								columnWidths[name] = len(z)
							}
							row[name] = z
						case float64, float32:
							ff := fmt.Sprintf("%.2f", z)
							if len(ff) > max {
								columnWidths[name] = len(ff)
							}
							row[name] = ff
							// ff, err := strconv.ParseFloat(z)
							// _ = ff
							// _ = err
						case int, int8, int16, int32, int64:
							ff := fmt.Sprintf("%d", z)
							if len(ff) > max {
								columnWidths[name] = len(ff)
							}
							row[name] = ff
							// ff, err := strconv.ParseFloat(z)
							// _ = ff
							// _ = err
						}
						// fmt.Println(name, value)
					}
					rows = append(rows, row)
				}
			}
			var sb strings.Builder
			sb.WriteString("      |")
			for _, c := range columns {
				_ = c
				sb.WriteString(fmt.Sprintf(" %-"+strconv.Itoa(columnWidths[c])+"s ", c))
				//
				sb.WriteString("|")
			}
			sb.WriteString("\n")

			for _, r := range rows {
				_ = r
				sb.WriteString("      |")
				for _, c := range columns {
					sb.WriteString(fmt.Sprintf(" %-"+strconv.Itoa(columnWidths[c])+"s ", r[c]))
					_ = c

					sb.WriteString("|")
				}
				sb.WriteString("\n")
			}

			fmt.Printf("%s", sb.String())
			// fmt.Printf("ROWS: %#v\n", rows)
			// fmt.Printf("COLUMN WIDTHS: %#v\n", columnWidths)
			// fmt.Println()
		},
	})
}

func (s *FeatureSuite) World() *World {
	return s.world
}

func (s *FeatureSuite) Start() {
	// fmt.Printf("START....\n")
	printFeature(s)
	// debugFeature(s)
}

func debugFeature(f *FeatureSuite) {
	fmt.Printf("FEATURE: %#v\n", f)
	for _, s := range f.steps {
		fmt.Printf("STEP: %#v\n", s)
	}
}

func printFeature(s *FeatureSuite) {
	// var (
	// 	inBackground bool
	// 	inScenario   bool
	// )
	var (
		lastStep featureStepKind
	)

	for _, s := range s.steps {
		if s.kind == isFeature {
			fmt.Printf("Feature: %s\n", s.title)
			lastStep = s.kind
			continue
		}

		if s.kind == isBackground {
			fmt.Printf("\n  Background: %s\n", s.title)
			// inBackground = true
			lastStep = s.kind
			continue
		}

		if s.kind == isScenario {
			fmt.Printf("\n  Scenario: %s\n", s.title)
			// inScenario = true
			// inBackground = false
			lastStep = s.kind
			continue
		}

		if s.kind == isGiven {
			if lastStep == s.kind {
				fmt.Printf("    And %s\n", s.title)
				continue
			}
			fmt.Printf("    Given %s\n", s.title)
			lastStep = s.kind
			continue
			// ..
		}

		if s.kind == isWhen {
			if lastStep == s.kind {
				fmt.Printf("    And %s\n", s.title)
				continue
			}
			fmt.Printf("    When %s\n", s.title)
			lastStep = s.kind
			continue
			// ..
		}

		if s.kind == isThen {
			if lastStep == s.kind {
				fmt.Printf("    And %s\n", s.title)
				continue
			}
			fmt.Printf("    Then %s\n", s.title)
			lastStep = s.kind
			continue
			// ..
		}

		if s.kind == isTable {
			// fmt.Printf("      TABLE HERE...\n")
			s.cb()
			// if lastStep == s.kind {
			// 	fmt.Printf("    And %s\n", s.title)
			// 	continue
			// }
			// fmt.Printf("    Then %s\n", s.title)
			// lastStep = s.kind
			continue
			// ..
		}

		// fmt.Printf("STEP: %#v\n", s)
	}
	fmt.Println()
}
