package model

import (
	"strconv"
	"strings"
	"time"

	"github.com/garyburd/redigo/redis"
)

// User model
type User struct {
	id        string
	role      *UserRole
	Username  string
	Password  string
	Name      string
	Email     string
	AboutMe   string
	CreatedAt time.Time
	UpdatedAt time.Time
}

// UserRole type
type UserRole struct {
	Admin      bool
	Instructor bool
}

// ID returns user id
func (x *User) ID() string {
	return x.id
}

// Role returns user role
func (x *User) Role() *UserRole {
	if x.role == nil {
		x.role = &UserRole{}
	}
	return x.role
}

// Save saves user
func (x *User) Save(c redis.Conn) error {
	if len(x.id) == 0 {
		id, err := redis.Int64(c.Do("INCR", key("id", "u")))
		if err != nil {
			return err
		}
		x.id = strconv.FormatInt(id, 10)
	}

	c.Send("MULTI")
	c.Send("SADD", key("u", "all"), x.id)

	x.UpdatedAt = time.Now()
	if x.CreatedAt.IsZero() {
		x.CreatedAt = x.UpdatedAt
		c.Send("ZADD", key("u", "t0"), x.CreatedAt.UnixNano(), x.id)
	} else {
		c.Send("ZADD", key("u", "t0"), x.CreatedAt.UnixNano(), x.id) // TODO: remove after migrate
	}

	c.Send("ZADD", key("u", "t1"), x.UpdatedAt.UnixNano(), x.id)
	c.Send("HSET", key("u"), x.id, enc(x))

	if x.role != nil {
		if x.role.Admin {
			c.Send("SADD", key("u", "admin"), x.id)
		} else {
			c.Send("SREM", key("u", "admin"), x.id)
		}
		if x.role.Instructor {
			c.Send("SADD", key("u", "instructor"), x.id)
		} else {
			c.Send("SREM", key("u", "instructor"), x.id)
		}
	}

	_, err := c.Do("EXEC")
	if err != nil {
		return err
	}
	return nil
}

// GetUsers gets users
func GetUsers(c redis.Conn, userIDs []string) ([]*User, error) {
	xs := make([]*User, len(userIDs))
	for _, userID := range userIDs {
		c.Send("SISMEMBER", key("u", "all"), userID)
		c.Send("HGET", key("u"), userID)
		c.Send("SISMEMBER", key("u", "admin"), userID)
		c.Send("SISMEMBER", key("u", "instructor"), userID)
	}
	c.Flush()
	for i := range userIDs {
		exists, _ := redis.Bool(c.Receive())
		if !exists {
			c.Receive()
			c.Receive()
			c.Receive()
			continue
		}
		var x User
		b, err := redis.Bytes(c.Receive())
		if err != nil {
			return nil, err
		}
		err = dec(b, &x)
		if err != nil {
			return nil, err
		}
		x.role = &UserRole{}
		x.role.Admin, _ = redis.Bool(c.Receive())
		x.role.Instructor, _ = redis.Bool(c.Receive())
		xs[i] = &x
	}
	return xs, nil
}

// GetUser gets user from id
func GetUser(c redis.Conn, userID string) (*User, error) {
	xs, err := GetUsers(c, []string{userID})
	if err != nil {
		return nil, err
	}
	return xs[0], nil
}

// GetUserFromUsername gets user from username
func GetUserFromUsername(c redis.Conn, username string) (*User, error) {
	userID, err := redis.String(c.Do("HGET", key("u", "username"), username))
	if err == redis.ErrNil {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return GetUser(c, userID)
}

// GetUserFromEmail gets user from email
func GetUserFromEmail(c redis.Conn, email string) (*User, error) {
	userID, err := redis.String(c.Do("HGET", key("u", "email"), email))
	if err == redis.ErrNil {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return GetUser(c, userID)
}

// GetUserFromEmailOrUsername gets user from email if given input contains '@'
// otherwise get user from username
func GetUserFromEmailOrUsername(c redis.Conn, user string) (*User, error) {
	if strings.Contains(user, "@") {
		return GetUserFromEmail(c, user)
	}
	return GetUserFromUsername(c, user)
}

// GetUserFromProvider gets user from provider
func GetUserFromProvider(c redis.Conn, provider string, providerUserID string) (*User, error) {
	userID, err := redis.String(c.Do("HGET", key("u", "provider", provider), providerUserID))
	if err == redis.ErrNil {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return GetUser(c, userID)
}

// IsUserAdmin returns true if given user id is an admin
// func IsUserAdmin(c redis.Conn, userID string) (bool, error) {
// 	if len(userID) == 0 {
// 		return false, nil
// 	}
// 	return redis.Bool(c.Do("SISMEMBER", key("u", "admin"), userID))
// }

// // IsUserInstructor returns true if given user id is an instructor
// func IsUserInstructor(c redis.Conn, userID string) (bool, error) {
// 	if len(userID) == 0 {
// 		return false, nil
// 	}
// 	return redis.Bool(c.Do("SISMEMBER", key("u", "instructor"), userID))
// }
