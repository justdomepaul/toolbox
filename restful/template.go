package restful

const (
	ErrPageKey = "error"
)

var ErrorPageTmpl = `
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <title>{{ .Title }}</title>
</head>
<style>
    @import url('https://fonts.googleapis.com/css2?family=Montserrat:wght@500&family=Stick+No+Bills:wght@700&family=Abril+Fatface&display=swap');
    body {
        margin: 0;
    }
    .container {
        width: 100%;
        min-height: 100vh;
        display: flex;
        justify-content: center;
        align-items: center;
        flex-direction: column;
        background-color: #7d8492;
        position: relative;
        background-image: radial-gradient(40% 45%, rgba(255, 255, 255, 0.35), transparent 100%);
        background-position: center -50vh;
        background-repeat: no-repeat;
        overflow: hidden;
    }
    .container .errorCode {
        font-size: 10rem;
        line-height: 10rem;
        margin-bottom: 20px;
        font-family: 'Stick No Bills', sans-serif;
        box-shadow: 0 0 10px rgba(0, 0, 0, .5);
        z-index: 10;
        background-color: white;
        padding: 10px 20px 0;
        color: #3b4151;
    }
    .container .errorCode span {
        background: #ec008c;
        background: -webkit-linear-gradient(to right, #fc6767, #ec008c);
        background: linear-gradient(to bottom, #fc6767, #ec008c);
        background-clip: text;
        -webkit-background-clip: text;
        color: transparent;
        user-select: none;
    }

    .container .line {
        position: absolute;
        top: 30%;
        left: 0;
        width: 100%;
        height: 80px;
        background-color: #dcba0f;
        z-index: 1;
        transform: skew(0, 10deg);
        background-image: -webkit-repeating-linear-gradient(45deg, #fb3 0, #fb3 30px, #fff 30px, #fff 60px);
        box-shadow: 0 0 10px rgba(0, 0, 0, .2);
    }

    .container .line2 {
        position: absolute;
        left: 0;
        width: 100%;
        background-color: #dcba0f;
        z-index: 1;
        background-image: -webkit-repeating-linear-gradient(45deg, #fb3 0, #fb3 30px, #fff 30px, #fff 60px);
        box-shadow: 0 0 10px rgba(0, 0, 0, .2);
        transform: skew(0, -2deg);
        top: unset;
        bottom: 15%;
        height: 50px;
    }

    .container .line3 {
        position: absolute;
        left: 0;
        width: 100%;
        background-color: #dcba0f;
        z-index: 1;
        background-image: -webkit-repeating-linear-gradient(45deg, #fb3 0, #fb3 30px, #fff 30px, #fff 60px);
        box-shadow: 0 0 10px rgba(0, 0, 0, .2);
        transform: skew(0, -15deg);
        top: 25%;
        height: 50px;
    }

    .container .errorMsg {
        box-sizing: border-box;
        font-size: 5rem;
        line-height: 5rem;
        font-family: 'Stick No Bills', sans-serif;
        background-color: #ffda54;
        padding: 10px 40px;
        transform: skew(-10deg, -5deg);
        box-shadow: 0 0 10px rgba(0, 0, 0, .3);
        margin-top: -20px;
        color: #20242b;
        z-index: 9;
        max-width: 90%;
        text-align: center;
        user-select: none;
    }

    .container .footer {
        position: fixed;
        bottom: 0;
        display: flex;
        align-items: center;
        justify-content: center;
        padding-bottom: 5px;
        flex-wrap: wrap;
    }
    .footer img {
        width: 18px;
        height: 18px;
        margin: 0 5px;
    }

    @media screen and (max-width: 768px) {
        .errorMsg {
            font-size: 3rem;
            line-height: 3rem;
            padding: 10px 20px;
        }

        .line2 {
            bottom: 15%;
        }

        .line3 {
            top: 33%;
        }
    }
</style>
<body>
<div class="container">
    <div class="line"></div>
    <div class="line2"></div>
    <div class="line3"></div>
    <div class="errorCode"><span>{{ .Code }}</span></div>
    <div class="errorMsg">{{ .Title }}</div>
    <div class="footer">
        copyright © 2021
    </div>
</div>
</body>
</html>
`

type ErrPage struct {
	Title string
	Code  int
}
