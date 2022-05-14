# Babyfood Finder

> The year is 2022 and Baby Food has become the hard to find item.
> While the pandemic has hit *many* different products supply lines(2020 Toilet Paper shortage remembers), this one hits in a different way.
> 
> I was very fortunate to be not sweat the early shortages, but baby food hits in a slightly different way when you have a 9 month old daugther.
> This is my effort to make it slightly easier to find baby food while a pandemic is raging.

The goal of this project is to make it easier to find out when baby food is available to buy. 
It should optimize for making it easy to *find* stores with baby food.
It should *not* attempt to automate the purchase of the product.

Up Next:
- implement Target.com scraping
- ~~Configure a filter for search results~~
- ~~Container Application~~
- Run the application on a scheduled basis

## Running

Prerequisites:
- Copy `.env.sample` to `.env` and fill out Twilio account [settings](https://console.twilio.com/?frameUrl=%2Fconsole%3Fx-target-region%3Dus1)

```sh
docker build -t babyfood-finder:main .
docker run --env-file=.env -it babyfood-scraper:main -to +18675309
```