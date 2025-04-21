package api

import (
	"gastrobar-backend/cmd/app"
	"gastrobar-backend/internal/models"
	"gastrobar-backend/pkg/middleware"

	"github.com/gorilla/mux"
)

func SetupRoutes(app *app.App) *mux.Router {
    router := mux.NewRouter()

    // Rutas públicas (sin middleware)

    // POST /login: Inicia sesión usando username y password
    router.HandleFunc("/login", app.AuthHandler.LoginHandler()).Methods("POST")

    // GET /menu-items: Lista todos los ítems (público)
    router.HandleFunc("/menu-items", app.MenuItemHandler.ListMenuItemsHandler()).Methods("GET")


    // Order Details (público)
    router.HandleFunc("/order-details", app.OrderDetailHandler.CreateOrderDetailHandler()).Methods("POST")
    router.HandleFunc("/order-details/{id}", app.OrderDetailHandler.GetOrderDetailByIDHandler()).Methods("GET")
    router.HandleFunc("/orders/{order_id}/details", app.OrderDetailHandler.GetOrderDetailsByOrderIDHandler()).Methods("GET")
    router.HandleFunc("/order-details/{id}", app.OrderDetailHandler.UpdateOrderDetailHandler()).Methods("PUT")
    router.HandleFunc("/order-details/{id}", app.OrderDetailHandler.DeleteOrderDetailHandler()).Methods("DELETE")
    router.HandleFunc("/orders/{order_id}/complete", app.CustomerOrderHandler.CompleteOrderHandler()).Methods("POST")
    //------------------------------------------------------------------------------->>>
    //Rutas Protegidas (con middleware)

    // Rutas del módulo de business (restringidas a administradores y dueños)
    router.Handle("/business", middleware.AuthMiddleware([]models.EmployeeRole{models.EmployeeRoleAdmin, models.EmployeeRoleOwner})(app.BusinessHandler.GetBusinessHandler())).Methods("GET")
    router.Handle("/business", middleware.AuthMiddleware([]models.EmployeeRole{models.EmployeeRoleAdmin, models.EmployeeRoleOwner})(app.BusinessHandler.UpdateBusinessHandler())).Methods("PUT")

    // Rutas del módulo de employees (restringidas a administradores y dueños)
    router.Handle("/employees", middleware.AuthMiddleware([]models.EmployeeRole{models.EmployeeRoleAdmin, models.EmployeeRoleOwner})(app.EmployeeHandler.CreateEmployeeHandler())).Methods("POST")
    router.Handle("/employees/{id}", middleware.AuthMiddleware([]models.EmployeeRole{models.EmployeeRoleAdmin, models.EmployeeRoleOwner})(app.EmployeeHandler.GetEmployeeHandler())).Methods("GET")
    router.Handle("/employees", middleware.AuthMiddleware([]models.EmployeeRole{models.EmployeeRoleAdmin, models.EmployeeRoleOwner})(app.EmployeeHandler.ListEmployeesHandler())).Methods("GET")
    router.Handle("/employees/role/employee", middleware.AuthMiddleware([]models.EmployeeRole{models.EmployeeRoleAdmin, models.EmployeeRoleOwner})(app.EmployeeHandler.ListEmployeesByRoleHandler(models.EmployeeRoleEmployee))).Methods("GET")
    router.Handle("/employees/{id}", middleware.AuthMiddleware([]models.EmployeeRole{models.EmployeeRoleAdmin, models.EmployeeRoleOwner})(app.EmployeeHandler.UpdateEmployeeHandler())).Methods("PUT")
    router.Handle("/employees/{id}/password", middleware.AuthMiddleware([]models.EmployeeRole{models.EmployeeRoleAdmin, models.EmployeeRoleOwner})(app.EmployeeHandler.UpdateEmployeePasswordHandler())).Methods("PUT")

    // Rutas del módulo de employee_tasks
    router.Handle("/tasks", middleware.AuthMiddleware([]models.EmployeeRole{models.EmployeeRoleAdmin, models.EmployeeRoleOwner})(app.EmployeeTaskHandler.CreateTaskHandler())).Methods("POST")
    router.Handle("/tasks/{id}", middleware.AuthMiddleware([]models.EmployeeRole{models.EmployeeRoleAdmin, models.EmployeeRoleOwner, models.EmployeeRoleEmployee})(app.EmployeeTaskHandler.GetTaskHandler())).Methods("GET")
    router.Handle("/tasks", middleware.AuthMiddleware([]models.EmployeeRole{models.EmployeeRoleAdmin, models.EmployeeRoleOwner})(app.EmployeeTaskHandler.ListTasksHandler())).Methods("GET")
    router.Handle("/my-tasks", middleware.AuthMiddleware([]models.EmployeeRole{models.EmployeeRoleAdmin, models.EmployeeRoleOwner, models.EmployeeRoleEmployee})(app.EmployeeTaskHandler.ListTasksByEmployeeHandler())).Methods("GET")
    router.Handle("/tasks/{id}", middleware.AuthMiddleware([]models.EmployeeRole{models.EmployeeRoleAdmin, models.EmployeeRoleOwner})(app.EmployeeTaskHandler.UpdateTaskHandler())).Methods("PUT")
    router.Handle("/tasks/{id}/status", middleware.AuthMiddleware([]models.EmployeeRole{models.EmployeeRoleAdmin, models.EmployeeRoleOwner})(app.EmployeeTaskHandler.UpdateTaskStatusHandler())).Methods("PUT")
    router.Handle("/tasks/{id}", middleware.AuthMiddleware([]models.EmployeeRole{models.EmployeeRoleAdmin, models.EmployeeRoleOwner})(app.EmployeeTaskHandler.DeleteTaskHandler())).Methods("DELETE")

    // Rutas del módulo de tables
    router.Handle("/tables", middleware.AuthMiddleware([]models.EmployeeRole{models.EmployeeRoleAdmin, models.EmployeeRoleOwner})(app.TableHandler.CreateTableHandler())).Methods("POST")
    router.Handle("/tables/{id}", middleware.AuthMiddleware([]models.EmployeeRole{models.EmployeeRoleAdmin, models.EmployeeRoleOwner, models.EmployeeRoleEmployee})(app.TableHandler.GetTableHandler())).Methods("GET")
    router.Handle("/tables", middleware.AuthMiddleware([]models.EmployeeRole{models.EmployeeRoleAdmin, models.EmployeeRoleOwner, models.EmployeeRoleEmployee})(app.TableHandler.ListTablesHandler())).Methods("GET")
    router.Handle("/tables/{id}", middleware.AuthMiddleware([]models.EmployeeRole{models.EmployeeRoleAdmin, models.EmployeeRoleOwner, models.EmployeeRoleEmployee})(app.TableHandler.UpdateTableHandler())).Methods("PUT")

    // Rutas del módulo de menu_items
    router.Handle("/menu-items", middleware.AuthMiddleware([]models.EmployeeRole{models.EmployeeRoleAdmin})(app.MenuItemHandler.CreateMenuItemHandler())).Methods("POST")
    router.Handle("/menu-items/{id}", middleware.AuthMiddleware([]models.EmployeeRole{models.EmployeeRoleAdmin})(app.MenuItemHandler.GetMenuItemHandler())).Methods("GET")
    router.Handle("/menu-items/{id}", middleware.AuthMiddleware([]models.EmployeeRole{models.EmployeeRoleAdmin})(app.MenuItemHandler.UpdateMenuItemHandler())).Methods("PUT")
    router.Handle("/menu-items/{id}", middleware.AuthMiddleware([]models.EmployeeRole{models.EmployeeRoleAdmin})(app.MenuItemHandler.DeleteMenuItemHandler())).Methods("DELETE")


    // Rutas del módulo de ordenes
    router.Handle("/orders/{order_id}", middleware.AuthMiddleware([]models.EmployeeRole{models.EmployeeRoleAdmin, models.EmployeeRoleEmployee})(app.CustomerOrderHandler.GetOrderWithDetailsHandler())).Methods("GET")
    router.Handle("/orders/{order_id}/complete-by-employee", middleware.AuthMiddleware([]models.EmployeeRole{models.EmployeeRoleAdmin, models.EmployeeRoleEmployee})(app.CustomerOrderHandler.CompleteOrderByEmployeeHandler())).Methods("POST")

    return router
}