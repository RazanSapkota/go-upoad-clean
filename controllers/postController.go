package controllers

import (
	initialize "example/go-api/Initialize"
	"example/go-api/models"
	"fmt"

	"github.com/gin-gonic/gin"
)

func CreatePost(c *gin.Context) {

var post models.Post

c.Bind(&post)

result := initialize.DB.Create(&post) 


if result.Error!=nil{
	c.Status(400)
	return
}
	
	c.JSON(200, gin.H{
		"post":post,
	})
}

func GetPosts(c *gin.Context){

	// Get all records
	var posts []models.Post

   result := initialize.DB.Find(&posts)

   //Selecting a posts with amount > avg of all amount
    // subQuery:=initialize.DB.Table("posts").Select("AVG(amount)")
    // result := initialize.DB.Where("amount > (?)", subQuery).Find(&posts)
	

	//selecting multiples intable
	//result := initialize.DB.Where("(title, amount) IN ?", [][]interface{}{{"test4", 4}, {"test3", 2}}).Find(&posts)
if result.Error!=nil{
	c.Status(400)
	return
}
c.JSON(200, gin.H{
	"posts":posts,
})
}

func FetchOnePost(c *gin.Context){

	id:=c.Param("id")
	// Get all records
	var post models.Post
	
    result := initialize.DB.First(&post, id)
	

if result.Error!=nil{
	c.Status(400)
	fmt.Println(result.Error)
	return
}
c.JSON(200, gin.H{
	"post":post,
})
}

func UpdatePost(c *gin.Context){
	var post models.Post
//Getting an Id
	id:=c.Param("id")

	//finding post from database by id
    result := initialize.DB.First(&post, id)

	c.Bind(&post)
	
	if result.Error!=nil{
		c.Status(400)
		return
	}

	result = initialize.DB.Save(&post)
	//Update
	if result.Error!=nil{
		c.Status(400)
		return
	}
//

c.JSON(200, gin.H{
	"post":post,
})
}
