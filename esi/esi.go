package esi

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/arkezar/evehub/model"
	"github.com/spf13/viper"

	"github.com/antihax/goesi"
	"github.com/bradfitz/gomemcache/memcache"
	"github.com/gregjones/httpcache"
	httpmemcache "github.com/gregjones/httpcache/memcache"
)

type Client struct {
	client *goesi.APIClient
}

func New(config *viper.Viper) *Client {
	return &Client{client: createClient(config)}
}

func createClient(config *viper.Viper) *goesi.APIClient {
	cache := memcache.New(config.GetString("memcache.host"))

	transport := httpcache.NewTransport(httpmemcache.NewWithClient(cache))
	transport.Transport = &http.Transport{Proxy: http.ProxyFromEnvironment}
	client := &http.Client{Transport: transport}

	return goesi.NewAPIClient(client, "EVEHUB Maintainer: <arkezar@gmail.com>")
}

func (c Client) GetCorpNames(corpIds ...int32) map[int32]string {
	corpNames := make(map[int32]string, len(corpIds))
	response, _, err := c.client.ESI.UniverseApi.PostUniverseNames(context.Background(), []int32(corpIds), nil)
	if err != nil {
		fmt.Println(err)
	}
	for _, item := range response {
		corpNames[item.Id] = item.Name
	}
	return corpNames
}

func (c Client) GetKillmail(errs chan int32, killmailIds ...model.KillId) []model.Killmail {
	kms := make([]model.Killmail, len(killmailIds))
	for i, id := range killmailIds {
		response, _, err := c.client.ESI.KillmailsApi.GetKillmailsKillmailIdKillmailHash(context.Background(), id.EsiHash, int32(id.EsiID), nil)
		if err != nil {
			errs <- 1
		}
		data, _ := response.MarshalJSON()
		var km model.Killmail
		json.Unmarshal(data, &km)
		kms[i] = km
	}
	return kms
}
