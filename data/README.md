# Data

Prepared font family, background image, and data.json to generate fake ticket.

## File: background.png

Default background image

## File: NotoSansTC-Regular.ttf

Default font family

## File: data.json

Generating ticket information.

```json
{
  "width":6, // required
  "height":9, // required
  "background":"background.png", // required
  "font_family":"NotoSansTC-Regular.ttf", // required
  "cinema":{
    "name":"臺灣電影院" // required
  },
  "movie":{
    "name":"電影名稱", // required
    "eng_name":"Movie Name", // required
    "time":"2024/01/01 00:00" // required
  },
  "ticket":{
    "room":"1廳", // required
    "seat":"1徘1號", // required
    "type":"全票", // required
    "price":240, // required
    "sales_time":"2024/01/01 00:00" // optional
  }
}
```
