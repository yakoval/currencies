package update

import (
	"encoding/xml"
	"math"
	"strconv"
	"strings"
)

// Список валют источника.
type CurrencyList struct {
	//Служебный элемент.
	XMLName xml.Name `xml:"ValCurs"`

	// Дата актуализации.
	Date string `xml:"Date,attr"`

	//Список курсов валют.
	List []Currency `xml:"Valute"`
}

// Валюта источника.
type Currency struct {
	//Служебный элемент.
	XMLName xml.Name `xml:"Valute"`

	// Идентификатор во внешнем источнике.
	ID string `xml:"ID,attr"`

	// Числовой код.
	NumCode uint `xml:"NumCode"`

	// Буквенный код валюты.
	CharCode string `xml:"CharCode"`

	// Номинал - количество единиц данной валюты.
	Nominal uint `xml:"Nominal"`

	// Название.
	Name string `xml:"Name"`

	// Количество рублей, содержащихся в номинале.
	Value string `xml:"Value"`
}

func (c *Currency) Rate() (float64, error) {
	value, err := strconv.ParseFloat(strings.Replace(c.Value, ",", ".", 1), 32)
	if err != nil {
		return 0.0, err
	}
	return round(value/float64(c.Nominal), 4), nil
}

func round(n float64, prec int) float64 {
	var rounder float64
	pow := math.Pow(10, float64(prec))
	intermed := n * pow
	_, frac := math.Modf(intermed)
	if frac >= 0.5 {
		rounder = math.Ceil(intermed)
	} else {
		rounder = math.Floor(intermed)
	}

	return rounder / pow
}
