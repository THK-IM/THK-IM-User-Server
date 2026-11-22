package handler

import (
	"bytes"
	"encoding/json"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/thk-im/thk-im-base-server/server"
	msg "github.com/thk-im/thk-im-msgapi-server/pkg/dto"
	"github.com/thk-im/thk-im-user-server/pkg/app"
	"github.com/thk-im/thk-im-user-server/pkg/model"
)

type presenceRecorder struct {
	model.UserOnlineRecordModel
	events []msg.UserOnlineStatusReq
}

func (r *presenceRecorder) CreateUserOnlineRecord(uid, conn int64, platform string, ts int64) error {
	r.events = append(r.events, msg.UserOnlineStatusReq{UserId: uid, ConnId: conn, Platform: platform, TimestampMs: ts, IsOnline: true})
	return nil
}
func (r *presenceRecorder) UpdateUserOnlineRecord(uid, conn int64, platform string, ts int64) error {
	r.events = append(r.events, msg.UserOnlineStatusReq{UserId: uid, ConnId: conn, Platform: platform, TimestampMs: ts})
	return nil
}
func TestPresenceCallbackAcceptsMsgAPIPayload(t *testing.T) {
	recorder := &presenceRecorder{}
	ctx := &app.Context{Context: &server.Context{ModelMap: map[string]interface{}{"user_online_record": recorder}}}
	router := gin.New()
	router.POST("/user/online_status", postUserOnlineStatus(ctx))
	for i, online := range []bool{true, false} {
		event := msg.UserOnlineStatusReq{UserId: 12, ConnId: 25, Platform: "ios", TimestampMs: int64(100 + i), IsOnline: online}
		data, err := json.Marshal(event)
		if err != nil {
			t.Fatal(err)
		}
		w := httptest.NewRecorder()
		router.ServeHTTP(w, httptest.NewRequest("POST", "/user/online_status", bytes.NewReader(data)))
		if w.Code != 200 || len(recorder.events) != i+1 || recorder.events[i] != event {
			t.Fatalf("status %d, received %+v", w.Code, recorder.events)
		}
	}
}
