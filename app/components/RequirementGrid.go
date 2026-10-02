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

// RequirementGrid renders the Requirements view. It owns
// `[]*app_types.Requirement`, groups them into package tiles and renders their
// checkbox list items. Its search/collapse controls live in the paired
// RequirementFooter.
type RequirementGrid struct {
	Component    *components.Component    `json:"component"`
	Collapsed    bool                     `json:"collapsed"`
	requirements []*app_types.Requirement `json:"-"`
	packages     *dom.Element             `json:"-"`
	Footer       *RequirementFooter       `json:"-"`
	initialized  bool                     `json:"-"`
}

func NewRequirementGrid() *RequirementGrid {

	element := dom.GetDocument().CreateElement("requirement-grid")
	component := components.NewComponent(element)

	return &RequirementGrid{
		Component:    &component,
		requirements: make([]*app_types.Requirement, 0),
	}

}

func ToRequirementGrid(element *dom.Element) *RequirementGrid {

	component := components.NewComponent(element)

	return &RequirementGrid{
		Component:    &component,
		requirements: make([]*app_types.Requirement, 0),
	}

}

func (grid *RequirementGrid) SetRequirements(requirements []*app_types.Requirement, collapsed bool) {

	grid.requirements = requirements
	grid.Collapsed = collapsed

	grid.Render()

}

func (grid *RequirementGrid) Disable() bool {
	return false
}

func (grid *RequirementGrid) Enable() bool {
	return false
}

func (grid *RequirementGrid) SetScheduler(scheduler interfaces.Scheduler) {

	if grid.Component != nil {
		grid.Component.SetScheduler(scheduler)
	}

	if grid.Footer != nil {
		grid.Footer.SetScheduler(scheduler)
	}

}

func (grid *RequirementGrid) placeholder() string {

	if grid.Component != nil && grid.Component.Element != nil {

		tmp := strings.TrimSpace(grid.Component.Element.GetAttribute("data-placeholder"))

		if tmp != "" {
			return tmp
		}

	}

	return "Search ..."

}

func (grid *RequirementGrid) ensure() {

	if grid.initialized == true {
		return
	}

	if grid.Component == nil || grid.Component.Element == nil {
		return
	}

	grid.Component.Element.SetInnerHTML("<div class=\"packages\"></div>")

	grid.packages = grid.Component.Element.QuerySelector("div.packages")

	grid.Footer = NewRequirementFooter(grid.placeholder())
	grid.Component.Element.Append(grid.Footer.Component.Element)
	grid.Footer.Mount()

	grid.initialized = true

}

func (grid *RequirementGrid) Mount() bool {

	if grid.Component == nil || grid.Component.Element == nil {
		return false
	}

	grid.ensure()
	grid.Render()

	return true

}

func (grid *RequirementGrid) Render() *dom.Element {

	if grid.Component == nil || grid.Component.Element == nil {
		return nil
	}

	grid.ensure()

	if grid.packages != nil {

		grouped := make(map[string][]*app_types.Requirement)
		package_paths := make([]string, 0)

		for _, requirement := range grid.requirements {

			if requirement == nil {
				continue
			}

			package_path := requirement.PackagePath()

			if grid.Collapsed == true {
				package_path = utils_path.Collapse(package_path)
			}

			if _, ok := grouped[package_path]; ok == false {
				grouped[package_path] = make([]*app_types.Requirement, 0)
				package_paths = append(package_paths, package_path)
			}

			grouped[package_path] = append(grouped[package_path], requirement)

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

			for _, requirement := range grouped[package_path] {

				state := ""

				if requirement.IsImplemented == true {
					state = " checked=\"\""
				}

				label := strings.TrimPrefix(requirement.File, "./") + "#" + requirement.Symbol + ": " + requirement.Declaration

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

func (grid *RequirementGrid) Query(query string) interfaces.Component {

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

func (grid *RequirementGrid) String() string {

	if grid.Component == nil || grid.Component.Element == nil {
		return ""
	}

	tag := strings.ToLower(grid.Component.Element.TagName)

	return "<" + tag + ">" + grid.Component.Element.InnerHTML + "</" + tag + ">"

}

func (grid *RequirementGrid) Unmount() bool {

	if grid.Footer != nil {
		grid.Footer.Unmount()
	}

	return true

}

