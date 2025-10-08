package main

import (
   "net/http"
   "github.com/asccclass/sherryserver"
)

func NewRouter(srv *SherryServer.Server, documentRoot string)(*http.ServeMux) {
   router := http.NewServeMux()

   // 靜態檔案服務
   staticfileserver := SherryServer.StaticFileServer{documentRoot, "index.html"}
   staticfileserver.AddRouter(router)

   // API 路由
   // router.HandleFunc("/", indexHandler)
   router.HandleFunc("/api/locations", locationsHandler)
   router.HandleFunc("/api/routes", routesHandler)
   router.HandleFunc("/api/location-photos", locationPhotosHandler)
   router.HandleFunc("/api/days", daysHandler)
   return router
}
