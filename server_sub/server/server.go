package server

import (
	"github.com/gorilla/mux"
	"html/template"
	"log"
	"net/http"
	"wbL0/server_sub/entity"
	"wbL0/server_sub/store"
)

type Server struct {
	Srv     *http.Server
	storage store.StoreService
	Addr    string
}

func InitServer(store store.StoreService, addr string) *Server {
	server := Server{
		storage: store,
		Addr:    addr,
	}
	return &server
}

func (s *Server) Start() error {
	router := mux.NewRouter()
	router.HandleFunc("/orders/{o_id}", s.ordersHandler)
	s.Srv = &http.Server{Addr: s.Addr, Handler: router}
	log.Println("Server is starting")
	err := s.Srv.ListenAndServe()
	if err != nil {
		return err
	}
	return err
}

func (s *Server) Stop() error {
	log.Println("Server stops")
	return s.Srv.Close()
}

func (s *Server) ordersHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["o_id"]
	od := s.storage.GetFromCacheByUID(id)
	if od.OrderUid == "" {
		w.WriteHeader(400)
		parsedTemplate, _ := template.ParseFiles("./server/404.html")
		err := parsedTemplate.Execute(w, struct{ Id string }{Id: id})
		if err != nil {
			_, err = w.Write([]byte("no data with id " + id))
			if err != nil {
				return
			}
			log.Printf("Error occurred while executing the template : ", id)
			return
		}
		return
	}

	dataItem := entity.DataItem{
		ID:        id,
		OrderData: od,
	}
	parsedTemplate, _ := template.ParseFiles("./server/id.html")
	err := parsedTemplate.Execute(w, dataItem)
	if err != nil {
		w.Write([]byte("error while executing template"))
		log.Printf("Error occurred while executing the template : ", dataItem)
		return
	}
	w.WriteHeader(200)
}
