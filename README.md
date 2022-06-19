<div  align="center">
	<h1 align="center">Stock Game</h1>
	<h3 align="center">A stock market game with real prices.</h3>
</div>

## About the project
Trade stocks and cryptocurrencies, whether you are looking to become better at investing money or just have fun. In Stock Game, your objective is to make as much money as you can, by buying and selling assets at real-world prices.

The app features:
* Information about historical value of assets.
* An account system, so you can play on whatever device you are using at the moment.
* A view of all your assets with related information in one place.
* session maintained with cookies

## Built with
* Go
* Go Fiber
* HTML
* CSS
* JavaScript
* Bootstrap 5
* MongoDB

## UML Diagram
<img width="1404" alt="image" src="https://user-images.githubusercontent.com/49531832/173229995-40bc80ab-26d2-401c-b679-e4d151755ba7.png">

## Running the app
Prerequisites:
* docker
* pulled project
* environment variables file `.env` in project root directory
To run the project run following commands in project root directory:
```
docker build -t "stock_game" .
```
```
docker run -p 3000:3000 --env-file ".env" "stock_game"
```

## Testing
In project root directory run:
```
go test ./tests
```
