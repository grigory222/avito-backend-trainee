package dto

type ReceptionDto struct {
	Id       string `json:"id"`       // Идентификатор приемки
	DateTime string `json:"dateTime"` // Дата и время создания приемки
	PVZId    string `json:"pvzId"`    // Идентификатор ПВЗ
	Status   string `json:"status"`   // Статус приемки
}

type ReceptionWithProductsDto struct {
	Reception ReceptionDto `json:"reception"` // Информация о приемке
	Products  []ProductDto `json:"products"`  // Информация о всех товарах в приемке
}

type CreateReceptionRequestDto struct {
	PVZId string `json:"pvzId"` // Идентификатор ПВЗ
}

type CreateReceptionResponseDto struct {
	Id       string `json:"id"`       // Идентификатор приемки
	DateTime string `json:"dateTime"` // Дата и время приемки
	PVZId    string `json:"pvzId"`    // Идентификатор ПВЗ
	Status   string `json:"status"`   // Статус приемки
}

type CloseReceptionResponseDto struct {
	Id       string `json:"id"`       // Идентификатор приемки
	DateTime string `json:"dateTime"` // Дата и время приемки
	PVZId    string `json:"pvzId"`    // Идентификатор ПВЗ
	Status   string `json:"status"`   // Статус приемки
}
