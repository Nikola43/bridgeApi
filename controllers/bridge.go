package controllers

import (
	"bridgeApi/utils"
	"fmt"
)

func ON(context *fiber.Ctx) error {
	fmt.Println("ON")
	return utils.ReturnSuccessResponse(context)
}
