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
        copyright © 2021 <img src="data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAAGUAAABkCAYAAACfIP5qAAAACXBIWXMAAAsSAAALEgHS3X78AAAE/klEQVR4nO2czVXcMBRGX6ggC++8COe4AVJBoAOoIEwFYSogqQBSAXQAqQBSQaYBn7DxzotJBzkizzBYeGz5fXrSkHeXtgdm/Em6+rPJMAzDMAzDMAzDMAzDMAzDMAzDMAzDeIu828Xf1JbVAREdE9EnInpPRAe9S+6JaE1EP4notmjqB++PZMzOhNKWlbv5Z0T0mYj2vQu2syKi70VTX2+9KhOyD2UjjC9cKyS4GvMt93CyDoWbqZsZNWMM17ydFE29HrkuCdmG0pbVKRFdeSdwuECOiqZeaf+2MbIMRSGQjiyDyS4UxUA6sgsmq1DYIXcAoYey4mCycMyedyQR3MtyUtcOhHicc+EdTUQ2oXCThe5lhXDKTWdysgilLaszHqGn5oKb0KQkd0pbVofskVxI7pekNWXDIzmR3C9Ja0pbVq6GHHonwnAl2k2b/OFPfXB+AHy9RarpmGShtGX1lYjOvRNhuJu27Dc1bVntcw2U+CHZ+CVJKCCPuCn5E+/o8/9wTeMvYY8uiV/UnQLyiLtZC+/oBnwjt14zgSR+SSF66QBxze39aOktmtrNBi+9E2GccpddDdXmC+SRYAG3ZXUDGAd91PKLWiggj1wXTR3cJIH88sDBRPeLSvMF9Mispohv5GCnYCL7WrPXWk5R88gQ3PRI/XKs4ZfozVdbVhe8xi4BNpDbBb9ErSltWR0DArkEj6wX7AcJN9wkRyFaKDyqlrbBq6KppU3OC3bBLzFrCsIj0pv3Krn7JUoo7BHpusQi5s7Goqkv3VSNdyKMKOsv8FCAHpHesCksuKst4Q7tF2goII8Q7wGOzsb8mGRACF8TQtcU1MaHKw44OiC/HPIUEgRYKCCPdKiuSHKXW9rtPuepJDGQUEAe6XPAQWuxBPgFMn4Rh8K9j1h99jMOPDo5+UUUCpeKq8gb6K5ijp43ycUv0pqC9MgQ/51fZofCuwm1dhRCezcTSOqXWaGwR7TXrmG9mzFS+yU4FCWPDBF1dnYToF+CC++cmqLhkSFS+OXSOxFGcA8yKBSQR+69I2Go+oWXDqR+CZqhmBwKyCP3RVMfAUqfml+YE02/TAoF5JGn9RFQ6dP0ywNiY99Uv0ytKQiP9B+RVi19UngpQcUvo6Hw6prUI0verfgEqPS9Sb9sDQXkkVte5fMAlb4355fBUEAb6EZrA7D0vRm/DIYCejB06qs2pKVPbfciKfjl1VBAD4Yup25YA5U+ld2LHcAa7nWgvB2SfNEv7+NhbH2gZwjQbkq13fHcZP4WDhW8B5Ne1BQtjwyxg+MXxN40rzPVb740PTKEdHZW2y+uq//NOxHGixcrPIXCB9U8MsQu7Y7vKJr6K2BO76Ibv2zWFOkTVoPjkVCAq3+aywuI8Yt7+9+/UFju0qecpL2nPlK/vAc8oz8ZkF8ev29XU8QTjejHzkCrf153MyYAvzx+31fHKYGIPTIEyC+qIPzShTJ3d/t17FdlgPyizVy/PIa5R88j6tDSrlmK5/pFY+e+h8AvP6jXfIW0heIHQ0OY6Zf7lO+FnOGX7sU/z6HwJNvUZiKaR4YI9Ms6BxexX6bc0/XmVMsL0fOLA5ZbSuQDfzhJG8//d6zrndXbU/mebqsxq/739SYk6XkO7Lg3dlkpPV01Co98z7lf333HFbfJlzm+mTv3e2oYhmEYhmEYhmEYhmEYhmEYhmEYhmHEhYj+AsaJR2j2SmV+AAAAAElFTkSuQmCC" /> Flashaim Inc. (7551)
    </div>
</div>
</body>
</html>
`

type ErrPage struct {
	Title string
	Code  int
}
