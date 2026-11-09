package tools

import "encoding/json"
import "os"
import "os/exec"
import "path/filepath"
import "strings"
import "testing"

func newTestContainers(t *testing.T, methods []string, allowed_images []string) (*Containers, string) {

	t.Helper()

	playground, _ := os.MkdirTemp("/tmp", "exocomp-test-containers-*")
	sandbox := filepath.Join(playground, "containers")

	err0 := os.MkdirAll(sandbox, 0755)

	if err0 != nil {
		t.Fatalf("Expected %v to be nil", err0)
	}

	tool := NewContainers(methods, playground, sandbox, allowed_images)

	t.Setenv("EXOCOMP_PODMAN", buildFakePodman(t))
	t.Setenv("EXOCOMP_FAKE_PODMAN_ARGS", filepath.Join(playground, "arguments.json"))
	t.Setenv("EXOCOMP_FAKE_PODMAN_LABEL", tool.label())

	t.Cleanup(func() {

		if t.Failed() == true {
			t.Logf("Preserving folder %s for debugging.", playground)
		} else {
			os.RemoveAll(playground)
		}

	})

	return tool, sandbox

}

func readRecordedArguments(t *testing.T) []string {

	t.Helper()

	path := os.Getenv("EXOCOMP_FAKE_PODMAN_ARGS")
	data, err0 := os.ReadFile(path)

	arguments := make([]string, 0)

	if err0 == nil {

		err1 := json.Unmarshal(data, &arguments)

		if err1 != nil {
			t.Fatalf("Expected %v to be nil", err1)
		}

	} else {
		t.Fatalf("Expected %v to be nil", err0)
	}

	return arguments

}

func TestContainers_Call_ArgumentValidation(t *testing.T) {

	tool, _ := newTestContainers(t, []string{"Create", "Exec", "List", "Logs"}, []string{"alpine:latest"})

	if tool != nil {

		_, err1 := tool.Call("Start", map[string]interface{}{"name": "web"})

		if err1 == nil {
			t.Errorf("Expected %v to be not nil", err1)
		} else if strings.Contains(err1.Error(), "Method not allowed") == false {
			t.Errorf("Expected method not allowed error, got %v", err1)
		}

		_, err2 := tool.Call("Create", map[string]interface{}{})

		if err2 == nil {
			t.Errorf("Expected %v to be not nil", err2)
		} else if strings.Contains(err2.Error(), "Invalid parameters") == false {
			t.Errorf("Expected invalid parameters error, got %v", err2)
		}

		_, err3 := tool.Call("Create", map[string]interface{}{"image": 123, "name": "web"})

		if err3 == nil {
			t.Errorf("Expected %v to be not nil", err3)
		} else if strings.Contains(err3.Error(), "Invalid parameter \"image\"") == false {
			t.Errorf("Expected invalid image error, got %v", err3)
		}

		_, err4 := tool.Call("Create", map[string]interface{}{"image": "alpine:latest", "name": "web", "arguments": "nope"})

		if err4 == nil {
			t.Errorf("Expected %v to be not nil", err4)
		} else if strings.Contains(err4.Error(), "Invalid parameter \"arguments\"") == false {
			t.Errorf("Expected invalid arguments error, got %v", err4)
		}

		_, err5 := tool.Call("Exec", map[string]interface{}{"name": "web"})

		if err5 == nil {
			t.Errorf("Expected %v to be not nil", err5)
		} else if strings.Contains(err5.Error(), "Invalid parameter \"command\"") == false {
			t.Errorf("Expected invalid command error, got %v", err5)
		}

		_, err6 := tool.Call("Logs", map[string]interface{}{"name": "web", "tail": "ten"})

		if err6 == nil {
			t.Errorf("Expected %v to be not nil", err6)
		} else if strings.Contains(err6.Error(), "Invalid parameter \"tail\"") == false {
			t.Errorf("Expected invalid tail error, got %v", err6)
		}

	} else {
		t.Errorf("Expected tool to be not nil")
	}

}

