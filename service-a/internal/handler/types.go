package handler

// CEPRequest representa o request para validação de CEP
type CEPRequest struct {
	CEP string `json:"cep"`
}

// WeatherResponse representa a resposta com dados de clima
type WeatherResponse struct {
	City  string  `json:"city"`
	TempC float64 `json:"temp_C"`
	TempF float64 `json:"temp_F"`
	TempK float64 `json:"temp_K"`
}

// ErrorResponse representa uma resposta de erro
type ErrorResponse struct {
	Error string `json:"error"`
}
