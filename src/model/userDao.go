package model

import (
	"common"
	"encoding/json"
	"fmt"
	"github.com/garyburd/redigo/redis"
	logger "github.com/shengkehua/xlog4go"
	"reflect"
	"utils"
)

/*
	在服务器启动时，初始化一个userDao实例
*/
var (
	MyUserDao *UserDao
)

const (
	USER_REDIS_PREFIX_KEY    = "chat_service_user_"
	USRR_ONLINE_KEY          = "chat_service_online_user"
	USER_OFFLINE_MESSAGE_KEY = "chat_service_offline_message"
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

func GetRedisUserKey(userId string) string {
	return fmt.Sprintf(USER_REDIS_PREFIX_KEY + userId)
}

//根据传入的key，返回user实例对象或错误
func (up *UserDao) GetDataByKey(key string) (user *common.User, err error) {
	//从redis连接池中获取一个连接
	redisClient := up.pool.Get()
	defer redisClient.Close()
	//通过传入的id去redis查询
	//res, err := redis.String(redisClient.Do("HGet", "users", userId))
	res, err := redis.String(redisClient.Do("Get", key))

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

//根据传入的key和data，往redis string类型中塞入数据
func (up *UserDao) SetDataByKey(key string, data string) (err error) {
	//从redis连接池中获取一个连接
	redisClient := up.pool.Get()
	defer redisClient.Close()

	_, err = redisClient.Do("Set", key, data)
	if err != nil {
		logger.Error("save rigister mes err, err=%s", err.Error())
		return
	}
	return
}

// HSetDataByName 根据传入的username，储存发给该用户的离线消息到redis中 hash类型
func (up *UserDao) HSetDataByName(userName string, mesResp common.Message) (err error) {
	//从redis连接池中获取一个连接
	redisClient := up.pool.Get()
	defer redisClient.Close()
	//先将redis中目标的离线消息拿到手，然后再增加
	var data []common.Message
	data, err = up.HGetDataByName(userName)
	if err != nil && err != redis.ErrNil {
		logger.Error("err=%v||type=%v", err, reflect.TypeOf(err))
		return
	}
	data = append(data, mesResp)
	dataByte, _ := json.Marshal(data)
	_, err = redisClient.Do("HSet", USER_OFFLINE_MESSAGE_KEY, userName, string(dataByte))
	if err != nil {
		logger.Error("hset offline message to redis err, err=%v", err)
		return
	}
	return
}

// HGetDataByName 根据传入的username，返回该用户目前的离线消息
func (up *UserDao) HGetDataByName(userName string) (data []common.Message, err error) {
	//从redis连接池中获取一个连接
	redisClient := up.pool.Get()
	defer redisClient.Close()
	//先将redis中目标的离线消息拿到手，然后再增加
	content, err := redis.String(redisClient.Do("HGet", USER_OFFLINE_MESSAGE_KEY, userName))
	if err != nil {
		logger.Error("HGet data from chat_service_offline_message err, err=%v", err)
		return
	}
	err = json.Unmarshal([]byte(content), &data)
	if err != nil {
		logger.Error("offline message unmarshal err, err=%v", err)
		return
	}
	return
}

// HDelDataByName 根据传入的用户名，删除该用户的离线消息
func (up *UserDao) HDelDataByName(userName string) (err error) {
	//从redis连接池中获取一个连接
	redisClient := up.pool.Get()
	defer redisClient.Close()
	_, err = redisClient.Do("HDel", USER_OFFLINE_MESSAGE_KEY, userName)
	if err != nil {
		logger.Error("hdel offline message err, err=%v", err)
		return
	}
	return
}

/*
	根据传入的用户名和密码返回登录的结果
	如果登录成功，返回一个User对象，登陆失败，返回错误码
*/
func (up *UserDao) Login(userId string, userPwd string) (user *common.User, err error) {

	//从redis连接池中获取一个连接
	redisClient := up.pool.Get()
	defer redisClient.Close()

	user, err = up.GetDataByKey(USER_REDIS_PREFIX_KEY + userId)
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

//根据传入的User实例对象返回注册的结果
func (up *UserDao) RegisterUser(user common.User) (err error) {
	logger.Debug("register||user=%s", utils.Struct2String(user))
	_, err = up.GetDataByKey(USER_REDIS_PREFIX_KEY + user.UserId)
	if err == nil { //证明在数据库中存在这个id的用户
		err = ERROR_USER_EXISTS
		return
	}
	//走到这里证明数据库中不存在此id的用户，将user转为string
	data, err := json.Marshal(user)
	err = up.SetDataByKey(USER_REDIS_PREFIX_KEY+user.UserId, string(data))
	if err != nil {
		logger.Error("save register mes err, err=%s", err.Error())
		return
	}
	//走到这里证明用户名正确，返回err=nil
	logger.Info("register||save user message success")
	return
}
