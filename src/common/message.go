package common

/*
	此包客户端和服务端功用一套
*/

type MesType string

const (
	LoginMesType               MesType = "LoginMes"
	LoginRespMesType           MesType = "LoginRespMes"
	RegisterMesType            MesType = "RegisterMes"
	RegisterRespMesType        MesType = "RegisterRespMes"
	NotifyUserStatusMesType    MesType = "NotifyUserStatusMes"
	SmsMesType                 MesType = "SmsMes"            // 广播通信
	RecvSmsMesType             MesType = "RecvSmsMes"        // 在线收到的广播通信
	OfflineRecvSmsMesType      MesType = "OfflineRecvSmsMes" // 离线收到的广播通信
	QueryAllOnlineType         MesType = "QueryAllOnline"
	AllOnlineRespType          MesType = "AllOnlineResp"
	SmsToOneMesType            MesType = "SmsToOneMes"            // 1对1通信
	RecvSmsToOneMesType        MesType = "RecvSmsToOneMes"        // 在线收到的1对1通信
	OfflineRecvSmsToOneMesType MesType = "OfflineRecvSmsToOneMes" // 离线收到的1对1通信
	SmsRespMesType             MesType = "SmsRespMes"             // 发送1对1与群发消息的结果 ,失败或成功
)

//定义几个用户在线状态
const (
	UserOnline = iota
	UserOffline
	UserBusy
)

type Message struct {
	Type MesType `json:"type"` //定义消息类型
	Data string  `json:"data"` //消息内容
}

// LoginMes 客户端发送的登录消息
type LoginMes struct {
	UserId   string `json:"user_id"`
	UserPwd  string `json:"user_pwd"`
	UserName string `json:"user_name"`
}

// LoginRespMes 服务器端返回的登录的结果消息
type LoginRespMes struct {
	RespCode int    `json:"resp_code"`
	Error    string `json:"error"`
	User     User   `json:"user"`
}

// RegisterMes 客户端发送的注册消息
type RegisterMes struct {
	User User `json:"user"`
}

// User 定义一个用户的结构体,包含一个用户的所有信息
type User struct {
	UserId     string   `json:"user_id"`
	UserPwd    string   `json:"user_pwd"`
	UserName   string   `json:"user_name"`
	UserStatus int      `json:"user_status"` // 用户的在线状态
	UserFriend []string `json:"user_friend"` // 该用户的好友列表
	UserGroup  []string `json:"user_group"`  // 该用户的群组列表
}

// OnlineUserInfo 储存在线用户信息
type OnlineUserInfo struct {
	OnlineUserId   string `json:"online_user_id"`
	OnLineUserName string `json:"on_line_user_name"`
}

// StatusRespMes 服务器端返回操作的结果消息
type StatusRespMes struct {
	RespCode int    `json:"resp_code"`
	Error    string `json:"error"`
}

// NotifyUserStatusMes 服务器用于推送用户状态变化的消息
type NotifyUserStatusMes struct {
	UserId     string `json:"userid"`
	UserStatus int    `json:"user_status"`
	UserName   string `json:"user_name"`
}

// SmsMes 客户端发送聊天消息的结构体
type SmsMes struct {
	User                //匿名结构体
	Content      string `json:"content"`
	SmsMesTarget string `json:"sms_mes_target"` // 1对1消息发送的对象
}

// SmsRespMes 服务器转发聊天消息的结构体
type SmsRespMes struct {
	User               //匿名结构体
	Content     string `json:"content"`
	SmsRespFrom string `json:"sms_resp_from"` // 群发消息的来源
}
