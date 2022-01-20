package feed

import (
	"encoding/json"
	"fmt"
	"github.com/bradfitz/gomemcache/memcache"
	"github.com/deltamc/otus-social-networks-backend/db"
	_ "github.com/bradfitz/gomemcache/memcache"
	post2 "github.com/deltamc/otus-social-networks-backend/models/post"
	"os"
)

type userList struct {
	id int64 `db:"id"`
	users string `db:"users"`
}
var mc  *memcache.Client

func FeedSetCache(userId ...int64) (err error) {
	dbPool := db.OpenDB()

	where := ""
	if len(userId) > 0 {
		where = fmt.Sprintf("WHERE users.id=%d OR friends.friend_id =%d", userId[0], userId[0])
	}

	sqlStmt := fmt.Sprintf(`SELECT users.id, IFNULL(CONCAT(users.id, ",", GROUP_CONCAT(friends.friend_id)),users.id) as users FROM users
LEFT JOIN friends ON users.id= friends.user_id %s GROUP BY users.id`, where)

	rows, err := dbPool.Query(sqlStmt)

	if err != nil {
		return
	}

	defer rows.Close()

	if mc == nil {
		mc = memcache.New(fmt.Sprintf("%s:%s", os.Getenv("MEMCACHE_HOST"), os.Getenv("MEMCACHE_PORT")))
	}


	for rows.Next() {
		var users userList

		if err = rows.Err(); err != nil {
			return
		}

		err = rows.Scan(&users.id,&users.users)
		if err != nil {
			return
		}

		err = feedSetCache(users, mc)
		if err != nil {
			return
		}
	}

	return
}

func feedSetCache(users userList, mc *memcache.Client) (err error){
	dbPool := db.OpenDB()


	sqlStmt := fmt.Sprintf(
		`SELECT posts.id, user_id, body, created_at, IFNULL(CONCAT(users.first_name," ", users.last_name),"") as user_name  
				FROM posts 
				LEFT JOIN users ON users.id=user_id
				WHERE user_id in(%s) ORDER BY created_at DESC LIMIT 1000`,
		users.users)

	rows, err := dbPool.Query(sqlStmt)


	if err != nil {
		return
	}

	defer rows.Close()

	var posts []post2.Post

	for rows.Next() {
		var post post2.Post

		if err = rows.Err(); err != nil {
			return
		}

		err = rows.Scan(&post.Id, &post.UserId, &post.Body, &post.CreatedAt, &post.UserName)
		if err != nil {
			return
		}

		posts = append(posts, post)
	}

	//fmt.Println(users.id, len(posts), users.users)

	b, err := json.Marshal(posts)

	if err != nil {
		return err
	}

	mc.Set(&memcache.Item{Key: fmt.Sprintf("feed_%d", users.id), Value: []byte(b)})
	return
}
