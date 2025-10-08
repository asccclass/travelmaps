package main

import (
	"encoding/json"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
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
	data, err := ioutil.ReadFile("travel_data.json")
	if err != nil {
		log.Println("找不到 travel_data.json，使用預設資料")
		createDefaultData()
		return
	}

	err = json.Unmarshal(data, &travelData)
	if err != nil {
		log.Println("JSON 解析失敗，使用預設資料:", err)
		createDefaultData()
		return
	}
	log.Println("成功從 travel_data.json 載入資料")
}

func createDefaultData() {
	travelData = TravelData{
		TotalDays: 3,
		Locations: []Location{
			{
				ID:          1,
				Name:        "台北101",
				Lat:         25.0340,
				Lng:         121.5645,
				Description: "台北最著名的地標建築",
				Day:         1,
				Photos: []Photo{
					{
						Thumbnail: "https://images.unsplash.com/photo-1566022671514-a2e75cec97eb?w=150&h=150&fit=crop",
						Full:      "https://images.unsplash.com/photo-1566022671514-a2e75cec97eb?w=800&h=600&fit=crop",
						Caption:   "台北101外觀",
					},
					{
						Thumbnail: "https://images.unsplash.com/photo-1570994728901-3ad0cf27bafb?w=150&h=150&fit=crop",
						Full:      "https://images.unsplash.com/photo-1570994728901-3ad0cf27bafb?w=800&h=600&fit=crop",
						Caption:   "台北101夜景",
					},
				},
			},
			{
				ID:          2,
				Name:        "中正紀念堂",
				Lat:         25.0408,
				Lng:         121.5188,
				Description: "台灣重要的歷史文化地標",
				Day:         1,
				Photos: []Photo{
					{
						Thumbnail: "https://images.unsplash.com/photo-1587139223577-f48c6cd8b6c4?w=150&h=150&fit=crop",
						Full:      "https://images.unsplash.com/photo-1587139223577-f48c6cd8b6c4?w=800&h=600&fit=crop",
						Caption:   "中正紀念堂正面",
					},
				},
			},
			{
				ID:          3,
				Name:        "西門町",
				Lat:         25.0421,
				Lng:         121.5069,
				Description: "台北熱鬧的購物商圈",
				Day:         2,
				Photos: []Photo{
					{
						Thumbnail: "https://images.unsplash.com/photo-1536098561742-ca998e48cbcc?w=150&h=150&fit=crop",
						Full:      "https://images.unsplash.com/photo-1536098561742-ca998e48cbcc?w=800&h=600&fit=crop",
						Caption:   "西門町街景",
					},
				},
			},
			{
				ID:          4,
				Name:        "士林夜市",
				Lat:         25.0878,
				Lng:         121.5241,
				Description: "台北最大的夜市之一",
				Day:         2,
				Photos: []Photo{
					{
						Thumbnail: "https://images.unsplash.com/photo-1555939594-58d7cb561ad1?w=150&h=150&fit=crop",
						Full:      "https://images.unsplash.com/photo-1555939594-58d7cb561ad1?w=800&h=600&fit=crop",
						Caption:   "士林夜市美食",
					},
				},
			},
			{
				ID:          5,
				Name:        "淡水老街",
				Lat:         25.1677,
				Lng:         121.4425,
				Description: "欣賞淡水河夕陽的好地方",
				Day:         3,
				Photos: []Photo{
					{
						Thumbnail: "https://images.unsplash.com/photo-1590736969955-71cc94901144?w=150&h=150&fit=crop",
						Full:      "https://images.unsplash.com/photo-1590736969955-71cc94901144?w=800&h=600&fit=crop",
						Caption:   "淡水夕陽",
					},
				},
			},
		},
		Routes: []Route{
			{
				Name:  "第一天 - 市區景點",
				Color: "#3388ff",
				Day:   1,
				Points: []Point{
					{Lat: 25.0340, Lng: 121.5645},
					{Lat: 25.0408, Lng: 121.5188},
				},
			},
			{
				Name:  "第二天 - 購物美食",
				Color: "#ff6b6b",
				Day:   2,
				Points: []Point{
					{Lat: 25.0421, Lng: 121.5069},
					{Lat: 25.0878, Lng: 121.5241},
				},
			},
			{
				Name:  "第三天 - 淡水之旅",
				Color: "#4ecdc4",
				Day:   3,
				Points: []Point{
					{Lat: 25.0878, Lng: 121.5241},
					{Lat: 25.1677, Lng: 121.4425},
				},
			},
		},
	}

	// 儲存預設資料到檔案
	saveDefaultData()
}

func saveDefaultData() {
	data, err := json.MarshalIndent(travelData, "", "  ")
	if err != nil {
		log.Println("無法序列化資料:", err)
		return
	}
	err = ioutil.WriteFile("travel_data.json", data, 0644)
	if err != nil {
		log.Println("無法寫入 travel_data.json:", err)
		return
	}
	log.Println("已建立預設的 travel_data.json 檔案")
}

func main() {
	// 靜態檔案服務
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	
	// API 路由
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/api/locations", locationsHandler)
	http.HandleFunc("/api/routes", routesHandler)
	http.HandleFunc("/api/location-photos", locationPhotosHandler)
	http.HandleFunc("/api/days", daysHandler)

	log.Println("伺服器啟動於 http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	tmpl := `<!DOCTYPE html>
<html lang="zh-TW">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>旅行地圖</title>
    <link rel="stylesheet" href="https://unpkg.com/leaflet@1.9.4/dist/leaflet.css" />
    <link rel="stylesheet" href="/static/css/style.css" />
</head>
<body>
    <div id="map"></div>
    <div id="photo-detail"></div>
    <div id="route-info" class="route-info">正在規劃路徑...</div>
    
    <div class="controls">
        <div class="control-group">
            <div class="control-label">地圖類型</div>
            <div class="map-type-buttons">
                <button class="map-type-btn active" onclick="switchMapType('street')">街道地圖</button>
                <button class="map-type-btn" onclick="switchMapType('satellite')">衛星地圖</button>
            </div>
        </div>
        <div class="control-group">
            <div class="control-label">選擇日期</div>
            <div class="day-buttons" id="day-buttons"></div>
        </div>
    </div>
    
    <div class="legend">
        <div class="legend-title">圖例</div>
        <div id="legend-routes"></div>
        <div class="legend-item">
            <div style="width: 20px; height: 20px; background: #ff4444; border-radius: 50%; margin-right: 8px;"></div>
            <span>景點位置</span>
        </div>
    </div>

    <script src="https://unpkg.com/leaflet@1.9.4/dist/leaflet.js"></script>
    <script src="https://unpkg.com/htmx.org@1.9.10"></script>
    <script src="/static/js/app.js"></script>
</body>
</html>`

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	t, _ := template.New("index").Parse(tmpl)
	t.Execute(w, nil)
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