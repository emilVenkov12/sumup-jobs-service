package tasks

import "errors"

type Task struct {
	Name     string   `json:"name"`
	Command  string   `json:"command"`
	Requires []string `json:"requires,omitempty"`
}

type Tasks struct {
	Tasks []Task `json:"tasks"`
}

// Order the topologically based on the requires array
func TopologicallyOrderTasks(tasks []Task) ([]Task, error) {
	if len(tasks) == 0 || len(tasks) == 1 {
		return tasks, nil
	}

	requires := make(map[string][]string)
	tasksMap := make(map[string]Task)

	for _, task := range tasks {
		tasksMap[task.Name] = task
		for _, requiredTask := range task.Requires {
			requires[requiredTask] = append(requires[requiredTask], task.Name)
		}
	}

	if len(requires) == 0 {
		return tasks, nil
	}

	dependsOnCounts := make(map[string]int, len(tasksMap))

	for u := range requires {
		for _, v := range requires[u] {
			dependsOnCounts[v] += 1
		}
	}

	q := make([]string, 0, len(tasksMap))

	for taskName := range tasksMap {
		if dependsOnCounts[taskName] == 0 {
			q = append(q, taskName)
		}
	}

	orderedTasks := make([]Task, 0, len(tasksMap))
	cnt := 0
	for len(q) > 0 {
		u := q[0]
		t := Task{Name: tasksMap[u].Name, Command: tasksMap[u].Command}
		orderedTasks = append(orderedTasks, t)
		q = q[1:]

		for _, v := range requires[u] {
			dependsOnCounts[v] -= 1
			if dependsOnCounts[v] == 0 {
				q = append(q, v)
			}
		}
		cnt++
	}
	if cnt != len(tasksMap) {
		return nil, errors.New("incorrect task requires - cyclic dependencies")
	}
	return orderedTasks, nil
}
