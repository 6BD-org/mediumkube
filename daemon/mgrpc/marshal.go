package mgrpc

import "mediumkube/pkg/models"

func Marshal(d models.Domain) *DomainResp {
	res := DomainResp{
		Name:   d.Name,
		Status: d.Status,
		Ip:     d.IP,
		Reason: d.Reason,
	}
	return &res
}

func MarshalList(ds []models.Domain) []*DomainResp {
	res := make([]*DomainResp, 0)
	for _, d := range ds {
		res = append(res, Marshal(d))
	}
	return res
}
