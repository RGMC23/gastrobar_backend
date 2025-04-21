-- GASTROBAR --

-- =====================================================================
-- CREACIÓN DE TABLAS
-- =====================================================================

-- Crear la tabla para la información del negocio (business)
CREATE TABLE business (
    id               SERIAL PRIMARY KEY,
    business_name    VARCHAR(100) NOT NULL,
    address          TEXT         NOT NULL,
    phone_number     VARCHAR(20),
    email            VARCHAR(255),
    corporate_reason VARCHAR(20)  NOT NULL,
    created_at       TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at       TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Crear la tabla para almacenar los logs de cambios en business (business_info_logs)
CREATE TABLE business_info_logs (
    id               SERIAL PRIMARY KEY,
    business_info_id INTEGER     NOT NULL,
    operation        VARCHAR(10) NOT NULL,
    old_data         JSONB,
    new_data         JSONB,
    changed_at       TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Crear un tipo ENUM para los roles
CREATE TYPE employee_role AS ENUM ('dueño', 'administrador', 'empleado');

-- Crear la tabla para los empleados (employees, con role como ENUM)
CREATE TABLE employees (
    id            SERIAL PRIMARY KEY,
    employee_name VARCHAR(100)  NOT NULL,
    email         VARCHAR(255),
    phone_number  VARCHAR(20),
    role          employee_role NOT NULL,
    username      VARCHAR(200)  NOT NULL,
    password      VARCHAR(200)  NOT NULL,
    created_at    TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Crear la tabla para las mesas (tables)
CREATE TABLE tables (
    id         SERIAL PRIMARY KEY,
    table_name VARCHAR(50),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Crear la tabla menu_items para agregar el campo description
CREATE TABLE menu_items (
    id          SERIAL PRIMARY KEY,
    item_name   VARCHAR(100)   NOT NULL,
    category    VARCHAR(50),
    price       NUMERIC(10, 2) NOT NULL,
    stock       INTEGER        NOT NULL DEFAULT 0,
    description TEXT,
    created_at  TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Crear la tabla customer_orders para asociar pedidos con mesas
CREATE TABLE customer_orders (
    id           SERIAL PRIMARY KEY,
    table_id     INTEGER        NOT NULL REFERENCES tables(id) ON DELETE CASCADE,
    total_amount NUMERIC(10, 2) NOT NULL DEFAULT 0.0,
    status       VARCHAR(20)    NOT NULL CHECK (status IN ('pending', 'completed')),
    created_at   TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Crear la tabla order_details (eliminamos price_at_time)
CREATE TABLE order_details (
    id           SERIAL PRIMARY KEY,
    order_id     INTEGER        NOT NULL REFERENCES customer_orders(id) ON DELETE CASCADE,
    menu_item_id INTEGER        NOT NULL REFERENCES menu_items(id),
    quantity     INTEGER        NOT NULL CHECK (quantity > 0),
    created_at   TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Crear la tabla para las tareas de los empleados (employee_tasks)
CREATE TABLE employee_tasks (
    id               SERIAL PRIMARY KEY,
    employee_id      INTEGER NOT NULL REFERENCES employees(id),
    task_description TEXT    NOT NULL,
    status           VARCHAR(20) DEFAULT 'pending',
    assigned_at      TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    completed_at     TIMESTAMP WITH TIME ZONE,
    created_at       TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- =====================================================================
-- FUNCIONES Y TRIGGERS (ORDENADOS POR TABLA AFECTADA)
-- =====================================================================

-- --------------------- Triggers para business -------------------------

-- Crear la función para registrar los logs
CREATE OR REPLACE FUNCTION log_business_info_changes()
RETURNS TRIGGER AS $$
BEGIN
    IF TG_OP = 'INSERT' THEN
        INSERT INTO business_info_logs (business_info_id, operation, old_data, new_data, changed_at)
        VALUES (NEW.id, TG_OP, NULL, row_to_json(NEW), CURRENT_TIMESTAMP);
        RETURN NEW;
    ELSIF TG_OP = 'UPDATE' THEN
        INSERT INTO business_info_logs (business_info_id, operation, old_data, new_data, changed_at)
        VALUES (NEW.id, TG_OP, row_to_json(OLD), row_to_json(NEW), CURRENT_TIMESTAMP);
        RETURN NEW;
    ELSIF TG_OP = 'DELETE' THEN
        INSERT INTO business_info_logs (business_info_id, operation, old_data, new_data, changed_at)
        VALUES (OLD.id, TG_OP, row_to_json(OLD), NULL, CURRENT_TIMESTAMP);
        RETURN OLD;
    END IF;
    RETURN NULL;
END;
$$ LANGUAGE plpgsql;

-- Crear los triggers para INSERT, UPDATE y DELETE en business
CREATE TRIGGER business_info_insert_trigger
    AFTER INSERT ON business
    FOR EACH ROW
    EXECUTE FUNCTION log_business_info_changes();

CREATE TRIGGER business_info_update_trigger
    AFTER UPDATE ON business
    FOR EACH ROW
    EXECUTE FUNCTION log_business_info_changes();

CREATE TRIGGER business_info_delete_trigger
    AFTER DELETE ON business
    FOR EACH ROW
    EXECUTE FUNCTION log_business_info_changes();

-- --------------------- Triggers para menu_items (a través de order_details) -------------------------

-- Crear la función para actualizar el stock automáticamente
CREATE OR REPLACE FUNCTION update_menu_item_stock()
RETURNS TRIGGER AS $$
DECLARE
    item_name     VARCHAR(100);
    current_stock INTEGER;
BEGIN
    -- Obtener el nombre y el stock actual del ítem
    SELECT menu_items.item_name, stock INTO item_name, current_stock
    FROM menu_items
    WHERE id = NEW.menu_item_id;

    -- Reducir el stock del ítem en menu_items basado en la quantity de order_details
    UPDATE menu_items
    SET stock = stock - NEW.quantity
    WHERE id = NEW.menu_item_id
      AND stock >= NEW.quantity;

    -- Verificar si la actualización fue exitosa (stock suficiente)
    IF NOT FOUND THEN
        RAISE EXCEPTION 'Insufficient stock for item % (ID: %). Available: %, Requested: %', item_name, NEW.menu_item_id, current_stock, NEW.quantity;
    END IF;

    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Crear el trigger para actualizar el stock después de insertar en order_details
CREATE TRIGGER update_stock_after_order
    AFTER INSERT ON order_details
    FOR EACH ROW
    EXECUTE FUNCTION update_menu_item_stock();

-- --------------------- Triggers para customer_orders -------------------------

-- Crear la función para actualizar total_amount en customer_orders
CREATE OR REPLACE FUNCTION update_customer_order_total_amount()
RETURNS TRIGGER AS $$
DECLARE
    new_total NUMERIC(10, 2);
    order_status VARCHAR(20);
BEGIN
    -- Determinar el order_id dependiendo de la operación
    IF TG_OP = 'DELETE' THEN
        -- Para DELETE, usamos OLD.order_id
        SELECT status INTO order_status
        FROM customer_orders
        WHERE id = OLD.order_id;

        -- Calcular el nuevo total_amount sumando los subtotales de los order_details restantes
        SELECT COALESCE(SUM(od.quantity * mi.price), 0)
        INTO new_total
        FROM order_details od
        JOIN menu_items mi ON od.menu_item_id = mi.id
        WHERE od.order_id = OLD.order_id;
    ELSE
        -- Para INSERT y UPDATE, usamos NEW.order_id
        SELECT status INTO order_status
        FROM customer_orders
        WHERE id = NEW.order_id;

        -- Calcular el nuevo total_amount sumando los subtotales de los order_details
        SELECT COALESCE(SUM(od.quantity * mi.price), 0)
        INTO new_total
        FROM order_details od
        JOIN menu_items mi ON od.menu_item_id = mi.id
        WHERE od.order_id = NEW.order_id;
    END IF;

    -- Validar el estado de la customer_order
    IF order_status IS NULL THEN
        RAISE EXCEPTION 'Customer order with ID % not found.',
            CASE WHEN TG_OP = 'DELETE' THEN OLD.order_id ELSE NEW.order_id END;
    END IF;

    IF order_status != 'pending' THEN
        RAISE NOTICE 'Customer order with ID % is in state %, skipping total_amount update.',
            CASE WHEN TG_OP = 'DELETE' THEN OLD.order_id ELSE NEW.order_id END, order_status;
        RETURN NULL;
    END IF;

    -- Actualizar el total_amount en customer_orders
    UPDATE customer_orders
    SET total_amount = new_total
    WHERE id = CASE WHEN TG_OP = 'DELETE' THEN OLD.order_id ELSE NEW.order_id END;

    RETURN NULL;
END;
$$ LANGUAGE plpgsql;

-- Trigger para después de insertar un order_detail
CREATE TRIGGER update_total_amount_after_insert
    AFTER INSERT ON order_details
    FOR EACH ROW
    EXECUTE FUNCTION update_customer_order_total_amount();

-- Trigger para después de actualizar quantity en un order_detail (eliminamos price_at_time)
CREATE TRIGGER update_total_amount_after_update
    AFTER UPDATE OF quantity ON order_details
    FOR EACH ROW
    EXECUTE FUNCTION update_customer_order_total_amount();

-- Trigger para después de eliminar un order_detail
CREATE TRIGGER update_total_amount_after_delete
    AFTER DELETE ON order_details
    FOR EACH ROW
    EXECUTE FUNCTION update_customer_order_total_amount();

-- =====================================================================
-- ÍNDICES PARA MEJORAR EL RENDIMIENTO
-- =====================================================================

CREATE INDEX idx_customer_orders_table_id ON customer_orders(table_id);
CREATE INDEX idx_order_details_order_id ON order_details(order_id);
CREATE INDEX idx_order_details_menu_item_id ON order_details(menu_item_id);
CREATE INDEX idx_employee_tasks_employee_id ON employee_tasks(employee_id);
CREATE INDEX idx_employee_tasks_status ON employee_tasks(status);

-- =====================================================================
-- DATOS INICIALES PARA PRUEBAS
-- =====================================================================

-- Datos para la información del negocio
INSERT INTO business (business_name, address, phone_number, email, corporate_reason)
VALUES ('Gastrobar El Sabor', 'Calle Falsa 123, Ciudad Ejemplo', '+1234567890', 'contacto@gastrobar.com', 'S.A.S.');

-- Datos para los empleados (usando los roles del ENUM)
INSERT INTO employees (employee_name, email, phone_number, role, username, password)
VALUES ('Juan Pérez', 'juan@gastrobar.com', '+1234567891', 'dueño', 'juanperez', '$2a$10$Ef6Y5TC3WGfxvLrAvnKQIuq1e7gp/JteObA.mBEfw3Mi68I4ldBNG'),
       ('María Gómez', 'maria@gastrobar.com', '+1234567892', 'administrador', 'mariagomez', '$2a$10$Ef6Y5TC3WGfxvLrAvnKQIuq1e7gp/JteObA.mBEfw3Mi68I4ldBNG'),
       ('Carlos López', 'carlos@gastrobar.com', '+1234567893', 'empleado', 'carloslopez', '$2a$10$Ef6Y5TC3WGfxvLrAvnKQIuq1e7gp/JteObA.mBEfw3Mi68I4ldBNG');

-- Datos para las mesas
INSERT INTO tables (table_name)
VALUES ('Mesa 1 - Para 4'),
       ('Mesa 2 - Para 2'),
       ('Mesa 3 - Para 6');

-- Datos para el menú (actualizados con description)
INSERT INTO menu_items (item_name, category, price, stock, description)
VALUES ('Cerveza artesanal', 'Bebidas', 5.50, 100, 'Cerveza artesanal de cebada y lúpulo'),
       ('Ensalada César', 'Entradas', 8.00, 50, 'Ensalada con lechuga, pollo y aderezo César'),
       ('Hamburguesa clásica', 'Platos fuertes', 12.00, 30, 'Hamburguesa con carne y vegetales frescos');

-- Datos para las tareas de los empleados
INSERT INTO employee_tasks (employee_id, task_description, status)
VALUES (1, 'Limpiar Mesa 1', 'pending'),
       (2, 'Atender Mesa 2', 'pending'),
       (3, 'Preparar pedido de Hamburguesa', 'pending');

-- Datos para customer_orders
INSERT INTO customer_orders (table_id, total_amount, status)
VALUES (1, 0.0, 'pending');

-- Datos para order_details (eliminamos price_at_time)
INSERT INTO order_details (order_id, menu_item_id, quantity)
VALUES (1, 1, 2), -- 2 cervezas artesanales
       (1, 2, 1); -- 1 ensalada César