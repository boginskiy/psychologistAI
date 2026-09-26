package jwtservice

type ClaimConf struct {
	TimeLiveToken int
	HostName      string
}

func NewClaimConfig(timeLiveToken int, hostName string) *ClaimConf {
	return &ClaimConf{}
}

func (c *ClaimConf) GetTimeLiveToken() int {
	return c.TimeLiveToken
}

func (c *ClaimConf) GetHostName() string {
	return c.HostName
}
