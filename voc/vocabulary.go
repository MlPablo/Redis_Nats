package voc

const (
	SubjectGet    = "service.get"
	SubjectDelete = "service.delete"
	SubjectUpdate = "service.update"
	SubjectCreate = "service.create"

	SubjectCreateOrder       = "service.create.order"
	SubjectStatusCreateOrder = "service.create.order.status"

	NatsToOrderServicesQueue = "order_service_queue"
	NatsToCrudServicesQueue  = "crud_service_queue"
	NatsOrderStreamName      = "OrdersStream"
)
