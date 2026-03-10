package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"

	"task-api/internal/domain"
	"task-api/internal/httpapi/gen"
	"task-api/internal/service"
)

type TaskHandler struct {
	svc *service.TaskService
}

// TaskHandler - это структура, которая содержит ссылку на сервис TaskService.
// Она отвечает за обработку HTTP-запросов, связанных с задачами.
// Внутри TaskHandler определены методы для создания задачи, получения списка задач, обновления задачи и удаления задачи.
// Эти методы обрабатывают входящие HTTP-запросы, взаимодействуют с сервисом и возвращают соответствующие HTTP-ответы.

func NewTaskHandler(svc *service.TaskService) *TaskHandler {
	return &TaskHandler{svc: svc}
}

// NewTaskHandler - это функция-конструктор для TaskHandler.
// Она принимает ссылку на TaskService и возвращает новый экземпляр TaskHandler, который будет использовать
// этот сервис для обработки запросов.

type createTaskRequest struct {
	Task   string `json:"task"`
	IsDone bool   `json:"is_done"`
}

// createTaskRequest - это структура, которая представляет собой тело запроса для создания новой задачи.

type patchTaskRequest struct {
	Task   *string `json:"task"`
	IsDone *bool   `json:"is_done"`
}

// patchTaskRequest - это структура, которая представляет собой тело запроса для обновления существующей задачи.

func (h *TaskHandler) Tasks(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		h.createTask(w, r)
	case http.MethodGet:
		h.listTasks(w)
	default:
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
	}
}

// Tasks - это метод, который обрабатывает HTTP-запросы на пути "/tasks".
// В зависимости от HTTP-метода, он вызывает соответствующий метод для создания новой задачи (POST) или получения списка задач (GET).
// Если HTTP-метод не поддерживается, возвращается ошибка "method not allowed".

func (h *TaskHandler) TaskByID(w http.ResponseWriter, r *http.Request) {
	id, err := parseIDFromPath(r.URL.Path, "/tasks/")
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid id")
		return
	}
	uid := uint(id)

	switch r.Method {
	case http.MethodPatch:
		h.patchTask(w, r, uid)
	case http.MethodDelete:
		h.deleteTask(w, uid)
	default:
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
	}
}

// TaskByID - это метод, который обрабатывает HTTP-запросы на пути "/tasks/{id}".
// Он извлекает идентификатор задачи из URL и в зависимости от HTTP-метода вызывает соответствующий метод
//  для обновления задачи (PATCH или PUT) или удаления задачи (DELETE).
// Если HTTP-метод не поддерживается, возвращается ошибка "method not allowed".

func (h *TaskHandler) createTask(w http.ResponseWriter, r *http.Request) {
	var req gen.CreateTaskRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "bad json")
		return
	}

	t, err := h.svc.Create(req.Task, req.IsDone)
	if err != nil {
		if errors.Is(err, service.ErrInvalidInput) {
			writeError(w, http.StatusBadRequest, "task is required")
			return
		}
		writeError(w, http.StatusInternalServerError, "internal error")
		return
	}

	writeJSON(w, http.StatusCreated, toCreateTaskResponse(*t))
}

// createTask - это метод, который обрабатывает HTTP-запросы для создания новой задачи.
// Он декодирует JSON-тело запроса в структуру createTaskRequest, вызывает метод Create сервиса для создания новой задачи и
//  возвращает созданную задачу в формате JSON.
// Если JSON некорректный или входные данные недопустимые, возвращается соответствующая ошибка.

func (h *TaskHandler) listTasks(w http.ResponseWriter) {
	tasks, err := h.svc.List()
	if err != nil {
		writeError(w, http.StatusInternalServerError, "internal error")
		return
	}
	writeJSON(w, http.StatusOK, toListTasksResponse(tasks))
}

// listTasks - это метод, который обрабатывает HTTP-запросы для получения списка всех задач.
// Он вызывает метод List сервиса для получения всех задач и возвращает их в формате JSON.
// Если возникает ошибка при получении задач, возвращается ошибка "internal error".

func (h *TaskHandler) patchTask(w http.ResponseWriter, r *http.Request, id uint) {
	var req gen.PatchTaskRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "bad json")
		return
	}

	t, err := h.svc.Patch(id, req.Task, req.IsDone)
	if err != nil {
		if errors.Is(err, service.ErrNotFound) {
			writeError(w, http.StatusNotFound, "not found")
			return
		}
		if errors.Is(err, service.ErrInvalidInput) {
			writeError(w, http.StatusBadRequest, "invalid input")
			return
		}
		writeError(w, http.StatusInternalServerError, "internal error")
		return
	}

	writeJSON(w, http.StatusOK, toPatchTaskResponse(*t))
}

