package ddns

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"toolbox/common"
	"toolbox/config"

	cloudflare "github.com/cloudflare/cloudflare-go"
)

// 存储更新后的 DNS 记录
var recordAddr string

type cloudflareDnsApi struct {
	apiKey     string
	email      string
	domain     string
	recordName string
	dnsType    string
	content    string
	ttl        int // s
}

func newCloudflareDnsApi(content string) *cloudflareDnsApi {
	return &cloudflareDnsApi{
		apiKey:     config.Config.ApiKey,
		email:      config.Config.Email,
		domain:     config.Config.Domain,
		recordName: config.Config.RecordName,
		dnsType:    "A",
		content:    content,
		ttl:        60,
	}
}

func (c *cloudflareDnsApi) ddns(isInit bool) (newRecord *string, err error) {
	// Construct a new API object using a global API key
	api, err := cloudflare.New(c.apiKey, c.email)
	if err != nil {
		return nil, err
	}

	ctx := context.Background()

	// Fetch the zone ID for zone example.org
	zoneID, err := api.ZoneIDByName(c.domain)
	if err != nil {
		return nil, err
	}
	// 返回的 DNSRecord 是列表，但使用完全限定域名作为查询参数，因此只会返回唯一的一条 DNSRecord
	// 完全限定域名使用 recordName + domain 进行拼接
	listParams := cloudflare.ListDNSRecordsParams{
		Type: c.dnsType,
		Name: strings.Join([]string{c.recordName, c.domain}, "."),
	}
	records, _, err := api.ListDNSRecords(ctx, cloudflare.ZoneIdentifier(zoneID), listParams)
	if err != nil {
		return nil, err
	}

	// 如果是初始化调用，则将查询的记录赋值给全局变量后 return，不往下执行 Update
	if isInit {
		recordAddr = records[0].Content
		return &recordAddr, nil
	}

	// Update 必须使用 Record ID，而不是 Name
	// 根据上面定义的查询结果只有一条 DNSRecord，因此直接使用索引取 Record ID
	updateParams := cloudflare.UpdateDNSRecordParams{
		Type:    c.dnsType,
		Content: c.content,
		ID:      records[0].ID,
		TTL:     c.ttl,
	}
	newRecords, err := api.UpdateDNSRecord(ctx, cloudflare.ZoneIdentifier(zoneID), updateParams)
	if err != nil {
		return nil, err
	}

	// 将更新后的记录赋值给全局变量
	recordAddr = newRecords.Content

	return &recordAddr, nil
}

func getIP() (ip *string, err error) {
	var addr string
	res, err := http.Get(config.Config.Endpoint)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		return nil, errors.New("StatusCode error")
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	addr = string(body)
	return &addr, nil
}

// 执行初始化，将首次查询的 DNS 记录赋值给全局变量
// 该初始化调用 API 时需要依赖配置，因此必须在配置加载成功后调用
func IsInit() {
	api := newCloudflareDnsApi("")
	record, err := api.ddns(true)
	if err != nil {
		log.Println(err)
		return
	}
	log.Printf("ddns init successful: %s\n", *record)
}

func RunDDNS() (message *string, err error) {
	nowAddr, err := getIP()
	if err != nil {
		return nil, err
	}

	// 判断当前地址与全局变量中的 DNS 记录，如没有变更直接 return
	if *nowAddr == recordAddr {
		return nil, nil
	}

	api := newCloudflareDnsApi(*nowAddr)
	newRecord, err := api.ddns(false)
	if err != nil {
		return nil, err
	}

	text := fmt.Sprintf("DNS Record: %s", *newRecord)
	// 发送通知到 telegram
	notification := common.DefaultNotify()
	notification.SendToTelegram(common.TEXT, []byte(text))
	return &text, nil
}
