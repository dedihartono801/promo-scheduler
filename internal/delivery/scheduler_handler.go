package delivery

import (
	"github.com/dedihartono801/promo-scheduler/internal/app/usecase/scheduler"
)

type SchedulerHandler interface {
	Scheduler() error
}

type schedulerHandler struct {
	service scheduler.Service
}

func NewSchedulerHandler(service scheduler.Service) SchedulerHandler {
	return &schedulerHandler{service}
}

func (h *schedulerHandler) Scheduler() error {
	err := h.service.Scheduler()

	if err != nil {
		return err
	}
	return nil

}
