package dao

import (
	"fmt"
	"log"
	"wifidog-server/framework"
	"wifidog-server/model"
)

type UserDao struct {
}

func (p *UserDao) SelectUserByUsername(username string) []model.Useraccount {
	rows, err := framework.DB.Query("SELECT username, password, email, phnumber, department, humanname, role FROM useraccount WHERE username=?", username)
	if err != nil {
		fmt.Println("selectuserbyname error1", err)
		return nil
	}

	var useraccounts []model.Useraccount
	for rows.Next() {
		var useraccount model.Useraccount
		err := rows.Scan(&useraccount.Name, &useraccount.Password,
			&useraccount.Email, &useraccount.PhoneNumber, &useraccount.Department, &useraccount.Human, &useraccount.Role)
		if err != nil {
			fmt.Println("selectuserbyname error2")
			continue
		}
		useraccounts = append(useraccounts, useraccount)
	}
	rows.Close()
	return useraccounts
}

func (p *UserDao) SelectUserByPhoneNumber(phonenumber string) []model.Useraccount {
	rows, err := framework.DB.Query("SELECT username,password,email,phnumber,department,humanname FROM useraccount WHERE phnumber=?", phonenumber)
	if err != nil {
		fmt.Println("get user account by phone number from database error")
		return nil
	}

	var useraccounts []model.Useraccount
	for rows.Next() {
		var useraccount model.Useraccount
		err := rows.Scan(&useraccount.Name, &useraccount.Password,
			&useraccount.Email, &useraccount.PhoneNumber, &useraccount.Department, &useraccount.Human)
		if err != nil {
			fmt.Println("scan error in SelectUserByPhoneNumber")
			continue
		}
		useraccounts = append(useraccounts, useraccount)
	}
	rows.Close()
	return useraccounts
}
func (p *UserDao) GetUserAll() []model.Useraccount {
	rows, err := framework.DB.Query("SELECT id, username, email, phnumber, department, humanname, role FROM useraccount")
	if err != nil {
		fmt.Println("can not find user in table")
		return nil
	}
	var useraccounts []model.Useraccount
	for rows.Next() {
		var useraccount model.Useraccount
		err := rows.Scan(&useraccount.Id, &useraccount.Name, &useraccount.Email,
			&useraccount.PhoneNumber, &useraccount.Department, &useraccount.Human, &useraccount.Role)
		if err != nil {
			fmt.Println("get user error : ", err)
			continue
		}
		useraccounts = append(useraccounts, useraccount)
	}
	rows.Close()
	return useraccounts
}

func (p *UserDao) ModifyUserByPhoneNumber(user model.Useraccount) string {
	_, err := framework.DB.Exec("UPDATE useraccount SET password=? WHERE phnumber=?;", user.Password, user.PhoneNumber)
	if err != nil {
		log.Println("modify password failed", err)
		return "error"
	}
	return "ok"
}

func (p *UserDao) AddUser(user model.Useraccount) string {
	_, err := framework.DB.Exec("INSERT INTO useraccount (username, password, email, phnumber, department, humanname) VALUES (?, ?, ?, ?, ?, ?);",
		user.Name, user.Password, user.Email, user.PhoneNumber, user.Department, user.Human)
	if err != nil {
		log.Println("add user failed", err)
		return "error"
	}
	return "ok"
}

func (p *UserDao) ModifyUser(user model.Useraccount) string {
	_, err := framework.DB.Exec("UPDATE useraccount SET username=?, email=?, phnumber=?, department=?, humanname=?, role=? WHERE id=?;",
		user.Name, user.Email, user.PhoneNumber, user.Department, user.Human, user.Role, user.Id)
	if err != nil {
		log.Println("modify user failed", err)
		return "error"
	}
	return "ok"
}

func (p *UserDao) DeleteUser(id int) string {
	_, err := framework.DB.Exec("DELETE FROM useraccount WHERE id=?;", id)
	if err != nil {
		log.Println("delete user failed", err)
		return "error"
	}
	return "ok"
}

func (p *UserDao) ModifyUserByUsername(user model.Useraccount) string {
	_, err := framework.DB.Exec("UPDATE useraccount SET password=? WHERE username=?;", user.Password, user.Name)
	if err != nil {
		log.Println("modify password failed", err)
		return "error"
	}
	return "ok"
}
