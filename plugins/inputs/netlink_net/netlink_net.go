// +build linux
package netlink_net

/*
Had to change something in docker/libnetwork because when I used the latest "github.com/vishvananda/netlink", the ipvs/netlink.go would break

Change line 220 in /root/go/pkg/mod/github.com/docker/libnetwork/ipvs/netlink.go to this: to use latest netlink code
msgs, _, err := s.Receive()
*/

import (
	"github.com/influxdata/telegraf"
	"github.com/influxdata/telegraf/plugins/inputs"
	"github.com/influxdata/telegraf/filter"    
    "fmt"
    "github.com/vishvananda/netlink"
	"sync"
    "strconv"
)

type Netlink struct {
	// This is the list of interface names to include
	InterfaceInclude []string `toml:"interface_include"`

	// This is the list of interface names to ignore
	InterfaceExclude []string `toml:"interface_exclude"`

	Log telegraf.Logger `toml:"-"`
    
    ShowVirtualFunctions bool `toml:"show_virtual_functions"`
    Verbose bool `toml:"verbose`
}

const (
	pluginName    = "netlink_net"
	tagInterface  = "interface"
	tagDriverName = "driver"

	sampleConfig = `
  ## List of interfaces to pull metrics for
  # interface_include = ["eth0"]

  ## List of interfaces to ignore when pulling metrics.
  # interface_exclude = ["eth1"]
`
)


func NewNetlink() * Netlink {
    
    obj :=  new (Netlink)

    return obj
}

func (s *Netlink) SampleConfig() string {
	return sampleConfig
}

func (s *Netlink) Description() string {
	return "Gathers the CPU Frequency for each core"
}

