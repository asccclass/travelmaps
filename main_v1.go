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
    <script src="https://unpkg.com/leaflet@1.9.4/dist/leaflet.js"></script>
    <script src="https://unpkg.com/htmx.org@1.9.10"></script>
    <style>
        * {
            margin: 0;
            padding: 0;
            box-sizing: border-box;
        }
        body {
            font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif;
            height: 100vh;
            overflow: hidden;
        }
        #map {
            width: 100%;
            height: 100vh;
        }
        .photo-popup {
            text-align: center;
        }
        .photo-thumbnail {
            width: 80px;
            height: 80px;
            object-fit: cover;
            border-radius: 8px;
            cursor: pointer;
            box-shadow: 0 2px 8px rgba(0,0,0,0.3);
            transition: transform 0.2s;
        }
        .photo-thumbnail:hover {
            transform: scale(1.1);
        }
        .location-name {
            font-weight: bold;
            margin-bottom: 8px;
            font-size: 16px;
        }
        .location-day {
            font-size: 12px;
            color: #666;
            margin-bottom: 8px;
        }
        #photo-detail {
            position: fixed;
            right: 20px;
            top: 20px;
            width: 400px;
            max-height: 80vh;
            background: white;
            border-radius: 12px;
            box-shadow: 0 4px 20px rgba(0,0,0,0.3);
            padding: 20px;
            overflow-y: auto;
            display: none;
            z-index: 1000;
        }
        #photo-detail.show {
            display: block;
            animation: slideIn 0.3s ease-out;
        }
        @keyframes slideIn {
            from {
                opacity: 0;
                transform: translateX(50px);
            }
            to {
                opacity: 1;
                transform: translateX(0);
            }
        }
        .detail-header {
            display: flex;
            justify-content: space-between;
            align-items: center;
            margin-bottom: 16px;
            border-bottom: 2px solid #eee;
            padding-bottom: 12px;
        }
        .detail-title {
            font-size: 20px;
            font-weight: bold;
            color: #333;
        }
        .close-btn {
            background: #ff4444;
            color: white;
            border: none;
            border-radius: 50%;
            width: 32px;
            height: 32px;
            cursor: pointer;
            font-size: 20px;
            line-height: 1;
            transition: background 0.2s;
        }
        .close-btn:hover {
            background: #cc0000;
        }
        .photo-item {
            margin-bottom: 20px;
            border-radius: 8px;
            overflow: hidden;
            box-shadow: 0 2px 8px rgba(0,0,0,0.1);
        }
        .photo-item img {
            width: 100%;
            height: auto;
            display: block;
        }
        .photo-caption {
            padding: 12px;
            background: #f8f8f8;
            font-size: 14px;
            color: #666;
        }
        .location-description {
            padding: 12px;
            background: #fff3cd;
            border-radius: 8px;
            margin-bottom: 16px;
            font-size: 14px;
            color: #856404;
        }
        .controls {
            position: fixed;
            top: 20px;
            left: 60px;
            background: white;
            padding: 15px 20px;
            border-radius: 8px;
            box-shadow: 0 2px 10px rgba(0,0,0,0.2);
            z-index: 1000;
        }
        .control-group {
            margin-bottom: 12px;
        }
        .control-group:last-child {
            margin-bottom: 0;
        }
        .control-label {
            font-weight: bold;
            margin-bottom: 8px;
            display: block;
            color: #333;
            font-size: 14px;
        }
        .day-buttons {
            display: flex;
            flex-wrap: wrap;
            gap: 8px;
        }
        .day-btn {
            padding: 8px 16px;
            border: 2px solid #ddd;
            background: white;
            border-radius: 6px;
            cursor: pointer;
            font-size: 14px;
            transition: all 0.2s;
            font-weight: 500;
        }
        .day-btn:hover {
            background: #f0f0f0;
        }
        .day-btn.active {
            background: #3388ff;
            color: white;
            border-color: #3388ff;
        }
        .map-type-buttons {
            display: flex;
            gap: 8px;
        }
        .map-type-btn {
            padding: 8px 16px;
            border: 2px solid #ddd;
            background: white;
            border-radius: 6px;
            cursor: pointer;
            font-size: 14px;
            transition: all 0.2s;
            font-weight: 500;
        }
        .map-type-btn:hover {
            background: #f0f0f0;
        }
        .map-type-btn.active {
            background: #4CAF50;
            color: white;
            border-color: #4CAF50;
        }
        .legend {
            position: fixed;
            bottom: 20px;
            left: 20px;
            background: white;
            padding: 15px;
            border-radius: 8px;
            box-shadow: 0 2px 10px rgba(0,0,0,0.2);
            z-index: 1000;
            max-width: 250px;
        }
        .legend-title {
            font-weight: bold;
            margin-bottom: 8px;
            color: #333;
        }
        .legend-item {
            display: flex;
            align-items: center;
            margin: 5px 0;
            font-size: 14px;
        }
        .legend-color {
            width: 30px;
            height: 3px;
            margin-right: 8px;
            border-radius: 2px;
        }
    </style>
