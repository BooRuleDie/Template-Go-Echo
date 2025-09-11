package i18n

type Translation struct {
	// IsInternal determines whether this translation
	// should be exposed to the UI through the i18n handler.
	// If set to false, the translation is intended for internal use only.
	IsInternal bool
	Messages   map[Locale]string
}

// Translation codes follow this format:
// [CATEGORY][:MODULE]_[MESSAGE]
//
// CATEGORY can be one of:
//   - ERR:      Error messages
//   - VAL:      Validation errors
//   - UI:       User interface strings
//   - SUC:      Success messages
//
// The MODULE part is required for ERR, VAL, and SUC, but optional for UI codes.
// MODULE helps categorize translations by application domain, such as USER, AUTH, etc.
//
// Examples:
//   Error:            "ERR:USER_NOT_FOUND"
//   Validation Error: "VAL:FIELD_REQUIRED"
//   UI (with module): "UI:BUTTON_SAVE"
//   UI (no module):   "UI:LOADING"
//   Success:          "SUC:AUTH_LOGIN"

// Translation store: Code -> Translation -> Locale -> Message
var Translations = map[string]Translation{
	// ========== USER MODULE ==========
	"SUC:USER_CREATED": {
		IsInternal: true,
		Messages: map[Locale]string{
			EN_US: "User created successfully",
			TR_TR: "Kullanıcı başarıyla oluşturuldu",
		},
	},
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
	"ERR:USER_INVALID_ID": {
		IsInternal: true,
		Messages: map[Locale]string{
			EN_US: "Invalid user ID",
			TR_TR: "Geçersiz kullanıcı ID'si",
		},
	},
	"ERR:USER_INVALID_REQUEST_PAYLOAD": {
		IsInternal: true,
		Messages: map[Locale]string{
			EN_US: "Invalid request payload",
			TR_TR: "Geçersiz istek verisi",
		},
	},

	// ========== GENERIC ERROR MESSAGES ==========
	"ERR:INTERNAL_SERVER_ERROR": {
		IsInternal: true,
		Messages: map[Locale]string{
			EN_US: "Something went wrong",
			TR_TR: "Bir şeyler ters gitti",
		},
	},
}
