package main

import(
   "fmt"
   "io/ioutil"
   "encoding/json"
)


func saveDefaultData() {
   data, err := json.MarshalIndent(travelData, "", "  ")
   if err != nil {
      fmt.Println("無法序列化資料:", err)
      return
   }
   err = ioutil.WriteFile("travel_data.json", data, 0644)
   if err != nil {
      fmt.Println("無法寫入 travel_data.json:", err)
      return
   }
   fmt.Println("已建立預設的 travel_data.json 檔案")
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
   saveDefaultData()  // 儲存預設資料到檔案
}
