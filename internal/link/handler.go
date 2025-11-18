package link

import (
	"fmt"
	"gorm.io/gorm"
	"links-shortener/configs"
	"links-shortener/pkg/event"
	"links-shortener/pkg/middleware"
	"links-shortener/pkg/req"
	"links-shortener/pkg/res"
	"net/http"
	"strconv"
)

type LinkHandlerDeps struct {
	LinkRepository *LinkRepository
	Config         *configs.Config
	EventBus       *event.EventBus
}

type LinkHandler struct {
	LinkRepository *LinkRepository
	EventBus       *event.EventBus
}

func NewLinkHandler(router *http.ServeMux, deps LinkHandlerDeps) {
	handler := &LinkHandler{
		LinkRepository: deps.LinkRepository,
		EventBus:       deps.EventBus,
	}
	router.HandleFunc("POST /link", handler.Create())
	router.Handle("PATCH /link/{id}",
		middleware.IsAuthenticated(handler.Update(), deps.Config))
	router.HandleFunc("DELETE /link/{id}", handler.Delete())
	router.HandleFunc("GET /{hash}", handler.GoTo())

	router.Handle("GET /link",
		middleware.IsAuthenticated(handler.GetAll(), deps.Config))
}

func (handler *LinkHandler) Create() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		body, err := req.HandleBody[LinkCreateRequest](&writer, request)
		if err != nil {
			return
		}
		link := NewLink(body.Url)

		for {
			existedLink, _ := handler.LinkRepository.GetByHash(link.Hash)
			if existedLink == nil {
				break
			}
			link.GenerateHash()
		}

		createdLink, err := handler.LinkRepository.Create(link)
		if err != nil {
			http.Error(writer, err.Error(), http.StatusInternalServerError)
			return
		}
		res.JsonResp(writer, createdLink, http.StatusCreated)
	}
}

func (handler *LinkHandler) Update() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		email, ok := request.Context().Value(middleware.ContextEmailKey).(string)
		if ok {
			fmt.Println("@@@ ", email)
		}
		body, err := req.HandleBody[LinkUpdateRequest](&writer, request)
		if err != nil {
			return
		}
		idString := request.PathValue("id")
		id, err := strconv.ParseUint(idString, 10, 32)
		if err != nil {
			http.Error(writer, err.Error(), http.StatusBadRequest)
			return
		}
		link, err := handler.LinkRepository.Update(&Link{
			Model: gorm.Model{ID: uint(id)},
			Url:   body.Url,
			Hash:  body.Hash,
		})
		if err != nil {
			http.Error(writer, err.Error(), http.StatusBadRequest)
			return
		}
		res.JsonResp(writer, link, http.StatusOK)
	}
}

func (handler *LinkHandler) Delete() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		idString := request.PathValue("id")
		id, err := strconv.ParseUint(idString, 10, 32)
		if err != nil {
			http.Error(writer, err.Error(), http.StatusBadRequest)
			return
		}
		_, err = handler.LinkRepository.GetById(uint(id))
		if err != nil {
			http.Error(writer, err.Error(), http.StatusNotFound)
			return
		}
		err = handler.LinkRepository.Delete(uint(id))
		if err != nil {
			http.Error(writer, err.Error(), http.StatusInternalServerError)
			return
		}
		res.JsonResp(writer, nil, http.StatusOK)
	}
}

func (handler *LinkHandler) GoTo() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		hash := request.PathValue("hash")
		link, err := handler.LinkRepository.GetByHash(hash)
		if err != nil {
			http.Error(writer, err.Error(), http.StatusNotFound)
			return
		}
		go handler.EventBus.Publish(event.Event{
			Type: event.EventLinkVisited,
			Data: link.ID,
		})
		http.Redirect(writer, request, link.Url, http.StatusTemporaryRedirect)
	}
}

func (handler *LinkHandler) GetAll() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		limit, err := strconv.Atoi(request.URL.Query().Get("limit"))
		if err != nil {
			http.Error(writer, "Invalid limit", http.StatusBadRequest)
			return
		}

		offset, err := strconv.Atoi(request.URL.Query().Get("offset"))
		if err != nil {
			http.Error(writer, "Invalid offset", http.StatusBadRequest)
			return
		}
		links := handler.LinkRepository.GetAll(limit, offset)
		count := handler.LinkRepository.GetCount()
		res.JsonResp(writer, GetAllLinksResponse{
			Links: links,
			Count: count,
		}, http.StatusOK)
	}
}
