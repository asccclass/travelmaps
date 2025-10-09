package main

import (
   "os"
   "fmt"
   "strconv"
   "net/http"
   "io/ioutil"
   "encoding/json"
   // "html/template"
   "github.com/joho/godotenv"
   "github.com/asccclass/sherryserver"
)

// Location 代表地圖上的一個地點
type Location struct {
   ID          int       `json:"id"`
   Name        string    `json:"name"`
   Lat         float64   `json:"lat"`
   Lng         float64   `json:"lng"`
   Description string    `json:"description"`
   Photos      []Photo   `json:"photos"`
   Day         int       `json:"day"`
}

// Photo 代表一張照片
type Photo struct {
   Thumbnail string `json:"thumbnail"`
   Full      string `json:"full"`
   Caption   string `json:"caption"`
}

// Route 代表移動軌跡
type Route struct {
   Points []Point `json:"points"`
   Color  string  `json:"color"`
   Name   string  `json:"name"`
   Day    int     `json:"day"`
}

// Point 代表路線上的一個點
type Point struct {
   Lat float64 `json:"lat"`
   Lng float64 `json:"lng"`
   Routez int  `json:"routez"`
}

// TravelData 旅行資料結構
type TravelData struct {
   Locations []Location `json:"locations"`
   Routes    []Route    `json:"routes"`
   TotalDays int        `json:"totalDays"`
}

var travelData TravelData

func init() {
   // 從 JSON 檔案讀取資料
   data, err := ioutil.ReadFile("data/travel_data.json")
   if err != nil {
      fmt.Println("找不到 travel_data.json，使用預設資料")
      createDefaultData()
      return
   }

   err = json.Unmarshal(data, &travelData)
   if err != nil {
      fmt.Println("JSON 解析失敗，使用預設資料:", err)
      createDefaultData()
      return
   }
   fmt.Println("成功從 travel_data.json 載入資料")
}

func main() {
   if err := godotenv.Load("envfile"); err != nil {
      fmt.Println(err.Error())
      return
   }
   port := os.Getenv("PORT")
   if port == "" {
      port = "80"
   }
   documentRoot := os.Getenv("DocumentRoot")
   if documentRoot == "" {
      documentRoot = "www/html"
   }
   templateRoot := os.Getenv("TemplateRoot")
   if templateRoot == "" {
      templateRoot = "www/template"
   }

   server, err := SherryServer.NewServer(":" + port, documentRoot, templateRoot)
   if err != nil {
      panic(err)
   }
   router := NewRouter(server, documentRoot)
   if router == nil {
      fmt.Println("router return nil")
      return
   }
   server.Server.Handler = router  // server.CheckCROS(router)  // 需要自行implement, overwrite 預設的
   server.Start()
}

func locationsHandler(w http.ResponseWriter, r *http.Request) {
   w.Header().Set("Content-Type", "application/json")
   json.NewEncoder(w).Encode(travelData.Locations)
}

func routesHandler(w http.ResponseWriter, r *http.Request) {
   w.Header().Set("Content-Type", "application/json")
   json.NewEncoder(w).Encode(travelData.Routes)
}

func daysHandler(w http.ResponseWriter, r *http.Request) {
   w.Header().Set("Content-Type", "application/json")
   json.NewEncoder(w).Encode(map[string]int{
      "totalDays": travelData.TotalDays,
   })
}

func locationPhotosHandler(w http.ResponseWriter, r *http.Request) {
   locationID := r.URL.Query().Get("id")
   id, _ := strconv.Atoi(locationID)

   var location *Location
   for i := range travelData.Locations {
      if travelData.Locations[i].ID == id {
         location = &travelData.Locations[i]
         break
      }
   }

   if location == nil {
      http.Error(w, "Location not found", http.StatusNotFound)
      return
   }

   html := `<div class="detail-header">
      <div class="detail-title">` + location.Name + `</div>
      <button class="close-btn" onclick="closePhotoDetail()">×</button>
   </div>
   <div class="location-description">` + location.Description + ` (第 ` + strconv.Itoa(location.Day) + ` 天)</div>`

   for _, photo := range location.Photos {
      if photo.Full == "" {
         photo.Full = photo.Thumbnail
      }
      html += `<div class="photo-item">
         <img src="` + photo.Full + `" alt="` + photo.Caption + `">
         <div class="photo-caption">` + photo.Caption + `</div>
      </div>`
   }

   w.Header().Set("Content-Type", "text/html; charset=utf-8")
   w.Write([]byte(html))
}
