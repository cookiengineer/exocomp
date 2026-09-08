package tools

import "os"
import "path/filepath"
import "strings"
import "testing"

func createTestSkill(t *testing.T, playground string, name string, description string) {

	skill_folder := filepath.Join(playground, "skills", name)
	scripts_folder := filepath.Join(skill_folder, "scripts")

	err1 := os.MkdirAll(scripts_folder, 0755)

	if err1 != nil {
		t.Fatalf("Expected %v to be nil", err1)
	}

	skill_source := `{"name":"` + name + `","description":"` + description + `"}`

	err2 := os.WriteFile(filepath.Join(skill_folder, "SKILL.md"), []byte(skill_source), 0644)

	if err2 != nil {
		t.Fatalf("Expected %v to be nil", err2)
	}

	script_source := `package main

import "fmt"

func main() {
	fmt.Println("hello from skill")
}
`

	err3 := os.WriteFile(filepath.Join(scripts_folder, "main.go"), []byte(script_source), 0644)

	if err3 != nil {
		t.Fatalf("Expected %v to be nil", err3)
	}

}

func newTestSkills(t *testing.T, name string, description string) *Skills {

	playground, _ := os.MkdirTemp("/tmp", "exocomp-test-skills-*")
	sandbox       := playground

	createTestSkill(t, playground, name, description)

	return NewSkills([]string{"List", "Load", "Unload", "Execute"}, playground, sandbox, []string{"go"}, []string{})

}

func TestSkills_List(t *testing.T) {

	tool := newTestSkills(t, "scan", "Scan websites for links")

	result, err := tool.List()

	if err != nil {
		t.Errorf("Expected %v to be nil", err)
	}

	if strings.Contains(result, "skills.List: 1 skills available.") == false {
		t.Errorf("Expected one skill to be available, got %s", result)
	}

	if strings.Contains(result, "Skill: scan") == false {
		t.Errorf("Expected scan skill to be listed, got %s", result)
	}

	if strings.Contains(result, "Status: unloaded") == false {
		t.Errorf("Expected scan skill to be unloaded, got %s", result)
	}

	if strings.Contains(result, "Scripts: main.go") == false {
		t.Errorf("Expected main.go script to be listed, got %s", result)
	}

}

func TestSkills_Load_Unload(t *testing.T) {

	tool := newTestSkills(t, "scan", "Scan websites for links")

	result1, err1 := tool.Load("scan")

	if err1 != nil {
		t.Errorf("Expected %v to be nil", err1)
	}

	if strings.Contains(result1, "Skill \"scan\" got loaded") == false {
		t.Errorf("Expected skill to be loaded, got %s", result1)
	}

	result2, err2 := tool.Load("missing")

	if result2 != "" {
		t.Errorf("Expected %s to be empty", result2)
	}

	if err2 == nil {
		t.Errorf("Expected %v to be not nil", err2)
	} else if strings.Contains(err2.Error(), "doesn't exist") == false {
		t.Errorf("Expected missing skill error, got %v", err2)
	}

	result3, err3 := tool.Unload("scan")

	if err3 != nil {
		t.Errorf("Expected %v to be nil", err3)
	}

	if strings.Contains(result3, "got unloaded") == false {
		t.Errorf("Expected skill to be unloaded, got %s", result3)
	}

	result4, err4 := tool.Unload("scan")

	if result4 != "" {
		t.Errorf("Expected %s to be empty", result4)
	}

	if err4 == nil {
		t.Errorf("Expected %v to be not nil", err4)
	} else if strings.Contains(err4.Error(), "isn't loaded") == false {
		t.Errorf("Expected not loaded error, got %v", err4)
	}

}

func TestSkills_Execute(t *testing.T) {

	tool := newTestSkills(t, "scan", "Scan websites for links")

	result1, err1 := tool.Execute("scan", "main.go", []string{})

	if result1 != "" {
		t.Errorf("Expected %s to be empty", result1)
	}

	if err1 == nil {
		t.Errorf("Expected %v to be not nil", err1)
	} else if strings.Contains(err1.Error(), "isn't loaded") == false {
		t.Errorf("Expected not loaded error, got %v", err1)
	}

	_, err2 := tool.Load("scan")

	if err2 != nil {
		t.Errorf("Expected %v to be nil", err2)
	}

	result3, err3 := tool.Execute("scan", "missing.go", []string{})

	if result3 != "" {
		t.Errorf("Expected %s to be empty", result3)
	}

	if err3 == nil {
		t.Errorf("Expected %v to be not nil", err3)
	} else if strings.Contains(err3.Error(), "has no runtime") == false {
		t.Errorf("Expected missing runtime error, got %v", err3)
	}

	result4, err4 := tool.Execute("scan", "main.go", []string{})

	if err4 != nil {
		t.Errorf("Expected %v to be nil", err4)
	}

	if strings.Contains(result4, "hello from skill") == false {
		t.Errorf("Expected script output, got %s", result4)
	}

}
