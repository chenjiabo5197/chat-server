package common

const(
	LoginMesType 			= "LoginMes"
	LoginResMesType 		= "LoginResMes"
	RegisterMesType 		= "RegisterMes"
	RegisterResMesType 		= "RegisterResMes"
	NotifyUserStatusMesType = "NotifyUserStatusMes"
	SmsMesType				= "SmsMes"
	SmsResMesType			= "SmsResMes"
)

//定义几个用户在线状态
const(
	UserOnline    = iota
	UserOffline
	Userbusy
)

type Message struct{
	Type string `json:"type"` //定义消息类型
	Data string	`json:"data"` //消息内容
}

/*
	客户端发送的登录消息
*/
type LoginMes struct{
	UserId int		`json:"userid"`
	UserPwd string	`json:"userpwd"`
	UserName string `json:"username"`
}

/*
	服务器端返回的登录的结果消息
*/
type LoginResMes struct{
	ResCode int		`json:"rescode"`
	Error string	`json:"error"`
	UsersId []int	
}

/*
	客户端发送的注册消息
*/
type RegisterMes struct{
	User User `json:"user"`
}

/*
	定义一个用户的结构体，
*/
type User struct{
	UserId int 			`json:"userid"`
	UserPwd string		`json:"userpwd"`
	UserName string		`json:"username"`
	UserStatus int 		`json:"userstatus"`
}

/*
	服务器端返回的注册的结果消息
*/
type RegisterResMes struct{
	ResCode int		`json:"rescode"`
	Error string	`json:"error"`
}

/*
	服务器用于推送用户状态变化的消息
*/
type NotifyUserStatusMes struct {
	UserId int 		`json:"userid"`
	UserStatus int 	`json:"userstatus"`
	UserName string	`json:"username"`
}

/*
	客户端发送聊天消息的结构体
*/
type SmsMes struct {
	User  //匿名结构体
	Content string `json:"content"`
}

/*
	服务器转发聊天消息的结构体
*/
type SmsResMes struct {
	User  //匿名结构体
	Content string `json:"content"`
}
