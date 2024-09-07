package internal_sys

import (
	internalHttp "github.com/odycenter/std-library/app/internal/web/http"
	"github.com/odycenter/std-library/app/property"
	"github.com/odycenter/std-library/app/web/errors"
	"github.com/odycenter/std-library/nets"
	"net/http"
	"os"
	"strconv"
	"strings"
)

type PropertyController struct {
	propertyManager *property.Manager
	accessControl   *internalHttp.IPv4AccessControl
}

func NewPropertyController(propertyManager *property.Manager) *PropertyController {
	return &PropertyController{
		propertyManager: propertyManager,
		accessControl:   &internalHttp.IPv4AccessControl{},
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
	if c.propertyManager.DefaultHTTPPort > 0 {
		sb.WriteString("# No http listen port configured, using port to start HTTP server\n")
		sb.WriteString("http.port=")
		sb.WriteString(strconv.Itoa(c.propertyManager.DefaultHTTPPort) + " (can not override by ENV)\n")
		sb.WriteString("\n")
	}
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
