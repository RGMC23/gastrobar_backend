package services

import (
	"gastrobar-backend/internal/models"
	"gastrobar-backend/internal/repositories"
	"time"

	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
)

type OrderDetailService interface {
	CreateOrderDetail(orderDetail models.OrderDetail, tableID int) (models.OrderDetail, error)
	GetOrderDetailByID(id int) (models.OrderDetail, error)
	GetOrderDetailsByOrderID(orderID int) ([]models.OrderDetail, error)
	UpdateOrderDetail(orderDetail models.OrderDetail) (models.OrderDetail, error)
	DeleteOrderDetail(id int) error
}

type orderDetailService struct {
	orderDetailRepo   repositories.OrderDetailRepository
	customerOrderRepo repositories.CustomerOrderRepository
	menuItemRepo      repositories.MenuItemRepository
	tableRepo         repositories.TableRepository
}

func NewOrderDetailService(
	orderDetailRepo repositories.OrderDetailRepository,
	customerOrderRepo repositories.CustomerOrderRepository,
	menuItemRepo repositories.MenuItemRepository,
	tableRepo repositories.TableRepository,
) OrderDetailService {
	return &orderDetailService{
		orderDetailRepo:   orderDetailRepo,
		customerOrderRepo: customerOrderRepo,
		menuItemRepo:      menuItemRepo,
		tableRepo:         tableRepo,
	}
}

func (s *orderDetailService) CreateOrderDetail(orderDetail models.OrderDetail, tableID int) (models.OrderDetail, error) {
	// Validar que la mesa exista
	_, err := s.tableRepo.FindByID(tableID)
	if err != nil {
		return models.OrderDetail{}, errors.Wrap(err, "failed to find table")
	}

	var customerOrder models.CustomerOrder
	var menu_item models.MenuItem

	// Caso 1: No se proporciona order_id (o es 0), intentamos crear un nuevo customer_order
	if orderDetail.OrderID == 0 {
		// Verificar si ya existe un customer_order pendiente para esta mesa
		pendingOrder, err := s.customerOrderRepo.FindPendingByTableID(tableID)
		if err == nil && pendingOrder.ID != 0 {
			return models.OrderDetail{}, errors.New("cannot create new customer order: a pending customer order already exists for this table")
		}

		// Verificar si la última orden está completada (opcional, para confirmar el estado)
		completedOrder, completedErr := s.customerOrderRepo.FindCompletedByTableID(tableID)
		if completedErr == nil && completedOrder.ID != 0 {
			// La última orden está completada, podemos proceder a crear una nueva
		}

		// Crear un nuevo customer_order
		newOrder := models.CustomerOrder{
			TableID:     tableID,
			TotalAmount: decimal.NewFromFloat(0.0),
			Status:      "pending",
			CreatedAt:   time.Now(),
		}
		customerOrder, err = s.customerOrderRepo.Create(newOrder)
		if err != nil {
			return models.OrderDetail{}, errors.Wrap(err, "failed to create customer order")
		}
	} else {
		// Caso 2: Se proporciona un order_id, validamos que exista y esté pendiente
		customerOrder, err = s.customerOrderRepo.FindByID(orderDetail.OrderID)
		if err != nil {
			return models.OrderDetail{}, errors.Wrap(err, "failed to find customer order")
		}
		if customerOrder.TableID != tableID {
			return models.OrderDetail{}, errors.New("customer order does not belong to the specified table")
		}
		if customerOrder.Status == "completed" {
			return models.OrderDetail{}, errors.New("cannot add order detail: customer order is already completed")
		}
	}

	// Validar el menu_item y el stock
	menuItem, err := s.menuItemRepo.FindByID(orderDetail.MenuItemID)
	if err != nil {
		return models.OrderDetail{}, errors.Wrap(err, "failed to find menu item")
	}

	// Validar el stock
	if menuItem.Stock < orderDetail.Quantity {
		return models.OrderDetail{}, errors.Errorf("insufficient stock for item %s (ID: %d). Available: %d, Requested: %d",
			menuItem.ItemName, menuItem.ID, menuItem.Stock, orderDetail.Quantity)
	}

	// Asignar el order_id y created_at
	orderDetail.OrderID = customerOrder.ID
	orderDetail.CreatedAt = time.Now()

	// Crear el order_detail
	createdDetail, err := s.orderDetailRepo.Create(orderDetail)
	if err != nil {
		return models.OrderDetail{}, errors.Wrap(err, "failed to create order detail")
	}
	//cosultar el producto asociado a la orden y insertarlo en la orden
	menu_item, err = s.menuItemRepo.FindByID(orderDetail.MenuItemID)
	createdDetail.MenuItem = menu_item

	return createdDetail, nil
}

