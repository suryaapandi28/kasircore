package response

type Response struct {
	Meta Meta        `json:"meta"`
	Data interface{} `json:"data"`
}

type Responsefieldempty struct {
	Meta Meta_error `json:"Meta_error"`
}

type ResponseUnverified struct {
	Meta Meta_unverified `json:"Meta_unverified"`
}

type Responseduplicateemail struct {
	Meta Meta_error_duplicate `json:"Meta_email_duplicate"`
}

type Meta struct {
	RC      int    `json:"rc"`
	Message string `json:"message"`
}

type Meta_error struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Colum   string `json:"colum"`
}

type Meta_unverified struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Status  string `json:"status"`
}

type Meta_error_duplicate struct {
	RC      int    `json:"rc"`
	Message string `json:"message"`
}

func SuccessResponse(rc int, message string, data interface{}) Response {
	return Response{
		Meta: Meta{
			RC:      rc,
			Message: message,
		},
		Data: data,
	}
}

func ErrorResponse(rc int, message string) Response {
	return Response{
		Meta: Meta{
			RC:      rc,
			Message: message,
		},
		Data: nil,
	}
}

func UNVERIFIEDResponse(code int, status string, message string) ResponseUnverified {
	return ResponseUnverified{
		Meta: Meta_unverified{
			Code:    code,
			Message: message,
			Status:  status,
		},
	}
}
func Errorfieldempty(code int, colum string) Responsefieldempty {
	return Responsefieldempty{
		Meta: Meta_error{
			Code:    code,
			Message: "column cannot be empty",
			Colum:   colum,
		},
	}
}

func DuplicateEmailResponse(rc int, message string) Responseduplicateemail {
	return Responseduplicateemail{
		Meta: Meta_error_duplicate{
			RC:      rc,
			Message: message,
		},
	}
}
