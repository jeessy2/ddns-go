package dns

import (
	"fmt"
	"sort"
	"strconv"
	"strings"

	"github.com/jeessy2/ddns-go/v6/config"
	"github.com/jeessy2/ddns-go/v6/util"
)

const edgeoneOriginRecordType = "IP_DOMAIN"

type EdgeOneOriginGroup struct {
	GroupId string                `json:"GroupId"`
	Name    string                `json:"Name"`
	Records []EdgeOneOriginRecord `json:"Records"`
}

type EdgeOneOriginGroupResponse struct {
	Response struct {
		Error struct {
			Code    string `json:"Code"`
			Message string `json:"Message"`
		} `json:"Error"`
		OriginGroups []EdgeOneOriginGroup `json:"OriginGroups"`
		TotalCount   int                  `json:"TotalCount"`
	} `json:"Response"`
}

type EdgeOneOriginRecord struct {
	Record string `json:"Record"`
	Type   string `json:"Type"`
	Weight int    `json:"Weight,omitempty"`
}

func (eo *EdgeOne) addUpdateOriginGroups(domainCache config.DomainTuples) {
	for _, domainTuple := range domainCache {
		if domainTuple == nil || domainTuple.Primary == nil || !eo.isOriginGroupDomain(domainTuple.Primary) {
			continue
		}

		zoneId, err := eo.getZoneId(domainTuple.Primary)
		if err != nil {
			util.Log("查询 EdgeOne 站点信息发生异常! %s", err)
			domainTuple.SetUpdateStatus(config.UpdatedFailed)
			continue
		}

		records, err := eo.getDesiredOriginRecords(domainTuple)
		if err != nil {
			util.Log("整理 EdgeOne 源站组记录失败! %s", err)
			domainTuple.SetUpdateStatus(config.UpdatedFailed)
			continue
		}

		originGroup, err := eo.getOriginGroup(domainTuple.Primary, zoneId)
		if err != nil {
			util.Log("查询 EdgeOne 源站组信息发生异常! %s", err)
			domainTuple.SetUpdateStatus(config.UpdatedFailed)
			continue
		}

		eo.modifyOriginGroup(originGroup, domainTuple, zoneId, records)
	}
}

func (eo *EdgeOne) isOriginGroupDomain(domain *config.Domain) bool {
	if domain == nil {
		return false
	}
	params := domain.GetCustomParams()
	return params.Has("GroupId") || params.Has("OriginGroupName")
}

func (eo *EdgeOne) getZoneId(domain *config.Domain) (string, error) {
	params := domain.GetCustomParams()
	if params.Has("ZoneId") {
		return params.Get("ZoneId"), nil
	}

	zoneResult, err := eo.getZone(domain.DomainName)
	if err != nil {
		return "", err
	}
	if zoneResult.Response.TotalCount <= 0 {
		return "", fmt.Errorf("在 EdgeOne 中未找到站点: %s", domain.DomainName)
	}
	for _, zone := range zoneResult.Response.Zones {
		if zone.ZoneName == domain.DomainName {
			return zone.ZoneId, nil
		}
	}
	return "", fmt.Errorf("在 EdgeOne 中未找到站点: %s", domain.DomainName)
}

func (eo *EdgeOne) getDesiredOriginRecords(domainTuple *config.DomainTuple) ([]EdgeOneOriginRecord, error) {
	domain := domainTuple.Primary
	domainName := domain.String()
	weight := eo.getOriginWeight(domain)
	records := make([]EdgeOneOriginRecord, 0, 2)

	if eo.hasDomain(eo.Domains.Ipv4Domains, domainName) {
		if eo.Domains.Ipv4Addr == "" {
			return nil, fmt.Errorf("未能获取域名 %s 对应的 IPv4 地址", domain)
		}
		records = append(records, EdgeOneOriginRecord{
			Record: eo.Domains.Ipv4Addr,
			Type:   edgeoneOriginRecordType,
			Weight: weight,
		})
	}
	if eo.hasDomain(eo.Domains.Ipv6Domains, domainName) {
		if eo.Domains.Ipv6Addr == "" {
			return nil, fmt.Errorf("未能获取域名 %s 对应的 IPv6 地址", domain)
		}
		records = append(records, EdgeOneOriginRecord{
			Record: eo.Domains.Ipv6Addr,
			Type:   edgeoneOriginRecordType,
			Weight: weight,
		})
	}
	if len(records) == 0 {
		return nil, fmt.Errorf("域名 %s 未配置可更新的源站记录", domain)
	}

	return records, nil
}

func (eo *EdgeOne) hasDomain(domains []*config.Domain, fullDomain string) bool {
	for _, domain := range domains {
		if domain != nil && domain.String() == fullDomain {
			return true
		}
	}
	return false
}

