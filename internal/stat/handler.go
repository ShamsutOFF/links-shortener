package stat

import (
	"fmt"
	"links-shortener/configs"
	"links-shortener/pkg/event"
	"links-shortener/pkg/middleware"
	"net/http"
	"time"
)

const (
	FilterByDay   = "day"
	FilterByMonth = "month"
)

type StatHandlerDeps struct {
	StatRepository *StatRepository
	Config         *configs.Config
}

type StatHandler struct {
	StatRepository *StatRepository
	EventBus       *event.EventBus
}

func NewStatHandler(router *http.ServeMux, deps StatHandlerDeps) {
	handler := &StatHandler{
		StatRepository: deps.StatRepository,
	}
	router.Handle("GET /stat",
		middleware.IsAuthenticated(handler.GetStat(), deps.Config))
}

func (h *StatHandler) GetStat() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		from, err := time.Parse("2006-01-02", r.URL.Query().Get("from"))
		if err != nil {
			http.Error(w, "Invalid from format", http.StatusBadRequest)
			return
		}
		to, err := time.Parse("2006-01-02", r.URL.Query().Get("to"))
		if err != nil {
			http.Error(w, "Invalid to format", http.StatusBadRequest)
			return
		}
		by := r.URL.Query().Get("by")
		if by != FilterByDay && by != FilterByMonth {
			http.Error(w, "Invalid by format", http.StatusBadRequest)
			return
		}
		fmt.Println(from, to, by)
	}
}
