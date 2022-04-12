package controllers

import (
	"fmt"

	"github.com/nikola43/bridgeApi/utils"
)

func ON(context *fiber.Ctx) error {
	fmt.Println("ON")
	return utils.ReturnSuccessResponse(context)
}
