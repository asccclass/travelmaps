const base = '/travel';
// 初始化地圖
// const map = L.map('map').setView([25.0340, 121.5300], 12);
// 初始化地圖
// 設定新的中心點經度 (Longitude)
// 緯度 (lat): 0 (赤道附近，方便觀察全球)
// 經度 (lng): 175 (將地圖中心設在太平洋中央，靠近180度線)
const map = L.map('map').setView([0, -175], 4); 

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
var defaultz = 1;
let currentMapType = 'street';

if (defaultz == 0) {
   streetLayer.addTo(map);
} else {
   satelliteLayer.addTo(map);
   currentMapType = 'satellite';
}

let currentLocationId = null;
const photoDetail = document.getElementById('photo-detail');
const routeInfo = document.getElementById('route-info');
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
fetch(base + '/api/days')
    .then(response => response.json())
    .then(data => {
        totalDays = data.totalDays;
        const dayButtonsContainer = document.getElementById('day-buttons');
        
        // 全部顯示按鈕
        const allBtn = document.createElement('button');
        allBtn.className = 'day-btn active';
        allBtn.textContent = 'All';
        allBtn.onclick = () => filterByDay(0);
        dayButtonsContainer.appendChild(allBtn);
        
        // 各天按鈕
        for (let i = 1; i <= totalDays; i++) {
            const btn = document.createElement('button');
            btn.className = 'day-btn';
            btn.textContent = i;
            btn.onclick = () => filterByDay(i);
            dayButtonsContainer.appendChild(btn);
        }
    });

// 載入路線
fetch(base + '/api/routes')
    .then(response => response.json())
    .then(routes => {
        allRoutes = routes;
        updateLegend();
        renderRoutes();
    });

// 載入地點
fetch(base + '/api/locations')
    .then(response => response.json())
    .then(locations => {
        allLocations = locations;
        renderLocations();
    });

// 取得上一個需要的路徑
function getLastNeededRoute(points, i) {   
   var j = i - 1;
   if(points[i].routez == 0)  return {lat: points[j].lat, lng: points[j].lng};
   for(; j >= 0 && points[j].routez == 0; j--);
   if(j < 0)  return {lat: points[i-1].lat, lng: points[i-1].lng};
   else {
      return {lat: points[j].lat, lng: points[j].lng};
   }     
}

async function renderRoutes() {
    polylines.forEach(p => map.removeLayer(p));  // 清除現有路線
    polylines = [];

    const routesToShow = currentDay === 0 
        ? allRoutes 
        : allRoutes.filter(r => r.day === currentDay);

    if (routesToShow.length > 0) {
        routeInfo.classList.add('show');
    }

    // 為每條路線獲取實際路徑
    for(const route of routesToShow) {  // routes
        // 每天的路線
        for(var i = 1; i < route.points.length; i++) {  // 跳過第一個點，因為沒有前一個點可以連接
            const end = {lat: route.points[i].lat, lng: route.points[i].lng};
            const start = getLastNeededRoute(route.points, i);
            var routePoints = [start, end];
            if(route.points[i].routez == 0) {
                  const polyline = L.polyline(
                     routePoints.map(p => [p.lat, p.lng]),
                     {
                        color: route.color,
                        weight: 4,
                        opacity: 0.7,
                        smoothFactor: 1,
                        dashArray: '10, 10'
                     }
                  ).addTo(map);
                  polyline.bindPopup('<b>' + route.name + '</b>');
                  polylines.push(polyline);
                  continue;
            }
            try {
                  const routePath = await getRoutePath(routePoints);
		    console.log(routePath);
                  const polyline = L.polyline(
                     routePath,
                     {
                        color: route.color,
                        weight: 4,
                        opacity: 0.7,
                        smoothFactor: 1
                     }
                  ).addTo(map);
                  polyline.bindPopup('<b>' + route.name + '</b>');
                  polylines.push(polyline);
            } catch (error) {
                  console.error('路徑規劃失敗，使用直線連接:', error);
                  // 如果路徑規劃失敗，使用直線作為備案
                  const polyline = L.polyline(
                     routePoints.map(p => [p.lat, p.lng]),
                     {
                        color: route.color,
                        weight: 4,
                        opacity: 0.7,
                        smoothFactor: 1,
                        dashArray: '10, 10'
                     }
                  ).addTo(map);
                  polyline.bindPopup('<b>' + route.name + '</b>');
                  polylines.push(polyline);
            }
         }
    }

    routeInfo.classList.remove('show');

    // 自動調整地圖視野
    if (polylines.length > 0) {
        const group = new L.featureGroup(polylines);
        map.fitBounds(group.getBounds().pad(0.1));
    }
}

// 使用 OSRM 路徑規劃服務獲取實際路徑
async function getRoutePath(points) {
    if (points.length < 2) {
        return points.map(p => [p.lat, p.lng]);
    }

    // 構建 OSRM API 請求
    const coordinates = points.map(p => `${p.lng},${p.lat}`).join(';');
    const url = `https://router.project-osrm.org/route/v1/driving/${coordinates}?overview=full&geometries=geojson`;

    const response = await fetch(url);
    const data = await response.json();

    if (data.code === 'Ok' && data.routes && data.routes.length > 0) {
        // 將 GeoJSON 座標轉換為 Leaflet 格式 [lat, lng]
        return data.routes[0].geometry.coordinates.map(coord => [coord[1], coord[0]]);
    } else {
        throw new Error('路徑規劃失敗:' + coordinates);
    }
}

function renderLocations() {    
    markers.forEach(m => map.removeLayer(m));  // 清除現有標記
    markers = [];

    const locationsToShow = currentDay === 0 
        ? allLocations 
        : allLocations.filter(l => l.day === currentDay);

    locationsToShow.forEach(location => {
        const markerSize = getMarkerSize();
        
        const icon = L.divIcon({
            className: 'custom-marker',
            html: `<div class="marker-pin" style="background: #ff4444; width: ${markerSize}px; height: ${markerSize}px; border-radius: 50%; border: ${Math.max(2, markerSize / 10)}px solid white; box-shadow: 0 2px 8px rgba(0,0,0,0.3); transition: all 0.3s ease;"></div>`,
            iconSize: [markerSize, markerSize],
            iconAnchor: [markerSize / 2, markerSize / 2]
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

// 根據地圖縮放等級計算標記大小
function getMarkerSize() {
    const zoom = map.getZoom();
    // 縮放等級 10-19，標記大小 20-50px
    const minSize = 20;
    const maxSize = 50;
    const minZoom = 10;
    const maxZoom = 19;
    
    if (zoom <= minZoom) return minSize;
    if (zoom >= maxZoom) return maxSize;
    
    const size = minSize + ((zoom - minZoom) / (maxZoom - minZoom)) * (maxSize - minSize);
    return Math.round(size);
}

// 更新所有標記的大小
function updateMarkerSizes() {
    const markerSize = getMarkerSize();
    const borderWidth = Math.max(2, markerSize / 10);
    
    markers.forEach(marker => {
        const markerElement = marker.getElement();
        if (markerElement) {
            const pin = markerElement.querySelector('.marker-pin');
            if (pin) {
                pin.style.width = markerSize + 'px';
                pin.style.height = markerSize + 'px';
                pin.style.borderWidth = borderWidth + 'px';
            }
            
            // 更新圖標錨點
            marker.options.icon.options.iconSize = [markerSize, markerSize];
            marker.options.icon.options.iconAnchor = [markerSize / 2, markerSize / 2];
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
    
    fetch(base + '/api/location-photos?id=' + locationId)
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
