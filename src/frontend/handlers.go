// Copyright 2018 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

	pb "github.com/samcfinan/microservices-demo/src/frontend/genproto"
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

	// if err := templates.ExecuteTemplate(w, "home", map[string]interface{}{
	// 	"session_id":    sessionID(r),
	// 	"request_id":    r.Context().Value(ctxKeyRequestID{}),
	// 	"user_currency": currentCurrency(r),
	// 	"currencies":    currencies,
	// 	"products":      ps,
	// 	"cart_size":     len(cart),
	// 	"banner_color":  os.Getenv("BANNER_COLOR"), // illustrates canary deployments
	// 	"ad":            fe.chooseAd(r.Context(), []string{}, log),
	// }); err != nil {
	// 	log.Error(err)
	// }
	log.WithField("path", "/")
	json.NewEncoder(w).Encode(map[string]string{"Status": "200"})
}

func (fe *frontendServer) nameHandler(w http.ResponseWriter, r *http.Request) {
	log := r.Context().Value(ctxKeyLog{}).(logrus.FieldLogger)
	fmt.Println("Herehere")
	name := mux.Vars(r)["name"]
	resp, err := fe.getNameLength(r.Context(), name)
	if err != nil {
		renderHTTPError(log, r, w, errors.Wrap(err, "Could not calculate name length"), http.StatusInternalServerError)
	}

	fmt.Println(name)
	json.NewEncoder(w).Encode(resp)
}

func (fe *frontendServer) productHandler(w http.ResponseWriter, r *http.Request) {
	// log := r.Context().Value(ctxKeyLog{}).(logrus.FieldLogger)
	// id := mux.Vars(r)["id"]
	// if id == "" {
	// 	renderHTTPError(log, r, w, errors.New("product id not specified"), http.StatusBadRequest)
	// 	return
	// }
	// log.WithField("id", id).WithField("currency", currentCurrency(r)).
	// 	Debug("serving product page")

	// p, err := fe.getProduct(r.Context(), id)
	// if err != nil {
	// 	renderHTTPError(log, r, w, errors.Wrap(err, "could not retrieve product"), http.StatusInternalServerError)
	// 	return
	// }
	// currencies, err := fe.getCurrencies(r.Context())
	// if err != nil {
	// 	renderHTTPError(log, r, w, errors.Wrap(err, "could not retrieve currencies"), http.StatusInternalServerError)
	// 	return
	// }

	// cart, err := fe.getCart(r.Context(), sessionID(r))
	// if err != nil {
	// 	renderHTTPError(log, r, w, errors.Wrap(err, "could not retrieve cart"), http.StatusInternalServerError)
	// 	return
	// }

	// price, err := fe.convertCurrency(r.Context(), p.GetPriceUsd(), currentCurrency(r))
	// if err != nil {
	// 	renderHTTPError(log, r, w, errors.Wrap(err, "failed to convert currency"), http.StatusInternalServerError)
	// 	return
	// }

	// recommendations, err := fe.getRecommendations(r.Context(), sessionID(r), []string{id})
	// if err != nil {
	// 	renderHTTPError(log, r, w, errors.Wrap(err, "failed to get product recommendations"), http.StatusInternalServerError)
	// 	return
	// }

	// product := struct {
	// 	Item  *pb.Product
	// 	Price *pb.Money
	// }{p, price}

}

func (fe *frontendServer) addToCartHandler(w http.ResponseWriter, r *http.Request) {
	log := r.Context().Value(ctxKeyLog{}).(logrus.FieldLogger)
	quantity, _ := strconv.ParseUint(r.FormValue("quantity"), 10, 32)
	productID := r.FormValue("product_id")
	if productID == "" || quantity == 0 {
		renderHTTPError(log, r, w, errors.New("invalid form input"), http.StatusBadRequest)
		return
	}
	log.WithField("product", productID).WithField("quantity", quantity).Debug("adding to cart")

	p, err := fe.getProduct(r.Context(), productID)
	if err != nil {
		renderHTTPError(log, r, w, errors.Wrap(err, "could not retrieve product"), http.StatusInternalServerError)
		return
	}

	if err := fe.insertCart(r.Context(), sessionID(r), p.GetId(), int32(quantity)); err != nil {
		renderHTTPError(log, r, w, errors.Wrap(err, "failed to add to cart"), http.StatusInternalServerError)
		return
	}
	w.Header().Set("location", "/cart")
	w.WriteHeader(http.StatusFound)
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

func (fe *frontendServer) setCurrencyHandler(w http.ResponseWriter, r *http.Request) {
	log := r.Context().Value(ctxKeyLog{}).(logrus.FieldLogger)
	cur := r.FormValue("currency_code")
	log.WithField("curr.new", cur).WithField("curr.old", currentCurrency(r)).
		Debug("setting currency")

	if cur != "" {
		http.SetCookie(w, &http.Cookie{
			Name:   cookieCurrency,
			Value:  cur,
			MaxAge: cookieMaxAge,
		})
	}
	referer := r.Header.Get("referer")
	if referer == "" {
		referer = "/"
	}
	w.Header().Set("Location", referer)
	w.WriteHeader(http.StatusFound)
}

// chooseAd queries for advertisements available and randomly chooses one, if
// available. It ignores the error retrieving the ad since it is not critical.
func (fe *frontendServer) chooseAd(ctx context.Context, ctxKeys []string, log logrus.FieldLogger) *pb.Ad {
	ads, err := fe.getAd(ctx, ctxKeys)
	if err != nil {
		log.WithField("error", err).Warn("failed to retrieve ads")
		return nil
	}
	return ads[rand.Intn(len(ads))]
}

func renderHTTPError(log logrus.FieldLogger, r *http.Request, w http.ResponseWriter, err error, code int) {
	log.WithField("error", err).Error("request error")
	fmt.Sprintf("%+v", err)

	w.WriteHeader(code)
}

func currentCurrency(r *http.Request) string {
	c, _ := r.Cookie(cookieCurrency)
	if c != nil {
		return c.Value
	}
	return defaultCurrency
}

func sessionID(r *http.Request) string {
	v := r.Context().Value(ctxKeySessionID{})
	if v != nil {
		return v.(string)
	}
	return ""
}

func cartIDs(c []*pb.CartItem) []string {
	out := make([]string, len(c))
	for i, v := range c {
		out[i] = v.GetProductId()
	}
	return out
}
