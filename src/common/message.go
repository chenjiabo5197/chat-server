package common

/*
	此包客户端和服务端功用一套
*/

type MesType string

const (
	LoginMesType            MesType = "LoginMes"
	LoginRespMesType        MesType = "LoginRespMes"
	RegisterMesType         MesType = "RegisterMes"
	RegisterRespMesType     MesType = "RegisterRespMes"
	NotifyUserStatusMesType MesType = "NotifyUserStatusMes"
	SmsMesType              MesType = "SmsMes" // 广播通信
	SmsRespMesType          MesType = "SmsRespMes"
	QueryAllOnlineType      MesType = "QueryAllOnline"
	AllOnlineRespType       MesType = "AllOnlineResp"
	SmsToOneMesType         MesType = "SmsToOneMes" // 1对1通信
	SmsToOneRespMesType     MesType = "SmsToOneRespMes"
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
	UserId   int    `json:"userid"`
	UserPwd  string `json:"userpwd"`
	UserName string `json:"username"`
}

// LoginRespMes 服务器端返回的登录的结果消息
type LoginRespMes struct {
	RespCode int    `json:"respcode"`
	Error    string `json:"error"`
	UsersId  int
	UserName string
}

// RegisterMes 客户端发送的注册消息
type RegisterMes struct {
	User User `json:"user"`
}

/*
	定义一个用户的结构体，
*/
type User struct {
	UserId     int    `json:"userid"`
	UserPwd    string `json:"userpwd"`
	UserName   string `json:"username"`
	UserStatus int    `json:"userstatus"`
}

// RegisterRespMes 服务器端返回的注册的结果消息
type RegisterRespMes struct {
	RespCode int    `json:"respcode"`
	Error    string `json:"error"`
}

// NotifyUserStatusMes 服务器用于推送用户状态变化的消息
type NotifyUserStatusMes struct {
	UserId     int    `json:"userid"`
	UserStatus int    `json:"userstatus"`
	UserName   string `json:"username"`
}

// SmsMes 客户端发送聊天消息的结构体
type SmsMes struct {
	User           //匿名结构体
	Content string `json:"content"`
	Target  string `json:"target"` // 消息发送的对象
}

// SmsRespMes 服务器转发聊天消息的结构体
type SmsRespMes struct {
	User           //匿名结构体
	Content string `json:"content"`
}
