package websites

import "context"
import "net"
import "time"

func (tool *Websites) canConnect(host string) bool {

	if tool.Capabilities != nil && tool.Capabilities.HasNetwork() == true {

		if tool.Capabilities.CanConnect(host) == true {
			return true
		} else {

			if net.ParseIP(host) == nil {

				ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
				defer cancel()

				addresses, err0 := net.DefaultResolver.LookupIPAddr(ctx, host)

				if err0 == nil {

					for _, address := range addresses {

						if tool.Capabilities.CanConnectIP(address.IP) == true {
							return true
						}

					}

				}

			}

		}
	}

	return false

}
