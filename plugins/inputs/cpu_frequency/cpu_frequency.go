package cpu_frequency

import (
	"github.com/influxdata/telegraf"
	"github.com/influxdata/telegraf/plugins/inputs"
    "fmt"
	"os"
	"path"
    "io/ioutil"
    "strconv"
    "strings"
)

type CPU_Frequency struct {
    sysCpuDir string  
    cpuCount int
}


var CPU_FrequencyConfig = `
  ## nothing to setup for this plugin

`

func NewCPU_Frequency() * CPU_Frequency {
    
    obj := new (CPU_Frequency)
    obj.sysCpuDir = path.Join(getHostSys(), "/devices/system/cpu/")

    return obj
}

func (s *CPU_Frequency) SampleConfig() string {
	return CPU_FrequencyConfig
}

func (s *CPU_Frequency) Description() string {
	return "Gathers the CPU Frequency for each core"
}

func (s *CPU_Frequency) getFrequencies() map[string]interface{}  {
    
    checkCount := s.cpuCount
    
    if 0 == checkCount {
        checkCount = 250 // haven't run this yet, so give it a large #
    }
    retMap := make(map[string]interface{})    
    
    for coreNum :=0; coreNum < checkCount; coreNum++ { // loop through call cores
        fName := fmt.Sprintf("%s/cpu%d/cpufreq/scaling_cur_freq",s.sysCpuDir,coreNum)
        freqStr, err := ioutil.ReadFile(fName)
        if err != nil {
            if 0 == s.cpuCount && coreNum > 0 {
                s.cpuCount = coreNum // use it next time
            } else {
                 fmt.Printf("error reading file  %s", fName)
                 return nil
            }
            break
        }
        freqStrClean := strings.TrimSuffix(string(freqStr),"\n") // strip carriage return
        key := fmt.Sprintf("cpu%d",coreNum) // column 
        freq,err := strconv.ParseInt(freqStrClean,10,64) // convert to uint
        if err != nil {
            fmt.Println(err)
         }
        retMap[key] = freq
    }

    return retMap
}

/*  use environment variable for containerized version */
func getHostSys() string {
	procPath := "/sys"
	if os.Getenv("HOST_SYS") != "" {
		procPath = os.Getenv("HOST_SYS")
	}
	return procPath
}


func (s *CPU_Frequency) Gather(acc telegraf.Accumulator) error {
	
    freqMap := s.getFrequencies()

    //fmt.Println(freqMap)
    acc.AddGauge("CPU_Frequency", freqMap, nil)
	return nil
}

func init() {
	inputs.Add("cpu_frequency", func() telegraf.Input { return NewCPU_Frequency() })
}
