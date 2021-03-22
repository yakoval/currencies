package update

import (
	"immo-currencies/project"
)

// Менеджер обновления валют.
type Manager struct {
	// Счетчик валют, полученных из внешнего источника.
	ReadCounter int

	// Счетчик измененных валют.
	UpdateCounter int

	// Счетчик добавленных валют.
	InsertCounter int

	// Счетчик удаленных валют.
	DeleteCounter int

	db             *project.Database
	externalSource *externalSource
	existed        map[string]project.Currency
}

// Создает новый менеджер обновления валют.
func NewUpdater(db *project.Database, conf *project.UpdaterConfig) (*Manager, error) {
	m := &Manager{
		db: db,
		externalSource: &externalSource{
			httpSource: httpSource{uri: conf.URI},
		},
		existed: map[string]project.Currency{},
	}
	if err := m.fillExistedRates(); err != nil {
		return nil, err
	}

	return m, nil
}

func (m *Manager) Work() error {
	currencies, err := m.externalSource.data()
	if err != nil {
		return err
	}

	for _, currency := range currencies {
		m.ReadCounter++
		existedCurrency, isExist := m.existed[currency.Name]
		if isExist {
			if err = m.update(&currency, existedCurrency.Rate); err != nil {
				return err
			}
			delete(m.existed, currency.Name)
		} else {
			internalCurrency, err := convertToInternal(&currency)
			if err != nil {
				return err
			}
			if err = m.db.Insert(internalCurrency); err != nil {
				return err
			}
			m.InsertCounter++
		}
	}

	//Удаление старых записей.
	for _, oldCurrency := range m.existed {
		if err = m.db.DeleteByID(oldCurrency.ID); err != nil {
			return err
		}
		m.DeleteCounter++
	}

	return nil
}

func (m *Manager) update(currency *Currency, rateOld float64) error {
	rate, err := currency.Rate()
	if err != nil {
		return err
	}
	if rateOld != rate {
		internalCurrency, err := convertToInternal(currency)
		if err != nil {
			return err
		}
		err = m.db.Update(internalCurrency)
		if err != nil {
			return err
		}
		m.UpdateCounter++
	}
	return nil
}

func (m *Manager) fillExistedRates() error {
	existedList, err := m.db.GetAll()
	if err != nil {
		return err
	}

	for _, currency := range existedList {
		m.existed[currency.Name] = currency
	}
	return nil
}
