CREATE TABLE tasks (
    id SERIAL PRIMARY KEY,
    task VARCHAR(255) NOT NULL,
    is_done BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMP DEFAULT NULL
);


// Эта миграция создает таблицу "tasks" в базе данных с полями id, task, is_done, created_at, updated_at и deleted_at.
// Поле id является первичным ключом и автоматически увеличивается при добавлении новых записей. 
// Поле task хранит текст задачи и не может быть пустым. Поле is_done указывает, выполнена ли задача, 
// и по умолчанию имеет значение false. Поля created_at и updated_at автоматически устанавливаются на текущую дату и 
// время при создании и обновлении записи соответственно. Поле deleted_at используется для логического удаления задач, 
// позволяя сохранять информацию о том, когда задача была удалена, без фактического удаления записи из базы данных.