// patchTask - это метод, который обрабатывает HTTP-запросы для обновления существующей задачи.
// Он декодирует JSON-тело запроса в структуру patchTaskRequest, вызывает метод Patch сервиса для обновления задачи и
//  возвращает обновленную задачу в формате JSON.
// Если JSON некорректный, задача не найдена или входные данные недопустимые, возвращается соответствующая ошибка.

func (h *TaskHandler) deleteTask(w http.ResponseWriter, id uint) {
	if err := h.svc.Delete(id); err != nil {
		if errors.Is(err, service.ErrNotFound) {
			writeError(w, http.StatusNotFound, "not found")
			return
		}
		writeError(w, http.StatusInternalServerError, "internal error")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// deleteTask - это метод, который обрабатывает HTTP-запросы для удаления существующей задачи.
// Он вызывает метод Delete сервиса для удаления задачи по идентификатору.
// Если задача не найдена, возвращается ошибка "not found".
// Если возникает другая ошибка, возвращается ошибка "internal error".
// Если удаление прошло успешно, возвращается статус "204 No Content".

func parseIDFromPath(path, prefix string) (int, error) {
	if !strings.HasPrefix(path, prefix) {
		return 0, errors.New("bad path")
	}
	raw := strings.TrimPrefix(path, prefix)
	raw = strings.Trim(raw, "/")
	if raw == "" {
		return 0, errors.New("empty id")
	}
	return strconv.Atoi(raw)
}

// parseIDFromPath - это вспомогательная функция, которая извлекает идентификатор из URL-пути.
// Она проверяет, что путь начинается с указанного префикса, удаляет префикс и любые ведущие или завершающие слэши,
//  а затем пытается преобразовать оставшуюся строку в целое число.
// Если путь не соответствует ожидаемому формату или идентификатор не является числом, возвращается ошибка.

func toCreateTaskResponse(t domain.Task) gen.CreateTaskResponse {
	return gen.CreateTaskResponse{
		Id:     int64(t.ID),
		Task:   t.Task,
		IsDone: t.IsDone,
	}
}

// toCreateTaskResponse - это функция, которая преобразует структуру domain.Task в структуру gen.CreateTaskResponse.
// Она используется для формирования ответа на запрос создания новой задачи, возвращая только необходимые поля Id, Task и IsDone.

func toPatchTaskResponse(t domain.Task) gen.PatchTaskResponse {
	return gen.PatchTaskResponse{
		Id:     int64(t.ID),
		Task:   t.Task,
		IsDone: t.IsDone,
	}
}

// toPatchTaskResponse - это функция, которая преобразует структуру domain.Task в структуру gen.PatchTaskResponse.
// Она используется для формирования ответа на запрос обновления существующей задачи, возвращая только необходимые поля Id, Task и IsDone.

func toListTasksResponse(tasks []domain.Task) gen.ListTasksResponse {
	out := make(gen.ListTasksResponse, 0, len(tasks))
	for _, t := range tasks {
		out = append(out, gen.Task{
			Id:     int64(t.ID),
			Task:   t.Task,
			IsDone: t.IsDone,
		})
	}
	return out
}

// toListTasksResponse - это функция, которая преобразует срез структур domain.Task в срез структур gen.Task.
// Она используется для формирования ответа на запрос получения списка задач, возвращая только необходимые поля Id, Task и IsDone для каждой задачи.

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

// writeJSON - это вспомогательная функция, которая записывает данные в формате JSON в HTTP-ответ.
// Она устанавливает заголовок "Content-Type" в "application/json", устанавливает статус ответа и кодирует переданное значение в JSON.
// Если кодирование не удалось, ошибка игнорируется, так как это вспомогательная функция для упрощения записи JSON-ответов.

func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, gen.ErrorResponse{Message: message})
}

// writeJSON - это вспомогательная функция, которая записывает данные в формате JSON в HTTP-ответ.
// Она устанавливает заголовок "Content-Type" в "application/json", устанавливает статус ответа и кодирует переданное значение в JSON.
// Если кодирование не удалось, ошибка игнорируется, так как это вспомогательная функция для упрощения записи JSON-ответов.