func TestContainers_Create(t *testing.T) {

	tool, sandbox := newTestContainers(t, []string{"Create"}, []string{"alpine:latest"})

	if tool != nil {

		result1, err1 := tool.Create("alpine:latest", "web", "sh", []string{"-c", "sleep 60"}, "none")

		if err1 != nil {
			t.Errorf("Expected %v to be nil", err1)
		}

		if result1 != "containers.Create: Container \"web\" created from image \"alpine:latest\"." {
			t.Errorf("Expected container to be created, got %q", result1)
		}

		arguments1 := readRecordedArguments(t)
		line1 := strings.Join(arguments1, " ")

		if strings.Contains(line1, "--network none") == false {
			t.Errorf("Expected network mode in %q", line1)
		}

		if strings.Contains(line1, "--volume "+sandbox+":/workspace:Z") == false {
			t.Errorf("Expected sandbox volume in %q", line1)
		}

		if strings.Contains(line1, "--label io.exocomp.sandbox="+tool.label()) == false {
			t.Errorf("Expected sandbox label in %q", line1)
		}

		if strings.HasSuffix(line1, "alpine:latest sh -c sleep 60") == false {
			t.Errorf("Expected command suffix in %q", line1)
		}

		result2, err2 := tool.Create("alpine:latest", "worker", "", nil, "bridge")

		if err2 != nil {
			t.Errorf("Expected %v to be nil", err2)
		}

		if result2 != "containers.Create: Container \"worker\" created from image \"alpine:latest\"." {
			t.Errorf("Expected default container to be created, got %q", result2)
		}

		arguments2 := readRecordedArguments(t)
		line2 := strings.Join(arguments2, " ")

		if strings.HasSuffix(line2, "alpine:latest sleep infinity") == false {
			t.Errorf("Expected default keepalive command in %q", line2)
		}

	} else {
		t.Errorf("Expected tool to be not nil")
	}

}

func TestContainers_CreateValidation(t *testing.T) {

	tool, sandbox := newTestContainers(t, []string{"Create"}, []string{"alpine:latest"})

	if tool != nil {

		_, err1 := tool.Create("busybox:latest", "web", "", nil, "bridge")

		if err1 == nil {
			t.Errorf("Expected %v to be not nil", err1)
		} else if strings.Contains(err1.Error(), "Attempt to use unallowed image") == false {
			t.Errorf("Expected unallowed image error, got %v", err1)
		}

		_, err2 := tool.Create("alpine:latest", "-web", "", nil, "bridge")

		if err2 == nil {
			t.Errorf("Expected %v to be not nil", err2)
		} else if strings.Contains(err2.Error(), "Invalid container name") == false {
			t.Errorf("Expected invalid name error, got %v", err2)
		}

		_, err3 := tool.Create("alpine:latest", "web", "", nil, "host")

		if err3 == nil {
			t.Errorf("Expected %v to be not nil", err3)
		} else if strings.Contains(err3.Error(), "Invalid network") == false {
			t.Errorf("Expected invalid network error, got %v", err3)
		}

		_, err4 := tool.Create("alpine:latest", "web", "", []string{"infinity"}, "bridge")

		if err4 == nil {
			t.Errorf("Expected %v to be not nil", err4)
		} else if strings.Contains(err4.Error(), "requires \"command\"") == false {
			t.Errorf("Expected command requirement error, got %v", err4)
		}

		t.Setenv("EXOCOMP_FAKE_PODMAN", "run-image-missing")
		_, err5 := tool.Create("alpine:latest", "web", "", nil, "bridge")

		if err5 == nil {
			t.Errorf("Expected %v to be not nil", err5)
		} else if strings.Contains(err5.Error(), "Image \"alpine:latest\" does not exist") == false {
			t.Errorf("Expected missing image error, got %v", err5)
		}

		t.Setenv("EXOCOMP_FAKE_PODMAN", "run-fail")
		result6, err6 := tool.Create("alpine:latest", "web", "", nil, "bridge")

		if err6 == nil {
			t.Errorf("Expected %v to be not nil", err6)
		} else if strings.Contains(err6.Error(), "/home/") == true || strings.Contains(err6.Error(), sandbox) == true {
			t.Errorf("Expected sanitized error, got %v", err6)
		}

		if result6 != "" {
			t.Errorf("Expected %s to be empty", result6)
		}

	} else {
		t.Errorf("Expected tool to be not nil")
	}

}

