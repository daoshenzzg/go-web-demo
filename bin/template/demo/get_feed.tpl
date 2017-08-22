{
    "code":"{{.Code}}",
    "msg":"",
    "data":{
        {{if .Data.feed}}
            "id": "{{.Data.feed.Id}}",
            "uid": "{{.Data.feed.Uid}}",
            "title": "{{.Data.feed.Title}}"
        {{end}}
    }
}