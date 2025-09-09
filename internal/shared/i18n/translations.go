package i18n

// Translation store: Code -> Locale -> Message
var translations = map[string]map[Locale]string{
	"ERR:USER_NOT_FOUND": {
		EN_US: "User with ID %v not found",
		TR_TR: "ID'si %v olan kullanıcı bulunamadı",
	},
	"ERR:USER_ALREADY_EXISTS": {
		EN_US: "user already exists",
		TR_TR: "kullanıcı kaydı yapılmış",
	},
	"ERR:INVALID_ID": {
		EN_US: "Invalid user ID",
		TR_TR: "Geçersiz kullanıcı ID'si",
	},
	"ERR:INVALID_REQUEST_PAYLOAD": {
		EN_US: "Invalid request payload",
		TR_TR: "Geçersiz istek verisi",
	},
	"ERR:INTERNAL_SERVER_ERROR": {
		EN_US: "Something went wrong",
		TR_TR: "Bir şeyler ters gitti",
	},
}
