package app

import (
	"database/sql"

	"gastrobar-backend/internal/handlers"
	"gastrobar-backend/internal/repositories"
	"gastrobar-backend/internal/services"
)

// App contiene todas las dependencias de la aplicación
type App struct {
	EmployeeRepo            repositories.EmployeeRepository
	BusinessRepo            repositories.BusinessRepository
	EmployeeTaskRepo        repositories.EmployeeTaskRepository
	TableRepo               repositories.TableRepository
	MenuItemRepo            repositories.MenuItemRepository
	OrderDetailRepo         repositories.OrderDetailRepository
	CustomerOrderRepo       repositories.CustomerOrderRepository
	AuthSvc                 services.AuthService
	BusinessSvc             services.BusinessService
	EmployeeSvc             services.EmployeeService
	EmployeeTaskSvc         services.EmployeeTaskService
	TableSvc                services.TableService
	MenuItemSvc             services.MenuItemService
	CustomerOrderSvc        services.CustomerOrderService
	OrderDetailSvc          services.OrderDetailService
	AuthHandler             *handlers.AuthHandler
	BusinessHandler         *handlers.BusinessHandler
	EmployeeHandler         *handlers.EmployeeHandler
	EmployeeTaskHandler     *handlers.EmployeeTaskHandler
	TableHandler            *handlers.TableHandler
	MenuItemHandler         *handlers.MenuItemHandler
	CustomerOrderHandler    *handlers.CustomerOrderHandler
	OrderDetailHandler      *handlers.OrderDetailHandler
}

// NewApp inicializa todas las dependencias de la aplicación
func NewApp(db *sql.DB) *App {
	// Inicializar repositorios
	employeeRepo := repositories.NewEmployeeRepository(db)
	businessRepo := repositories.NewBusinessRepository(db)
	employeeTaskRepo := repositories.NewEmployeeTaskRepository(db)
	tableRepo := repositories.NewTableRepository(db)
	menuItemRepo := repositories.NewMenuItemRepository(db)
	customerOrderRepo := repositories.NewCustomerOrderRepository(db)
	orderDetailRepo := repositories.NewOrderDetailRepository(db)

	// Inicializar servicios
	authSvc := services.NewAuthService(employeeRepo)
	businessSvc := services.NewBusinessService(businessRepo)
	employeeSvc := services.NewEmployeeService(employeeRepo)
	employeeTaskSvc := services.NewEmployeeTaskService(employeeTaskRepo, employeeRepo)
	tableSvc := services.NewTableService(tableRepo)
	menuItemSvc := services.NewMenuItemService(menuItemRepo)
	customerOrderSvc := services.NewCustomerOrderService(customerOrderRepo, menuItemRepo, orderDetailRepo)
	orderDetailSvc := services.NewOrderDetailService(orderDetailRepo, customerOrderRepo, menuItemRepo, tableRepo)

	// Inicializar manejadores
	authHandler := handlers.NewAuthHandler(authSvc)
	businessHandler := handlers.NewBusinessHandler(businessSvc)
	employeeHandler := handlers.NewEmployeeHandler(employeeSvc)
	employeeTaskHandler := handlers.NewEmployeeTaskHandler(employeeTaskSvc)
	tableHandler := handlers.NewTableHandler(tableSvc)
	menuItemHandler := handlers.NewMenuItemHandler(menuItemSvc)
	customerOrderHandler := handlers.NewCustomerOrderHandler(customerOrderSvc)
	orderDetailHandler := handlers.NewOrderDetailHandler(orderDetailSvc)

	return &App{
		EmployeeRepo:            employeeRepo,
		BusinessRepo:            businessRepo,
		EmployeeTaskRepo:        employeeTaskRepo,
		TableRepo:               tableRepo,
		MenuItemRepo:            menuItemRepo,
		CustomerOrderRepo:       customerOrderRepo,
		OrderDetailRepo:         orderDetailRepo,
		AuthSvc:                 authSvc,
		BusinessSvc:             businessSvc,
		EmployeeSvc:             employeeSvc,
		EmployeeTaskSvc:         employeeTaskSvc,
		TableSvc:                tableSvc,
		MenuItemSvc:             menuItemSvc,
		CustomerOrderSvc:        customerOrderSvc,
		OrderDetailSvc:          orderDetailSvc,
		AuthHandler:             authHandler,
		BusinessHandler:         businessHandler,
		EmployeeHandler:         employeeHandler,
		EmployeeTaskHandler:     employeeTaskHandler,
		TableHandler:            tableHandler,
		MenuItemHandler:         menuItemHandler,
		CustomerOrderHandler:    customerOrderHandler,
		OrderDetailHandler:      orderDetailHandler,
	}
}
