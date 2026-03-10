
package gen


type CreateTaskRequest struct {
	IsDone bool   `json:"is_done"`
	Task   string `json:"task"`
}
// CreateTaskRequest - это структура, которая представляет собой тело запроса для создания новой задачи.
// Она содержит поля IsDone и Task, которые являются булевым значением и строкой соответственно.

type CreateTaskResponse = Task
// CreateTaskResponse - это тип, который представляет собой ответ на запрос создания новой задачи. 
// Он является псевдонимом для структуры Task, которая содержит поля Id, IsDone и Task.

type ErrorResponse struct {
	Message string `json:"message"`
}
// ErrorResponse - это структура, которая представляет собой тело ответа в случае ошибки.
// Она содержит поле Message, которое содержит текст сообщения об ошибке.

type ListTasksResponse = []Task
// ListTasksResponse - это тип, который представляет собой ответ на запрос получения списка задач.
// Он является псевдонимом для среза структур Task, каждая из которых содержит поля Id, IsDone и Task.

type PatchTaskRequest struct {
	IsDone *bool   `json:"is_done,omitempty"`
	Task   *string `json:"task,omitempty"`
}
// PatchTaskRequest - это структура, которая представляет собой тело запроса для обновления существующей задачи.
// Она содержит поля IsDone и Task, которые являются указателями на булевое значение и строку соответственно. 
// Это позволяет отличать случаи, когда поле не было указано в запросе (nil) от случаев, когда поле было указано с нулевым значением (false для IsDone и "" для Task).

type PatchTaskResponse = Task
// PatchTaskResponse - это тип, который представляет собой ответ на запрос обновления существующей задачи.
// Он является псевдонимом для структуры Task, которая содержит поля Id, IsDone и Task. 
// Ответ будет содержать обновленные данные задачи после успешного обновления.

type Task struct {
	Id     int64  `json:"id"`
	IsDone bool   `json:"is_done"`
	Task   string `json:"task"`
}
// Task - это структура, которая представляет собой задачу в ответах API.
// Она содержит поля Id, IsDone и Task, которые являются целым числом, булевым значением и строкой соответственно. 
// Это структура, которая используется в ответах API для представления задач.

type CreateTaskJSONRequestBody = CreateTaskRequest
// CreateTaskJSONRequestBody - это тип, который представляет собой тело запроса в формате JSON для создания новой задачи.
// Он является псевдонимом для структуры CreateTaskRequest, которая содержит поля IsDone и Task.

type PatchTaskJSONRequestBody = PatchTaskRequest
// PatchTaskJSONRequestBody - это тип, который представляет собой тело запроса в формате JSON для обновления существующей задачи.
// Он является псевдонимом для структуры PatchTaskRequest, которая содержит поля IsDone и Task.