func TestContainers_List(t *testing.T) {

	tool, _ := newTestContainers(t, []string{"List"}, []string{"alpine:latest"})

	if tool != nil {

		result1, err1 := tool.List()

		if err1 != nil {
			t.Errorf("Expected %v to be nil", err1)
		}

		if strings.Contains(result1, "containers.List: 1 containers available.") == false {
			t.Errorf("Expected one container, got %q", result1)
		}

		if strings.Contains(result1, "Name: web, Image: alpine:latest, Status: Up 2 minutes") == false {
			t.Errorf("Expected container details, got %q", result1)
		}

		arguments := readRecordedArguments(t)
		line := strings.Join(arguments, " ")

		if strings.Contains(line, "--filter label=io.exocomp.sandbox="+tool.label()) == false {
			t.Errorf("Expected sandbox label filter in %q", line)
		}

		t.Setenv("EXOCOMP_FAKE_PODMAN_EMPTY", "1")
		result2, err2 := tool.List()

		if err2 != nil {
			t.Errorf("Expected %v to be nil", err2)
		}

		if strings.Contains(result2, "containers.List: 0 containers available.") == false {
			t.Errorf("Expected no containers, got %q", result2)
		}

	} else {
		t.Errorf("Expected tool to be not nil")
	}

}

func TestContainers_Images(t *testing.T) {

	tool, _ := newTestContainers(t, []string{"Images"}, []string{"alpine:latest"})

	if tool != nil {

		result1, err1 := tool.Images()

		if err1 != nil {
			t.Errorf("Expected %v to be nil", err1)
		}

		if strings.Contains(result1, "containers.Images: 1 images available.") == false {
			t.Errorf("Expected one allowed image, got %q", result1)
		}

		if strings.Contains(result1, "Name: alpine:latest") == false {
			t.Errorf("Expected alpine image, got %q", result1)
		}

		if strings.Contains(result1, "busybox:latest") == true {
			t.Errorf("Expected busybox image to be filtered, got %q", result1)
		}

	} else {
		t.Errorf("Expected tool to be not nil")
	}

}

func TestContainers_Stat(t *testing.T) {

	tool, _ := newTestContainers(t, []string{"Stat"}, []string{"alpine:latest"})

	if tool != nil {

		result1, err1 := tool.Stat("web")

		if err1 != nil {
			t.Errorf("Expected %v to be nil", err1)
		}

		if strings.Contains(result1, "containers.Stat: web") == false {
			t.Errorf("Expected container name, got %q", result1)
		}

		if strings.Contains(result1, "ID: 0123456789ab") == false {
			t.Errorf("Expected short container ID, got %q", result1)
		}

		if strings.Contains(result1, "Status: running") == false {
			t.Errorf("Expected running status, got %q", result1)
		}

		t.Setenv("EXOCOMP_FAKE_PODMAN_FOREIGN", "1")
		_, err2 := tool.Stat("web")

		if err2 == nil {
			t.Errorf("Expected %v to be not nil", err2)
		} else if strings.Contains(err2.Error(), "does not exist") == false {
			t.Errorf("Expected foreign container rejection, got %v", err2)
		}

		t.Setenv("EXOCOMP_FAKE_PODMAN", "inspect-missing")
		_, err3 := tool.Stat("web")

		if err3 == nil {
			t.Errorf("Expected %v to be not nil", err3)
		} else if strings.Contains(err3.Error(), "/") == true {
			t.Errorf("Expected sanitized error, got %v", err3)
		}

	} else {
		t.Errorf("Expected tool to be not nil")
	}

}

