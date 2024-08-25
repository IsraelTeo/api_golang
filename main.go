package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type task struct {
	ID      int    `json:"ID"`
	Name    string `json:"Name"`
	Content string `json:"Content"`
}

type allTask []task

var tasks = allTask{
	{
		ID:      1,
		Name:    "Task One",
		Content: "Some Content",
	},
}

func getTask(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	taskID, err := strconv.Atoi(vars["id"])
	if err != nil {
		fmt.Fprintf(w, "Invalid ID")
		return
	}
	for _, task := range tasks {
		if task.ID == taskID {
			w.Header().Set("Content-type", "aplication/json")
			json.NewEncoder(w).Encode(task)
		}
	}
}

func getAllTasks(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(tasks)
}

func createTask(w http.ResponseWriter, r *http.Request) {
	var newTask task
	reqBody, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Insert a valid task", http.StatusBadRequest)
		return
	}
	json.Unmarshal(reqBody, &newTask)
	newTask.ID = len(tasks) + 1
	tasks = append(tasks, newTask)
	w.Header().Set("Content-type", "aplication/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(tasks)
}

func updateTask(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	var updatedTask task

	taskID, err := strconv.Atoi(vars["id"])
	if err != nil {
		fmt.Fprintf(w, "Invalid ID")
	}

	reqBody, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Fprintf(w, "Please enter valid data")
	}
	json.Unmarshal(reqBody, &updatedTask)

	for i, task := range tasks {
		if task.ID == taskID {
			tasks = append(tasks[:i], tasks[i+1:]...)
			updatedTask.ID = taskID
			tasks = append(tasks, updatedTask)
			fmt.Fprintf(w, "The task with ID %v has updated successfuslly", taskID)
		}
	}
}

func deleteTask(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	taskID, err := strconv.Atoi(vars["id"])
	if err != nil {
		fmt.Fprintf(w, "Invalid ID")
		return
	}

	for i, task := range tasks {
		if task.ID == taskID {
			tasks = append(tasks[:i], tasks[i+1:]...)
			fmt.Fprintf(w, "The task with ID %v has been moved succesfully", taskID)
		}
	}
}

func main() {
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/tasks", getAllTasks).Methods("GET")
	router.HandleFunc("/task", createTask).Methods("POST")
	router.HandleFunc("/task/{id}", getTask).Methods("GET")
	router.HandleFunc("/task/{id}", deleteTask).Methods("DELETE")
	router.HandleFunc("/task/{id}", updateTask).Methods("PUT")

	log.Fatal(http.ListenAndServe(":8080", router))
}
