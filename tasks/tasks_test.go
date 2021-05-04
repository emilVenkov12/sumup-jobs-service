package tasks

import (
	"testing"
)

func TestCorrectTaskOrder(t *testing.T) {
	tsks := []Task{{Name: "task-1", Command: "cat /tmp/file1", Requires: []string{"task-2"}}, {Name: "task-2", Command: "touch /tmp/file1"}}
	orderedTasks, _ := TopologicallyOrderTasks(tsks)

	if orderedTasks[0].Name != "task-2" {
		t.Error("Incorect tasks order")
	}
}

func TestErrosOnCyclicDependencies(t *testing.T) {
	tsks := []Task{{Name: "task-1", Command: "cat /tmp/file1", Requires: []string{"task-2"}}, {Name: "task-2", Command: "touch /tmp/file1", Requires: []string{"task-1"}}}
	_, err := TopologicallyOrderTasks(tsks)

	if err == nil {
		t.Error("Dependecy cycle not detected")
	}
}

func TestMissingTaskRequired(t *testing.T) {
	tsks := []Task{{Name: "task-1", Command: "ls -lah"}, {Name: "task-2", Command: "pwd", Requires: []string{"task-3"}}}

	_, err := TopologicallyOrderTasks(tsks)

	t.Log(err)
	if err == nil {
		t.Error("Missing tasks required not detected")
	}
}
