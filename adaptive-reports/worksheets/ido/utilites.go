package ido

import (
	"github.com/adaptiveteam/adaptive/adaptive-reports/models"
)

func completedString(b bool) (rv string) {
	if b {
		rv = "Closed"
	} else {
		rv = "Active"
	}
	return rv
}

func getUpdateRow(ido *models.IDO, allIDOs models.IDOs) (rv int) {
	for i := 0; i < len(allIDOs) && ido.Name() != allIDOs[i].Name(); i++ {
		if ido.Completed() == allIDOs[i].Completed() {
			rv = rv + len(allIDOs[i].Updates())
		}
	}
	return rv
}

func getStatusStyle(state bool, status string) (rv string) {
	if state {
		rv = "Closed Background"
	} else {
		rv = status
	}
	return rv
}