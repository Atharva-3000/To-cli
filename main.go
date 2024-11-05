package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"
//  one external package is used for the prompt ui
	"github.com/manifoldco/promptui"
)

type Task struct {
	ID        int
	Name      string
	Status    string
	CreatedAt string
	UpdatedAt string
}

type TaskList struct {
	tasks []Task
}

func (t *TaskList) loadAllTasks() error {
	file, err := os.ReadFile("./tasks.json")
	if err != nil {
		if os.IsNotExist(err) {
			fmt.Println("No tasks yet 🫡!")
			return nil
		}
		return err
	}
	if len(file) == 0 {
		fmt.Println("No tasks yet 🫡!")
		return nil
	}
	return json.Unmarshal(file, &t.tasks)
}

func (t *TaskList) saveTasks() error {
	file, err := json.MarshalIndent(t.tasks, "", " ")
	if err != nil {
		return err
	}
	return os.WriteFile("./tasks.json", file, 0644)
}

func (t *TaskList) ensureTasksFileExists() error {
	if _, err := os.Stat("./tasks.json"); os.IsNotExist(err) {
		emptyTasks := []Task{}
		file, err := json.MarshalIndent(emptyTasks, "", " ")
		if err != nil {
			return err
		}
		return os.WriteFile("./tasks.json", file, 0644)
	}
	return nil
}

func (t *TaskList) addTask(name string) {
	if strings.TrimSpace(name) == "" {
		fmt.Println("Task name cannot be empty!")
		return
	}

	currentTime := time.Now().Format("15:04 02, January, 2006")
	newTask := Task{
		ID:        len(t.tasks) + 1,
		Name:      name,
		Status:    "todo",
		CreatedAt: currentTime,
		UpdatedAt: currentTime,
	}

	t.tasks = append(t.tasks, newTask)
	fmt.Println("Task added:", newTask.Name)
	t.saveTasks()
}

func (t *TaskList) updateTask() {
	if len(t.tasks) == 0 {
		fmt.Println("No tasks to update.")
		return
	}

	taskNames := []string{}
	for _, task := range t.tasks {
		taskNames = append(taskNames, task.Name)
	}

	prompt := promptui.Select{
		Label: "Select a task to update",
		Items: taskNames,
	}

	_, selectedName, err := prompt.Run()
	if err != nil {
		fmt.Println("Prompt failed:", err)
		return
	}

	for i, task := range t.tasks {
		if task.Name == selectedName {
			statusPrompt := promptui.Select{
				Label: "Select new status",
				Items: []string{"todo", "in-progress", "done"},
			}

			_, status, err := statusPrompt.Run()
			if err != nil {
				fmt.Println("Prompt failed:", err)
				return
			}

			t.tasks[i].Status = status
			t.tasks[i].UpdatedAt = time.Now().Format("15:04 02, January, 2006")
			fmt.Println("Task updated:", selectedName)
			t.saveTasks()
			return
		}
	}
}

func (t *TaskList) deleteTask() {
	if len(t.tasks) == 0 {
		fmt.Println("No tasks to delete.")
		return
	}

	taskNames := []string{}
	for _, task := range t.tasks {
		taskNames = append(taskNames, task.Name)
	}

	prompt := promptui.Select{
		Label: "Select a task to delete",
		Items: taskNames,
	}

	_, selectedName, err := prompt.Run()
	if err != nil {
		fmt.Println("Prompt failed:", err)
		return
	}

	for i, task := range t.tasks {
		if task.Name == selectedName {
			t.tasks = append(t.tasks[:i], t.tasks[i+1:]...)
			fmt.Println("Task deleted:", selectedName)
			t.saveTasks()
			return
		}
	}
}

func (t *TaskList) listTasksWithFilter() {
	prompt := promptui.Select{
		Label: "Select task status to list",
		Items: []string{"All Tasks 📋", "Done Tasks ✅", "In-Progress Tasks 🔄", "Todo Tasks 📝"},
	}

	_, result, err := prompt.Run()
	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		return
	}

	var filter string
	switch result {
	case "All Tasks 📋":
		filter = ""
	case "Done Tasks ✅":
		filter = "done"
	case "In-Progress Tasks 🔄":
		filter = "in-progress"
	case "Todo Tasks 📝":
		filter = "todo"
	}

	t.displayTasks(filter)
}

func (t *TaskList) displayTasks(status string) {
	hasTasks := false
	fmt.Println("ID   | Task Name           | Status        | Created At            | Updated At            ")
	fmt.Println("-----|----------------------|---------------|-----------------------|-----------------------")
	for _, task := range t.tasks {
		if status == "" || task.Status == status {
			fmt.Printf("%-4d | %-20s | %-13s | %-21s | %-21s\n", task.ID, task.Name, task.Status, task.CreatedAt, task.UpdatedAt)
			hasTasks = true
		}
	}
	if !hasTasks {
		fmt.Println("No tasks found with the selected status.")
	}
	fmt.Println()
}

func main() {
	taskList := &TaskList{}

	err := taskList.ensureTasksFileExists()
	if err != nil {
		fmt.Println("Error ensuring tasks file exists:", err)
		return
	}

	err = taskList.loadAllTasks()
	if err != nil {
		fmt.Println("Error loading tasks:", err)
		return
	}

	reader := bufio.NewReader(os.Stdin)
	for {
		prompt := promptui.Select{
			Label: "Choose an option",
			Items: []string{"Add Task ➕", "Update Task 📝", "Delete Task 🗑️", "List Tasks 📋", "Exit 😺"},
		}

		_, result, err := prompt.Run()
		if err != nil {
			fmt.Printf("Prompt failed %v\n", err)
			return
		}

		switch result {
		case "Add Task ➕":
			fmt.Println("Enter TASK to be ADDED:")
			name, _ := reader.ReadString('\n')
			name = strings.TrimSpace(name)
			taskList.addTask(name)
		case "Update Task 📝":
			taskList.updateTask()
		case "Delete Task 🗑️":
			taskList.deleteTask()
		case "List Tasks 📋":
			taskList.listTasksWithFilter()
		case "Exit 😺":
			taskList.saveTasks()
			fmt.Println("Sayonara 🙇🏼‍♂️, Keep Building!")
			return
		}
	}
}