func TestContainers_StartStopRemove(t *testing.T) {

	tool, _ := newTestContainers(t, []string{"Start", "Stop", "Remove"}, []string{"alpine:latest"})

	if tool != nil {

		result1, err1 := tool.Start("web")

		if err1 != nil {
			t.Errorf("Expected %v to be nil", err1)
		}

		if result1 != "containers.Start: Container \"web\" started." {
			t.Errorf("Expected start result, got %q", result1)
		}

		result2, err2 := tool.Stop("web")

		if err2 != nil {
			t.Errorf("Expected %v to be nil", err2)
		}

		if result2 != "containers.Stop: Container \"web\" stopped." {
			t.Errorf("Expected stop result, got %q", result2)
		}

		result3, err3 := tool.Remove("web")

		if err3 != nil {
			t.Errorf("Expected %v to be nil", err3)
		}

		if result3 != "containers.Remove: Container \"web\" removed." {
			t.Errorf("Expected remove result, got %q", result3)
		}

		arguments := readRecordedArguments(t)
		line := strings.Join(arguments, " ")

		if line != "rm -f web" {
			t.Errorf("Expected remove command, got %q", line)
		}

	} else {
		t.Errorf("Expected tool to be not nil")
	}

}

func TestContainers_Exec(t *testing.T) {

	tool, _ := newTestContainers(t, []string{"Exec"}, []string{"alpine:latest"})

	if tool != nil {

		result1, err1 := tool.Exec("web", "sh", []string{"-c", "id"}, "/workspace")

		if err1 != nil {
			t.Errorf("Expected %v to be nil", err1)
		}

		if strings.Contains(result1, "containers.Exec: podman exec web sh -c id") == false {
			t.Errorf("Expected exec command line, got %q", result1)
		}

		arguments := readRecordedArguments(t)
		line := strings.Join(arguments, " ")

		if line != "exec --workdir /workspace web sh -c id" {
			t.Errorf("Expected exec arguments, got %q", line)
		}

		_, err2 := tool.Exec("web", "id", []string{}, "workspace")

		if err2 == nil {
			t.Errorf("Expected %v to be not nil", err2)
		} else if strings.Contains(err2.Error(), "Invalid workdir") == false {
			t.Errorf("Expected invalid workdir error, got %v", err2)
		}

		t.Setenv("EXOCOMP_FAKE_PODMAN", "exec-fail")
		result3, err3 := tool.Exec("web", "id", []string{}, "/workspace")

		if err3 == nil {
			t.Errorf("Expected %v to be not nil", err3)
		} else if strings.Contains(err3.Error(), "Permission denied") == false {
			t.Errorf("Expected permission denied error, got %v", err3)
		}

		if result3 != "" {
			t.Errorf("Expected %s to be empty", result3)
		}

		t.Setenv("EXOCOMP_FAKE_PODMAN", "exec-notfound")
		_, err4 := tool.Exec("web", "missing", []string{}, "/workspace")

		if err4 == nil {
			t.Errorf("Expected %v to be not nil", err4)
		} else if strings.Contains(err4.Error(), "does not exist") == true {
			t.Errorf("Expected command failure error, got %v", err4)
		}

	} else {
		t.Errorf("Expected tool to be not nil")
	}

}

func TestContainers_Logs(t *testing.T) {

	tool, _ := newTestContainers(t, []string{"Logs"}, []string{"alpine:latest"})

	if tool != nil {

		result1, err1 := tool.Logs("web", 5)

		if err1 != nil {
			t.Errorf("Expected %v to be nil", err1)
		}

		if strings.Contains(result1, "containers.Logs: web") == false {
			t.Errorf("Expected logs header, got %q", result1)
		}

		if strings.Contains(result1, "line one") == false || strings.Contains(result1, "line two") == false {
			t.Errorf("Expected log lines, got %q", result1)
		}

		arguments := readRecordedArguments(t)
		line := strings.Join(arguments, " ")

		if line != "logs --tail 5 web" {
			t.Errorf("Expected logs arguments, got %q", line)
		}

		_, err2 := tool.Logs("bad/name", 5)

		if err2 == nil {
			t.Errorf("Expected %v to be not nil", err2)
		} else if strings.Contains(err2.Error(), "Invalid container name") == false {
			t.Errorf("Expected invalid name error, got %v", err2)
		}

	} else {
		t.Errorf("Expected tool to be not nil")
	}

}

