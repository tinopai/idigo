# idigo
Idi**go**, yet another Minecraft screensharing tool made in Golang

### Disclaimer
⚠️ This is **__NOT__** production ready, you should add more features and fix some security issues like the (really old and bad) strings being on plaintext, you could use either AES-256-CBC encryption or decryption or just make a new PHP file with those which doesn't let you print the text and just do a require() from **pin.php** removing every <?php ?> thingy. **This uses [strings2](https://github.com/glmcdona/strings2), you should read [it's license](https://github.com/glmcdona/strings2/blob/master/license.txt), aswell as every other dependency one, like AES-Everywhere or Astilectron**

### Installation
##### Golang
1. `go get github.com/mervick/aes-everywhere/go/aes256` - Thanks [mervick and every contributor](https://github.com/mervick/aes-everywhere) 💖
2. `go get -u github.com/asticode/go-astilectron` - Thanks [electron](https://www.electronjs.org/), [asticode and every contributor](https://github.com/asticode/go-astilectron) 💖
3. `go get golang.org/x/sys/windows`
4. Modify every variable it says you should change (on `app.go`)
4. `go build -ldflags -H=windowsgui app.go` (You can use `go build app.go` if you want to show the console)
##### PHP
1. [Install PHP](https://www.php.net/downloads)
2. Move every file in the `/server/` folder to your htdocs (or wherever you place your website files)
3. Modify every variable it says you should change (on `config.php`)
##### Node.js
1. [Install Node.js](https://nodejs.org/en/download/)
2. `npm i`
3. Modify `config.json` to your **MySQL** account details and every variable it says you should change (on `index.js`)
4. `node index.js` (you should make a new screen if using Linux, or just use [forever](https://www.npmjs.com/package/forever))
##### MySQL
1. Install [MySQL](https://dev.mysql.com/downloads/)
2. Create a new database named `idigo`
3. Create a new table named `pin` using [this columns](https://gist.github.com/tinopai/01fa90e57bdb0cb4412efa63d43896a7#file-idigo-mysql-L14-L23)
4. Create a new table named `users` using [this columns](https://gist.github.com/tinopai/01fa90e57bdb0cb4412efa63d43896a7#file-idigo-mysql-L2-L10)