func generateVirtualFunctionInfo(vfInfo netlink.VfInfo, link netlink.Link) (map[string]interface{}, map[string]string) {
/*type VfInfo struct {
	ID        int
	Mac       net.HardwareAddr
	Vlan      int
	Qos       int
	TxRate    int // IFLA_VF_TX_RATE  Max TxRate
	Spoofchk  bool
	LinkState uint32
	MaxTxRate uint32 // IFLA_VF_RATE Max TxRate
	MinTxRate uint32 // IFLA_VF_RATE Min TxRate
}
*/
    tags := make(map[string]string)

    fields :=  make(map[string]interface{})    
    tags["pf-name"] = link.Attrs().Name
    fields["pf-mac_address"] = link.Attrs().HardwareAddr.String()
    tags["vf_number"] = strconv.Itoa(vfInfo.ID)
    tags["mac_address"] = vfInfo.Mac.String()
    fields["link"] = vfInfo.LinkState
    fields["spoof-check"] = vfInfo.Spoofchk
    fields["tx_rate"] = vfInfo.TxRate
    fields["max-tx_rate"] = vfInfo.MaxTxRate
    fields["min-tx_rate"] = vfInfo.MinTxRate
    fields["vlan"] = vfInfo.Vlan
    fields["qos"] = vfInfo.Qos
    fields["tx_packets"] = vfInfo.TxPackets
    fields["rx_packets"] = vfInfo.RxPackets
    fields["tx_bytes"] = vfInfo.TxBytes
    fields["rx_bytes"] = vfInfo.RxBytes
    fields["tx_multicast"] = vfInfo.Multicast
    fields["tx_broadcast"] = vfInfo.Broadcast
    fields["tx_dropped"] = vfInfo.TxDropped
    fields["rx_dropped"] = vfInfo.RxDropped    
    fields["trust"] = vfInfo.Trust

    return fields,tags
}
func (nl *Netlink) gatherNetlinkNetworkStats(link netlink.Link, acc telegraf.Accumulator) error {
        
        fields :=  make(map[string]interface{})    
        
        tags := make(map[string]string)
            
        if true == nl.Verbose {
            fields["mtu"] = link.Attrs().MTU
            fields["tx_queue_length"] = link.Attrs().TxQLen
            tags["mac address"] = link.Attrs().HardwareAddr.String()
            //fields["Flags"] = link.Attrs().Flags
//                fields["RawFlags"] = link.Attrs().RawFlags
//                fields["ParentIndex"] = link.Attrs().ParentIndex
//                fields["MasterIndex"] = link.Attrs().MasterIndex
            if nil != link.Attrs().Namespace {
                fields["namespace"] = link.Attrs().Namespace
            } else {
              fields["namespace"] = ""
            }
            fields["alias"] = link.Attrs().Alias
            fields["promisc"] = link.Attrs().Promisc
            //fields["Xdp"] = link.Attrs().Xdp
            fields["encap_type"] = link.Attrs().EncapType
            //fields["proto_info"] = link.Attrs().Protinfo
            fields["link"] = link.Attrs().OperState
            fields["NetNsID"] = link.Attrs().NetNsID
            
            fields["num_tx_queues"] = link.Attrs().NumTxQueues
            fields["num_rx_queues"] = link.Attrs().NumRxQueues
            fields["gso_max_size"] = link.Attrs().GSOMaxSize
            fields["gso_max_segments"] = link.Attrs().GSOMaxSegs
//                fields["Group"] = link.Attrs().Group
//                fields["Slave"] = link.Attrs().Slave
        }
        tags["name"] = link.Attrs().Name

        fields["rx_Packets"] = link.Attrs().Statistics.RxPackets                
        fields["tx_Packets"] = link.Attrs().Statistics.TxPackets
        fields["tx_bytes"] = link.Attrs().Statistics.RxBytes
        fields["tx_bytes"] = link.Attrs().Statistics.TxBytes
        fields["rx_errors"] = link.Attrs().Statistics.RxErrors
        fields["tx_errors"] = link.Attrs().Statistics.TxErrors
        fields["rx_dropped"] = link.Attrs().Statistics.RxDropped
        fields["tx_dropped"] = link.Attrs().Statistics.TxDropped
        fields["multicast"] = link.Attrs().Statistics.Multicast
        fields["collisions"] = link.Attrs().Statistics.Collisions
        fields["rx-length_errors"] = link.Attrs().Statistics.RxLengthErrors
        fields["rx-overflow_Errors"] = link.Attrs().Statistics.RxOverErrors
        fields["rx-crc_Errors"] = link.Attrs().Statistics.RxCrcErrors
        fields["rx-frame_errors"] = link.Attrs().Statistics.RxFrameErrors
        fields["rx-fifo_errors"] = link.Attrs().Statistics.RxFifoErrors
        fields["rx-missed_errors"] = link.Attrs().Statistics.RxMissedErrors
        fields["rx-aborted_errors"] = link.Attrs().Statistics.TxAbortedErrors
        fields["tx-carrier_errors"] = link.Attrs().Statistics.TxCarrierErrors
        fields["tx-fifo_errors"] = link.Attrs().Statistics.TxFifoErrors
        fields["tx-heartbeat_errors"] = link.Attrs().Statistics.TxHeartbeatErrors
        fields["tx-window_errors"] = link.Attrs().Statistics.TxWindowErrors
        fields["rx_compressed"] = link.Attrs().Statistics.RxCompressed
        fields["tx_compressed"] = link.Attrs().Statistics.TxCompressed
        
        fields["num_vfs"] = len(link.Attrs().Vfs)
        
  //      fmt.Println(fields)
        acc.AddFields("netlink", fields, tags)
        
        if nl.ShowVirtualFunctions {
            if len(link.Attrs().Vfs) > 0 {
                for _,vfInfo := range link.Attrs().Vfs {
                    vfData,vfTags := generateVirtualFunctionInfo(vfInfo,link)
//                    fmt.Println(vfData)
                    
                    acc.AddFields("netlink_vf",vfData,vfTags)
                }
            }
        }    
        return nil
    }
func (nl *Netlink) Gather(acc telegraf.Accumulator) error {
	
    interfaceFilter, err :=  filter.NewIncludeExcludeFilter(nl.InterfaceInclude, nl.InterfaceExclude)    

    links, err :=  netlink.LinkList()

    if err != nil {
        fmt.Println(err)
    }
    
	var wg sync.WaitGroup
    for _, link :=  range links {
        if interfaceFilter.Match(link.Attrs().Name) {
  			wg.Add(1)
			go func(i netlink.Link) {  // do each in a thread, in case there are a bunch
				nl.gatherNetlinkNetworkStats(i,acc)
				wg.Done()
			}(link)
            
        }
        
    }    
	wg.Wait()
    
	return nil
}

func init() {
	inputs.Add(pluginName, func() telegraf.Input { return NewNetlink() })
}
