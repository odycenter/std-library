package internal_sys

import (
	"net/http"
	"os"
	internal_http "std-library/app/internal/web/http"
	"std-library/app/property"
	"std-library/app/web/errors"
	"std-library/nets"
	"strings"
)

type PropertyController struct {
	propertyManager *property.Manager
	accessControl   *internal_http.IPv4AccessControl
}

func NewPropertyController(propertyManager *property.Manager) *PropertyController {
	return &PropertyController{
		propertyManager: propertyManager,
		accessControl:   &internal_http.IPv4AccessControl{},
	}
}

func (c *PropertyController) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	err := c.accessControl.Validate(nets.IP(r).String())
	if err != nil {
		errors.Forbidden("access denied", "IP_ACCESS_DENIED")
	}

	if r.Method != http.MethodGet {
		w.WriteHeader(404)
		return
	}
	w.WriteHeader(200)
	w.Write([]byte(c.propertiesString()))
}

func (c *PropertyController) propertiesString() string {
	var sb strings.Builder
	sb.WriteString("# properties\n")
	for _, key := range c.propertyManager.GetKeysInOrder() {
		sb.WriteString(key)
		sb.WriteString("=")
		sb.WriteString(property.MaskValue(key, c.propertyManager.Get(key)))
		sb.WriteString("\n")
	}
	sb.WriteString("\n# env variables\n")
	for _, env := range os.Environ() {
		splitEnv := strings.SplitN(env, "=", 2)
		key := splitEnv[0]
		sb.WriteString(key)
		sb.WriteString("=")
		sb.WriteString(property.MaskValue(key, os.Getenv(key)))
		sb.WriteString("\n")
	}
	return sb.String()
}
