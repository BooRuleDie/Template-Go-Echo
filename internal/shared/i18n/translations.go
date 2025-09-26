package i18n

type Messages map[Locale]string

type Translation struct {
	// IsInternal determines whether this translation
	// should be exposed to the UI through the i18n handler.
	// If set to false, the translation is intended for internal use only.
	IsInternal bool
	Messages   Messages
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
// TODO: this is outdated fix this later on
// Examples:
//   Error:            "ERR:USER_NOT_FOUND"
//   Validation Error: "VAL:FIELD_REQUIRED"
//   UI (with module): "UI:BUTTON_SAVE"
//   UI (no module):   "UI:LOADING"
//   Success:          "SUC:AUTH_LOGIN"

// Translation store: Code -> Translation -> Locale -> Message
var Translations = map[string]Translation{
	// ========== VALIDATION FIELDS ==========
	"FIELD:ID": {
		IsInternal: true,
		Messages: map[Locale]string{
			EN_US: "ID",
			TR_TR: "ID",
		},
	},
	"FIELD:NAME": {
		IsInternal: true,
		Messages: map[Locale]string{
			EN_US: "Name",
			TR_TR: "İsim",
		},
	},
	"FIELD:EMAIL": {
		IsInternal: true,
		Messages: map[Locale]string{
			EN_US: "Email",
			TR_TR: "E-posta",
		},
	},
	"FIELD:PHONE": {
		IsInternal: true,
		Messages: map[Locale]string{
			EN_US: "Phone",
			TR_TR: "Telefon",
		},
	},
	"FIELD:TAGS": {
		IsInternal: true,
		Messages: map[Locale]string{
			EN_US: "Tags",
			TR_TR: "Etiketler",
		},
	},
	"FIELD:PASSWORD": {
		IsInternal: true,
		Messages: map[Locale]string{
			EN_US: "Password",
			TR_TR: "Şifre",
		},
	},
	"FIELD:AGE": {
		IsInternal: true,
		Messages: map[Locale]string{
			EN_US: "Age",
			TR_TR: "Yaş",
		},
	},
	"FIELD:ROLE": {
		IsInternal: true,
		Messages: map[Locale]string{
			EN_US: "Role",
			TR_TR: "Rol",
		},
	},

	// ========== VALIDATION MESSAGES ==========
	"VAL:VALIDATION_ERR": {
		IsInternal: true,
		Messages: map[Locale]string{
			EN_US: "Validation error",
			TR_TR: "Doğrulama hatası",
		},
	},
	"VAL:REQUIRED": {
		IsInternal: true,
		Messages: map[Locale]string{
			EN_US: "%v is required",
			TR_TR: "%v alanı zorunludur",
		},
	},
	"VAL:EMAIL": {
		IsInternal: true,
		Messages: map[Locale]string{
			EN_US: "%v must be a valid email",
			TR_TR: "%v geçerli bir e-posta olmalıdır",
		},
	},
	"VAL:PHONE": {
		IsInternal: true,
		Messages: map[Locale]string{
			EN_US: "%v must be a valid Turkish phone number (+905XXXXXXXXX)",
			TR_TR: "%v geçerli bir Türk GSM numarası olmalıdır (+905XXXXXXXXX)",
		},
	},
	"VAL:PASSWORD": {
		IsInternal: true,
		Messages: map[Locale]string{
			EN_US: "%v must be at least 8 characters long and include 1 uppercase letter, 1 lowercase letter, 1 digit, and 1 special character",
			TR_TR: "%v en az 8 karakter olmalı, 1 büyük harf, 1 küçük harf, 1 rakam ve 1 özel karakter içermelidir",
		},
	},

	// min / max / len -- string
	"VAL:MIN_STRING": {
		IsInternal: true,
		Messages: map[Locale]string{
			EN_US: "%v must be at least %v characters",
			TR_TR: "%v en az %v karakter olmalıdır",
		},
	},
	"VAL:MAX_STRING": {
		IsInternal: true,
		Messages: map[Locale]string{
			EN_US: "%v must be at most %v characters",
			TR_TR: "%v en fazla %v karakter olmalıdır",
		},
	},
	"VAL:LEN_STRING": {
		IsInternal: true,
		Messages: map[Locale]string{
			EN_US: "%v must be exactly %v characters",
			TR_TR: "%v tam olarak %v karakter olmalıdır",
		},
	},

	// min / max / len -- numeric
	"VAL:MIN_NUMBER": {
		IsInternal: true,
		Messages: map[Locale]string{
			EN_US: "%v must be greater than or equal to %v",
			TR_TR: "%v en az %v olmalıdır",
		},
	},
	"VAL:MAX_NUMBER": {
		IsInternal: true,
		Messages: map[Locale]string{
			EN_US: "%v must be less than or equal to %v",
			TR_TR: "%v en fazla %v olmalıdır",
		},
	},
	"VAL:LEN_NUMBER": {
		IsInternal: true,
		Messages: map[Locale]string{
			EN_US: "%v must be %v",
			TR_TR: "%v %v olmalıdır",
		},
	},

	// min / max / len -- slice
	"VAL:MIN_SLICE": {
		IsInternal: true,
		Messages: map[Locale]string{
			EN_US: "%v must contain at least %v items",
			TR_TR: "%v en az %v öğe içermelidir",
		},
	},
	"VAL:MAX_SLICE": {
		IsInternal: true,
		Messages: map[Locale]string{
			EN_US: "%v must contain at most %v items",
			TR_TR: "%v en fazla %v öğe içermelidir",
		},
	},
	"VAL:LEN_SLICE": {
		IsInternal: true,
		Messages: map[Locale]string{
			EN_US: "%v must contain exactly %v items",
			TR_TR: "%v tam olarak %v öğe içermelidir",
		},
	},

	// gte / lte numeric messages
	"VAL:GTE_NUMBER": {
		IsInternal: true,
		Messages: map[Locale]string{
			EN_US: "%v must be greater than or equal to %v",
			TR_TR: "%v en az %v olmalıdır",
		},
	},
	"VAL:LTE_NUMBER": {
		IsInternal: true,
		Messages: map[Locale]string{
			EN_US: "%v must be less than or equal to %v",
			TR_TR: "%v en fazla %v olmalıdır",
		},
	},

	// others
	"VAL:ALPHA": {
		IsInternal: true,
		Messages: map[Locale]string{
			EN_US: "%v must contain letters only",
			TR_TR: "%v yalnızca harf içermelidir",
		},
	},
	"VAL:ALPHANUM": {
		IsInternal: true,
		Messages: map[Locale]string{
			EN_US: "%v must contain only letters and numbers",
			TR_TR: "%v yalnızca harf ve rakam içermelidir",
		},
	},
	"VAL:CONTAINS": {
		IsInternal: true,
		Messages: map[Locale]string{
			EN_US: "%v must contain '%v'",
			TR_TR: "%v '%v' içermelidir",
		},
	},
	"VAL:ONEOF": {
		IsInternal: true,
		Messages: map[Locale]string{
			EN_US: "%v must be one of [%v]",
			TR_TR: "%v şu değerlerden biri olmalıdır: [%v]",
		},
	},

	// ========== GENERIC ERROR MESSAGES ==========
	"ERR:HTTP_500": {
		IsInternal: true,
		Messages: map[Locale]string{
			EN_US: "Internal Server Error",
			TR_TR: "Bir Şeyler Ters Gitti",
		},
	},
	"ERR:HTTP_503": {
		IsInternal: true,
		Messages: map[Locale]string{
			EN_US: "Service Unavailable or Timed Out",
			TR_TR: "Servis Kullanılamıyor veya Zaman Aşımına Uğradı",
		},
	},
	"ERR:HTTP_400": {
		IsInternal: true,
		Messages: map[Locale]string{
			EN_US: "Bad Request",
			TR_TR: "Geçersiz İstek",
		},
	},
	"ERR:HTTP_401": {
		IsInternal: true,
		Messages: map[Locale]string{
			EN_US: "Unauthorized",
			TR_TR: "Yetkisiz Erişim",
		},
	},
	"ERR:HTTP_404": {
		IsInternal: true,
		Messages: map[Locale]string{
			EN_US: "Not Found",
			TR_TR: "Bulunamadı",
		},
	},
	"ERR:HTTP_405": {
		IsInternal: true,
		Messages: map[Locale]string{
			EN_US: "Method Not Allowed",
			TR_TR: "İzin Verilmeyen Yöntem",
		},
	},
	"ERR:HTTP_413": {
		IsInternal: true,
		Messages: map[Locale]string{
			EN_US: "Payload Too Large",
			TR_TR: "İstek Boyutu Çok Büyük",
		},
	},
	"ERR:HTTP_415": {
		IsInternal: true,
		Messages: map[Locale]string{
			EN_US: "Unsupported Media Type",
			TR_TR: "Desteklenmeyen Medya Türü",
		},
	},
	"ERR:HTTP_429": {
		IsInternal: true,
		Messages: map[Locale]string{
			EN_US: "Too Many Requests",
			TR_TR: "Çok Fazla İstek",
		},
	},
}
