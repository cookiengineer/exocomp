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
import "time"

// ChangelogGrid renders the Changelog view. It owns
// `[]*app_types.ChangelogEntry`, groups them into package tiles and renders
// their dated list items. Its search/collapse controls live in the paired
// ChangelogFooter.
type ChangelogGrid struct {
	Component   *components.Component       `json:"component"`
	Collapsed   bool                        `json:"collapsed"`
	entries     []*app_types.ChangelogEntry `json:"-"`
	packages    *dom.Element                `json:"-"`
	Footer      *ChangelogFooter            `json:"-"`
	initialized bool                        `json:"-"`
}

func NewChangelogGrid() *ChangelogGrid {

	element := dom.GetDocument().CreateElement("changelog-grid")
	component := components.NewComponent(element)

	return &ChangelogGrid{
		Component: &component,
		entries:   make([]*app_types.ChangelogEntry, 0),
	}

}

func ToChangelogGrid(element *dom.Element) *ChangelogGrid {

	component := components.NewComponent(element)

	return &ChangelogGrid{
		Component: &component,
		entries:   make([]*app_types.ChangelogEntry, 0),
	}

}

func (grid *ChangelogGrid) SetChangelog(entries []*app_types.ChangelogEntry, collapsed bool) {

	grid.entries = entries
	grid.Collapsed = collapsed

	grid.Render()

}

func (grid *ChangelogGrid) Disable() bool {
	return false
}

func (grid *ChangelogGrid) Enable() bool {
	return false
}

func (grid *ChangelogGrid) SetScheduler(scheduler interfaces.Scheduler) {

	if grid.Component != nil {
		grid.Component.SetScheduler(scheduler)
	}

	if grid.Footer != nil {
		grid.Footer.SetScheduler(scheduler)
	}

}

func (grid *ChangelogGrid) placeholder() string {

	if grid.Component != nil && grid.Component.Element != nil {

		tmp := strings.TrimSpace(grid.Component.Element.GetAttribute("data-placeholder"))

		if tmp != "" {
			return tmp
		}

	}

	return "Search ..."

}

func (grid *ChangelogGrid) ensure() {

	if grid.initialized == true {
		return
	}

	if grid.Component == nil || grid.Component.Element == nil {
		return
	}

	grid.Component.Element.SetInnerHTML("<div class=\"packages\"></div>")

	grid.packages = grid.Component.Element.QuerySelector("div.packages")

	grid.Footer = NewChangelogFooter(grid.placeholder())
	grid.Component.Element.Append(grid.Footer.Component.Element)
	grid.Footer.Mount()

	grid.initialized = true

}

func (grid *ChangelogGrid) Mount() bool {

	if grid.Component == nil || grid.Component.Element == nil {
		return false
	}

	grid.ensure()
	grid.Render()

	return true

}

func (grid *ChangelogGrid) Render() *dom.Element {

	if grid.Component == nil || grid.Component.Element == nil {
		return nil
	}

	grid.ensure()

	if grid.packages != nil {

		grouped := make(map[string][]*app_types.ChangelogEntry)
		package_paths := make([]string, 0)

		for _, entry := range grid.entries {

			if entry == nil {
				continue
			}

			package_path := entry.PackagePath()

			if grid.Collapsed == true {
				package_path = utils_path.Collapse(package_path)
			}

			if _, ok := grouped[package_path]; ok == false {
				grouped[package_path] = make([]*app_types.ChangelogEntry, 0)
				package_paths = append(package_paths, package_path)
			}

			grouped[package_path] = append(grouped[package_path], entry)

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

			for _, entry := range grouped[package_path] {

				date_text := ""
				datetime := ""

				if entry.Date.IsZero() == false {
					date_text = entry.Date.Format("2006-01-02")
					datetime = entry.Date.Format(time.RFC3339)
				}

				marker := ""

				if entry.Type != "" {
					marker = "[" + entry.Type + "] "
				}

				label := marker + strings.TrimPrefix(entry.File, "./") + "#" + entry.Symbol + ": " + entry.Description

				html += "<li><time datetime=\"" + utils_fmt.SanitizeContent(datetime) + "\">" + utils_fmt.SanitizeContent(date_text) + "</time><span>" + utils_fmt.SanitizeContent(label) + "</span></li>"

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

func (grid *ChangelogGrid) Query(query string) interfaces.Component {

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

func (grid *ChangelogGrid) String() string {

	if grid.Component == nil || grid.Component.Element == nil {
		return ""
	}

	tag := strings.ToLower(grid.Component.Element.TagName)

	return "<" + tag + ">" + grid.Component.Element.InnerHTML + "</" + tag + ">"

}

func (grid *ChangelogGrid) Unmount() bool {

	if grid.Footer != nil {
		grid.Footer.Unmount()
	}

	return true

}
