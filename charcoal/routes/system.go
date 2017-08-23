package routes

import "fmt"
import "strconv"
import "github.com/albrow/forms"

import "github.com/dadleyy/charcoal.api/charcoal/net"
import "github.com/dadleyy/charcoal.api/charcoal/defs"
import "github.com/dadleyy/charcoal.api/charcoal/models"

// FindSystemEmailDomains lists all of the email address domains allowed for login.
func FindSystemEmailDomains(runtime *net.RequestRuntime) *net.ResponseBucket {
	blueprint := runtime.Blueprint()
	var domains []models.SystemEmailDomain

	total, err := blueprint.Apply(&domains)

	if err != nil {
		runtime.Errorf("[system domains find] blueprint err: %s", err.Error())
		return runtime.LogicError("bad-request")
	}

	return runtime.SendResults(total, domains)
}

func CreateSystemEmailDomain(runtime *net.RequestRuntime) *net.ResponseBucket {
	body, err := forms.Parse(runtime.Request)

	if err != nil {
		return runtime.LogicError("invalid-request")
	}

	validator := body.Validator()
	validator.Require("domain")
	validator.MinLength("domain", 2)

	// if the validator picked up errors, add them to the request
	// runtime and then return
	if validator.HasErrors() == true {
		errors := []error{}

		for _, m := range validator.Messages() {
			errors = append(errors, fmt.Errorf("field:%s", m))
		}

		return runtime.SendErrors(errors...)
	}

	domain := models.SystemEmailDomain{Domain: body.Get("domain")}
	cursor := runtime.Model(&domain)
	existing := 0

	if err := cursor.Where("domain = ?", domain.Domain).Count(&existing).Error; err != nil || existing >= 1 {
		return runtime.LogicError("duplicate")
	}

	if err := cursor.Create(&domain).Error; err != nil {
		runtime.Errorf("[create system domain] unable to create: %s", err.Error())
		return runtime.ServerError()
	}

	return runtime.SendResults(1, []models.SystemEmailDomain{domain})
}

func DestroySystemEmailDomain(runtime *net.RequestRuntime) *net.ResponseBucket {
	id, ok := runtime.IntParam("id")

	if ok != true {
		return runtime.LogicError("invalid-id")
	}

	var domain models.SystemEmailDomain

	if err := runtime.First(&domain, id).Error; err != nil {
		return runtime.LogicError("not-found")
	}

	if err := runtime.Delete(&domain).Error; err != nil {
		runtime.Errorf("[delete sys domain] unable to delete domain: %s", err.Error())
		return runtime.ServerError()
	}

	return nil
}

func UpdateSystem(runtime *net.RequestRuntime) *net.ResponseBucket {
	var settings models.SystemSettings

	if err := runtime.FirstOrCreate(&settings, models.SystemSettings{}).Error; err != nil {
		runtime.Errorf("[update system] error: %s", err.Error())
		return runtime.LogicError("not-found")
	}

	body, err := forms.Parse(runtime.Request)

	if err != nil {
		return runtime.LogicError("bad-request")
	}

	updates := make(map[string]interface{})
	restriction := body.KeyExists(defs.EmailDomainRestriction)

	if value := body.Get(defs.EmailDomainRestriction); restriction && (value != "true" && value != "false") {
		return runtime.LogicError("invalid-restriction-value")
	}

	if restriction {
		value, _ := strconv.ParseBool(body.Get(defs.EmailDomainRestriction))
		updates[defs.EmailDomainRestriction] = value
	}

	if len(updates) == 0 {
		return nil
	}

	cursor := runtime.Model(&settings).Where("id = ?", settings.ID)

	if err := cursor.Updates(updates).Error; err != nil {
		runtime.Errorf("[update system] unable to update: %s", err.Error())
		return runtime.ServerError()
	}

	return runtime.SendResults(1, []models.SystemSettings{settings})
}

func PrintSystem(runtime *net.RequestRuntime) *net.ResponseBucket {
	var settings models.SystemSettings

	if admin := runtime.IsAdmin(); admin != true {
		runtime.Warnf("non-admin access of system route: %d", runtime.User.ID)
		return runtime.SendResults(1, "ok")
	}

	if err := runtime.First(&settings).Error; err != nil {
		runtime.Errorf("[system] unable to find settings: %s", err.Error())
		return runtime.LogicError("not-found")
	}

	return runtime.SendResults(1, []models.SystemSettings{settings})
}
