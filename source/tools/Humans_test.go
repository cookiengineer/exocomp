package tools

import "strings"
import "testing"
import "time"

func waitForQuestion(t *testing.T, tool *Humans, timeout time.Duration) string {

	deadline := time.Now().Add(timeout)

	for time.Now().Before(deadline) {

		tool.mutex.Lock()

		for id, _ := range tool.contents {
			tool.mutex.Unlock()
			return id
		}

		tool.mutex.Unlock()

		time.Sleep(5 * time.Millisecond)

	}

	t.Fatalf("Expected a question to be asked within %s", timeout)
	return ""

}

func TestHumans_Ask_Await_Blocks_Until_Answer(t *testing.T) {

	tool := NewHumans()

	result_ch := make(chan string, 1)
	err_ch    := make(chan error, 1)

	go func() {
		result, err := tool.Ask("What is your name?")
		result_ch <- result
		err_ch    <- err
	}()

	id := waitForQuestion(t, tool, 1*time.Second)

	select {
	case result := <-result_ch:
		t.Fatalf("Expected Ask to block until answered, but it returned early: %s", result)
	case <-time.After(50 * time.Millisecond):
		// still blocked, expected
	}

	err0 := tool.Answer(id, "Alice")

	if err0 != nil {
		t.Fatalf("Expected %v to be nil", err0)
	}

	select {
	case err := <-err_ch:
		if err != nil {
			t.Fatalf("Expected %v to be nil", err)
		}
	case <-time.After(1 * time.Second):
		t.Fatalf("Expected Ask to return after the answer")
	}

	result := <-result_ch

	if !strings.Contains(result, "Alice") {
		t.Errorf("Expected the answer to contain \"Alice\", got %s", result)
	}

	question, ok := tool.contents[id]

	if ok == false {
		t.Fatalf("Expected question %s to exist", id)
	}

	if question.Answer != "Alice" {
		t.Errorf("Expected answer to be %q, got %q", "Alice", question.Answer)
	}

}

func TestHumans_Choice_Persists_Options_And_Multiple(t *testing.T) {

	tool := NewHumans()

	result_ch := make(chan string, 1)
	err_ch    := make(chan error, 1)

	go func() {
		result, err := tool.Choice("Which database?", []string{"PostgreSQL", "SQLite", "MongoDB"}, true)
		result_ch <- result
		err_ch    <- err
	}()

	id := waitForQuestion(t, tool, 1*time.Second)

	tool.mutex.Lock()
	question, ok := tool.contents[id]
	tool.mutex.Unlock()

	if ok == false {
		t.Fatalf("Expected question %s to exist", id)
	}

	if question.Type != "Choice" {
		t.Errorf("Expected type to be %q, got %q", "Choice", question.Type)
	}

	if len(question.Options) != 3 {
		t.Errorf("Expected 3 options, got %d", len(question.Options))
	}

	if question.Multiple != true {
		t.Errorf("Expected multiple to be true, got %v", question.Multiple)
	}

	err0 := tool.Answer(id, "PostgreSQL\nSQLite")

	if err0 != nil {
		t.Fatalf("Expected %v to be nil", err0)
	}

	select {
	case err := <-err_ch:
		if err != nil {
			t.Fatalf("Expected %v to be nil", err)
		}
	case <-time.After(1 * time.Second):
		t.Fatalf("Expected Choice to return after the answer")
	}

	result := <-result_ch

	if !strings.Contains(result, "PostgreSQL") {
		t.Errorf("Expected the answer to contain \"PostgreSQL\", got %s", result)
	}

}

func TestHumans_Await_Timeout(t *testing.T) {

	tool := NewHumans()
	tool.Timeout = 50 * time.Millisecond

	start := time.Now()
	_, err0 := tool.Ask("Are you still there?")
	elapsed := time.Since(start)

	if err0 == nil {
		t.Errorf("Expected a non-nil error for an unanswered question")
	}

	if !strings.Contains(err0.Error(), "never answered") {
		t.Errorf("Expected never answered error, got %v", err0)
	}

	if elapsed > 5*time.Second {
		t.Errorf("Expected Ask to return within ~5s, took %s", elapsed)
	}

}

func TestHumans_Await_NeverAsked(t *testing.T) {

	tool := NewHumans()

	_, err0 := tool.Await("question-999")

	if err0 == nil {
		t.Errorf("Expected a non-nil error for a question that was never asked")
	}

	if !strings.Contains(err0.Error(), "never asked") {
		t.Errorf("Expected never asked error, got %v", err0)
	}

}

func TestHumans_Answer_InvalidID(t *testing.T) {

	tool := NewHumans()

	err0 := tool.Answer("question-999", "answer")

	if err0 == nil {
		t.Errorf("Expected a non-nil error for an invalid question id")
	}

	if !strings.Contains(err0.Error(), "Invalid question id") {
		t.Errorf("Expected invalid question id error, got %v", err0)
	}

}

func TestHumans_Call_ArgumentValidation(t *testing.T) {

	tool := NewHumans()

	_, err0 := tool.Call("Ask", map[string]interface{}{})

	if err0 == nil || !strings.Contains(err0.Error(), "question") {
		t.Errorf("Expected missing question error, got %v", err0)
	}

	_, err1 := tool.Call("Choose", map[string]interface{}{
		"question": "Pick one",
	})

	if err1 == nil || !strings.Contains(err1.Error(), "options") {
		t.Errorf("Expected missing options error, got %v", err1)
	}

	_, err2 := tool.Call("Choose", map[string]interface{}{
		"question": "Pick one",
		"options":  []interface{}{"a", 42},
	})

	if err2 == nil || !strings.Contains(err2.Error(), "options[1]") {
		t.Errorf("Expected non-string option error, got %v", err2)
	}

	_, err3 := tool.Call("Invalid", map[string]interface{}{})

	if err3 == nil || !strings.Contains(err3.Error(), "Invalid method") {
		t.Errorf("Expected invalid method error, got %v", err3)
	}

}
