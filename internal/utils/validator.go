package utils

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

var FieldMessages = map[string]map[string]string{
	"Name": {
		"required": "Name field is required",
		"min":      "Name field must be at least 2 characters",
		"max":      "Name field must be at most 255 characters",
	},
	"Status": {
		"required": "Status field is required",
		"oneof":    "Status field must be either 0 or 1",
	},
	"Image": {
		"required": "Image field is required",
	},
	"CategoryId": {
		"required": "Category field is required",
	},
}

// ValidateStruct binds & validates any DTO and returns friendly error messages
func ValidateStruct(c *gin.Context, s interface{}) (bool, map[string]string) {
	if err := c.ShouldBind(s); err != nil {
		errorsMap := make(map[string]string)
		if errs, ok := err.(validator.ValidationErrors); ok {
			for _, e := range errs {
				field := e.Field()
				tag := e.Tag()
				if fieldMsg, ok := FieldMessages[field]; ok {
					if msg, ok := fieldMsg[tag]; ok {
						errorsMap[field] = msg
					} else {
						errorsMap[field] = "Invalid value"
					}
				} else {
					errorsMap[field] = "Invalid value"
				}
			}
		} else {
			errorsMap["General"] = err.Error()
		}
		return false, errorsMap
	}
	return true, nil
}
