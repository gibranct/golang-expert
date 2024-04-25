package handlers

import (
	"encoding/json"
	"net/http"

	"github.com.br/gibranct/golang-course-api/internal/dto"
	"github.com.br/gibranct/golang-course-api/internal/entity"
	"github.com.br/gibranct/golang-course-api/internal/infra/database"
	entityPkg "github.com.br/gibranct/golang-course-api/pkg/entity"
	"github.com/go-chi/chi"
)

type ProductHandler struct {
	ProductDB database.ProductInterface
}

func NewProductHandler(db database.ProductInterface) *ProductHandler {
	return &ProductHandler{
		ProductDB: db,
	}
}

func (h *ProductHandler) CreateProduct(w http.ResponseWriter, r *http.Request) {
	var product dto.CreateProductInput
	err := json.NewDecoder(r.Body).Decode(&product)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	p, err := entity.NewProduct(product.Name, product.Price)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	err = h.ProductDB.Create(p)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

func (h *ProductHandler) FindByID(w http.ResponseWriter, r *http.Request) {
	pId := chi.URLParam(r, "id")
	if pId == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if _, err := entityPkg.FromString(pId); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	p, err := h.ProductDB.FindByID(pId)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	w.Header().Add("content-type", "application/json")
	json.NewEncoder(w).Encode(p)
}

func (h *ProductHandler) Update(w http.ResponseWriter, r *http.Request) {
	pId := chi.URLParam(r, "id")
	if pId == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if _, err := entityPkg.FromString(pId); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	p, err := h.ProductDB.FindByID(pId)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	var product dto.UpdateProductInput
	err = json.NewDecoder(r.Body).Decode(&product)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	p.Update(product.Name, product.Price)
	err = h.ProductDB.Update(p)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (h *ProductHandler) DeleteById(w http.ResponseWriter, r *http.Request) {
	pId := chi.URLParam(r, "id")
	if pId == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if _, err := entityPkg.FromString(pId); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	_, err := h.ProductDB.FindByID(pId)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	err = h.ProductDB.Delete(pId)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
