package security

type GeoSecurity interface {
	IsCountryChanged(currentIP, referenceIP string) bool
	IsProviderChanged(currentIP, referenceIP string) bool
	IsDataCenter(ipStr string) bool
}
