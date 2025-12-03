package cloakerredirect

import (
	"encoding/base64"
	"log"
	"strings"

	"github.com/houseme/mobiledetect/ua"
	"github.com/italoservio/serviosoftware_ads/internal/modules/cloakers/clients"
	"github.com/italoservio/serviosoftware_ads/internal/modules/cloakers/models"
	"github.com/italoservio/serviosoftware_ads/internal/modules/cloakers/repos"
	"github.com/italoservio/serviosoftware_ads/pkg/exception"
)

type RedirectCloakerInput struct {
	EncodedID string `bson:"encodedId" json:"encodedId" validate:"required"`
	UserIP    string `bson:"userIp" json:"userIp" validate:"required,ip"`
	UserAgent string `bson:"userAgent" json:"userAgent" validate:"required"`
}

type RedirectCloakerOutput struct {
	RedirectURL string `bson:"redirectUrl" json:"redirectUrl"`
}

type RedirectCloakerCmd struct {
	IPLookupRepo repos.IPLookupRepository
	CloakerRepo  repos.CloakerRepository
	NetifyClient clients.NetifyClient
}

func NewRedirectCloakerCmd(
	ipLookupRepo repos.IPLookupRepository,
	cloakerRepo repos.CloakerRepository,
	netifyClient clients.NetifyClient,
) *RedirectCloakerCmd {
	return &RedirectCloakerCmd{
		IPLookupRepo: ipLookupRepo,
		CloakerRepo:  cloakerRepo,
		NetifyClient: netifyClient,
	}
}

func (c *RedirectCloakerCmd) Exec(input *RedirectCloakerInput) (*RedirectCloakerOutput, error) {
	cloaker, err := c.getCloakerByEncodedID(input.EncodedID)
	if err != nil {
		return nil, err
	}

	if !cloaker.IsActive {
		return nil, exception.NewNotFoundException()
	}

	agent := ua.New(input.UserAgent)
	if agent.Bot() {
		return &RedirectCloakerOutput{RedirectURL: cloaker.WhiteURL}, nil
	}

	if cloaker.Config.AllowOnlyMobile && !agent.Mobile() {
		return &RedirectCloakerOutput{RedirectURL: cloaker.WhiteURL}, nil
	}

	ipPattern := ipToPattern(input.UserIP)
	ipLookup, err := c.IPLookupRepo.GetByIPPattern(ipPattern)
	if err != nil {
		return nil, exception.NewRepoException("nao foi possivel consultar o ip", err.Error())
	}

	if ipLookup != nil {
		go c.incrementAccessCount(ipLookup.StringID())
		return &RedirectCloakerOutput{RedirectURL: cloaker.WhiteURL}, nil
	}

	netifyResponse, err := c.NetifyClient.GetIPMetadata(input.UserIP)
	if err != nil {
		return nil, exception.NewClientException(err.Error())
	}

	netifyData := netifyResponse.Data
	var netifyApplications []string

	if len(netifyData.ApplicationList) > 0 {
		netifyApplications = make([]string, len(netifyData.ApplicationList))

		for i, app := range netifyData.ApplicationList {
			netifyApplications[i] = app.Label
		}

		if isFromKnownApplication(netifyApplications, input.UserAgent) || netifyData.SharedScore < 60 {
			go c.createIPLookupFromNetifyData(input.UserIP, ipPattern, netifyApplications, netifyData)

			return &RedirectCloakerOutput{RedirectURL: cloaker.WhiteURL}, nil
		}
	}

	return &RedirectCloakerOutput{RedirectURL: cloaker.BlackURL}, nil
}

func (c *RedirectCloakerCmd) getCloakerByEncodedID(encodedID string) (*models.Cloaker, error) {
	var cloakerID string

	decodedBytes, err := base64.URLEncoding.DecodeString(encodedID)
	if err != nil {
		return nil, exception.NewInternalException("falha ao decodificar id", err.Error())
	}

	cloakerID = string(decodedBytes)

	cloaker, err := c.CloakerRepo.GetByID(cloakerID)
	if err != nil {
		return nil, err
	}

	if cloaker == nil {
		return nil, exception.NewNotFoundException()
	}

	return cloaker, nil
}

func (c *RedirectCloakerCmd) incrementAccessCount(id string) {
	err := c.IPLookupRepo.IncrementAccessCount(id)
	if err != nil {
		log.Printf("falha ao incrementar contagem de acessos para o ip %s: %v", id, err)
	}
}

func (c *RedirectCloakerCmd) createIPLookupFromNetifyData(
	ip string,
	ipPattern string,
	applications []string,
	netifyData clients.NetifySuccessResponseData,
) {
	if len(applications) == 0 {
		applications = []string{"unknown"}
	}

	newIPLookup := &models.IPLookup{
		IP:           ip,
		IPPattern:    ipPattern,
		Applications: applications,
		SharedScore:  netifyData.SharedScore,
		GeoLocation: models.IPLookupGeoLocation{
			Country:   netifyData.Geolocation.Country.Label,
			Continent: netifyData.Geolocation.Continent.Label,
		},
	}

	_, err := c.IPLookupRepo.Create(newIPLookup)
	if err != nil {
		log.Printf("falha ao criar ip %s: %v", ip, err)
	}
}

func ipToPattern(ip string) string {
	// IPv4: 192.168.1.100 -> 192.168.1.*
	if strings.Contains(ip, ".") {
		if i := strings.LastIndex(ip, "."); i != -1 {
			return ip[:i] + ".*"
		}
	}

	// IPv6: 2001:0db8:85a3::8a2e:0370:7334 -> 2001:0db8:85a3:*
	if strings.Contains(ip, ":") {
		parts := strings.Split(ip, ":")
		if len(parts) > 3 {
			return strings.Join(parts[:3], ":") + ":*"
		}
	}

	return ip
}

func isFromKnownApplication(applications []string, userAgent string) bool {
	knownApplications := []string{
		"amazon", "apple", "facebook", "google", "instagram",
		"kwai", "tiktok", "whatsapp", "youtube",
	}

	lowerUserAgent := strings.ToLower(userAgent)

	for _, app := range knownApplications {
		if strings.Contains(lowerUserAgent, app) {
			return true
		}

		for _, application := range applications {
			if strings.Contains(strings.ToLower(application), app) {
				return true
			}
		}
	}

	return false
}
