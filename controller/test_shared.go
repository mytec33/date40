package controller

import "date_calculation/models"

type ResponseWrapper struct {
	Results models.OutputResults `json:"results"`
}
