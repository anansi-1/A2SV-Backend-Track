# API Documentation

### Get Tasks

  * Endpoint: `Get/tasks`
  * Description: retrieves a list of tasks.
  * Response:
    - Status Code: `200 Ok`
    - Body: A JSON response array of objects.
        ```json
           {
            "id": "string",
            "title": "string",
            "description": "string",
            "due_date": "time.Time",
            "status": "string"
            }
            ```

### Get Task By ID

  * Endpoint: `Get/task/:id`
  * Description: retrieves a task with the specific ID.
  * Response:
    - Status Code: `200 Ok`
    - Body: A JSON object representing the task.
           ```json
           {
            "id": "string",
            "title": "string",
            "description": "string",
            "due_date": "time.Time",
            "status": "string"
            }
            ```       
  * Possible Errors: 
    - `404 Not Found`: Task with specified Id dones not exist.

### Add Task 

  * Endpoint: `POST/tasks`
  * Description: Adds a new task.
  * Request Payload:
  ```json
  {
    "id": "string",
    "title": "string",
    "description": "string",
    "due_date": "time.Time",
    "status": "string"
  }
  ```
  * Response:
    - Status Code: `201 Created`
    - Body: JSON message `"Task created"`
  * Possible Errors: 
    - `400 Bad Request`: Invalid request payload.

### Remove task

  * Endpoint: `DELETE/tasks/:id`
  * Description: Deletes the task with the specified ID.
  * Request Parameters:
  `id`: The ID of the task to delete.
  * Response:
    - Status Code: `200 OK`
    - Body: JSON message `"Task deleted"`
  * Possible Errors: 
    - `404 Not Found` Task with the specified ID does not exist.

### Update task

  * Endpoint: `PUT/tasks/:id`
  * Description: Updates the task with the specificed ID.
  * Request Parameters:
  `id`: The ID of the task to be update.
  * Request Payload:
  ```json
  {
    "id": "string", 
    "title": "string", 
    "description": "string", 
    "due_date": "time.Time", 
    "status": "string"   }
  ```
  * Response:
    - Status Code: `200 OK`
    - Body: JSON message `"Task updated"`
  * Possible Errors: 
    - `404 Not Found` Invalid request payload
    - `400 Bad Request`: Task with the specified ID does not exist.