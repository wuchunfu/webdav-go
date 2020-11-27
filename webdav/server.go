package webdav

import (
	"log"
	"net"
	"net/http"
	"os"
	"strconv"
	"strings"

	wd "golang.org/x/net/webdav"
	"gopkg.in/yaml.v2"

	"webdav/config"
)

func StartServer(c_ *config.DavServer) net.Listener {
	c := &Config{
		DavServer: *c_,
		handler: &wd.Handler{
			Prefix:     "/dav",
			FileSystem: wd.Dir(c_.Scope),
			LockSystem: wd.NewMemLS(),
		},
	}
	// auto ip
	if c.Ip == "auto" {
		c.Ip = GetIP()
	}
	// auto choose port
	listener, err := net.Listen("tcp", c.Ip+":"+strconv.Itoa(int(c.Port)))
	for err != nil && c.Port-c_.Port <= 30 {
		log.Println(err.Error())
		c.Port += 1
		log.Printf("Change port to %d.\n", c.Port)
		listener, err = net.Listen("tcp", c.Ip+":"+strconv.Itoa(int(c.Port)))
	}
	go func() {
		log.Println("Listening on", listener.Addr().String())
		cc, _ := yaml.Marshal(c_)
		log.Println(strings.ReplaceAll(string(cc), "\n", ";"))
		if c.Tls {
			if c.Cert == "" {
				c.Cert = os.Getenv("TLS_CERT")
				if c.Cert == "" {
					c.Cert = "C:\\nginx-1.18.0\\conf\\nolva.pem"
				}
			}
			if c.Key == "" {
				c.Key = os.Getenv("TLS_KEY")
				if c.Key == "" {
					c.Key = "C:\\nginx-1.18.0\\conf\\nolva.key"
				}
			}
			if err = http.ServeTLS(listener, c, c.Cert, c.Key); err != nil {
				log.Println(err)
			}
		} else {
			if err = http.Serve(listener, c); err != nil {
				log.Println(err)
			}
		}
		log.Println("Finished listening...")
	}()
	return listener
}

// GetIP
func GetIP() string {
	ret := ""
	netInterfaces, err := net.Interfaces()
	if err != nil {
		return ret
	}

	for i := 0; i < len(netInterfaces); i++ {
		if (netInterfaces[i].Flags & net.FlagUp) != 0 {
			addrs, _ := netInterfaces[i].Addrs()
			for _, address := range addrs {
				if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
					if ipnet.IP.To4() != nil {
						ret = ipnet.IP.String()
						if strings.Index(ret, "10.") == 0 {
							return ret
						}
					}
				}
			}
		}
	}
	return ret
}
