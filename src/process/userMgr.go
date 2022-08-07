package processes

import (
	"fmt"
)

/*
	UserMgr实例在服务器上有且只有一个，且在很多地方会使用到，故定义为全局变量
*/
var (
	Usermgr *UserMgr
)

/*
	该结构体储存的map类型保存所有在线的用户，key为用户id，value为UserProcess的指针，使用该value可以拿到该用户对应的net.Conn
*/
type UserMgr struct {
	OnlineUsers map[int]*UserProcess
}

//完成对userMgr的初始化工作
func init() {
	Usermgr = &UserMgr{
		OnlineUsers: make(map[int]*UserProcess, 0),
	}
}

// AddOnlineUsers 完成对onlineUsers的添加
func (um *UserMgr) AddOnlineUsers(up *UserProcess) {
	um.OnlineUsers[up.UserId] = up
}

// DeleteOnlineUsers 完成对onlineUsers的删除
func (um *UserMgr) DeleteOnlineUsers(userId int) {
	delete(um.OnlineUsers, userId)
}

// GetAllOnlineUsers 返回所有的在线用户
func (um *UserMgr) GetAllOnlineUsers() map[int]*UserProcess {
	return um.OnlineUsers
}

// GetUserProcessById 根据传入的UserId返回当前的在线用户的UserProcess指针
func (um *UserMgr) GetUserProcessById(userId int) (up *UserProcess, err error) {
	up, ok := um.OnlineUsers[userId]

	if !ok {
		err = fmt.Errorf("用户%d不在线", userId)
		return
	}

	return
}
