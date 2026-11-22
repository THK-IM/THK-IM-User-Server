package handler

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/gin-gonic/gin"
	base "github.com/thk-im/thk-im-base-server/dto"
	msg "github.com/thk-im/thk-im-msgapi-server/pkg/dto"
	msgErr "github.com/thk-im/thk-im-msgapi-server/pkg/errorx"
	"github.com/thk-im/thk-im-msgapi-server/pkg/model"
	"github.com/thk-im/thk-im-msgapi-server/pkg/sdk"
)

type contactCheckFunc func(*msg.CheckMessageReq, base.ThkClaims) (*msg.CheckMessageRes, error)

func (f contactCheckFunc) CheckMessage(req *msg.CheckMessageReq, claims base.ThkClaims) (*msg.CheckMessageRes, error) {
	return f(req, claims)
}

func TestMessageCheckRouting(t *testing.T) {
	for _, tc := range []struct {
		name                     string
		sessionType, messageType int
		sender                   int64
		expectContact            bool
		contactError             error
		status                   int
	}{
		{"single", model.SingleSessionType, 5, 10, true, nil, http.StatusOK},
		{"blocked", model.SingleSessionType, 5, 10, true, msgErr.ErrUserReject, http.StatusBadRequest},
		{"contact unavailable", model.SingleSessionType, 1, 10, true, errors.New("contact unavailable"), http.StatusInternalServerError},
		{"group", model.GroupSessionType, 5, 10, false, nil, http.StatusOK},
		{"supergroup", model.SuperGroupSessionType, 5, 10, false, nil, http.StatusOK},
		{"system sender", model.SingleSessionType, 1, 0, false, nil, http.StatusOK},
		{"control message", model.SingleSessionType, -1, 10, false, nil, http.StatusOK},
	} {
		t.Run(tc.name, func(t *testing.T) {
			req := msg.CheckMessageReq{SessionType: tc.sessionType, SessionId: 123, EntityId: 20, FromUId: tc.sender, MessageType: tc.messageType, MessageContent: "test"}
			calls := 0
			contact := contactCheckFunc(func(got *msg.CheckMessageReq, claims base.ThkClaims) (*msg.CheckMessageRes, error) {
				calls++
				if *got != req || claims.GetTraceId() != "test-trace" {
					t.Fatalf("request or trace changed: %+v", got)
				}
				return &msg.CheckMessageRes{ExcludeSaveMsgUserIds: []int64{42}}, tc.contactError
			})
			w := postCheck(t, contact, req)
			if w.Code != tc.status || (tc.expectContact && calls != 1) || (!tc.expectContact && calls != 0) {
				t.Fatalf("status %d, contact calls %d, body %s", w.Code, calls, w.Body)
			}
			if tc.status == http.StatusOK {
				var result msg.CheckMessageRes
				if err := json.Unmarshal(w.Body.Bytes(), &result); err != nil {
					t.Fatal(err)
				}
				var want []int64
				if tc.expectContact {
					want = []int64{42}
				}
				if !reflect.DeepEqual(result.ExcludeSaveMsgUserIds, want) {
					t.Fatal(result)
				}
			}
			if tc.contactError == msgErr.ErrUserReject {
				var result struct{ Code int }
				if err := json.Unmarshal(w.Body.Bytes(), &result); err != nil || result.Code != msgErr.ErrUserReject.Code {
					t.Fatalf("lost rejection code: %s", w.Body)
				}
			}
		})
	}
}

func postCheck(t *testing.T, contact sdk.MsgCheckerApi, req msg.CheckMessageReq) *httptest.ResponseRecorder {
	t.Helper()
	router := gin.New()
	router.POST("/system/message/check", func(c *gin.Context) { c.Set(base.ClaimsKey, base.ThkClaims{base.TraceID: "test-trace"}); c.Next() }, checkMessage(contact))
	data, err := json.Marshal(req)
	if err != nil {
		t.Fatal(err)
	}
	w := httptest.NewRecorder()
	router.ServeHTTP(w, httptest.NewRequest(http.MethodPost, "/system/message/check", bytes.NewReader(data)))
	return w
}
