package services

import (
    "gastrobar-backend/internal/models"
    "gastrobar-backend/internal/repositories"

    "github.com/pkg/errors"
)

type CustomerOrderService interface {
    GetOrderWithDetails(orderID int) (models.CustomerOrder, error)
    CompleteOrder(orderID int) (models.CustomerOrder, error)
    CompleteOrderByEmployee(orderID int) (models.CustomerOrder, error)
}

type customerOrderService struct {
    customerOrderRepo repositories.CustomerOrderRepository
    menuItemRepo      repositories.MenuItemRepository
    orderDetailRepo   repositories.OrderDetailRepository
}

func NewCustomerOrderService(
    customerOrderRepo repositories.CustomerOrderRepository,
    menuItemRepo repositories.MenuItemRepository,
    orderDetailRepo repositories.OrderDetailRepository,
) CustomerOrderService {
    return &customerOrderService{
        customerOrderRepo: customerOrderRepo,
        menuItemRepo:      menuItemRepo,
        orderDetailRepo:   orderDetailRepo,
    }
}

func (s *customerOrderService) GetOrderWithDetails(orderID int) (models.CustomerOrder, error) {
    // Validar que la orden existe
    order, err := s.customerOrderRepo.FindByID(orderID)
    if err != nil {
        return models.CustomerOrder{}, errors.Wrap(err, "failed to find customer order")
    }

    // Obtener los order_details asociados
    orderDetails, err := s.orderDetailRepo.FindByOrderID(orderID)
    if err != nil {
        return models.CustomerOrder{}, errors.Wrap(err, "failed to find order details")
    }

    // Asignar los order_details a la orden
    order.OrderDetails = orderDetails

    return order, nil
}

func (s *customerOrderService) CompleteOrder(orderID int) (models.CustomerOrder, error) {
    // Validar que la orden existe
    order, err := s.customerOrderRepo.FindByID(orderID)
    if err != nil {
        return models.CustomerOrder{}, errors.Wrap(err, "failed to find customer order")
    }

    // Validar que la orden esté en estado 'pending'
    if order.Status != "pending" {
        return models.CustomerOrder{}, errors.New("order is not in 'pending' state")
    }

    // Actualizar el estado a 'completed'
    updatedOrder, err := s.customerOrderRepo.UpdateStatus(orderID, "completed")
    if err != nil {
        return models.CustomerOrder{}, errors.Wrap(err, "failed to complete order")
    }

    return updatedOrder, nil
}

func (s *customerOrderService) CompleteOrderByEmployee(orderID int) (models.CustomerOrder, error) {
    return s.CompleteOrder(orderID)
}