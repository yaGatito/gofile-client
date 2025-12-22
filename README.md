# DROPPED: no reason to continue workarounding the resource restrictions. 
Last restriction added `X-Website-Token` as header allowing to download resource.
<br/>
Response:
`This content has been moved to cold storage due to inactivity. To access it, you'll need to import it into a premium account first.`

<br/>
Requests:
```shell
curl.exe -X POST https://api.gofile.io/contents/createFolder -H "Authorization: Bearer <TOKEN>" -H "Content-Type: application/json" -d '{\"parentFolderId\":\"<ID>\",\"folderName\":\"main\"}'
```

GET file
```shell
curl.exe -X GET https://api.gofile.io/contents/{ID} -H "Authorization: Bearer <TOKEN>" -H "Content-Type: application/json"
```

GET file
```shell
curl.exe -X GET https://store-eu-par-1.gofile.io/download/web/{ID}/file.txt -H "Authorization: Bearer <TOKEN>"
```


