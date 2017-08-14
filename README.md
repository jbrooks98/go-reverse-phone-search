# go-reverse-phone-search
A web scraper for the site https://www.truepeoplesearch.com/ written in Go. 
The input from the user will be a phone number and a name.

# Basic requirements

* Do a "Reverse Phone" search by the phone number only (do not use the user-provided name)
* If multiple results are returned, use the user-provided name to determine the correct entry
* Display the user's full name and address

# Additional requirements

* Create a browser-based application in Go
* Any pre-existing packages (for example, from Github), may be used
* Store the results in a database, so future lookups are cached and do not require duplicate scraping
