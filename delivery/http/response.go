package delivery

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/azinudinachzab/scr-syky-tech-test/model"
	"github.com/azinudinachzab/scr-syky-tech-test/pkg/errs"
	"github.com/azinudinachzab/scr-syky-tech-test/pkg/strcase"
	"github.com/go-playground/validator/v10"
)

type keyCtxStatusCode struct{}
type httpResponse struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func responseData(rw http.ResponseWriter, r *http.Request, data httpResponse) {
	rw.Header().Set("Content-type", "application/json")
	writeHeaderAndContext(rw, r, http.StatusOK)
	toJson(rw, data)
}

func responseError(rw http.ResponseWriter, r *http.Request, err error) {
	rw.Header().Set("Content-type", "application/json")

	cerr, ok := err.(*errs.Error)
	if !ok {
		writeHeaderAndContext(rw, r, http.StatusInternalServerError)
		toJson(rw, map[string]string{"message": err.Error()})
		return
	}

	if cerr.Code == model.ECodeValidateFail && cerr.Err != nil && cerr.Attributes == nil {
		writeHeaderAndContext(rw, r, http.StatusBadRequest)
		data := convertValidatorErrToAttribute(cerr)
		toJson(rw, data)
		return
	}

	stsCode := http.StatusBadRequest
	if cerr.Code == model.ECodeInternal {
		stsCode = http.StatusInternalServerError
	}

	writeHeaderAndContext(rw, r, stsCode)
	toJson(rw, cerr)
}

func toJson(w http.ResponseWriter, data interface{}) {
	if err := json.NewEncoder(w).Encode(data); err != nil {
		log.Println("fail to encode to JSON", err)
	}
}

func convertValidatorErrToAttribute(cerr *errs.Error) *errs.Error {
	if _, ok := cerr.Err.(*validator.InvalidValidationError); ok {
		return errs.New(model.ECodeInternal, "unknown error")
	}

	attrs := make([]errs.Attribute, 0)
	for _, fe := range cerr.Err.(validator.ValidationErrors) {
		fld := strcase.ToSnakeCase(fe.Field())
		msg := tagToMsg(fld, fe)
		attrs = append(attrs, errs.Attribute{
			Field:   fld,
			Message: msg,
		})
	}

	return errs.NewWithAttribute(cerr.Code, cerr.Message, attrs)
}

func tagToMsg(field string, fe validator.FieldError) string {
	switch fe.Tag() {
	case "min":
		return "less than minimum length"
	case "max":
		return "over than maximum length max: " + fe.Param()
	case "required":
		return "cannot be empty."
	case "required_if":
		return "cannot be empty if field " + fe.Param()
	case "oneof":
		return "value must be one of: " + fe.Param()
	case "len":
		return "value must have length " + fe.Param()
	}
	// TODO: need to add more custom message for every tag validation

	return field + " is failed to validate"
}

func writeHeaderAndContext(w http.ResponseWriter, r *http.Request, code int) {
	w.WriteHeader(code)
	ctx := context.WithValue(r.Context(), keyCtxStatusCode{}, code)
	*r = *(r.WithContext(ctx))
}
