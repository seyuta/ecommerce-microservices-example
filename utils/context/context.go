package utilcontext

import (
	"context"
	"fmt"
)

// GetUsername from Context
func GetUsername(c context.Context) string {
	uname := fmt.Sprintf("%v", c.Value("username"))
	if uname == "" || uname == "<nil>" {
		uname = "anon"
	}
	return uname
}

// GetUserphone from Context
func GetUserphone(c context.Context) string {
	phone := fmt.Sprintf("%v", c.Value("phone"))
	if phone == "" || phone == "<nil>" {
		phone = ""
	}
	return phone
}

// GetUserEmail from Context
func GetUserEmail(c context.Context) string {
	email := fmt.Sprintf("%v", c.Value("email"))
	if email == "" || email == "<nil>" {
		email = ""
	}
	return email
}

// GetUserID from Context
func GetUserID(c context.Context) string {
	uid := fmt.Sprintf("%v", c.Value("userid"))
	if uid == "" || uid == "<nil>" {
		uid = ""
	}
	return uid
}

// GetDeviceID from Context
func GetDeviceID(c context.Context) string {
	dvcid := fmt.Sprintf("%v", c.Value("dvcid"))
	if dvcid == "" || dvcid == "<nil>" {
		dvcid = ""
	}
	return dvcid
}
