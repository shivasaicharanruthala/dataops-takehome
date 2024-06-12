package handler

import (
	"encoding/json"
	"github.com/shivasaicharanruthala/dataops-takehome/api/store"
	"github.com/shivasaicharanruthala/dataops-takehome/model"
	"net/http"
	"strconv"
)

type loginHandler struct {
	loginStore store.Login
}

func New(loginStore store.Login) *loginHandler {
	return &loginHandler{
		loginStore: loginStore,
	}
}

type responseErr struct {
	StatusCode int    `json:"code"`
	Err        string `json:"message"`
}

func (lh loginHandler) Get(w http.ResponseWriter, r *http.Request) {
	var filter model.Filter

	limit := r.URL.Query().Get("limit")
	page := r.URL.Query().Get("page")
	isEncrypted := r.URL.Query().Get("isEncrypted")
	groupDuplicates := r.URL.Query().Get("groupDuplicates")

	if groupDuplicates != "" {
		groupDuplicatesConv, _ := strconv.ParseBool(groupDuplicates)
		filter.GroupDuplicates = groupDuplicatesConv
	}

	if limit == "" || page == "" || isEncrypted == "" {
		errResp, _ := json.Marshal(responseErr{StatusCode: 400, Err: "query params limit or page or isEncrypted is missing."})

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write(errResp)
		return
	}

	limitConv, _ := strconv.Atoi(limit)
	pageConv, _ := strconv.Atoi(page)
	isEncryptedConv, _ := strconv.ParseBool(isEncrypted)

	filter.Limit = limitConv
	filter.Page = pageConv
	filter.IsEncrypted = isEncryptedConv

	resp, err := lh.loginStore.Get(&filter)
	if err != nil {
		errResp, _ := json.Marshal(responseErr{StatusCode: 400, Err: err.Error()})

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write(errResp)
		return
	}

	respJson, _ := json.Marshal(resp)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	_, _ = w.Write(respJson)
}
