package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

func (fe *frontendServer) homeHandler(w http.ResponseWriter, r *http.Request) {
	log := r.Context().Value(ctxKeyLog{}).(logrus.FieldLogger)
	// log.WithField("currency", currentCurrency(r)).Info("home")
	// currencies, err := fe.getCurrencies(r.Context())
	// if err != nil {
	// 	renderHTTPError(log, r, w, errors.Wrap(err, "could not retrieve currencies"), http.StatusInternalServerError)
	// 	return
	// }
	// products, err := fe.getProducts(r.Context())
	// if err != nil {
	// 	renderHTTPError(log, r, w, errors.Wrap(err, "could not retrieve products"), http.StatusInternalServerError)
	// 	return
	// }
	// cart, err := fe.getCart(r.Context(), sessionID(r))
	// if err != nil {
	// 	renderHTTPError(log, r, w, errors.Wrap(err, "could not retrieve cart"), http.StatusInternalServerError)
	// 	return
	// }

	// type productView struct {
	// 	Item  *pb.Product
	// 	Price *pb.Money
	// }
	// ps := make([]productView, len(products))
	// for i, p := range products {
	// 	price, err := fe.convertCurrency(r.Context(), p.GetPriceUsd(), currentCurrency(r))
	// 	if err != nil {
	// 		renderHTTPError(log, r, w, errors.Wrapf(err, "failed to do currency conversion for product %s", p.GetId()), http.StatusInternalServerError)
	// 		return
	// 	}
	// 	ps[i] = productView{p, price}
	// }

	log.WithField("path", "/")
	json.NewEncoder(w).Encode(map[string]string{"Status": "200"})
}

func (fe *frontendServer) nameHandler(w http.ResponseWriter, r *http.Request) {
	log := r.Context().Value(ctxKeyLog{}).(logrus.FieldLogger)
	name := mux.Vars(r)["name"]
	resp, err := fe.getNameLength(r.Context(), name)
	if err != nil {
		renderHTTPError(log, r, w, errors.Wrap(err, "Could not calculate name length"), http.StatusInternalServerError)
	}

	w.Header().Set("Content-type", "application/json")
	json.NewEncoder(w).Encode(resp)
}


func (fe *frontendServer) logoutHandler(w http.ResponseWriter, r *http.Request) {
	log := r.Context().Value(ctxKeyLog{}).(logrus.FieldLogger)
	log.Debug("logging out")
	for _, c := range r.Cookies() {
		c.Expires = time.Now().Add(-time.Hour * 24 * 365)
		c.MaxAge = -1
		http.SetCookie(w, c)
	}
	w.Header().Set("Location", "/")
	w.WriteHeader(http.StatusFound)
}



func renderHTTPError(log logrus.FieldLogger, r *http.Request, w http.ResponseWriter, err error, code int) {
	log.WithField("error", err).Error("request error")
	// fmt.Sprintf("%+v", err)

	w.WriteHeader(code)
}


func sessionID(r *http.Request) string {
	v := r.Context().Value(ctxKeySessionID{})
	if v != nil {
		return v.(string)
	}
	return ""
}
