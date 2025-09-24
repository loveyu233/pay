package pay

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/ArtisanCloud/PowerWeChat/v3/src/kernel/contract"
	"github.com/ArtisanCloud/PowerWeChat/v3/src/kernel/messages"
	models2 "github.com/ArtisanCloud/PowerWeChat/v3/src/kernel/models"
	"github.com/ArtisanCloud/PowerWeChat/v3/src/kernel/power"
	"github.com/ArtisanCloud/PowerWeChat/v3/src/officialAccount/broadcasting/request"
	"github.com/ArtisanCloud/PowerWeChat/v3/src/officialAccount/server/handlers/models"
	tRequest "github.com/ArtisanCloud/PowerWeChat/v3/src/officialAccount/templateMessage/request"
	"github.com/ArtisanCloud/PowerWeChat/v3/src/officialAccount/templateMessage/response"
	"github.com/gin-gonic/gin"
	"github.com/loveyu233/gb"
)

func (wx *wxOfficial) RegisterHandlers(r *gin.RouterGroup) {
	r.Use(gb.GinLogSetModuleName("微信公众号"))
	r.GET("/wx/callback", gb.GinLogSetOptionName("回调验证", wx.IsSaveHandlerLog), wx.oACallbackVerify)
	r.POST("/wx/callback", gb.GinLogSetOptionName("收到消息", wx.IsSaveHandlerLog), wx.oACallback)
	r.POST("/wx/push", gb.GinLogSetOptionName("推送消息", wx.IsSaveHandlerLog), wx.pushHand)
}

func (wx *wxOfficial) oACallbackVerify(c *gin.Context) {
	//回调验证
	rs, err := wx.OfficialAccountApp.Server.VerifyURL(c.Request)
	if err != nil {
		panic(err)
	}
	text, _ := io.ReadAll(rs.Body)
	c.String(http.StatusOK, string(text))
}

func (wx *wxOfficial) oACallback(c *gin.Context) {
	rs, err := wx.OfficialAccountApp.Server.Notify(c.Request, func(event contract.EventInterface) interface{} {
		fmt.Println("event", event)
		switch event.GetMsgType() {
		case models2.CALLBACK_MSG_TYPE_TEXT:
			//收到用户的消息
			msg := models.MessageText{}
			err := event.ReadMessage(&msg)
			if err != nil {
				println(err.Error())
				return "error"
			}

		case models2.CALLBACK_MSG_TYPE_EVENT:
			fmt.Println(event.GetToUserName(), event.GetFromUserName())
			rs, _ := wx.OfficialAccountApp.User.Get(context.Background(), event.GetFromUserName(), "zh_CN")
			switch event.GetEvent() {
			case "subscribe": // 关注
				if rs.OpenID != "" {
					wx.subscribe(rs, event)
				}
				// 这里回复success告诉微信我收到,后续需要回复用户信息可以主动调发消息接口
				return messages.NewText("感谢您的关注！")
			case "unsubscribe": // 取消关注
				wx.unSubscribe(rs, event)
			}
		}
		return ""
	})

	if err != nil {
		c.String(200, err.Error())
		return
	}

	text, _ := io.ReadAll(rs.Body)
	c.String(http.StatusOK, string(text))
}

func (wx *wxOfficial) pushHand(c *gin.Context) {
	users, message := wx.pushHandler(c)
	_, err := wx.push(users, message)
	if err != nil {
		gb.ResponseError(c, gb.ErrRequestWechat.WithMessage("推送消息失败:%s", err.Error()))
		return
	}
	gb.ResponseSuccess(c, nil)
}

func (wx *wxOfficial) push(toUser []string, message string) (interface{}, error) {
	if len(toUser) == 1 {
		//至少两个才能发送成功 添加一个空id
		toUser = append(toUser, "")
	}
	d, err := wx.OfficialAccountApp.Broadcasting.SendText(context.Background(), message, &request.Reception{
		ToUser: toUser,
		Filter: &request.Filter{
			IsToAll: false,
			TagID:   0,
		},
	}, &power.HashMap{})
	if err != nil {
		return nil, err
	}

	return d, err
}

func (wx *wxOfficial) PushTemplateMessage(toUser, templateId string, data interface{}) (*response.ResponseTemplateSend, error) {
	d := make(power.HashMap)
	dataMap, err := structToMapWithJSONTag(data)
	if err != nil {
		return nil, err
	}
	for i, v := range dataMap {
		d[i] = &power.HashMap{
			"value": v,
			"color": "#173177",
		}
	}
	send, err := wx.OfficialAccountApp.TemplateMessage.Send(context.Background(), &tRequest.RequestTemlateMessage{
		ToUser:     toUser,
		TemplateID: templateId,
		Data:       &d,
	})
	return send, err
}

func structToMapWithJSONTag(obj interface{}) (map[string]interface{}, error) {
	data, err := json.Marshal(obj)
	if err != nil {
		return nil, err
	}

	var m map[string]interface{}
	if err := json.Unmarshal(data, &m); err != nil {
		return nil, err
	}
	return m, nil
}
