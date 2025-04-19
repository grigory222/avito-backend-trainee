package dto

const (
	ProductTypeElectronics = "электроника"
	ProductTypeClothes     = "одежда"
	ProductTypeShoes       = "обувь"
)

type ProductDto struct {
	Id          string `json:"id"`          // Идентификатор
	DateTime    string `json:"dateTime"`    // Дата и время
	Type        string `json:"type"`        // Тип товара
	ReceptionId string `json:"receptionId"` // Идентификатор приемки
}

type AddProductRequestDto struct {
	Type  string `json:"type"`  // Тип товара
	PVZId string `json:"pvzId"` // Идентификатор ПВЗ, на который добавляется товар
}

type AddProductResponseDto struct {
	Id          string `json:"id"`          // Идентификатор
	DateTime    string `json:"dateTime"`    // Дата и время
	Type        string `json:"type"`        // Тип товара
	ReceptionId string `json:"receptionId"` // Идентификатор приемки
}