func (eo *EdgeOne) getOriginWeight(domain *config.Domain) int {
	weight := 100
	if domain == nil {
		return weight
	}

	if s := domain.GetCustomParams().Get("Weight"); s != "" {
		if v, err := strconv.Atoi(s); err == nil && v > 0 {
			weight = v
		}
	}
	return weight
}

func (eo *EdgeOne) getOriginGroup(domain *config.Domain, zoneId string) (EdgeOneOriginGroup, error) {
	params := domain.GetCustomParams()
	record := struct {
		ZoneId  string   `json:"ZoneId"`
		Filters []Filter `json:"Filters"`
	}{
		ZoneId: zoneId,
	}

	if params.Has("GroupId") {
		record.Filters = []Filter{{Name: "origin-group-id", Values: []string{params.Get("GroupId")}}}
	} else if params.Has("OriginGroupName") {
		record.Filters = []Filter{{Name: "origin-group-name", Values: []string{params.Get("OriginGroupName")}}}
	} else {
		return EdgeOneOriginGroup{}, fmt.Errorf("请在域名后追加 ?GroupId=xxx 或 ?OriginGroupName=xxx")
	}

	var result EdgeOneOriginGroupResponse
	if err := eo.request("DescribeOriginGroup", record, &result); err != nil {
		return EdgeOneOriginGroup{}, err
	}
	if result.Response.Error.Code != "" {
		return EdgeOneOriginGroup{}, fmt.Errorf("%s", result.Response.Error.Message)
	}
	if result.Response.TotalCount <= 0 || len(result.Response.OriginGroups) == 0 {
		return EdgeOneOriginGroup{}, fmt.Errorf("在 EdgeOne 中未找到源站组: %s", domain)
	}

	if params.Has("GroupId") {
		groupId := params.Get("GroupId")
		for _, group := range result.Response.OriginGroups {
			if group.GroupId == groupId {
				return group, nil
			}
		}
		return EdgeOneOriginGroup{}, fmt.Errorf("在 EdgeOne 中未找到源站组 GroupId=%s", groupId)
	}

	groupName := params.Get("OriginGroupName")
	for _, group := range result.Response.OriginGroups {
		if group.Name == groupName {
			return group, nil
		}
	}
	if len(result.Response.OriginGroups) == 1 {
		return result.Response.OriginGroups[0], nil
	}
	return EdgeOneOriginGroup{}, fmt.Errorf("找到多个名称匹配的源站组，请改用 GroupId 指定唯一源站组")
}

func (eo *EdgeOne) modifyOriginGroup(originGroup EdgeOneOriginGroup, domainTuple *config.DomainTuple, zoneId string, records []EdgeOneOriginRecord) {
	if sameEdgeOneOriginRecords(originGroup.Records, records) {
		util.Log("你的IP %s 没有变化, EdgeOne 源站组 %s", strings.Join(edgeOneOriginRecordValues(records), ","), originGroup.Name)
		return
	}

	var status EdgeOneStatus
	err := eo.request(
		"ModifyOriginGroup",
		struct {
			ZoneId  string                `json:"ZoneId"`
			GroupId string                `json:"GroupId"`
			Records []EdgeOneOriginRecord `json:"Records"`
		}{
			ZoneId:  zoneId,
			GroupId: originGroup.GroupId,
			Records: records,
		},
		&status,
	)
	if err != nil {
		util.Log("更新 EdgeOne 源站组 %s 失败! 异常信息: %s", originGroup.Name, err)
		domainTuple.SetUpdateStatus(config.UpdatedFailed)
		return
	}

	if status.Response.Error.Code == "" {
		util.Log("更新 EdgeOne 源站组 %s 成功! IP: %s", originGroup.Name, strings.Join(edgeOneOriginRecordValues(records), ","))
		domainTuple.SetUpdateStatus(config.UpdatedSuccess)
	} else {
		util.Log("更新 EdgeOne 源站组 %s 失败! 异常信息: %s", originGroup.Name, status.Response.Error.Message)
		domainTuple.SetUpdateStatus(config.UpdatedFailed)
	}
}

func sameEdgeOneOriginRecords(current []EdgeOneOriginRecord, desired []EdgeOneOriginRecord) bool {
	if len(current) != len(desired) {
		return false
	}

	currentKeys := edgeOneOriginRecordKeys(current)
	desiredKeys := edgeOneOriginRecordKeys(desired)
	for i := range currentKeys {
		if currentKeys[i] != desiredKeys[i] {
			return false
		}
	}
	return true
}

func edgeOneOriginRecordKeys(records []EdgeOneOriginRecord) []string {
	keys := make([]string, 0, len(records))
	for _, record := range records {
		keys = append(keys, fmt.Sprintf("%s|%s|%d", record.Record, record.Type, record.Weight))
	}
	sort.Strings(keys)
	return keys
}

func edgeOneOriginRecordValues(records []EdgeOneOriginRecord) []string {
	values := make([]string, 0, len(records))
	for _, record := range records {
		values = append(values, record.Record)
	}
	sort.Strings(values)
	return values
}