</head>
<body>
    <div id="map"></div>
    <div id="photo-detail"></div>
    
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

    <script>
        // 初始化地圖
        const map = L.map('map').setView([25.0340, 121.5300], 12);
        
        // 定義地圖圖層
        const streetLayer = L.tileLayer('https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png', {
            attribution: '© OpenStreetMap contributors',
            maxZoom: 19
        });

        const satelliteLayer = L.tileLayer('https://server.arcgisonline.com/ArcGIS/rest/services/World_Imagery/MapServer/tile/{z}/{y}/{x}', {
            attribution: '© Esri',
            maxZoom: 19
        });

        const satelliteLabelLayer = L.tileLayer('https://server.arcgisonline.com/ArcGIS/rest/services/Reference/World_Boundaries_and_Places/MapServer/tile/{z}/{y}/{x}', {
            attribution: '© Esri',
            maxZoom: 19
        });

        // 預設使用街道地圖
        streetLayer.addTo(map);
        let currentMapType = 'street';

        let currentLocationId = null;
        const photoDetail = document.getElementById('photo-detail');
        let currentDay = 0; // 0 表示顯示全部
        let allLocations = [];
        let allRoutes = [];
        let markers = [];
        let polylines = [];
        let totalDays = 0;

        // 切換地圖類型
        function switchMapType(type) {
            document.querySelectorAll('.map-type-btn').forEach(btn => btn.classList.remove('active'));
            event.target.classList.add('active');
            
            if (type === 'satellite' && currentMapType !== 'satellite') {
                map.removeLayer(streetLayer);
                map.addLayer(satelliteLayer);
                map.addLayer(satelliteLabelLayer);
                currentMapType = 'satellite';
            } else if (type === 'street' && currentMapType !== 'street') {
                map.removeLayer(satelliteLayer);
                map.removeLayer(satelliteLabelLayer);
                map.addLayer(streetLayer);
                currentMapType = 'street';
            }
        }

        // 載入天數資訊
        fetch('/api/days')
            .then(response => response.json())
            .then(data => {
                totalDays = data.totalDays;
                const dayButtonsContainer = document.getElementById('day-buttons');
                
                // 全部顯示按鈕
                const allBtn = document.createElement('button');
                allBtn.className = 'day-btn active';
                allBtn.textContent = '全部顯示';
                allBtn.onclick = () => filterByDay(0);
                dayButtonsContainer.appendChild(allBtn);
                
                // 各天按鈕
                for (let i = 1; i <= totalDays; i++) {
                    const btn = document.createElement('button');
                    btn.className = 'day-btn';
                    btn.textContent = '第 ' + i + ' 天';
                    btn.onclick = () => filterByDay(i);
                    dayButtonsContainer.appendChild(btn);
                }
            });

        // 載入路線
        fetch('/api/routes')
            .then(response => response.json())
            .then(routes => {
                allRoutes = routes;
                updateLegend();
                renderRoutes();
            });

        // 載入地點
        fetch('/api/locations')
            .then(response => response.json())
            .then(locations => {
                allLocations = locations;
                renderLocations();
            });

        function renderRoutes() {
            // 清除現有路線
            polylines.forEach(p => map.removeLayer(p));
            polylines = [];

            const routesToShow = currentDay === 0 
                ? allRoutes 
                : allRoutes.filter(r => r.day === currentDay);

            routesToShow.forEach(route => {
                const polyline = L.polyline(
                    route.points.map(p => [p.lat, p.lng]),
                    {
                        color: route.color,
                        weight: 4,
                        opacity: 0.7,
                        smoothFactor: 1
                    }
                ).addTo(map);
                
                polyline.bindPopup('<b>' + route.name + '</b>');
                polylines.push(polyline);
            });

            // 自動調整地圖視野
            if (polylines.length > 0) {
                const group = new L.featureGroup(polylines);
                map.fitBounds(group.getBounds().pad(0.1));
            }
        }

        function renderLocations() {
            // 清除現有標記
            markers.forEach(m => map.removeLayer(m));
            markers = [];

            const locationsToShow = currentDay === 0 
                ? allLocations 
                : allLocations.filter(l => l.day === currentDay);

            locationsToShow.forEach(location => {
                const icon = L.divIcon({
                    className: 'custom-marker',
                    html: '<div style="background: #ff4444; width: 30px; height: 30px; border-radius: 50%; border: 3px solid white; box-shadow: 0 2px 8px rgba(0,0,0,0.3);"></div>',
                    iconSize: [30, 30],
                    iconAnchor: [15, 15]
                });

                const marker = L.marker([location.lat, location.lng], {icon: icon}).addTo(map);
                
                let popupContent = '<div class="photo-popup">';
                popupContent += '<div class="location-name">' + location.name + '</div>';
                popupContent += '<div class="location-day">第 ' + location.day + ' 天</div>';
                if (location.photos && location.photos.length > 0) {
                    popupContent += '<img src="' + location.photos[0].thumbnail + '" class="photo-thumbnail" data-location-id="' + location.id + '">';
                }
                popupContent += '</div>';
                
                marker.bindPopup(popupContent);
                markers.push(marker);
            });

            // 使用事件委派處理縮圖點擊
            map.on('popupopen', function(e) {
                const thumbnail = e.popup._contentNode.querySelector('.photo-thumbnail');
                if (thumbnail) {
                    thumbnail.addEventListener('mouseenter', function() {
                        const locationId = this.getAttribute('data-location-id');
                        loadLocationPhotos(locationId);
                    });
                }
            });
        }

        function filterByDay(day) {
            currentDay = day;
            
            // 更新按鈕狀態
            document.querySelectorAll('.day-btn').forEach(btn => btn.classList.remove('active'));
            event.target.classList.add('active');
            
            renderRoutes();
            renderLocations();
        }

        function updateLegend() {
            const legendRoutes = document.getElementById('legend-routes');
            legendRoutes.innerHTML = '';
            
            allRoutes.forEach(route => {
                const item = document.createElement('div');
                item.className = 'legend-item';
                item.innerHTML = '<div class="legend-color" style="background: ' + route.color + ';"></div><span>' + route.name + '</span>';
                legendRoutes.appendChild(item);
            });
        }

        function loadLocationPhotos(locationId) {
            if (currentLocationId === locationId && photoDetail.classList.contains('show')) {
                return;
            }
            
            currentLocationId = locationId;
            
            fetch('/api/location-photos?id=' + locationId)
                .then(response => response.text())
                .then(html => {
                    photoDetail.innerHTML = html;
                    photoDetail.classList.add('show');
                });
        }

        function closePhotoDetail() {
            photoDetail.classList.remove('show');
            currentLocationId = null;
        }

        // 點擊地圖關閉照片詳情
        map.on('click', function() {
            closePhotoDetail();
        });
    </script>
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
		html += `<div class="photo-item">
			<img src="` + photo.Full + `" alt="` + photo.Caption + `">
			<div class="photo-caption">` + photo.Caption + `</div>
		</div>`
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write([]byte(html))
}