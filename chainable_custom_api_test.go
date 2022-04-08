package gorm

import (
	"fmt"
	"testing"
)

type SchoolPlanDto struct {
	Id       string `json:"id"`
	Name     string `json:"name"`
	Describe string `json:"describe"`
	Sort     int    `json:"sort"`
	Status   int    `json:"status"`
}

func TestName(t *testing.T) {
	fmt.Println(structToTag((*SchoolPlanDto)(nil)))

}
