package security

import (
	"errors"
	"net"
	"strings"

	"github.com/oschwald/geoip2-golang"
)

// Список самых популярных облачных провайдеров.
var DataCenters = map[string]struct{}{
	"hetzner":         struct{}{},
	"digitalocean":    struct{}{},
	"amazon aws":      struct{}{},
	"google cloud":    struct{}{},
	"microsoft azure": struct{}{},
	"oracle cloud":    struct{}{},
	"linode":          struct{}{},
	"vultr":           struct{}{},
	"ovh":             struct{}{},
	"choopa":          struct{}{},
	"leaseweb":        struct{}{},
	"contina":         struct{}{},
	"quadranet":       struct{}{},
	"psychz":          struct{}{},
	"akamai":          struct{}{},
	"cloudflare":      struct{}{}}

type GeoChecker struct {
	GeoCountryDB *geoip2.Reader
	GeoAsnDB     *geoip2.Reader
}

func NewGeoChecker(pathCountryDB, pathAsnDB string) (*GeoChecker, error) {
	db1, err := geoip2.Open(pathCountryDB)
	if err != nil {
		return nil, err
	}

	db2, err := geoip2.Open(pathAsnDB)
	if err != nil {
		return nil, err
	}

	return &GeoChecker{
		GeoCountryDB: db1,
		GeoAsnDB:     db2,
	}, nil
}

func (c *GeoChecker) Close() error {
	c.GeoCountryDB.Close()
	c.GeoAsnDB.Close()
	return nil
}

// isDatacenter возвращает true, если IP принадлежит известному облаку
func (c *GeoChecker) IsProviderChanged(currentIP, referenceIP string) bool {
	if currentIP == "" || referenceIP == "" {
		return false
	}

	asnCurrentIP, err1 := c.getASN(currentIP)
	asnReferenceIP, err2 := c.getASN(referenceIP)

	if err1 != nil || err2 != nil {
		// + logger
		return false
	}

	return asnCurrentIP != asnReferenceIP
}

func (c *GeoChecker) getASN(ipStr string) (string, error) {
	ip := net.ParseIP(ipStr)
	if ip == nil {
		return "", errors.New("IP is not valid")
	}

	record, err := c.GeoAsnDB.ASN(ip)
	if err != nil {
		return "", err
	}

	return strings.ToLower(record.AutonomousSystemOrganization), nil
}

func (c *GeoChecker) IsCountryChanged(currentIP, referenceIP string) bool {
	// Если один из адресов пустой, считаем что смены нет (или ошибка ввода)
	if currentIP == "" || referenceIP == "" {
		return false
	}

	//
	isoCodeCurrentIP, err1 := c.getISOCode(currentIP)
	isoCodeReferenceIP, err2 := c.getISOCode(referenceIP)
	if err1 != nil || err2 != nil {
		// + logger
		return false
	}
	return isoCodeCurrentIP != isoCodeReferenceIP
}

func (c *GeoChecker) getISOCode(ipStr string) (string, error) {
	ip := net.ParseIP(ipStr)
	if ip == nil {
		return "", errors.New("IP is not valid")
	}

	record, err := c.GeoCountryDB.Country(ip)
	if err != nil {
		return "", err
	}

	// Возвращаем код, например "RU"
	return record.Country.IsoCode, nil
}
