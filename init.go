package afs

import (
	"github.com/viant/afs/file"
	"github.com/viant/afs/http"
	"github.com/viant/afs/mem"
	"github.com/viant/afs/scp"
)

func init() {
	registry := GetRegistry()
	registry.Register(file.Scheme, file.Provider)
	registry.Register(mem.Scheme, mem.Provider)
	registry.Register(http.Scheme, http.Provider)
	registry.Register(http.SecureScheme, http.Provider)
	registry.Register(scp.Scheme, scp.Provider)
}
