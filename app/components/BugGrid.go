//go:build wasm

package components

import app_types "exocomp-app/types"
import utils_fmt "exocomp-app/utils/fmt"
import utils_path "exocomp-app/utils/path"
import "github.com/cookiengineer/gooey/bindings/dom"
import "github.com/cookiengineer/gooey/components"
import "github.com/cookiengineer/gooey/components/interfaces"
import "github.com/cookiengineer/gooey/components/utils"
import std_html "html"
import "sort"
import "strings"

// BugGrid renders the Bugs view. It owns `[]*app_types.Bug`, groups them into
// package tiles and renders their checkbox list items. Its search/collapse
// controls live in the paired BugFooter.
type BugGrid struct {
	Component   *components.Component `json:"component"`
	Collapsed   bool                  `json:"collapsed"`
	bugs        []*app_types.Bug      `json:"-"`
	packages    *dom.Element          `json:"-"`
	Footer      *BugFooter            `json:"-"`
	initialized bool                  `json:"-"`
}

func NewBugGrid() *BugGrid {

	element := dom.GetDocument().CreateElement("bug-grid")
	component := components.NewComponent(element)

	return &BugGrid{
		Component: &component,
		bugs:      make([]*app_types.Bug, 0),
	}

}

func ToBugGrid(element *dom.Element) *BugGrid {

	component := components.NewComponent(element)

	return &BugGrid{
		Component: &component,
		bugs:      make([]*app_types.Bug, 0),
	}

}

func (grid *BugGrid) SetBugs(bugs []*app_types.Bug, collapsed bool) {

	grid.bugs = bugs
	grid.Collapsed = collapsed

	grid.Render()

}

func (grid *BugGrid) Disable() bool {
	return false
}

func (grid *BugGrid) Enable() bool {
	return false
}

func (grid *BugGrid) SetScheduler(scheduler interfaces.Scheduler) {

	if grid.Component != nil {
		grid.Component.SetScheduler(scheduler)
	}

	if grid.Footer != nil {
		grid.Footer.SetScheduler(scheduler)
	}

}

func (grid *BugGrid) placeholder() string {

	if grid.Component != nil && grid.Component.Element != nil {

		tmp := strings.TrimSpace(grid.Component.Element.GetAttribute("data-placeholder"))

		if tmp != "" {
			return tmp
		}

	}

	return "Search ..."

}

func (grid *BugGrid) ensure() {

	if grid.initialized == true {
		return
	}

	if grid.Component == nil || grid.Component.Element == nil {
		return
	}

	grid.Component.Element.SetInnerHTML("<div class=\"packages\"></div>")

	grid.packages = grid.Component.Element.QuerySelector("div.packages")

	grid.Footer = NewBugFooter(grid.placeholder())
	grid.Component.Element.Append(grid.Footer.Component.Element)
	grid.Footer.Mount()

	grid.initialized = true

}

func (grid *BugGrid) Mount() bool {

	if grid.Component == nil || grid.Component.Element == nil {
		return false
	}

	grid.ensure()
	grid.Render()

	return true

}

func (grid *BugGrid) Render() *dom.Element {

	if grid.Component == nil || grid.Component.Element == nil {
		return nil
	}

	grid.ensure()

	if grid.packages != nil {

		grouped := make(map[string][]*app_types.Bug)
		package_paths := make([]string, 0)

		for _, bug := range grid.bugs {

			if bug == nil {
				continue
			}

			package_path := bug.PackagePath()

			if grid.Collapsed == true {
				package_path = utils_path.Collapse(package_path)
			}

			if _, ok := grouped[package_path]; ok == false {
				grouped[package_path] = make([]*app_types.Bug, 0)
				package_paths = append(package_paths, package_path)
			}

			grouped[package_path] = append(grouped[package_path], bug)

		}

		sort.Strings(package_paths)

		html := ""

		for _, package_path := range package_paths {

			label := package_path

			if label == "" {
				label = "."
			}

			html += "<article class=\"package\" data-package=\"" + std_html.EscapeString(package_path) + "\">"
			html += "<h3>" + std_html.EscapeString(label) + "</h3>"
			html += "<ul>"

			for _, bug := range grouped[package_path] {

				state := ""

				if bug.IsFixed == true {
					state = " checked=\"\""
				}

				label := strings.TrimPrefix(bug.File, "./") + "#" + bug.Symbol + ": " + bug.Description

				html += "<li><input type=\"checkbox\" disabled=\"\"" + state + "/><span>" + utils_fmt.SanitizeContent(label) + "</span></li>"

			}

			html += "</ul>"
			html += "</article>"

		}

		grid.packages.SetInnerHTML(html)

	}

	if grid.Footer != nil {
		grid.Footer.SetCollapsed(grid.Collapsed)
	}

	return grid.Component.Element

}

func (grid *BugGrid) Query(query string) interfaces.Component {

	selectors := utils.SplitQuery(query)

	if len(selectors) == 1 {

		if grid.Component != nil && grid.Component.Element != nil {

			if utils.MatchesQuery(grid.Component.Element, selectors[0]) == true {
				return grid
			}

		}

	}

	return nil

}

func (grid *BugGrid) String() string {

	if grid.Component == nil || grid.Component.Element == nil {
		return ""
	}

	tag := strings.ToLower(grid.Component.Element.TagName)

	return "<" + tag + ">" + grid.Component.Element.InnerHTML + "</" + tag + ">"

}

func (grid *BugGrid) Unmount() bool {

	if grid.Footer != nil {
		grid.Footer.Unmount()
	}

	return true

}