func TestContainers_Schemas(t *testing.T) {

	methods := []string{"Create", "Exec", "Images", "List", "Logs", "Remove", "Start", "Stat", "Stop"}
	tool := NewContainers(methods, "", "", []string{})

	schemas := tool.Schemas()

	if len(schemas) == len(methods) {

		names := make([]string, 0)

		for _, schema := range schemas {
			names = append(names, schema.Function.Name)
		}

		for _, method := range methods {

			if strings.Contains(strings.Join(names, " "), "containers."+method) == false {
				t.Errorf("Expected schema for containers.%s", method)
			}

		}

	} else {
		t.Errorf("Expected %d schemas to be %d", len(schemas), len(methods))
	}

}

func TestContainers_PodmanMissing(t *testing.T) {

	tool, playground := newTestContainers(t, []string{"List"}, []string{"alpine:latest"})

	if tool != nil {

		t.Setenv("EXOCOMP_PODMAN", filepath.Join(playground, "missing-podman"))
		_, err1 := tool.List()

		if err1 == nil {
			t.Errorf("Expected %v to be not nil", err1)
		} else if strings.Contains(err1.Error(), "Program \"podman\" doesn't exist.") == false {
			t.Errorf("Expected missing podman error, got %v", err1)
		}

	} else {
		t.Errorf("Expected tool to be not nil")
	}

}

func TestContainers_Integration(t *testing.T) {

	if os.Getenv("EXOCOMP_PODMAN_TEST") != "1" {
		t.Skip("set EXOCOMP_PODMAN_TEST=1 to run against a real podman installation")
	}

	_, err0 := exec.LookPath("podman")

	if err0 != nil {
		t.Skip("podman is not installed")
	}

	playground, _ := os.MkdirTemp("/tmp", "exocomp-test-containers-integration-*")
	sandbox := filepath.Join(playground, "containers")
	err1 := os.MkdirAll(sandbox, 0755)

	if err1 != nil {
		t.Fatalf("Expected %v to be nil", err1)
	}

	t.Setenv("EXOCOMP_PODMAN", "")

	image := "docker.io/library/alpine:latest"
	name := "exocomp-integration"
	tool := NewContainers([]string{"Create", "Exec", "List", "Remove", "Stat"}, playground, sandbox, []string{image})

	if tool != nil {

		t.Cleanup(func() {

			if t.Failed() == true {
				t.Logf("Preserving folder %s for debugging.", playground)
				tool.Remove(name)
			} else {
				tool.Remove(name)
				os.RemoveAll(playground)
			}

		})

		result1, err2 := tool.Create(image, name, "", nil, "none")

		if err2 != nil {

			t.Errorf("Expected %v to be nil", err2)

		} else {

			if strings.Contains(result1, "created from image") == false {
				t.Errorf("Expected create result, got %q", result1)
			}

			result2, err3 := tool.Exec(name, "sh", []string{"-c", "echo exocomp-integration"}, "/workspace")

			if err3 != nil {
				t.Errorf("Expected %v to be nil", err3)
			} else if strings.Contains(result2, "exocomp-integration") == false {
				t.Errorf("Expected exec output, got %q", result2)
			}

			result3, err4 := tool.List()

			if err4 != nil {
				t.Errorf("Expected %v to be nil", err4)
			} else if strings.Contains(result3, "Name: "+name) == false {
				t.Errorf("Expected container in list, got %q", result3)
			}

		}

	} else {
		t.Errorf("Expected tool to be not nil")
	}

}
