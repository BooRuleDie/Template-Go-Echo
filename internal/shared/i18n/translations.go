package i18n

type Translation struct {
	// IsInternal determines whether this translation
	// should be exposed to the UI through the i18n handler.
	// If set to false, the translation is intended for internal use only.
	IsInternal bool
	Messages   map[Locale]string
}

// Translation store: Code -> Locale -> Message
var Translations = map[string]Translation{
	"ERR:USER_NOT_FOUND": {
		IsInternal: true,
		Messages: map[Locale]string{
			EN_US: "User with ID %v not found",
			TR_TR: "ID'si %v olan kullanıcı bulunamadı",
		},
	},
	"ERR:USER_ALREADY_EXISTS": {
		IsInternal: true,
		Messages: map[Locale]string{
			EN_US: "User already exists",
			TR_TR: "Kullanıcı kaydı yapılmış",
		},
	},
	"ERR:INVALID_ID": {
		IsInternal: true,
		Messages: map[Locale]string{
			EN_US: "Invalid user ID",
			TR_TR: "Geçersiz kullanıcı ID'si",
		},
	},
	"ERR:INVALID_REQUEST_PAYLOAD": {
		IsInternal: true,
		Messages: map[Locale]string{
			EN_US: "Invalid request payload",
			TR_TR: "Geçersiz istek verisi",
		},
	},
	"ERR:INTERNAL_SERVER_ERROR": {
		IsInternal: true,
		Messages: map[Locale]string{
			EN_US: "Something went wrong",
			TR_TR: "Bir şeyler ters gitti",
		},
	},
}
