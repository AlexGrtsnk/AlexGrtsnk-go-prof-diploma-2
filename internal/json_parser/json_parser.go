package jsonparser

// Тип, содержащий логин и пароль пользователя
type Auth struct {
	// Login поле с логином пользователя
	Login string `json:"login"`
	// Password поле с паролем пользователя
	Password string `json:"password"`
}

// Масиив текстов пользователя
type TextList []TextAnswer

// Тип, содержащий текст пользователя
type TextAnswer struct {
	// UserText поле с текстом пользователя
	UserText string `json:"user_text"`
}

// Масиив бинарников пользователя
type BinList []BinAnswer

// Тип, содержащий бинарные данные пользователя
type BinAnswer struct {
	// UserBin поле с бинарными данными пользователя
	UserBin []byte `json:"user_bin"`
}

// Масиив текстов пользователя
type CardList []CardAnswer

// Тип, содержащий данные о карте пользователя
type CardAnswer struct {
	// UserCardNum поле с номер карты пользователя
	UserCardNum string `json:"user_card_num"`
	// UserCardName поле с номером держателя карты пользователя
	UserCardName string `json:"user_card_name"`
	// UserCardCVV поле с CVV карты пользователя
	UserCardCVV string `json:"user_card_cvv"`
}

// Setting - тип, используюшийся для чтения файла формата JSON с параметрами конфигурации exe
type Setting struct {
	// аналог переменной окружения SERVER_ADDRESS или флага -a
	ServerAddress string `json:"server_address"`
	// аналог переменной окружения BASE_URL или флага -b
	BaseURL string `json:"base_url"`
	// LongURL изначальный отправленный пользователем url
	FileStoragePath string `json:"file_storage_path"`
	// аналог переменной окружения DATABASE_DSN или флага -d
	DataBaseDSN string `json:"database_dsn"`
	// DelFlag флаг мягкого удаления из базы данныз
	EnableHTTPS bool `json:"enable_https"`
	// подсеть
	TrustedSubnet string `json:"trusted_subnet"`
}
