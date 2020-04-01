package beater

import (
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"time"

	"github.com/elastic/beats/libbeat/common"
	"github.com/elastic/beats/libbeat/logp"
	"github.com/julienschmidt/httprouter"
)

type hapiMessage struct {
	Host   string `json:"host"`
	Schema string `json:"schema"`
	// TODO: If http-good is not aggregating events
	// this could be better to use instead of trying
	// to access the "timestamp" field
	//Timestamp int64             `json:"timeStamp"`
	Events []json.RawMessage `json"events"`
}

type HapiMessage struct {
	Host   string `json:"host"`
	Schema string `json:"schema"`
	//Timestamp time.Time        `json:"timeStamp"`
	Events []*common.MapStr `json"events"`
}

func (m *HapiMessage) UnmarshalJSON(raw []byte) error {
	msg := &hapiMessage{}
	err := json.Unmarshal(raw, msg)
	if err != nil {
		return err
	}
	m.Host = msg.Host
	m.Schema = msg.Schema
	m.Events = []*common.MapStr{}
	for _, eventRaw := range msg.Events {
		evt := &common.MapStr{}
		err = json.Unmarshal([]byte(eventRaw), evt)
		if err != nil {
			break
		}
		m.Events = append(m.Events, evt)
	}
	return err
}

type Server struct {
	auth   TokenAuth
	events chan *common.MapStr
}

type HandleFunc func(http.ResponseWriter, *http.Request, httprouter.Params) error

func HandleWrapper(s *Server, fn HandleFunc) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		// verify auth token, if enabled
		if s.auth.Enabled() && !s.auth.Verify(r) {
			logp.Warn("request token verification failed")
			w.WriteHeader(http.StatusForbidden)
			return
		}
		err := fn(w, r, p)
		if err != nil {
			logp.Warn("handler responded with err %s", err.Error())
			http.Error(w, err.Error(), 500)
		}
		logp.Info("%s %s %s", r.URL.String(), r.Method, r.Header.Get("Host"))
	}
}

func (s *Server) Events(w http.ResponseWriter, r *http.Request, p httprouter.Params) error {
	msg := &HapiMessage{}
	err := json.NewDecoder(r.Body).Decode(msg)
	if err != nil {
		return err
	}
	for _, payload := range msg.Events {
		// ensure data field is always an object, if exists
		d, err := payload.GetValue("data")
		if err == nil {
			if dStr, ok := d.(string); ok {
				logp.Info("enclosing string field in object: %s", dStr)
				payload.Put("data", map[string]string{"string": dStr})
			}
		}

		// parse timestamp
		t, err := payload.GetValue("timestamp")
		if err != nil {
			return fmt.Errorf("cannot find timestamp")
		}
		v, ok := t.(float64)
		if !ok {
			return fmt.Errorf("cannot cast timestamp to float")
		}
		payload.Put("@timestamp", common.Time(readTime(v)))

		logp.Debug("parsed event payload: %s", payload.StringToPrint())
		s.events <- payload
	}
	return nil
}

func readTime(v float64) time.Time {
	// NodeJS Date has millisecond precision, so round to the second
	v = v / 1000
	return time.Unix(int64(v+math.Copysign(0.5, v)), 0)
}

func Run(auth TokenAuth, events chan *common.MapStr) error {
	server := &Server{auth, events}
	if auth.Enabled() {
		logp.Info("token auth enabled: %s", auth.token)
	}
	router := httprouter.New()
	router.POST("/events", HandleWrapper(server, server.Events))
	logp.Info("listening @0.0.0.0:9090")
	return http.ListenAndServe("0.0.0.0:9090", router)
}
