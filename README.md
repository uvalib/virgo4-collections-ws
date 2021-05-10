# Virgo4 Collections Context Service

This is a web service provides collections context servives to the Virgo4 Client

Requires Go 1.16.0+

### Current API

* GET /version - get version information
* GET /collections/:name - get collection context for a named collection. Ex: Daily Progress Digitized Microfilm
* GET /collections/:name/dates?year=YYYY - get publication dates for a collection year
* GET /collections/:name/items/:date/next - get the next published item; date format=yyyy-mm-dd
* GET /collections/:name/items/:date/previous - get the previous published item; date format=yyyy-mm-dd
