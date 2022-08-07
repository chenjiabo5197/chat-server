package model

import (
	"common"
	"encoding/json"
	"github.com/garyburd/redigo/redis"
	logger "github.com/shengkehua/xlog4go"
	"utils"
)

/*
	在服务器启动时，初始化一个userDao实例
*/
var (
	MyUserDao *UserDao
)

/*
	定义一个UserDao结构体
	完成对User结构体的crud操作
*/
type UserDao struct {
	pool *redis.Pool
}

/*
	使用工厂模式，创造一个UserDao实例
*/
func NewUserDao(pool *redis.Pool) (userDao *UserDao) {
	userDao = &UserDao{
		pool: pool,
	}
	return
}

/*
	定义一个方法，根据传入的redis连接和userid，返回user实例对象或错误
*/
func (up *UserDao) getUserById(conn redis.Conn, userId int) (user *common.User, err error) {

	//通过传入的id去redis查询
	res, err := redis.String(conn.Do("HGet", "users", userId))

	if err != nil {
		if err == redis.ErrNil { //在redis中没有找到相应对象
			err = ERROR_USER_NOTEXISTS
		}
		return
	}
	//将找到的对象反序列化
	err = json.Unmarshal([]byte(res), &user)

	if err != nil {
		logger.Error("redis data unmarshal err, err=%s", err.Error())
		return
	}
	return
}

/*
	根据传入的用户名和密码返回登录的结果
	如果登录成功，返回一个User对象，登陆失败，返回错误码
*/
func (up *UserDao) Login(userId int, userPwd string) (user *common.User, err error) {

	//从redis连接池中获取一个连接
	conn := up.pool.Get()

	defer conn.Close()

	user, err = up.getUserById(conn, userId)
	if err != nil { //证明数据库中不存在这个id的用户
		return
	}

	if user.UserPwd != userPwd { //用户密码错误F
		err = ERROR_USER_PWD
		return
	}
	//走到这里证明用户名正确，返回err=nil
	return
}

/*
	根据传入的User实例对象返回注册的结果
*/
func (up *UserDao) RegisterUser(user common.User) (err error) {

	//从redis连接池中获取一个连接
	conn := up.pool.Get()

	defer conn.Close()

	//user := message.User{}
	//err = json.Unmarshal([]byte(userString), &user)
	//if err != nil {
	//	fmt.Println("userString反序列化失败")
	//	return
	//}
	logger.Debug("user=%s", utils.Struct2String(user))

	_, err = up.getUserById(conn, user.UserId)

	if err == nil { //证明在数据库中存在这个id的用户
		err = ERROR_USER_EXISTS
		return
	}

	//走到这里证明数据库中不存在此id的用户，将user转为string
	data, err := json.Marshal(user)
	_, err = conn.Do("HSet", "users", user.UserId, string(data))
	if err != nil {
		logger.Error("save rigister mes err, err=%s", err.Error())
		return
	}

	//走到这里证明用户名正确，返回err=nil
	return
}
