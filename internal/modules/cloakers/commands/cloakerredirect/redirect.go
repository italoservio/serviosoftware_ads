package cloakerredirect

import (
	"encoding/base64"
	"fmt"
	"log"
	"net"
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

		if ipLookup.IsBlacklisted {
			return &RedirectCloakerOutput{RedirectURL: cloaker.WhiteURL}, nil
		}

		return &RedirectCloakerOutput{RedirectURL: cloaker.BlackURL}, nil
	}

	netifyResponse, err := c.NetifyClient.GetIPMetadata(input.UserIP)
	if err != nil {
		return nil, exception.NewClientException(err.Error())
	}

	netifyData := (*netifyResponse).Data
	netifyApplications := extractNetifyApplicationsLabels(netifyData)

	if isFromKnownApplication(netifyApplications, input.UserAgent) || netifyData.SharedScore < 60 {
		isBlacklisted := true
		go c.createIPLookupFromNetifyData(
			input.UserIP,
			ipPattern,
			netifyApplications,
			netifyData,
			isBlacklisted,
		)

		return &RedirectCloakerOutput{RedirectURL: cloaker.WhiteURL}, nil
	}

	isBlacklisted := false
	go c.createIPLookupFromNetifyData(
		input.UserIP,
		ipPattern,
		netifyApplications,
		netifyData,
		isBlacklisted,
	)

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

func extractNetifyApplicationsLabels(netifyData clients.NetifySuccessResponseData) []string {
	labels := make([]string, len(netifyData.ApplicationList))
	for i, app := range netifyData.ApplicationList {
		labels[i] = app.Label
	}
	return labels
}

func (c *RedirectCloakerCmd) createIPLookupFromNetifyData(
	ip string,
	ipPattern string,
	applications []string,
	netifyData clients.NetifySuccessResponseData,
	isBlacklisted bool,
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
		IsBlacklisted: isBlacklisted,
	}

	_, err := c.IPLookupRepo.Create(newIPLookup)
	if err != nil {
		log.Printf("falha ao criar ip %s: %v", ip, err)
	}
}

func ipToPattern(ip string) string {
	// remove port ([IPv6]:port)
	if strings.HasPrefix(ip, "[") {
		if i := strings.Index(ip, "]"); i != -1 {
			ip = ip[1:i]
		}
	}

	parsedIP := net.ParseIP(ip)
	if parsedIP == nil {
		return ip
	}

	// IPv4: 192.168.1.100 -> 192.168.1.*
	if parsedIP.To4() != nil {
		ipStr := parsedIP.String()
		if i := strings.LastIndex(ipStr, "."); i != -1 {
			return ipStr[:i] + ".*"
		}
	}

	// IPv6: get first 48 bits (3 hextets)
	// parsedIP.To16() give us the full representation with 16 bytes
	ipv6 := parsedIP.To16()
	if ipv6 != nil {
		// format first 6 bytes (48 bits) as default
		// each hextet has 2 bytes
		return fmt.Sprintf("%02x%02x:%02x%02x:%02x%02x:*",
			ipv6[0], ipv6[1], ipv6[2], ipv6[3], ipv6[4], ipv6[5])
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
