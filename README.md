# stores-server

Data is from OpenStreetMap.


Testing with `bobatool`:

```
bobatool stores find-stores \
	--address localhost:8080 \
	--insecure \
	--bounds.max.latitude 40 \
	--bounds.min.latitude 39.9 \
	--bounds.max.longitude -79.9 \
	--bounds.min.longitude -80 \
	--json
{
  "count": 1,
  "stores": [
    {
      "name": "USPS Isabella",
      "type": "office",
      "title": "USPS Isabella",
      "location": {
        "latitude": 39.9438,
        "longitude": -79.937996
      },
      "address": {
        "street": "1st Street",
        "city": "Isabella",
        "state": "PA",
        "zipCode": 15447,
        "regionCode": "us"
      }
    }
  ]
}
```
