package helpers

import (
	"github.com/Clementol/restur-manag/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func GetUserDetails(userObj primitive.M, user models.User) primitive.M {

	userObj["user_id"] = user.User_id
	userObj["avatar"] = user.Avatar
	userObj["email"] = user.Email
	userObj["first_name"] = user.First_Name
	userObj["last_name"] = user.Last_Name
	userObj["created_at"] = user.Created_at
	userObj["updated_at"] = user.Updated_at
	userObj["phone"] = user.Phone
	userObj["user_id"] = user.User_id
	userObj["token"] = user.Token

	return userObj

}
