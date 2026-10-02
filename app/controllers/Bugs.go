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

// Bugs is the standalone controller of the Bugs grid view. It mirrors the
// legacy `_old_frontend/ui/grids/Bugs.mjs`, owns the fetched reports and the
// collapse/filter state, and drives the `BugGrid` component.
type Bugs struct {
	Main       *app.Main
	Grid       *app_components.BugGrid
	Reports    []*app_types.Bug
	Collapsed  bool
	Filter     string
	running    bool
	generation int
}

func NewBugs(main *app.Main, view interfaces.View) *Bugs {

	var grid *app_components.BugGrid = nil

	if view != nil {

		if component := view.Query("section > bug-grid"); component != nil {
			grid, _ = components.UnwrapComponent[*app_components.BugGrid](component)
		}

	}

	return &Bugs{
		Main:    main,
		Grid:    grid,
		Reports: make([]*app_types.Bug, 0),
	}

}

func (bugs *Bugs) Name() string {
	return "bugs"
}

func (bugs *Bugs) Enter() bool {

	if bugs.Grid != nil && bugs.Grid.Footer != nil {
		bugs.attach(bugs.Grid.Footer.Component)
	}

	bugs.Render()

	bugs.running = true
	bugs.generation = bugs.generation + 1

	go bugs.loop(bugs.generation)

	return true

}

func (bugs *Bugs) Leave() bool {

	bugs.running = false

	if bugs.Grid != nil && bugs.Grid.Footer != nil {
		bugs.Grid.Footer.Component.RemoveEventListener("toggle-collapse", nil)
		bugs.Grid.Footer.Component.RemoveEventListener("change-filter", nil)
	}

	return true

}

func (bugs *Bugs) Update() {

	bugs.fetchReports()
	bugs.Render()

}

func (bugs *Bugs) Render() {

	if bugs.Grid == nil {
		return
	}

	bugs.Grid.SetBugs(bugs.visible(), bugs.Collapsed)

}

func (bugs *Bugs) loop(generation int) {

	time.Sleep(500 * time.Millisecond)

	for bugs.running == true && bugs.generation == generation {

		bugs.Update()
		time.Sleep(5 * time.Second)

	}

}

func (bugs *Bugs) attach(component *components.Component) {

	if component == nil {
		return
	}

	component.AddEventListener("toggle-collapse", components.ToEventListener(func(event string, attributes map[string]any) {

		bugs.Collapsed = bugs.Collapsed != true

		bugs.Render()

	}, false))

	component.AddEventListener("change-filter", components.ToEventListener(func(event string, attributes map[string]any) {

		if value, ok := attributes["value"].(string); ok == true {
			bugs.Filter = value
			bugs.Render()
		}

	}, false))

}

func (bugs *Bugs) fetchReports() {

	response, err := bugs.Main.Client.Read("/api/session/bugs")

	if err != nil || response == nil {
		return
	}

	raw := make(map[string]map[string]app_types.Bug)

	if json.Unmarshal(response.Body, &raw) != nil {
		return
	}

	reports := make([]*app_types.Bug, 0)

	for _, symbols := range raw {

		for _, bug := range symbols {

			report := bug
			reports = append(reports, &report)

		}

	}

	bugs.Reports = reports

}

func (bugs *Bugs) visible() []*app_types.Bug {

	query := strings.ToLower(strings.TrimSpace(bugs.Filter))

	if query == "" {
		return bugs.Reports
	}

	result := make([]*app_types.Bug, 0)

	for _, report := range bugs.Reports {

		if report == nil {
			continue
		}

		haystack := strings.ToLower(report.PackagePath() + " " + report.File + " " + report.Symbol + " " + report.Description)

		if strings.Contains(haystack, query) == true {
			result = append(result, report)
		}

	}

	return result

}
