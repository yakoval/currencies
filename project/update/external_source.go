package update

import (
	"encoding/xml"
	"fmt"
	"io"

	"golang.org/x/text/encoding/charmap"
)

type externalSource struct {
	httpSource httpSource
}

func (es *externalSource) data() ([]Currency, error) {
	rawDataReader, err := es.httpSource.content()
	if err != nil {
		return nil, err
	}

	decoder := xml.NewDecoder(rawDataReader)
	decoder.CharsetReader = func(charset string, input io.Reader) (io.Reader, error) {
		switch charset {
		case "windows-1251":
			return charmap.Windows1251.NewDecoder().Reader(input), nil
		default:
			return nil, fmt.Errorf("unknown charset: %s", charset)
		}
	}

	currencyList := &CurrencyList{}
	err = decoder.Decode(currencyList)
	if err != nil {
		return nil, err
	}

	return currencyList.List, nil
}
