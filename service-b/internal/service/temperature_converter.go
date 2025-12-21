package service

// TemperatureConverter interface para conversão de temperaturas
type TemperatureConverter interface {
	ConvertFromCelsius(celsius float64) (float64, float64, float64) // Retorna C, F, K
}

// temperatureConverter implementação da conversão de temperaturas
type temperatureConverter struct{}

// NewTemperatureConverter cria uma nova instância do conversor de temperaturas
func NewTemperatureConverter() TemperatureConverter {
	return &temperatureConverter{}
}

// ConvertFromCelsius converte temperatura de Celsius para Fahrenheit e Kelvin
// Retorna: Celsius, Fahrenheit, Kelvin
func (c *temperatureConverter) ConvertFromCelsius(celsius float64) (float64, float64, float64) {
	// Celsius permanece o mesmo
	tempC := celsius

	// Celsius para Fahrenheit: °F = (°C × 9/5) + 32
	tempF := (celsius * 9.0 / 5.0) + 32.0

	// Celsius para Kelvin: K = °C + 273.15
	tempK := celsius + 273.15

	return tempC, tempF, tempK
}
