# Energy Consumption

## Description
My energy supplier is Octopus Energy. They provide API Access to Smart Data readings via the [Developer Dashboard](https://octopus.energy/dashboard/developer/). 

I wanted to grab the API data and display it in a Grafana dashboard so we could see visually what our Energy and Gas consumption is. 

It also seemed like a good opportunity to learn some Go :) 

## Usage

### Octopus Details
Below is the config layout, you'll need your API key from the [Developer Dashboard](https://octopus.energy/dashboard/developer/). You'll also need your Gas/Electricity meter numbers available on the same page. 

Finally, you'll need the p/kWh for your tarriff this is available on your [Account Dashboard](https://octopus.energy/dashboard/).

*NOTE: Keep this all secret!* 

### Influx
You'll need the account's access token, bucket, org and the URL where influx is located. 

Config.json:
``` JSON
{
    "octopus":{
        "apikey": "<api key>",
        "electricity": {
            "mpan": "#############",
            "serial": "##########",
            "cost": 0.0000 
        },
        "gas": {
            "mprn": "##########",
            "serial": "#############",
            "cost": 0.0000
        }
    },
    "influx":{
        "token": "<influx token>",
        "bucket": "bukcet",
        "org":"org",
        "url": "http://example.org"
    }
}
```