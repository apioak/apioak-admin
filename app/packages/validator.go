package packages

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
)

var trans *ut.Translator

func SetTranslator(translator *ut.Translator)  {
	trans = translator
}

func Translate(errs validator.ValidationErrors) string {
	var errMsg string
	for _, e := range errs {
		errMsg = e.Translate(*trans)
	}
	return errMsg
}

func ParseRequestParams(c *gin.Context, request interface{}) (string, error) {
	if err := c.ShouldBind(request); err != nil {

		var errStr string
		switch err.(type) {
		case validator.ValidationErrors:
			errStr = Translate(err.(validator.ValidationErrors))
		case *json.UnmarshalTypeError:
			unmarshalTypeError := err.(*json.UnmarshalTypeError)
			errStr = fmt.Errorf("[%s]类型错误，期望类型:%s", unmarshalTypeError.Field, unmarshalTypeError.Type.String()).Error()
		default:
			errStr = err.Error()
		}

		return errStr, err
	}
	return "", nil
}
