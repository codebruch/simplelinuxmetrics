package main
//go get -v github.com/prometheus/client_golang/prometheus
import (
	"log"
	"net/http"
	"strings"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"strconv"
)

//import "os"
import "os/exec"

import "fmt"


type CPUInfo struct {
	ProcessorName string
	ProcessorNumber int 
}

var CPUs []CPUInfo




var (


/*
	cpuTempOpts = prometheus.GaugeOpts{
		Name: "cpu_temperature_celsius",
		Help: "Current temperature of the CPU.",
	}

	cpuTemp = prometheus.NewGaugeVec(
		cpuTempOpts
		,
	[]string{"CPU Type","Cores"},
        )

*/

	cpuTempOpts = prometheus.GaugeOpts{
                Name: "cpu_temperature_celsius",
                Help: "Current temperature of the CPU.",
        }



        cpuTemp = prometheus.NewGaugeVec(
                cpuTempOpts,
        []string{"CPU_Type","Cores"},
        )



	hdFailures = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "hd_errors_total",
			Help: "Number of hard-disk errors.",
		},
		[]string{"device","hdname"},
	)
	
	hdTemp = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
	        	Name: "hd_temperature_celsius",
                	Help: "Current temperature of the HardDiscs in the system",
        },
	[]string{"device","hdname"},
	)

)


func detectCPUs()  *CPUInfo {


	binary, lookErr := exec.LookPath("cat")
    	if lookErr != nil {
        	panic(lookErr)
    	}

	args := []string{"cat", "/proc/cpuinfo"}
       // env := os.Environ()

       cmd := exec.Command(binary, args[1])
       output, err := cmd.Output()       


       if err != nil {
          panic(err)
       }
       outputstr := fmt.Sprintf("%s", output)
       var CPUCount = strings.Count(outputstr,"processor")
       outputstr = strings.Trim(outputstr," ")

       var outputstrArray = strings.Split(outputstr,":")
       var CPUModel  = ""

       for i := 0; i < len(outputstrArray); i++  {

       //   fmt.Printf("CPU %s\n", strings.Trim(outputstrArray[i]," " ))


	 if(strings.Contains(outputstrArray[i], "model name")){

		 CPUModel = strings.Trim(outputstrArray[i+1]," " )
		 fmt.Printf("CPU model %s\n", CPUModel)
		}
	}
       fmt.Printf("No CPUs: %d\n", CPUCount)


	cpus := &CPUInfo{

	ProcessorName: CPUModel,
	ProcessorNumber: CPUCount,

       }




  return cpus

}



func readCPUtemperature(corenum string) float64 { 

   var CPUTemp = ""
   var temperature  = -1.0 
   binary, lookErr := exec.LookPath("sensors")
   if lookErr != nil {
           panic(lookErr)
   }

   args := []string{"sensors", "-u"}

   cmd := exec.Command(binary, args[1])
   output, err := cmd.Output()


   if err != nil {
       panic(err)
   }
   outputstr := fmt.Sprintf("%s", output)

   var outputstrArray = strings.Split(outputstr,":")

   for i := 0; i < len(outputstrArray); i++  {

       //   fmt.Printf("CPU %s\n", strings.Trim(outputstrArray[i]," " ))

         if(strings.Contains(outputstrArray[i], "temp"+corenum+"_input")){

                 CPUTemp = strings.Trim(outputstrArray[i+1]," " )
                 fmt.Printf("CPU model %s\n", CPUTemp)
                }
   }

 return temperature

}


func init() {

	// Metrics have to be registered to be exposed:
	prometheus.MustRegister(cpuTemp)
	prometheus.MustRegister(hdFailures)
	prometheus.MustRegister(hdTemp)
}

func main() {

	var cpuinfo = detectCPUs()

       

	cpuTemp.With(prometheus.Labels{"CPU_Type":cpuinfo.ProcessorName, "Cores":strconv.Itoa(cpuinfo.ProcessorNumber)}).Set(65.3)


	//detect hdd and write device
	hdFailures.With(prometheus.Labels{"device":"/dev/sda", "hdname":"Hitachi HDS721050CLA362"}).Inc()

	hdTemp.With(prometheus.Labels{"device":"/dev/sda", "hdname":"Hitachi HDS721050CLA362"}).Set(38)

	// The Handler function provides a default handler to expose metrics
	// via an HTTP server. "/metrics" is the usual endpoint for that.
	http.Handle("/metrics", promhttp.Handler())
	log.Fatal(http.ListenAndServe(":1212", nil))
}
