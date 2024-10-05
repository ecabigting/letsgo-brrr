package main

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/ecabigting/letsgo-brrr/devicemonitor/internal/server"
)

func main() {
	fmt.Println(".. Loading devicemonitor")

	svr := server.NewServer(10, "6969")
	go func(s *server.Server) {
		for {
			timeStamp := time.Now().Format("2006-01-02 15:04:05")
			sectionSys, sectionDisc, sectionCPU := s.GetHardwareData()

			html := `
      <div hx-swap-oob="innerHTML:#timestamps">
          <table class='table table-striped table-hover table-sm'>
            <tbody>
              <tr>
                <td> Current System time:</td>
                <td>` + timeStamp + `</td>
              </tr>
            </tbody>
          </table>
      </div>
      <div hx-swap-oob="innerHTML:#system">` + sectionSys + `</div>
      <div hx-swap-oob="innerHTML:#disk">` + sectionDisc + `</div>
      <div hx-swap-oob="innerHTML:#cpu">` + sectionCPU + `</div>
      `

			s.Broadcast([]byte(html))

			time.Sleep(3 * time.Second)
		}
	}(svr)
	err := http.ListenAndServe(":"+svr.Port, &svr.Mux)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