func (s *orderDetailService) GetOrderDetailByID(id int) (models.OrderDetail, error) {
	if id <= 0 {
		return models.OrderDetail{}, errors.New("invalid order detail ID")
	}

	orderDetail, err := s.orderDetailRepo.FindByID(id)
	if err != nil {
		return models.OrderDetail{}, errors.Wrap(err, "failed to get order detail")
	}

	return orderDetail, nil
}

func (s *orderDetailService) GetOrderDetailsByOrderID(orderID int) ([]models.OrderDetail, error) {
	if orderID <= 0 {
		return nil, errors.New("invalid order ID")
	}

	// Validar que la orden existe
	_, err := s.customerOrderRepo.FindByID(orderID)
	if err != nil {
		return nil, errors.Wrap(err, "failed to find customer order")
	}

	orderDetails, err := s.orderDetailRepo.FindByOrderID(orderID)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get order details by order id")
	}

	return orderDetails, nil
}

func (s *orderDetailService) UpdateOrderDetail(orderDetail models.OrderDetail) (models.OrderDetail, error) {
	if orderDetail.ID <= 0 {
		return models.OrderDetail{}, errors.New("invalid order detail ID")
	}

	// Validar que la orden existe
	customerOrder, err := s.customerOrderRepo.FindByID(orderDetail.OrderID)
	if err != nil {
		return models.OrderDetail{}, errors.Wrap(err, "failed to find customer order")
	}

	// Validar que la customer_order no esté completada
	if customerOrder.Status == "completed" {
		return models.OrderDetail{}, errors.New("cannot update order detail: customer order is already completed")
	}

	// Validar el menu_item y el stock
	menuItem, err := s.menuItemRepo.FindByID(orderDetail.MenuItemID)
	if err != nil {
		return models.OrderDetail{}, errors.Wrap(err, "failed to find menu item")
	}

	// Validar el stock
	if menuItem.Stock < orderDetail.Quantity {
		return models.OrderDetail{}, errors.Errorf("insufficient stock for item %s (ID: %d). Available: %d, Requested: %d",
			menuItem.ItemName, menuItem.ID, menuItem.Stock, orderDetail.Quantity)
	}

	// Actualizar el order_detail
	updatedDetail, err := s.orderDetailRepo.Update(orderDetail)
	if err != nil {
		return models.OrderDetail{}, errors.Wrap(err, "failed to update order detail")
	}

	return updatedDetail, nil
}

func (s *orderDetailService) DeleteOrderDetail(id int) error {
	if id <= 0 {
		return errors.New("invalid order detail ID")
	}

	// Obtener el order_detail para encontrar su order_id
	orderDetail, err := s.orderDetailRepo.FindByID(id)
	if err != nil {
		return errors.Wrap(err, "failed to find order detail")
	}

	// Validar que la customer_order no esté completada
	customerOrder, err := s.customerOrderRepo.FindByID(orderDetail.OrderID)
	if err != nil {
		return errors.Wrap(err, "failed to find customer order")
	}
	if customerOrder.Status == "completed" {
		return errors.New("cannot delete order detail: customer order is already completed")
	}

	// Eliminar el order_detail
	err = s.orderDetailRepo.Delete(id)
	if err != nil {
		return errors.Wrap(err, "failed to delete order detail")
	}

	return nil
}
