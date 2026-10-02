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

// Requirements is the standalone controller of the Requirements grid view. It
// mirrors the legacy `_old_frontend/ui/grids/Requirements.mjs`, owns the fetched
// reports and the collapse/filter state, and drives the `RequirementGrid`
// component.
type Requirements struct {
	Main    *app.Main
	Grid    *app_components.RequirementGrid
	Reports []*app_types.Requirement

	Collapsed bool
	Filter    string

	running    bool
	generation int
}

func NewRequirements(main *app.Main, view interfaces.View) *Requirements {

	var grid *app_components.RequirementGrid = nil

	if view != nil {

		if component := view.Query("section > requirement-grid"); component != nil {
			grid, _ = components.UnwrapComponent[*app_components.RequirementGrid](component)
		}

	}

	return &Requirements{
		Main:    main,
		Grid:    grid,
		Reports: make([]*app_types.Requirement, 0),
	}

}

func (requirements *Requirements) Name() string {
	return "requirements"
}

func (requirements *Requirements) Enter() bool {

	if requirements.Grid != nil && requirements.Grid.Footer != nil {
		requirements.attach(requirements.Grid.Footer.Component)
	}

	requirements.Render()

	requirements.running = true
	requirements.generation = requirements.generation + 1

	go requirements.loop(requirements.generation)

	return true

}

func (requirements *Requirements) Leave() bool {

	requirements.running = false

	if requirements.Grid != nil && requirements.Grid.Footer != nil {
		requirements.Grid.Footer.Component.RemoveEventListener("toggle-collapse", nil)
		requirements.Grid.Footer.Component.RemoveEventListener("change-filter", nil)
	}

	return true

}

func (requirements *Requirements) Update() {

	requirements.fetchReports()
	requirements.Render()

}

func (requirements *Requirements) Render() {

	if requirements.Grid == nil {
		return
	}

	requirements.Grid.SetRequirements(requirements.visible(), requirements.Collapsed)

}

func (requirements *Requirements) loop(generation int) {

	time.Sleep(500 * time.Millisecond)

	for requirements.running == true && requirements.generation == generation {

		requirements.Update()
		time.Sleep(5 * time.Second)

	}

}

func (requirements *Requirements) attach(component *components.Component) {

	if component == nil {
		return
	}

	component.AddEventListener("toggle-collapse", components.ToEventListener(func(event string, attributes map[string]any) {

		requirements.Collapsed = requirements.Collapsed != true

		requirements.Render()

	}, false))

	component.AddEventListener("change-filter", components.ToEventListener(func(event string, attributes map[string]any) {

		if value, ok := attributes["value"].(string); ok == true {
			requirements.Filter = value
			requirements.Render()
		}

	}, false))

}

func (requirements *Requirements) fetchReports() {

	response, err := requirements.Main.Client.Read("/api/session/requirements")

	if err != nil || response == nil {
		return
	}

	raw := make(map[string]map[string]app_types.Requirement)

	if json.Unmarshal(response.Body, &raw) != nil {
		return
	}

	reports := make([]*app_types.Requirement, 0)

	for _, symbols := range raw {

		for _, requirement := range symbols {

			report := requirement
			reports = append(reports, &report)

		}

	}

	requirements.Reports = reports

}

func (requirements *Requirements) visible() []*app_types.Requirement {

	query := strings.ToLower(strings.TrimSpace(requirements.Filter))

	if query == "" {
		return requirements.Reports
	}

	result := make([]*app_types.Requirement, 0)

	for _, report := range requirements.Reports {

		if report == nil {
			continue
		}

		haystack := strings.ToLower(report.PackagePath() + " " + report.File + " " + report.Symbol + " " + report.Declaration + " " + report.Behavior + " " + report.Type)

		if strings.Contains(haystack, query) == true {
			result = append(result, report)
		}

	}

	return result

}
