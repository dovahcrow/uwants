package models

import (
	"github.com/astaxie/beego/orm"
	_ "github.com/mattn/go-sqlite3"
	. "tools"
)

type User struct {
	Id       int
	Username string
	Password string
}

var O orm.Ormer

func init() {
	orm.RegisterDataBase("default", "sqlite3", "database.db")
	orm.RegisterModel(new(User), new(Threads))
	orm.RunCommand()
	O = orm.NewOrm()
}

func GetAllUsers() (ret []*User, err error) {
	defer func() {
		if e := recover(); e != nil {
			err = e.(error)
			return
		}
	}()
	_, err = O.QueryTable(new(User)).All(&ret)
	E(`获取所有用户错误`, err)
	return
}
func GetUserById(id int) (ret *User, err error) {
	defer func() {
		if e := recover(); e != nil {
			err = e.(error)
			return
		}
	}()
	ret = &User{}
	ret.Id = id
	err = O.Read(ret)
	E(`根据id读取用户错误`, err)
	return
}

func CreateUser(username, password string) (err error) {
	defer func() {
		if e := recover(); e != nil {
			err = e.(error)
			return
		}
	}()
	ret := &User{}
	ret.Password = password
	ret.Username = username
	_, err = O.Insert(ret)
	E(`新建用户错误`, err)
	return
}

func UpdateUser(id int, username, password string) (err error) {
	defer func() {
		if e := recover(); e != nil {
			err = e.(error)
			return
		}
	}()
	ret := &User{}
	ret.Password = password
	ret.Username = username
	ret.Id = id
	_, err = O.Update(ret)
	E(`更新用户错误`, err)
	return
}
func DeleteUser(id int) (err error) {
	defer func() {
		if e := recover(); e != nil {
			err = e.(error)
			return
		}
	}()
	ret := &User{}
	ret.Id = id
	_, err = O.Delete(ret)
	E(`删除用户错误`, err)
	return
}

type Threads struct {
	Id    int
	Title string
	Text  string
}

func GetAllThreads() (ret []*Threads, err error) {
	defer func() {
		if e := recover(); e != nil {
			err = e.(error)
			return
		}
	}()
	_, err = O.QueryTable(new(Threads)).All(&ret)
	E(`获取所有帖子错误`, err)
	return
}
func GetThreadById(id int) (ret *Threads, err error) {
	defer func() {
		if e := recover(); e != nil {
			err = e.(error)
			return
		}
	}()
	ret = &Threads{}
	ret.Id = id
	err = O.Read(ret)
	E(`根据id读取帖子错误`, err)
	return
}

func CreateThread(title, text string) (err error) {
	defer func() {
		if e := recover(); e != nil {
			err = e.(error)
			return
		}
	}()
	ret := &Threads{}
	ret.Text = text
	ret.Title = title
	_, err = O.Insert(ret)
	E(`新建帖子错误`, err)
	return
}

func UpdateThread(id int, title, text string) (err error) {
	defer func() {
		if e := recover(); e != nil {
			err = e.(error)
			return
		}
	}()
	ret := &Threads{}
	ret.Text = text
	ret.Title = title
	ret.Id = id
	_, err = O.Update(ret)
	E(`更新帖子错误`, err)
	return
}
func DeleteThread(id int) (err error) {
	defer func() {
		if e := recover(); e != nil {
			err = e.(error)
			return
		}
	}()
	ret := &Threads{}
	ret.Id = id
	_, err = O.Delete(ret)
	E(`删除帖子错误`, err)
	return
}
