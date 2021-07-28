package packages

import ut "github.com/go-playground/universal-translator"

var trans *ut.Translator

func SetTranslator(translator *ut.Translator)  {
	trans = translator
}
