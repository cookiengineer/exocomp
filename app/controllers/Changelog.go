//go:build wasm

package controllers

import app_components "exocomp-app/components"
import app_types "exocomp-app/types"
import "github.com/cookiengineer/gooey/components"
import "github.com/cookiengineer/gooey/components/app"
import "github.com/cookiengineer/gooey/components/interfaces"
import "encoding/json"
import "strings"
import "time"

// Changelog is the standalone controller of the Changelog grid view. It mirrors
// the legacy `_old_frontend/ui/grids/Changelog.mjs`, flattens the nested
// `map[path]map[symbol][]ChangelogEntry` payload, owns the fetched reports and
// the collapse/filter state, and drives the `ChangelogGrid` component.
type Changelog struct {
	Main    *app.Main
	Grid    *app_components.ChangelogGrid
	Reports []*app_types.ChangelogEntry

	Collapsed bool
	Filter    string

	running    bool
	generation int

	toggleCollapseListener *components.EventListener
	changeFilterListener   *components.EventListener
}

func NewChangelog(main *app.Main, view interfaces.View) *Changelog {

	var grid *app_components.ChangelogGrid = nil

	if view != nil {

		if component := view.Query("section > changelog-grid"); component != nil {
			grid, _ = components.UnwrapComponent[*app_components.ChangelogGrid](component)
		}

	}

	return &Changelog{
		Main:    main,
		Grid:    grid,
		Reports: make([]*app_types.ChangelogEntry, 0),
	}

}

func (changelog *Changelog) Name() string {
	return "changelog"
}

func (changelog *Changelog) Enter() bool {

	if changelog.Grid != nil && changelog.Grid.Footer != nil {
		changelog.attach(changelog.Grid.Footer.Component)
	}

	changelog.Render()

	changelog.running = true
	changelog.generation = changelog.generation + 1

	go changelog.loop(changelog.generation)

	return true

}

func (changelog *Changelog) Leave() bool {

	changelog.running = false

	if changelog.Grid != nil && changelog.Grid.Footer != nil {

		if changelog.toggleCollapseListener != nil {
			changelog.Grid.Footer.Component.RemoveEventListener("toggle-collapse", changelog.toggleCollapseListener)
			changelog.toggleCollapseListener = nil
		}

		if changelog.changeFilterListener != nil {
			changelog.Grid.Footer.Component.RemoveEventListener("change-filter", changelog.changeFilterListener)
			changelog.changeFilterListener = nil
		}

	}

	return true

}

func (changelog *Changelog) Update() {

	changelog.fetchReports()
	changelog.Render()

}

func (changelog *Changelog) Render() {

	if changelog.Grid == nil {
		return
	}

	changelog.Grid.SetChangelog(changelog.visible(), changelog.Collapsed)

}

func (changelog *Changelog) loop(generation int) {

	time.Sleep(500 * time.Millisecond)

	for changelog.running == true && changelog.generation == generation {

		changelog.Update()
		time.Sleep(5 * time.Second)

	}

}

func (changelog *Changelog) attach(component *components.Component) {

	if component == nil {
		return
	}

	changelog.toggleCollapseListener = components.ToEventListener(func(event string, attributes map[string]any) {

		changelog.Collapsed = changelog.Collapsed != true

		changelog.Render()

	}, false)

	component.AddEventListener("toggle-collapse", changelog.toggleCollapseListener)

	changelog.changeFilterListener = components.ToEventListener(func(event string, attributes map[string]any) {

		if value, ok := attributes["value"].(string); ok == true {
			changelog.Filter = value
			changelog.Render()
		}

	}, false)

	component.AddEventListener("change-filter", changelog.changeFilterListener)

}

func (changelog *Changelog) fetchReports() {

	response, err := changelog.Main.Client.Read("/api/session/changelog")

	if err != nil || response == nil {
		return
	}

	raw := make(map[string]map[string][]app_types.ChangelogEntry)

	if json.Unmarshal(response.Body, &raw) != nil {
		return
	}

	reports := make([]*app_types.ChangelogEntry, 0)

	for _, symbols := range raw {

		for _, entries := range symbols {

			for _, entry := range entries {

				report := entry
				reports = append(reports, &report)

			}

		}

	}

	changelog.Reports = reports

}

func (changelog *Changelog) visible() []*app_types.ChangelogEntry {

	query := strings.ToLower(strings.TrimSpace(changelog.Filter))

	if query == "" {
		return changelog.Reports
	}

	result := make([]*app_types.ChangelogEntry, 0)

	for _, report := range changelog.Reports {

		if report == nil {
			continue
		}

		haystack := strings.ToLower(report.PackagePath() + " " + report.Type + " " + report.File + " " + report.Symbol + " " + report.Description)

		if strings.Contains(haystack, query) == true {
			result = append(result, report)
		}

	}

	return result

}
