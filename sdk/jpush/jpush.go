// Package jpush 极光推送SDK
package jpush

import (
	"errors"
	"github.com/laiyinghate18/jpush-api-go-client"
)

// SendMessage 极光推送发送消息
// jpushAppKey 对应极光推送的AppKey
// jpushSecret 对应极光推送的Master Secret
// msgContent 消息标题
// msgTitle 消息内容
// registerIDs 可以指定推送的设备ID
// isSendAll 如果不想指定指定设备id，这里可以直接配置为true，那就是默认全发送
func SendMessage(jpushAppKey, jpushSecret, msgTitle, msgContent string, registerIDs []string, isSendAll bool, extras map[string]any) (string, error) {
	//Platform
	var pf jpushclient.Platform
	_ = pf.Add(jpushclient.ANDROID)
	_ = pf.Add(jpushclient.IOS)
	_ = pf.Add(jpushclient.WINPHONE)
	//pf.All()

	//Audience
	var ad jpushclient.Audience
	/*
		s := []string{"1", "2", "3"}
		ad.SetTag(s)
		ad.SetAlias(s)
		ad.SetID(s)
		ad.All()
	*/
	if isSendAll == true {
		// 管理员专用
		ad.All()
	} else {
		if len(registerIDs) > 0 {
			ad.SetID(registerIDs)
		} else {
			return "", errors.New("registerIDs Limit")
		}
	}

	//Notice
	var notice jpushclient.Notice
	notice.SetAlert(msgContent)
	notice.SetAndroidNotice(&jpushclient.AndroidNotice{Alert: msgContent, Title: msgTitle, Extras: extras, BadgeAddNum: 1, BadgeClass: "com.game.web.MainActivity"})

	// 针对ios单独设置
	iosAPNs := map[string]string{
		"title": msgTitle, //可选设置
		//"subtitle": "子标题:" + msgTitle, //可选设置
		"body": msgContent, //必填，否则通知将不展示，在不设置 title 和 subtitle 时直接对 alert 传值即可，不需要特地写 body 字段
	}
	notice.SetIOSNotice(&jpushclient.IOSNotice{Alert: iosAPNs, Badge: 1, Extras: extras}) // 这里设置1
	notice.SetWinPhoneNotice(&jpushclient.WinPhoneNotice{Alert: msgContent, Title: msgTitle, Extras: extras})

	var msg jpushclient.Message
	msg.Title = msgContent
	msg.Content = msgTitle
	msg.Extras = extras

	payload := jpushclient.NewPushPayLoad()
	payload.SetPlatform(&pf)
	payload.SetAudience(&ad)
	payload.SetMessage(&msg)
	payload.SetNotice(&notice)

	/**
	如果目标平台为 iOS 平台，推送 Notification 时需要在 options 中通过 apns_production 字段来设定推送环境。
	True 表示推送生产环境，False 表示要推送开发环境；

	如果不指定则为推送生产环境；一次只能推送给一个环境。
	*/
	options := jpushclient.Option{}
	options.ApnsProduction = true
	payload.SetOptions(&options)

	bytes, _ := payload.ToBytes()

	//push
	c := jpushclient.NewPushClient(jpushSecret, jpushAppKey)
	return c.Send(bytes)
}
