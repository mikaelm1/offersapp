package main

import "github.com/gin-gonic/gin"

func main() {

	router := gin.Default()

	// usersGroup = router.Group("users")

	router.Run(":3000")

}
